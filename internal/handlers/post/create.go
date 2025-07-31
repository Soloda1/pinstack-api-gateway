package post_handler

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/middlewares"
	"pinstack-api-gateway/internal/models"
	"pinstack-api-gateway/internal/utils"
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
	Position int32  `json:"position" validate:"gte=1,lte=9"`
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

// Create godoc
// @Summary Create a new post
// @Description Create a new post with title, content, tags and media
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreatePostRequest true "Post creation data"
// @Success 201 {object} CreatePostResponse "Post created successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts [post]
func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Debug("Failed to decode create post request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	claims, err := middlewares.GetClaimsFromContext(r.Context())
	if err != nil {
		h.log.Debug("No user claims in context", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrUnauthenticated.Error())
		return
	}
	h.log.Debug("requested model", slog.Any("model", req))

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
		modelReq.MediaItems = make([]*models.PostMediaInput, 0, len(req.MediaItems))
		for _, item := range req.MediaItems {
			modelReq.MediaItems = append(modelReq.MediaItems, &models.PostMediaInput{
				URL:      item.URL,
				Type:     models.MediaType(item.Type),
				Position: item.Position,
			})
		}
	}

	h.log.Debug("creating post", slog.Any("model", modelReq))

	post, err := h.postClient.CreatePost(r.Context(), modelReq)
	if err != nil {
		h.log.Error("create post failed", slog.String("error", err.Error()))

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.SendError(w, http.StatusBadRequest, "invalid input data")
				return
			case codes.Unauthenticated:
				utils.SendError(w, http.StatusUnauthorized, "unauthenticated")
				return
			case codes.AlreadyExists:
				utils.SendError(w, http.StatusConflict, "post already exists")
				return
			case codes.PermissionDenied:
				utils.SendError(w, http.StatusForbidden, "access denied")
				return
			case codes.Unavailable:
				utils.SendError(w, http.StatusServiceUnavailable, "service unavailable")
				return
			}
		}

		utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		return
	}
	h.log.Debug("created post", slog.Any("post", post))
	author, err := h.userClient.GetUser(r.Context(), post.Post.AuthorID)
	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrUserNotFound):
			h.log.Warn("author not found, using placeholder data", slog.Int64("authorID", post.Post.AuthorID))
			author = utils.GenerateUnknownAuthor()
		default:
			h.log.Error("Failed to get user", slog.Int64("id", post.Post.AuthorID), slog.String("error", err.Error()))
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
			return
		}
	}
	resp := CreatePostResponse{
		ID:        post.Post.ID,
		Title:     post.Post.Title,
		Content:   post.Post.Content,
		CreatedAt: post.Post.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: post.Post.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		AuthorID:  post.Post.AuthorID,
	}

	resp.AuthorEmail = author.Email
	resp.AuthorAvatarURL = author.AvatarURL
	resp.AuthorBio = author.Bio
	resp.AuthorFullName = author.FullName
	resp.AuthorUsername = author.Username

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
