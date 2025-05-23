package auth_handler

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/models"
	"pinstack-api-gateway/internal/utils"
)

type RegisterRequest struct {
	Username  string  `json:"username" validate:"required,min=3,max=32"`
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=6"`
	FullName  *string `json:"full_name,omitempty"`
	Bio       *string `json:"bio,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type RegisterResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Debug("Failed to decode register request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		h.log.Debug("Failed to validate register request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
		return
	}

	// Convert to models struct
	modelReq := &models.RegisterRequest{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		FullName:  req.FullName,
		Bio:       req.Bio,
		AvatarURL: req.AvatarURL,
	}

	tokens, err := h.authClient.Register(r.Context(), modelReq)
	if err != nil {
		h.log.Error("register failed", slog.String("error", err.Error()))
		switch err {
		case custom_errors.ErrUsernameExists:
			utils.SendError(w, http.StatusConflict, custom_errors.ErrUsernameExists.Error())
		case custom_errors.ErrEmailExists:
			utils.SendError(w, http.StatusConflict, custom_errors.ErrEmailExists.Error())
		case custom_errors.ErrInvalidUsername:
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidUsername.Error())
		case custom_errors.ErrInvalidEmail:
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidEmail.Error())
		case custom_errors.ErrInvalidPassword:
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidPassword.Error())
		case custom_errors.ErrUserAlreadyExists:
			utils.SendError(w, http.StatusConflict, custom_errors.ErrUserAlreadyExists.Error())
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		}
		return
	}

	response := RegisterResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	utils.Send(w, http.StatusCreated, response)
}
