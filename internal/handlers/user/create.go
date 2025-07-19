package user_handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/models"
	"pinstack-api-gateway/internal/utils"

	"github.com/go-playground/validator/v10"
)

type CreateUserRequest struct {
	Username  string  `json:"username" validate:"required,min=3,max=32"`
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=6"`
	FullName  *string `json:"full_name,omitempty" validate:"omitempty,max=100"`
	Bio       *string `json:"bio,omitempty" validate:"omitempty,max=500"`
	AvatarURL *string `json:"avatar_url,omitempty" validate:"omitempty,url"`
}

type CreateUserResponse struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	FullName  *string `json:"full_name,omitempty"`
	Bio       *string `json:"bio,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with provided data
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateUserRequest true "User creation data"
// @Success 201 {object} CreateUserResponse "User created successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 409 {object} map[string]string "User already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
		return
	}

	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		FullName:  req.FullName,
		Bio:       req.Bio,
		AvatarURL: req.AvatarURL,
	}

	createdUser, err := h.userClient.CreateUser(r.Context(), user)
	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrUsernameExists):
			utils.SendError(w, http.StatusConflict, err.Error())
		case errors.Is(err, custom_errors.ErrEmailExists):
			utils.SendError(w, http.StatusConflict, err.Error())
		case errors.Is(err, custom_errors.ErrInvalidUsername):
			utils.SendError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, custom_errors.ErrInvalidEmail):
			utils.SendError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, custom_errors.ErrInvalidPassword):
			utils.SendError(w, http.StatusBadRequest, err.Error())
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		}
		return
	}

	response := CreateUserResponse{
		ID:        createdUser.ID,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		FullName:  createdUser.FullName,
		Bio:       createdUser.Bio,
		AvatarURL: createdUser.AvatarURL,
		CreatedAt: createdUser.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: createdUser.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	utils.Send(w, http.StatusCreated, response)
}
