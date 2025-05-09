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

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Debug("Failed to decode login request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		h.log.Debug("Failed to validate login request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
		return
	}

	tokens, err := h.authClient.Login(r.Context(), &req)
	if err != nil {
		h.log.Error("login failed", slog.String("error", err.Error()))
		switch err {
		case custom_errors.ErrInvalidCredentials:
			utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrInvalidCredentials.Error())
		case custom_errors.ErrUserNotFound:
			utils.SendError(w, http.StatusNotFound, custom_errors.ErrUserNotFound.Error())
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		}
		return
	}

	utils.Send(w, http.StatusOK, tokens)
}
