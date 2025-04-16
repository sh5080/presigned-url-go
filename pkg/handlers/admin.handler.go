package handler

import (
	"context"
	"encoding/json"
	appCtx "lambda-go/pkg/contexts"
	"lambda-go/pkg/models"
	"lambda-go/pkg/utils"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// GetRestaurantRequests는 매장 생성 요청 목록을 조회합니다.
func (h *Handler) GetRestaurantRequests(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// OPTIONS 처리 코드 제거 (라우터에서 처리)
	
	resp, err := h.adminService.GetRestaurantRequests(ctx)
	if err != nil {
		return h.handleAppError(utils.InternalServerError("매장 생성 요청 목록 조회 중 오류가 발생했습니다", err)), nil
	}
	
	return h.successResponse(http.StatusOK, resp), nil
}

// ProcessRestaurantRequest는 매장 생성 요청을 처리합니다.
func (h *Handler) ProcessRestaurantRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// URL 파라미터에서 요청 ID 추출
	requestID := appCtx.GetParam(ctx, "id")
	if requestID == "" {
		return h.handleAppError(utils.BadRequest("유효하지 않은 요청 ID입니다")), nil
	}
	
	var payload models.ProcessRestaurantRequest
	
	err := json.Unmarshal([]byte(request.Body), &payload)
	if err != nil {
		return h.handleAppError(utils.BadRequest("잘못된 요청 형식입니다: " + err.Error())), nil
	}
	
	if err := utils.Validate(&payload); err != nil {
		return h.handleAppError(utils.BadRequest(err.Error())), nil
	}
	
	result, err := h.adminService.ProcessRestaurantRequest(ctx, requestID, &payload)
	if err != nil {
		return h.handleAppError(utils.InternalServerError("매장 생성 요청 처리 중 오류가 발생했습니다: " + err.Error())), nil
	}
	
	return h.successResponse(http.StatusOK, result), nil
}