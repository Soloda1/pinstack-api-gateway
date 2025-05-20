package post_handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/middlewares"
	"pinstack-api-gateway/internal/models"
	"pinstack-api-gateway/internal/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdatePostRequest struct {
	Title      *string           `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
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

func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.log.Debug("Missing post id in query params")
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Debug("Invalid post id format", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	claimsRaw := r.Context().Value("claims")
	claims, ok := claimsRaw.(*middlewares.Claims)
	if !ok {
		h.log.Debug("No claims found in context")
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
		h.log.Error("update post failed", slog.String("error", err.Error()))

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
				return
			case codes.NotFound:
				utils.SendError(w, http.StatusNotFound, custom_errors.ErrPostNotFound.Error())
				return
			case codes.PermissionDenied:
				utils.SendError(w, http.StatusForbidden, custom_errors.ErrForbidden.Error())
				return
			case codes.Internal:
				utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
				return
			}
		}

		utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		return
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
	if updatedPost.Author != nil {
		resp.Author = &UpdatePostAuthor{
			ID:        updatedPost.Author.ID,
			Username:  updatedPost.Author.Username,
			FullName:  updatedPost.Author.FullName,
			AvatarURL: updatedPost.Author.AvatarURL,
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
