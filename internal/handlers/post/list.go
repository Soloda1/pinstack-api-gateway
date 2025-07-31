package post_handler

import (
	"errors"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/models"
	"pinstack-api-gateway/internal/utils"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListPostsResponse struct {
	Posts []ListPostItem `json:"posts"`
	Total int64          `json:"total"`
}

type ListPostItem struct {
	ID        int64               `json:"id"`
	Title     string              `json:"title"`
	Content   *string             `json:"content,omitempty"`
	CreatedAt string              `json:"created_at"`
	UpdatedAt string              `json:"updated_at"`
	Author    *ListPostAuthor     `json:"author,omitempty"`
	Media     []PostMediaResponse `json:"media,omitempty"`
	Tags      []TagResponse       `json:"tags,omitempty"`
}

type ListPostAuthor struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	FullName  *string `json:"full_name,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// List godoc
// @Summary List posts with filters
// @Description Get a list of posts with optional filtering by author and date range
// @Tags posts
// @Accept json
// @Produce json
// @Param author_id query int false "Filter by author ID"
// @Param created_after query string false "Filter posts created after this time (RFC3339 format)"
// @Param created_before query string false "Filter posts created before this time (RFC3339 format)"
// @Param offset query int false "Pagination offset"
// @Param limit query int false "Pagination limit"
// @Success 200 {object} ListPostsResponse "List of posts"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts/list [get]
func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	h.log.Debug("query info", slog.Any("query", query))
	var authorID *int64
	if authorIDStr := query.Get("author_id"); authorIDStr != "" {
		id, err := strconv.ParseInt(authorIDStr, 10, 64)
		if err != nil {
			h.log.Debug("Invalid author_id format", slog.String("error", err.Error()))
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
			return
		}
		authorID = &id
	}

	var createdAfterTime *time.Time
	if createdAfterStr := query.Get("created_after"); createdAfterStr != "" {
		t, err := time.Parse(time.RFC3339, createdAfterStr)
		if err != nil {
			h.log.Debug("Invalid created_after format", slog.String("error", err.Error()))
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
			return
		}
		createdAfterTime = &t
	}

	var createdBeforeTime *time.Time
	if createdBeforeStr := query.Get("created_before"); createdBeforeStr != "" {
		t, err := time.Parse(time.RFC3339, createdBeforeStr)
		if err != nil {
			h.log.Debug("Invalid created_before format", slog.String("error", err.Error()))
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
			return
		}
		createdBeforeTime = &t
	}

	var offset *int
	if offsetStr := query.Get("offset"); offsetStr != "" {
		offsetVal, err := strconv.Atoi(offsetStr)
		if err != nil {
			h.log.Debug("Invalid offset format", slog.String("error", err.Error()))
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
			return
		}
		offset = &offsetVal
	}

	var limit *int
	if limitStr := query.Get("limit"); limitStr != "" {
		limitVal, err := strconv.Atoi(limitStr)
		if err != nil {
			h.log.Debug("Invalid limit format", slog.String("error", err.Error()))
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
			return
		}
		limit = &limitVal
	}

	filters := models.PostFilters{}
	if authorID != nil {
		filters.AuthorID = authorID
	}
	if createdAfterTime != nil {
		filters.CreatedAfter = createdAfterTime
	}
	if createdBeforeTime != nil {
		filters.CreatedBefore = createdBeforeTime
	}
	if offset != nil {
		filters.Offset = offset
	}
	if limit != nil {
		filters.Limit = limit
	}

	posts, total, err := h.postClient.ListPosts(r.Context(), &filters)
	if err != nil {
		h.log.Error("list posts failed", slog.String("error", err.Error()))
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
				return
			case codes.Internal:
				utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
				return
			}
		}
		utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		return
	}

	resp := ListPostsResponse{
		Posts: make([]ListPostItem, len(posts)),
		Total: total,
	}
	for i, p := range posts {
		item := ListPostItem{
			ID:        p.Post.ID,
			Title:     p.Post.Title,
			Content:   p.Post.Content,
			CreatedAt: p.Post.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: p.Post.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		author, err := h.userClient.GetUser(r.Context(), p.Post.AuthorID)
		if err != nil {
			switch {
			case errors.Is(err, custom_errors.ErrUserNotFound):
				h.log.Warn("author not found, using placeholder", slog.Int64("authorID", p.Post.AuthorID))
				author = &models.User{
					ID:        0,
					Username:  "unknown",
					FullName:  utils.StringPtr("Unknown Author"),
					AvatarURL: utils.StringPtr("http://unknown.unknown"),
				}
			default:
				h.log.Error("Failed to get user", slog.Int64("id", p.Post.AuthorID), slog.String("error", err.Error()))
				utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
				return
			}
		}

		item.Author = &ListPostAuthor{
			ID:        author.ID,
			Username:  author.Username,
			FullName:  author.FullName,
			AvatarURL: author.AvatarURL,
		}

		if len(p.Media) > 0 {
			item.Media = make([]PostMediaResponse, len(p.Media))
			for j, m := range p.Media {
				item.Media[j] = PostMediaResponse{
					ID:       m.ID,
					URL:      m.URL,
					Type:     string(m.Type),
					Position: m.Position,
				}
			}
		}

		if len(p.Tags) > 0 {
			item.Tags = make([]TagResponse, len(p.Tags))
			for k, t := range p.Tags {
				item.Tags[k] = TagResponse{
					ID:   t.ID,
					Name: t.Name,
				}
			}
		}
		resp.Posts[i] = item
	}
	utils.Send(w, http.StatusOK, resp)
}
