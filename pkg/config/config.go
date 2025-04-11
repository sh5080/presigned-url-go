package config

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Config는 애플리케이션 설정 값을 관리하는 구조체입니다.
type Config struct {
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	DefaultBucket      string
	DefaultMaxFileSize int64
	Environment        string
}

// 환경 변수로부터 기본값을 가져오는 함수
func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// NewConfig는 환경 변수에서 설정을 로드하여 Config 구조체를 반환합니다.
func NewConfig() *Config {
	defaultMaxFileSizeStr := GetEnvOrDefault("DEFAULT_MAX_FILE_SIZE", "10485760")
	defaultMaxFileSize, err := strconv.ParseInt(defaultMaxFileSizeStr, 10, 64)
	if err != nil {
		defaultMaxFileSize = 10485760 // 10MB 기본값
	}

	return &Config{
		AWSRegion:          GetEnvOrDefault("AWS_REGION", "ap-northeast-2"),
		AWSAccessKeyID:     GetEnvOrDefault("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: GetEnvOrDefault("AWS_SECRET_ACCESS_KEY", ""),
		DefaultBucket:      GetEnvOrDefault("DEFAULT_BUCKET", ""),
		DefaultMaxFileSize: defaultMaxFileSize,
		Environment:        GetEnvOrDefault("ENV", "dev"),
	}
}

// NewS3Client는 AWS S3 클라이언트와 Presign 클라이언트를 생성합니다.
func NewS3Client(ctx context.Context, cfg *Config) (*s3.Client, *s3.PresignClient, error) {
	// 기본 옵션 슬라이스 생성
	options := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.AWSRegion),
	}

	// 명시적인 AWS 자격 증명이 제공된 경우 사용
	if cfg.AWSAccessKeyID != "" && cfg.AWSSecretAccessKey != "" {
		options = append(options, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWSAccessKeyID,
			cfg.AWSSecretAccessKey,
			"", // Session token (선택적)
		)))
	}

	// AWS 설정 로드
	awsCfg, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return nil, nil, fmt.Errorf("AWS 설정 로드 실패: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg)
	presignClient := s3.NewPresignClient(s3Client)
	
	return s3Client, presignClient, nil
} 