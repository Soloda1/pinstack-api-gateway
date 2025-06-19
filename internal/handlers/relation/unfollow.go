package relation_handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/utils"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UnfollowRequest struct {
	FollowerID int64 `json:"follower_id" validate:"required,gt=0"`
	FolloweeID int64 `json:"followee_id" validate:"required,gt=0"`
}

type UnfollowResponse struct {
	Message string `json:"message"`
}

// Unfollow godoc
// @Summary Unfollow user
// @Description Unfollow another user
// @Tags relation
// @Accept json
// @Produce json
// @Param request body UnfollowRequest true "Unfollow request"
// @Success 200 {object} UnfollowResponse "Unfollowed successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /relation/unfollow [post]
func (h *RelationHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	var req UnfollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Debug("Failed to decode unfollow request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrInvalidInput.Error())
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		h.log.Debug("Failed to validate unfollow request", slog.String("error", err.Error()))
		utils.SendError(w, http.StatusBadRequest, custom_errors.ErrValidationFailed.Error())
		return
	}

	err := h.relationClient.Unfollow(r.Context(), req.FollowerID, req.FolloweeID)
	if err != nil {
		h.log.Error("unfollow failed", slog.String("error", err.Error()))

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

	response := UnfollowResponse{
		Message: "Unfollowed successfully",
	}
	utils.Send(w, http.StatusOK, response)
}
