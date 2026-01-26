# Multi-stage build를 사용하여 최종 이미지 크기 최소화
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 의존성 파일 복사
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 복사
COPY . .

# 바이너리 빌드
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o health-checker .

# 최종 이미지 (Alpine Linux 사용)
FROM alpine:latest

# ca-certificates는 HTTPS 요청에 필요
RUN apk --no-cache add ca-certificates

WORKDIR /app

# 빌드된 바이너리 복사
COPY --from=builder /app/health-checker .

# 실행 권한 부여
RUN chmod +x health-checker

# 기본 명령어 (환경 변수로 오버라이드 가능)
CMD ["./health-checker", "run", "--url", "${URL}", "--interval", "${INTERVAL:-60s}", "--timeout", "${TIMEOUT:-5s}"]
