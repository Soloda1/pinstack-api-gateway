package user_handler

import (
	"encoding/json"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/models"
	"pinstack-api-gateway/internal/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type UpdateUserRequest struct {
	ID       int64   `json:"id" validate:"required"`
	Username *string `json:"username,omitempty" validate:"omitempty,min=3,max=32"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	FullName *string `json:"full_name,omitempty" validate:"omitempty,max=100"`
	Bio      *string `json:"bio,omitempty" validate:"omitempty,max=500"`
}

type UpdateUserResponse struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	FullName  *string `json:"full_name,omitempty"`
	Bio       *string `json:"bio,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// UpdateUser godoc
// @Summary Update user information
// @Description Update user fields by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body UpdateUserRequest true "User update data"
// @Success 200 {object} UpdateUserResponse "User updated successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 409 {object} map[string]string "Username or email already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	currentUser, err := h.userClient.GetUser(r.Context(), id)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, custom_errors.ErrUserNotFound.Error())
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	req.ID = id

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
		return
	}

	updateUser := &models.User{
		ID:       id,
		Username: currentUser.Username,
		Email:    currentUser.Email,
		FullName: currentUser.FullName,
		Bio:      currentUser.Bio,
	}

	if req.Username != nil {
		updateUser.Username = *req.Username
	}
	if req.Email != nil {
		updateUser.Email = *req.Email
	}
	if req.FullName != nil {
		updateUser.FullName = req.FullName
	}
	if req.Bio != nil {
		updateUser.Bio = req.Bio
	}

	updatedUser, err := h.userClient.UpdateUser(r.Context(), updateUser)
	if err != nil {
		switch err {
		case custom_errors.ErrUsernameExists:
			utils.SendError(w, http.StatusConflict, err.Error())
		case custom_errors.ErrEmailExists:
			utils.SendError(w, http.StatusConflict, err.Error())
		case custom_errors.ErrInvalidUsername:
			utils.SendError(w, http.StatusBadRequest, err.Error())
		case custom_errors.ErrInvalidEmail:
			utils.SendError(w, http.StatusBadRequest, err.Error())
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		}
		return
	}

	response := UpdateUserResponse{
		ID:        updatedUser.ID,
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		FullName:  updatedUser.FullName,
		Bio:       updatedUser.Bio,
		AvatarURL: updatedUser.AvatarURL,
		CreatedAt: updatedUser.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: updatedUser.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	utils.Send(w, http.StatusOK, response)
}
