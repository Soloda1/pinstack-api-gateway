package post_handler

import (
	"errors"
	"github.com/soloda1/pinstack-proto-definitions/custom_errors"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type GetPostRequest struct {
	ID int64 `json:"id"`
}

type GetPostResponse struct {
	ID        int64           `json:"id"`
	Author    *GetPostUser    `json:"author,omitempty"`
	Title     string          `json:"title"`
	Content   *string         `json:"content,omitempty"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
	Media     []*GetPostMedia `json:"media,omitempty"`
	Tags      []*GetPostTag   `json:"tags,omitempty"`
}

type GetPostUser struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	FullName  *string `json:"full_name,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type GetPostMedia struct {
	ID       int64  `json:"id"`
	URL      string `json:"url"`
	Type     string `json:"type"`
	Position int32  `json:"position"`
}

type GetPostTag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Get godoc
// @Summary Get post by ID
// @Description Get detailed information about a specific post
// @Tags posts
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} GetPostResponse "Post information"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Post not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts/{id} [get]
func (h *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		h.log.Debug("Missing post id in path params")
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Debug("Invalid post id format", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	post, err := h.postClient.GetPostByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrPostNotFound):
			h.log.Debug("get post failed", slog.String("error", err.Error()))
			utils.SendError(w, http.StatusNotFound, custom_errors.ErrPostNotFound.Error())
			return
		case errors.Is(err, custom_errors.ErrPostValidation):
			h.log.Debug("post validation failed", slog.String("error", err.Error()))
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrPostValidation.Error())
			return
		default:
			h.log.Error("Failed to get post", slog.Int64("id", id), slog.String("error", err.Error()))
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
			return
		}

	}

	resp := &GetPostResponse{
		ID:        post.Post.ID,
		Title:     post.Post.Title,
		Content:   post.Post.Content,
		CreatedAt: post.Post.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: post.Post.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	author, err := h.userClient.GetUser(r.Context(), post.Post.AuthorID)
	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrUserNotFound):
			h.log.Warn("Author not found for post", slog.Int64("author_id", post.Post.AuthorID))
			author = utils.GenerateUnknownAuthor()
		default:
			h.log.Error("Failed to get user", slog.Int64("id", post.Post.AuthorID), slog.String("error", err.Error()))
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
			return
		}
	}

	resp.Author = &GetPostUser{
		ID:        author.ID,
		Username:  author.Username,
		FullName:  author.FullName,
		AvatarURL: author.AvatarURL,
	}

	if post.Media != nil {
		media := make([]*GetPostMedia, 0, len(post.Media))
		for _, m := range post.Media {
			media = append(media, &GetPostMedia{
				ID:       m.ID,
				URL:      m.URL,
				Type:     string(m.Type),
				Position: m.Position,
			})
		}
		resp.Media = media
	}
	if post.Tags != nil {
		tags := make([]*GetPostTag, 0, len(post.Tags))
		for _, t := range post.Tags {
			tags = append(tags, &GetPostTag{
				ID:   t.ID,
				Name: t.Name,
			})
		}
		resp.Tags = tags
	}
	utils.Send(w, http.StatusOK, resp)
}
