package auth_handler

import (
	"encoding/json"
	"errors"
	"github.com/soloda1/pinstack-proto-definitions/custom_errors"
	"log/slog"
	"net/http"

	"pinstack-api-gateway/internal/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-playground/validator/v10"
)

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Logout godoc
// @Summary User logout
// @Description Logout user and invalidate refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LogoutRequest true "Logout data"
// @Success 200 {object} map[string]string "Successful logout"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Invalid refresh token"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/logout [post]
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

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrInvalidRefreshToken.Error())
				return
			case codes.Unauthenticated:
				utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrUnauthenticated.Error())
				return
			case codes.NotFound:
				utils.SendError(w, http.StatusNotFound, custom_errors.ErrUserNotFound.Error())
				return
			case codes.Internal:
				utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
				return
			}
		}

		switch {
		case errors.Is(err, custom_errors.ErrInvalidRefreshToken):
			utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrInvalidRefreshToken.Error())
		case errors.Is(err, custom_errors.ErrTokenExpired):
			utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrTokenExpired.Error())
		case errors.Is(err, custom_errors.ErrUnauthenticated):
			utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrUnauthenticated.Error())
		case errors.Is(err, custom_errors.ErrUserNotFound):
			utils.SendError(w, http.StatusNotFound, custom_errors.ErrUserNotFound.Error())
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		}
		return
	}

	utils.Send(w, http.StatusOK, nil)
}
