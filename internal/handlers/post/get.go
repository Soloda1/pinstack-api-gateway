package post_handler

import (
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/utils"
	"strconv"
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

func (h *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
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
	post, err := h.postClient.GetPostByID(r.Context(), id)
	if err != nil {
		h.log.Error("get post failed", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusNotFound, custom_errors.ErrPostNotFound.Error())
		return
	}
	resp := &GetPostResponse{
		ID:        post.Post.ID,
		Title:     post.Post.Title,
		Content:   post.Post.Content,
		CreatedAt: post.Post.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: post.Post.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if post.Author != nil {
		resp.Author = &GetPostUser{
			ID:        post.Author.ID,
			Username:  post.Author.Username,
			FullName:  post.Author.FullName,
			AvatarURL: post.Author.AvatarURL,
		}
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
