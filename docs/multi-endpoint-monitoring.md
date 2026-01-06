# 다중 엔드포인트 모니터링

## 개요

현재 health-checker는 단일 URL만 모니터링할 수 있습니다. 다중 엔드포인트 모니터링 기능을 추가하여 여러 URL을 동시에 모니터링하고, 각 엔드포인트별로 독립적인 헬스 체크를 수행할 수 있습니다.

## 기능 요구사항

### 1. 기본 기능
- 여러 URL을 동시에 모니터링
- 각 엔드포인트별 독립적인 헬스 체크 주기 설정
- 각 엔드포인트별 독립적인 타임아웃 설정
- 각 엔드포인트별 독립적인 알림 설정

### 2. CLI 인터페이스

#### 옵션 1: 명령줄에서 여러 URL 지정
```bash
health-checker monitor \
  --url https://api.example.com \
  --url https://web.example.com \
  --url https://db.example.com/health \
  --interval 30s
```

#### 옵션 2: 설정 파일 사용 (권장)
```bash
health-checker monitor --config endpoints.yaml
```

### 3. 설정 파일 형식 (YAML)

```yaml
# endpoints.yaml
endpoints:
  - name: "API Server"
    url: https://api.example.com
    interval: 30s
    timeout: 5s
    latency_threshold: 1s
    slack_webhook: https://hooks.slack.com/services/API/WEBHOOK
    discord_webhook: https://discord.com/api/webhooks/API/WEBHOOK
    enabled: true

  - name: "Web Server"
    url: https://web.example.com
    interval: 60s
    timeout: 10s
    latency_threshold: 2s
    slack_webhook: https://hooks.slack.com/services/WEB/WEBHOOK
    enabled: true

  - name: "Database Health"
    url: https://db.example.com/health
    interval: 15s
    timeout: 3s
    enabled: false  # 일시적으로 비활성화

# 전역 설정 (모든 엔드포인트에 적용, 개별 설정으로 오버라이드 가능)
global:
  interval: 60s
  timeout: 5s
  slack_webhook: https://hooks.slack.com/services/GLOBAL/WEBHOOK
  discord_webhook: https://discord.com/api/webhooks/GLOBAL/WEBHOOK
  latency_threshold: 3s
```

### 4. JSON 형식 지원

```json
{
  "global": {
    "interval": "60s",
    "timeout": "5s",
    "slack_webhook": "https://hooks.slack.com/services/GLOBAL/WEBHOOK"
  },
  "endpoints": [
    {
      "name": "API Server",
      "url": "https://api.example.com",
      "interval": "30s",
      "timeout": "5s",
      "latency_threshold": "1s",
      "enabled": true
    },
    {
      "name": "Web Server",
      "url": "https://web.example.com",
      "interval": "60s",
      "enabled": true
    }
  ]
}
```

## 구현 계획

### 1. 새로운 서브커맨드 추가
- `health-checker monitor`: 다중 엔드포인트 모니터링 모드

### 2. 설정 파일 파싱
- YAML 파서: `gopkg.in/yaml.v3`
- JSON 파서: 표준 라이브러리 `encoding/json`
- 설정 검증 로직

### 3. 병렬 모니터링
- 각 엔드포인트별로 독립적인 goroutine 실행
- 각 엔드포인트별로 독립적인 ticker 사용
- 동시성 안전성 보장 (sync 패키지 활용)

### 4. 로깅 개선
- 각 엔드포인트별로 구분된 로그
- 엔드포인트 이름 포함한 로그 메시지
- 색상 코딩 (선택 사항)

## 사용 예제

### 예제 1: 설정 파일 사용
```bash
# endpoints.yaml 생성
cat > endpoints.yaml << EOF
endpoints:
  - name: "Production API"
    url: https://api.prod.example.com
    interval: 30s
    timeout: 5s
    latency_threshold: 1s
  
  - name: "Staging API"
    url: https://api.staging.example.com
    interval: 60s
    timeout: 10s
EOF

# 모니터링 시작
health-checker monitor --config endpoints.yaml
```

### 예제 2: 환경 변수와 함께 사용
```bash
export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK"
health-checker monitor --config endpoints.yaml
```

### 예제 3: 명령줄에서 직접 지정
```bash
health-checker monitor \
  --url https://api1.example.com --interval 30s \
  --url https://api2.example.com --interval 60s \
  --slack-webhook https://hooks.slack.com/services/YOUR/WEBHOOK
```

## 출력 예시

```
Starting multi-endpoint monitoring...
Monitoring 3 endpoints:
  - API Server (https://api.example.com) [interval: 30s]
  - Web Server (https://web.example.com) [interval: 60s]
  - Database Health (https://db.example.com/health) [interval: 15s]

[2025-01-15 10:00:00] ✅ [API Server] Status: 200 (latency: 150ms)
[2025-01-15 10:00:00] ✅ [Web Server] Status: 200 (latency: 200ms)
[2025-01-15 10:00:00] ✅ [Database Health] Status: 200 (latency: 50ms)
[2025-01-15 10:00:15] ✅ [Database Health] Status: 200 (latency: 45ms)
[2025-01-15 10:00:30] ✅ [API Server] Status: 200 (latency: 160ms)
[2025-01-15 10:00:30] ❌ [Web Server] Status: 500 (latency: 2.5s)
🚨 사이트 장애 감지: Web Server (https://web.example.com)
```

## 고려사항

### 1. 성능
- 많은 엔드포인트를 모니터링할 때 리소스 사용량
- 각 엔드포인트별 goroutine 관리
- 메모리 사용량 최적화

### 2. 설정 관리
- 설정 파일 변경 시 동적 리로드 (선택 사항)
- 설정 파일 검증 및 에러 처리
- 기본값 처리 로직

### 3. 확장성
- 수백 개의 엔드포인트 모니터링 지원
- 설정 파일 크기 제한
- 효율적인 스케줄링

## 향후 확장 가능성

1. **엔드포인트 그룹화**: 여러 엔드포인트를 그룹으로 묶어서 관리
2. **의존성 체크**: 엔드포인트 간 의존성 정의 및 체크
3. **동적 추가/제거**: 런타임에 엔드포인트 추가/제거
4. **엔드포인트 태그**: 태그 기반 필터링 및 그룹화


