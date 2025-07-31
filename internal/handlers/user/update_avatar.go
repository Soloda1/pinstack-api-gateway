package user_handler

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/middlewares"
	"pinstack-api-gateway/internal/utils"
)

type UpdateAvatarRequest struct {
	AvatarURL string `json:"avatar_url" validate:"required,url"`
}

// UpdateAvatar godoc
// @Summary Update user avatar
// @Description Update the avatar URL for a user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateAvatarRequest true "Avatar update data"
// @Success 200 {object} nil "Avatar updated successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Operation not allowed"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/avatar [put]
func (h *UserHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	claims, err := middlewares.GetClaimsFromContext(r.Context())
	if err != nil {
		h.log.Debug("No user claims in context", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrUnauthenticated.Error())
		return
	}

	var req UpdateAvatarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
		return
	}

	err = h.userClient.UpdateAvatar(r.Context(), claims.UserID, req.AvatarURL)
	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrUserNotFound):
			utils.SendError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, custom_errors.ErrValidationFailed):
			utils.SendError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, custom_errors.ErrOperationNotAllowed):
			utils.SendError(w, http.StatusForbidden, err.Error())
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		}
		return
	}

	utils.Send(w, http.StatusOK, nil)
}
