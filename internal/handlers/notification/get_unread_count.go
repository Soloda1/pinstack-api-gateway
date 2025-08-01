package notification_handler

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/middlewares"
	"pinstack-api-gateway/internal/utils"
)

type GetUnreadCountResponse struct {
	Count int32 `json:"count"`
}

// GetUnreadCount godoc
// @Summary Get unread notification count
// @Description Get the count of unread notifications for a user
// @Tags notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} GetUnreadCountResponse "Unread notification count"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /notification/unread-count [get]
func (h *NotificationHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	claims, err := middlewares.GetClaimsFromContext(r.Context())
	if err != nil {
		h.log.Debug("No user claims in context", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrUnauthenticated.Error())
		return
	}

	count, err := h.notificationClient.GetUnreadCount(r.Context(), claims.UserID)
	if err != nil {
		h.log.Error("Failed to get unread notification count", slog.Int64("user_id", claims.UserID), slog.String("error", err.Error()))

		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				utils.SendError(w, http.StatusNotFound, custom_errors.ErrUserNotFound.Error())
				return
			case codes.Internal:
				utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
				return
			}
		}

		switch {
		case errors.Is(err, custom_errors.ErrUserNotFound):
			utils.SendError(w, http.StatusNotFound, custom_errors.ErrUserNotFound.Error())
			return
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
			return
		}
	}

	response := GetUnreadCountResponse{
		Count: count,
	}
	utils.Send(w, http.StatusOK, response)
}
