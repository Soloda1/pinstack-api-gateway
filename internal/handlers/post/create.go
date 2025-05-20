package post_handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/middlewares"
	"pinstack-api-gateway/internal/models"
	"pinstack-api-gateway/internal/utils"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreatePostRequest struct {
	Title      string            `json:"title" validate:"required,min=1,max=255"`
	Content    *string           `json:"content,omitempty"`
	Tags       []string          `json:"tags,omitempty"`
	MediaItems []*MediaItemInput `json:"media_items,omitempty" validate:"max=9,dive"`
}

type MediaItemInput struct {
	URL      string `json:"url" validate:"required,url,max=512"`
	Type     string `json:"type" validate:"required,oneof=image video"`
	Position int32  `json:"position" validate:"gte=0,lte=100"`
}

type CreatePostResponse struct {
	ID              int64               `json:"id"`
	Title           string              `json:"title"`
	Content         *string             `json:"content,omitempty"`
	CreatedAt       string              `json:"created_at"`
	UpdatedAt       string              `json:"updated_at"`
	AuthorID        int64               `json:"author_id"`
	AuthorUsername  string              `json:"author_username"`
	AuthorEmail     string              `json:"author_email"`
	AuthorFullName  *string             `json:"author_full_name,omitempty"`
	AuthorBio       *string             `json:"author_bio,omitempty"`
	AuthorAvatarURL *string             `json:"author_avatar_url,omitempty"`
	Media           []PostMediaResponse `json:"media,omitempty"`
	Tags            []TagResponse       `json:"tags,omitempty"`
}

type PostMediaResponse struct {
	ID       int64  `json:"id"`
	URL      string `json:"url"`
	Type     string `json:"type"`
	Position int32  `json:"position"`
}

type TagResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Debug("Failed to decode create post request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	claimsRaw := r.Context().Value("claims")
	claims, ok := claimsRaw.(*middlewares.Claims)
	if !ok || claims == nil {
		h.log.Error("invalid token claims")
		utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrInvalidToken.Error())
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		h.log.Debug("Failed to validate create post request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
		return
	}

	modelReq := &models.CreatePostDTO{
		AuthorID: claims.UserID,
		Title:    req.Title,
		Content:  req.Content,
		Tags:     req.Tags,
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

	post, err := h.postClient.CreatePost(r.Context(), modelReq)
	if err != nil {
		h.log.Error("create post failed", slog.String("error", err.Error()))

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

	resp := CreatePostResponse{
		ID:              post.Post.ID,
		Title:           post.Post.Title,
		Content:         post.Post.Content,
		CreatedAt:       post.Post.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       post.Post.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		AuthorID:        post.Author.ID,
		AuthorUsername:  post.Author.Username,
		AuthorEmail:     post.Author.Email,
		AuthorFullName:  post.Author.FullName,
		AuthorBio:       post.Author.Bio,
		AuthorAvatarURL: post.Author.AvatarURL,
	}
	if len(post.Media) > 0 {
		resp.Media = make([]PostMediaResponse, len(post.Media))
		for i, m := range post.Media {
			resp.Media[i] = PostMediaResponse{
				ID:       m.ID,
				URL:      m.URL,
				Type:     string(m.Type),
				Position: m.Position,
			}
		}
	}
	if len(post.Tags) > 0 {
		resp.Tags = make([]TagResponse, len(post.Tags))
		for i, t := range post.Tags {
			resp.Tags[i] = TagResponse{
				ID:   t.ID,
				Name: t.Name,
			}
		}
	}
	utils.Send(w, http.StatusCreated, resp)
}
