package post_handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/models"
	"pinstack-api-gateway/internal/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListPostsRequest struct {
	AuthorID      *int64   `json:"author_id,omitempty"`
	TagNames      []string `json:"tag_names,omitempty"`
	CreatedAfter  *string  `json:"created_after,omitempty"`
	CreatedBefore *string  `json:"created_before,omitempty"`
	Offset        *int     `json:"offset,omitempty"`
	Limit         *int     `json:"limit,omitempty"`
}

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
// @Description Get a list of posts with optional filtering by author, tags, and date range
// @Tags posts
// @Accept json
// @Produce json
// @Param request body ListPostsRequest false "Filter parameters"
// @Success 200 {object} ListPostsResponse "List of posts"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts/list [post]
func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	var req ListPostsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Debug("Failed to decode list posts request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		h.log.Debug("Failed to validate list posts request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
		return
	}

	var createdAfterTime *time.Time
	if req.CreatedAfter != nil {
		t, err := time.Parse(time.RFC3339, *req.CreatedAfter)
		if err != nil {
			h.log.Debug("Invalid created_after format", slog.String("error", err.Error()))
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
			return
		}
		createdAfterTime = &t
	}
	var createdBeforeTime *time.Time
	if req.CreatedBefore != nil {
		t, err := time.Parse(time.RFC3339, *req.CreatedBefore)
		if err != nil {
			h.log.Debug("Invalid created_before format", slog.String("error", err.Error()))
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
			return
		}
		createdBeforeTime = &t
	}

	filters := models.PostFilters{}
	if req.AuthorID != nil {
		filters.AuthorID = req.AuthorID
	}
	if req.TagNames != nil {
		filters.TagNames = req.TagNames
	}
	if createdAfterTime != nil {
		filters.CreatedAfter = createdAfterTime
	}
	if createdBeforeTime != nil {
		filters.CreatedBefore = createdBeforeTime
	}
	if req.Offset != nil {
		filters.Offset = req.Offset
	}
	if req.Limit != nil {
		filters.Limit = req.Limit
	}

	posts, err := h.postClient.ListPosts(r.Context(), &filters)
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
		Total: int64(len(posts)),
	}
	for i, p := range posts {
		item := ListPostItem{
			ID:        p.Post.ID,
			Title:     p.Post.Title,
			Content:   p.Post.Content,
			CreatedAt: p.Post.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: p.Post.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		if p.Author != nil {
			item.Author = &ListPostAuthor{
				ID:        p.Author.ID,
				Username:  p.Author.Username,
				FullName:  p.Author.FullName,
				AvatarURL: p.Author.AvatarURL,
			}
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
