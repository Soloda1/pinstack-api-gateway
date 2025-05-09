package auth_handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/middlewares"
	"pinstack-api-gateway/internal/models"
	"pinstack-api-gateway/internal/utils"

	"github.com/go-playground/validator/v10"
)

func (h *AuthHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var req models.UpdatePasswordRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Debug("Failed to decode update password request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		h.log.Debug("Failed to validate update password request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
		return
	}

	claims, err := middlewares.GetClaimsFromContext(r.Context())
	if err != nil {
		h.log.Error("failed to get claims from context", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrUnauthenticated.Error())
		return
	}

	req.ID = claims.UserID

	err = h.authClient.UpdatePassword(r.Context(), &req)
	if err != nil {
		h.log.Error("update password failed", slog.String("error", err.Error()))
		switch err {
		case custom_errors.ErrInvalidPassword:
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidPassword.Error())
		case custom_errors.ErrInvalidCredentials:
			utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrInvalidCredentials.Error())
		case custom_errors.ErrUserNotFound:
			utils.SendError(w, http.StatusNotFound, custom_errors.ErrUserNotFound.Error())
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		}
		return
	}

	utils.Send(w, http.StatusOK, nil)
}
