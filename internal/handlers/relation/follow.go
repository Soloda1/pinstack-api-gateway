package relation_handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/middlewares"
	"pinstack-api-gateway/internal/utils"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FollowRequest struct {
	FolloweeID int64 `json:"followee_id" validate:"required,gt=0"`
}

type FollowResponse struct {
	Message string `json:"message"`
}

// Follow godoc
// @Summary Follow user
// @Description Follow another user
// @Tags relation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body FollowRequest true "Follow request"
// @Success 200 {object} FollowResponse "Followed successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /relation/follow [post]
func (h *RelationHandler) Follow(w http.ResponseWriter, r *http.Request) {
	var req FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Debug("Failed to decode follow request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	validate := validator.New()
	if err := validate.StructPartial(req, "FolloweeID"); err != nil {
		h.log.Debug("Failed to validate follow request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
		return
	}

	claims, err := middlewares.GetClaimsFromContext(r.Context())
	if err != nil {
		h.log.Debug("No user claims in context", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrUnauthenticated.Error())
		return
	}

	err = h.relationClient.Follow(r.Context(), claims.UserID, req.FolloweeID)
	if err != nil {
		h.log.Error("follow failed", slog.String("error", err.Error()))
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
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
		case errors.Is(err, custom_errors.ErrUserNotFound):
			utils.SendError(w, http.StatusNotFound, custom_errors.ErrUserNotFound.Error())
		default:
			utils.SendError(w, http.StatusInternalServerError, custom_errors.ErrExternalServiceError.Error())
		}
		return
	}

	response := FollowResponse{
		Message: "Followed successfully",
	}
	utils.Send(w, http.StatusOK, response)
}
