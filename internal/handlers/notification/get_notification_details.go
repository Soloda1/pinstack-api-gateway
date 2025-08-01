package notification_handler

import (
	"errors"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/middlewares"
	"pinstack-api-gateway/internal/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetNotificationDetails godoc
// @Summary Get notification details
// @Description Get notification details by ID
// @Tags notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param notification_id path int true "Notification ID"
// @Success 200 {object} models.NotificationSwagger "Notification details"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Notification not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /notification/{notification_id} [get]
func (h *NotificationHandler) GetNotificationDetails(w http.ResponseWriter, r *http.Request) {
	notificationIDStr := chi.URLParam(r, "notification_id")
	notificationID, err := strconv.ParseInt(notificationIDStr, 10, 64)
	if err != nil {
		h.log.Debug("Failed to parse notification ID", slog.String("notification_id", notificationIDStr), slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	claims, err := middlewares.GetClaimsFromContext(r.Context())
	if err != nil {
		h.log.Debug("No user claims in context", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrUnauthenticated.Error())
		return
	}

	notification, err := h.notificationClient.GetNotificationDetails(r.Context(), notificationID)
	if err != nil {
		h.log.Error("Failed to get notification details", slog.Int64("notification_id", notificationID), slog.String("error", err.Error()))
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				utils.SendError(w, http.StatusNotFound, custom_errors.ErrNotificationNotFound.Error())
				return
			case codes.InvalidArgument:
				utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
				return
			case codes.PermissionDenied:
				utils.SendError(w, http.StatusForbidden, custom_errors.ErrInsufficientRights.Error())
				return
			case codes.Internal:
				utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
				return
			}
		}

		switch {
		case errors.Is(err, custom_errors.ErrNotificationNotFound):
			utils.SendError(w, http.StatusNotFound, custom_errors.ErrNotificationNotFound.Error())
			return
		case errors.Is(err, custom_errors.ErrNotificationAccessDenied):
			utils.SendError(w, http.StatusForbidden, custom_errors.ErrNotificationAccessDenied.Error())
			return
		case errors.Is(err, custom_errors.ErrValidationFailed):
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
			return
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
			return
		}
	}

	if notification.UserID != claims.UserID {
		h.log.Debug(
			"User not authorized to view this notification",
			slog.Int64("notification_user_id", notification.UserID),
			slog.Int64("requester_user_id", claims.UserID),
		)
		utils.SendError(w, http.StatusForbidden, custom_errors.ErrNotificationAccessDenied.Error())
		return
	}

	utils.Send(w, http.StatusOK, notification)
}
