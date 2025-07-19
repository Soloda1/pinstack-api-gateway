package auth_handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-playground/validator/v10"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Refresh godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} RefreshTokenResponse "Successful token refresh"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Invalid refresh token"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Debug("Failed to decode refresh request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		h.log.Debug("Failed to validate refresh request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
		return
	}

	tokens, err := h.authClient.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		h.log.Error("refresh failed", slog.String("error", err.Error()))

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

	response := RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	utils.Send(w, http.StatusOK, response)
}
