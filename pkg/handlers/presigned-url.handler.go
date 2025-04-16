package handler

import (
	"context"
	"encoding/json"
	"lambda-go/pkg/models"
	"lambda-go/pkg/utils"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// PresignedURLHandleRequest는 Presigned URL 요청을 처리합니다.
func (h *Handler) PresignedURLHandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("요청 본문: %s", request.Body)

	// OPTIONS 메서드 처리 (CORS 프리플라이트 요청)
	if request.HTTPMethod == "OPTIONS" {
		return h.successResponse(http.StatusOK, ""), nil
	}

	// 요청 파싱
	var req models.PresignedURLRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		log.Printf("요청 파싱 오류: %v", err)
		return h.handleAppError(utils.BadRequest("잘못된 요청 형식입니다", err)), nil
	}

	// 요청 검증 및 전처리
	err = h.s3Service.ValidateAndPreprocessRequest(&req)
	if err != nil {
		log.Printf("요청 검증 오류: %v", err)
		return h.handleAppError(utils.BadRequest(err.Error())), nil
	}

	// Presigned URL 생성
	resp, err := h.s3Service.GeneratePresignedURL(ctx, &req)
	if err != nil {
		log.Printf("Presigned URL 생성 오류: %v", err)
		return h.handleAppError(utils.InternalServerError("Presigned URL 생성 중 오류가 발생했습니다", err)), nil
	}
	return h.successResponse(200, resp), nil
}