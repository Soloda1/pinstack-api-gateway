package post_handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/middlewares"
	"pinstack-api-gateway/internal/models"
	"pinstack-api-gateway/internal/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type UpdatePostRequest struct {
	Title      *string           `json:"title,omitempty" validate:"omitempty"`
	Content    *string           `json:"content,omitempty"`
	Tags       []string          `json:"tags,omitempty"`
	MediaItems []*MediaItemInput `json:"media_items,omitempty" validate:"max=9,dive"`
}

type UpdatePostResponse struct {
	ID        int64               `json:"id"`
	Title     string              `json:"title"`
	Content   *string             `json:"content,omitempty"`
	CreatedAt string              `json:"created_at"`
	UpdatedAt string              `json:"updated_at"`
	Author    *UpdatePostAuthor   `json:"author,omitempty"`
	Media     []PostMediaResponse `json:"media,omitempty"`
	Tags      []TagResponse       `json:"tags,omitempty"`
}

type UpdatePostAuthor struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	FullName  *string `json:"full_name,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// Update godoc
// @Summary Update a post
// @Description Update an existing post with new data
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Post ID"
// @Param request body UpdatePostRequest true "Post update data"
// @Success 200 {object} UpdatePostResponse "Post updated successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Post not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts/{id} [put]
func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	claims, err := middlewares.GetClaimsFromContext(r.Context())
	if err != nil {
		h.log.Debug("No user claims in context", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrUnauthenticated.Error())
		return
	}

	var req UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Debug("Failed to decode update post request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		h.log.Debug("Failed to validate update post request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
		return
	}

	modelReq := &models.UpdatePostDTO{
		UserID:  claims.UserID,
		Title:   req.Title,
		Content: req.Content,
		Tags:    req.Tags,
	}
	if len(req.MediaItems) > 0 {
		modelReq.MediaItems = make([]*models.PostMediaInput, len(req.MediaItems))
		for i, item := range req.MediaItems {
			modelReq.MediaItems[i] = &models.PostMediaInput{
				URL:      item.URL,
				Type:     models.MediaType(item.Type),
				Position: item.Position,
			}
		}
	}

	err = h.postClient.UpdatePost(r.Context(), id, modelReq)
	if err != nil {
		h.log.Error("Update post failed", slog.String("error", err.Error()))
		switch {
		case errors.Is(err, custom_errors.ErrPostNotFound):
			h.log.Debug("Post not found", slog.String("error", err.Error()))
			utils.SendError(w, http.StatusNotFound, custom_errors.ErrPostNotFound.Error())
			return
		case errors.Is(err, custom_errors.ErrForbidden):
			h.log.Debug("Forbidden", slog.String("error", err.Error()))
			utils.SendError(w, http.StatusForbidden, custom_errors.ErrForbidden.Error())
			return
		case errors.Is(err, custom_errors.ErrPostValidation):
			h.log.Debug("Post validation failed", slog.String("error", err.Error()))
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrPostValidation.Error())
			return
		default:
			h.log.Error("Update post failed", slog.String("error", err.Error()))
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
			return
		}
	}

	updatedPost, err := h.postClient.GetPostByID(r.Context(), id)
	if err != nil {
		h.log.Error("get post after update failed", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		return
	}

	resp := UpdatePostResponse{
		ID:        updatedPost.Post.ID,
		Title:     updatedPost.Post.Title,
		Content:   updatedPost.Post.Content,
		CreatedAt: updatedPost.Post.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: updatedPost.Post.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	author, err := h.userClient.GetUser(r.Context(), updatedPost.Post.AuthorID)
	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrUserNotFound):
			h.log.Warn("Author not found, setting author to null", slog.Int64("authorID", updatedPost.Post.AuthorID))
			resp.Author = &UpdatePostAuthor{
				ID:        0,
				Username:  "unknown",
				FullName:  utils.StringPtr("Unknown Author"),
				AvatarURL: utils.StringPtr("http://unknown.unknown"),
			}
		default:
			h.log.Error("Failed to get user", slog.Int64("id", updatedPost.Post.AuthorID), slog.String("error", err.Error()))
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
			return
		}
	}

	if author != nil {
		resp.Author = &UpdatePostAuthor{
			ID:        author.ID,
			Username:  author.Username,
			FullName:  author.FullName,
			AvatarURL: author.AvatarURL,
		}
	}

	if len(updatedPost.Media) > 0 {
		resp.Media = make([]PostMediaResponse, len(updatedPost.Media))
		for i, m := range updatedPost.Media {
			resp.Media[i] = PostMediaResponse{
				ID:       m.ID,
				URL:      m.URL,
				Type:     string(m.Type),
				Position: m.Position,
			}
		}
	}
	if len(updatedPost.Tags) > 0 {
		resp.Tags = make([]TagResponse, len(updatedPost.Tags))
		for i, t := range updatedPost.Tags {
			resp.Tags[i] = TagResponse{
				ID:   t.ID,
				Name: t.Name,
			}
		}
	}
	utils.Send(w, http.StatusOK, resp)
}
