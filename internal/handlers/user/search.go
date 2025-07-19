package user_handler

import (
	"errors"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/utils"
	"strconv"
)

type SearchUsersResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
}

type UserResponse struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	FullName  *string `json:"full_name,omitempty"`
	Bio       *string `json:"bio,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// SearchUsers godoc
// @Summary Search users
// @Description Search users by query string
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param query query string true "Search query"
// @Param page query int false "Page number"
// @Param limit query int false "Results per page (max 100)"
// @Success 200 {object} SearchUsersResponse "Search results"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/search [get]
func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidSearchQuery.Error())
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := h.userClient.SearchUsers(r.Context(), query, page, limit)
	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrInvalidSearchQuery):
			utils.SendError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, custom_errors.ErrSearchFailed):
			utils.SendError(w, http.StatusInternalServerError, err.Error())
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		}
		return
	}

	response := SearchUsersResponse{
		Users: make([]UserResponse, 0, len(users)),
		Total: total,
	}

	for _, user := range users {
		response.Users = append(response.Users, UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			Bio:       user.Bio,
			AvatarURL: user.AvatarURL,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	utils.Send(w, http.StatusOK, response)
}
