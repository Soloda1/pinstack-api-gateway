package notification_handler

import (
	"errors"
	"github.com/soloda1/pinstack-proto-definitions/custom_errors"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/middlewares"
	"pinstack-api-gateway/internal/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ReadNotificationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ReadNotification godoc
// @Summary Mark notification as read
// @Description Mark a notification as read
// @Tags notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param notification_id path int true "Notification ID"
// @Success 200 {object} ReadNotificationResponse "Notification marked as read"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Notification not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /notification/{notification_id}/read [put]
func (h *NotificationHandler) ReadNotification(w http.ResponseWriter, r *http.Request) {
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
		h.log.Error("Failed to get notification details for read", slog.Int64("notification_id", notificationID), slog.String("error", err.Error()))

		switch {
		case errors.Is(err, custom_errors.ErrNotificationNotFound):
			utils.SendError(w, http.StatusNotFound, custom_errors.ErrNotificationNotFound.Error())
			return
		case errors.Is(err, custom_errors.ErrValidationFailed):
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
			return
		case errors.Is(err, custom_errors.ErrNotificationAccessDenied):
			utils.SendError(w, http.StatusForbidden, custom_errors.ErrNotificationAccessDenied.Error())
			return
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
			return
		}
	}

	if notification.UserID != claims.UserID {
		h.log.Debug(
			"User not authorized to mark this notification as read",
			slog.Int64("notification_user_id", notification.UserID),
			slog.Int64("requester_user_id", claims.UserID),
		)
		utils.SendError(w, http.StatusForbidden, custom_errors.ErrNotificationAccessDenied.Error())
		return
	}

	if notification.IsRead {
		response := ReadNotificationResponse{
			Success: true,
			Message: "Notification already marked as read",
		}
		utils.Send(w, http.StatusOK, response)
		return
	}

	err = h.notificationClient.ReadNotification(r.Context(), notificationID)
	if err != nil {
		h.log.Error("Failed to mark notification as read", slog.Int64("notification_id", notificationID), slog.String("error", err.Error()))

		switch {
		case errors.Is(err, custom_errors.ErrNotificationNotFound):
			utils.SendError(w, http.StatusNotFound, custom_errors.ErrNotificationNotFound.Error())
			return
		case errors.Is(err, custom_errors.ErrValidationFailed):
			utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
			return
		case errors.Is(err, custom_errors.ErrNotificationAccessDenied):
			utils.SendError(w, http.StatusForbidden, custom_errors.ErrNotificationAccessDenied.Error())
			return
		case errors.Is(err, custom_errors.ErrExternalServiceTimeout):
			utils.SendError(w, http.StatusGatewayTimeout, custom_errors.ErrExternalServiceTimeout.Error())
			return
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
			return
		}
	}

	response := ReadNotificationResponse{
		Success: true,
		Message: "Notification marked as read",
	}
	utils.Send(w, http.StatusOK, response)
}
