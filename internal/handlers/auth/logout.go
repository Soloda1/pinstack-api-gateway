package auth_handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/utils"

	"github.com/go-playground/validator/v10"
)

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Debug("Failed to decode logout request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		h.log.Debug("Failed to validate logout request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
		return
	}

	err := h.authClient.Logout(r.Context(), req.RefreshToken)
	if err != nil {
		h.log.Error("logout failed", slog.String("error", err.Error()))
		switch err {
		case custom_errors.ErrInvalidRefreshToken:
			utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrInvalidRefreshToken.Error())
		case custom_errors.ErrTokenExpired:
			utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrTokenExpired.Error())
		case custom_errors.ErrUnauthenticated:
			utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrUnauthenticated.Error())
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		}
		return
	}

	utils.Send(w, http.StatusOK, nil)
}
