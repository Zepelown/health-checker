# 상세 메트릭 및 통계

## 개요

헬스 체크 결과를 단순히 로그로만 남기는 것이 아니라, 메트릭을 수집하고 통계를 제공하여 장기적인 모니터링과 분석이 가능하도록 합니다.

## 기능 요구사항

### 1. 수집할 메트릭

#### 기본 메트릭
- **응답 시간 (Latency)**: 평균, 최소, 최대, P50, P95, P99
- **상태 코드 분포**: 각 HTTP 상태 코드별 발생 횟수
- **가용성 (Availability)**: 시간대별 가용성 퍼센트
- **에러율 (Error Rate)**: 시간대별 에러 발생 비율
- **체크 횟수**: 총 체크 횟수, 성공 횟수, 실패 횟수

#### 고급 메트릭
- **연속 실패 횟수**: 현재 연속으로 실패한 횟수
- **복구 시간**: 장애 발생 후 복구까지 걸린 시간
- **MTTR (Mean Time To Recovery)**: 평균 복구 시간
- **MTBF (Mean Time Between Failures)**: 평균 장애 간격

### 2. 메트릭 저장 방식

#### 옵션 1: 인메모리 저장 (기본)
- 프로그램 실행 중에만 메트릭 유지
- 간단하고 빠름
- 재시작 시 데이터 손실

#### 옵션 2: 파일 기반 저장
- JSON 또는 CSV 형식으로 메트릭 저장
- 프로그램 재시작 후에도 데이터 유지
- 제한된 쿼리 기능

#### 옵션 3: 데이터베이스 저장 (향후)
- SQLite, PostgreSQL 등
- 복잡한 쿼리 및 분석 가능
- 장기 데이터 보관

### 3. 통계 출력 형식

#### CLI 출력
```bash
health-checker stats --url https://example.com
```

출력 예시:
```
Statistics for https://example.com
================================
Total Checks: 1440
Successful: 1435 (99.65%)
Failed: 5 (0.35%)

Response Time:
  Average: 150ms
  Min: 45ms
  Max: 2.5s
  P50: 140ms
  P95: 280ms
  P99: 450ms

Availability: 99.65%
Uptime: 23h 45m
Downtime: 15m

Status Code Distribution:
  200: 1435 (99.65%)
  500: 3 (0.21%)
  503: 2 (0.14%)

Recent Failures:
  2025-01-15 14:30:00 - Status: 500, Latency: 2.5s
  2025-01-15 14:25:00 - Status: 503, Latency: 1.8s
```

#### JSON 출력
```bash
health-checker stats --url https://example.com --format json
```

```json
{
  "url": "https://example.com",
  "period": {
    "start": "2025-01-15T00:00:00Z",
    "end": "2025-01-15T23:59:59Z"
  },
  "summary": {
    "total_checks": 1440,
    "successful": 1435,
    "failed": 5,
    "availability": 99.65
  },
  "response_time": {
    "average": 150000000,
    "min": 45000000,
    "max": 2500000000,
    "p50": 140000000,
    "p95": 280000000,
    "p99": 450000000
  },
  "status_codes": {
    "200": 1435,
    "500": 3,
    "503": 2
  }
}
```

### 4. Prometheus 메트릭 노출

#### HTTP 엔드포인트
```
GET /metrics
```

#### 메트릭 형식
```prometheus
# HELP health_checker_response_time_seconds Response time in seconds
# TYPE health_checker_response_time_seconds histogram
health_checker_response_time_seconds_bucket{url="https://example.com",le="0.1"} 1200
health_checker_response_time_seconds_bucket{url="https://example.com",le="0.5"} 1430
health_checker_response_time_seconds_bucket{url="https://example.com",le="1.0"} 1435
health_checker_response_time_seconds_bucket{url="https://example.com",le="+Inf"} 1440
health_checker_response_time_seconds_sum{url="https://example.com"} 216.5
health_checker_response_time_seconds_count{url="https://example.com"} 1440

# HELP health_checker_up Whether the endpoint is up (1) or down (0)
# TYPE health_checker_up gauge
health_checker_up{url="https://example.com"} 1

# HELP health_checker_status_code HTTP status code
# TYPE health_checker_status_code gauge
health_checker_status_code{url="https://example.com"} 200

# HELP health_checker_checks_total Total number of health checks
# TYPE health_checker_checks_total counter
health_checker_checks_total{url="https://example.com",status="success"} 1435
health_checker_checks_total{url="https://example.com",status="failure"} 5
```

#### 사용 예제
```bash
# Prometheus 메트릭 서버 시작
health-checker run \
  --url https://example.com \
  --metrics-port 9090 \
  --metrics-path /metrics

# Prometheus에서 스크랩
# scrape_configs:
#   - job_name: 'health-checker'
#     static_configs:
#       - targets: ['localhost:9090']
```

### 5. 실시간 통계 대시보드

#### CLI 대시보드
```bash
health-checker dashboard --url https://example.com
```

실시간 업데이트되는 대시보드:
```
┌─────────────────────────────────────────────────────────┐
│ Health Checker Dashboard - https://example.com         │
├─────────────────────────────────────────────────────────┤
│ Status: ✅ UP                                           │
│ Uptime: 99.65% (Last 24h)                              │
│ Current Response Time: 145ms                           │
├─────────────────────────────────────────────────────────┤
│ Response Time (Last Hour)                              │
│   Avg: 150ms  Min: 45ms  Max: 280ms                    │
│   P95: 280ms  P99: 450ms                               │
├─────────────────────────────────────────────────────────┤
│ Status Codes (Last Hour)                               │
│   200: ████████████████████████████████ 99.65%         │
│   500: █ 0.21%                                         │
│   503: █ 0.14%                                         │
├─────────────────────────────────────────────────────────┤
│ Recent Events                                           │
│   10:30:00 ✅ 200 (145ms)                              │
│   10:29:30 ✅ 200 (150ms)                              │
│   10:29:00 ❌ 500 (2.5s)                               │
└─────────────────────────────────────────────────────────┘
```

## 구현 계획

### 1. 메트릭 수집 구조

```go
type Metrics struct {
    URL           string
    StartTime     time.Time
    TotalChecks   int64
    Successful    int64
    Failed        int64
    ResponseTimes []time.Duration
    StatusCodes   map[int]int64
    Errors        []CheckError
    LastCheck     time.Time
    LastStatus    int
    LastLatency   time.Duration
}

type CheckError struct {
    Timestamp time.Time
    Status    int
    Latency   time.Duration
    Error     error
}
```

### 2. 통계 계산 함수

```go
type Statistics struct {
    Availability    float64
    AverageLatency  time.Duration
    MinLatency      time.Duration
    MaxLatency      time.Duration
    P50Latency      time.Duration
    P95Latency      time.Duration
    P99Latency      time.Duration
    StatusCodes     map[int]int64
    ErrorRate       float64
    MTTR            time.Duration
    MTBF            time.Duration
}
```

### 3. 데이터 저장

#### 파일 기반 저장 (JSON)
```json
{
  "url": "https://example.com",
  "checks": [
    {
      "timestamp": "2025-01-15T10:00:00Z",
      "status": 200,
      "latency_ms": 150,
      "success": true
    },
    {
      "timestamp": "2025-01-15T10:00:30Z",
      "status": 500,
      "latency_ms": 2500,
      "success": false
    }
  ]
}
```

### 4. 시간 윈도우 기반 통계

- **실시간**: 마지막 1분, 5분, 15분
- **단기**: 마지막 1시간, 6시간, 12시간
- **장기**: 마지막 24시간, 7일, 30일

## 사용 예제

### 예제 1: 기본 통계 확인
```bash
health-checker run --url https://example.com --save-metrics
health-checker stats --url https://example.com
```

### 예제 2: 특정 기간 통계
```bash
health-checker stats \
  --url https://example.com \
  --from "2025-01-15T00:00:00Z" \
  --to "2025-01-15T23:59:59Z"
```

### 예제 3: Prometheus 메트릭 노출
```bash
health-checker run \
  --url https://example.com \
  --metrics-server \
  --metrics-port 9090
```

### 예제 4: 실시간 대시보드
```bash
health-checker dashboard --url https://example.com --refresh 1s
```

## 고려사항

### 1. 성능
- 메트릭 수집 오버헤드 최소화
- 대량의 데이터 처리 효율성
- 메모리 사용량 관리

### 2. 데이터 보관
- 데이터 보관 기간 설정
- 오래된 데이터 자동 삭제
- 데이터 압축 및 최적화

### 3. 정확도
- 시간 윈도우 계산 정확도
- 백분위수 계산 알고리즘
- 통계 계산 성능

## 향후 확장 가능성

1. **Grafana 연동**: Prometheus 메트릭을 Grafana에서 시각화
2. **알림 규칙**: 메트릭 기반 알림 (예: 가용성이 99% 이하로 떨어지면)
3. **트렌드 분석**: 장기적인 트렌드 분석 및 예측
4. **비교 분석**: 여러 엔드포인트 간 성능 비교
5. **SLA 모니터링**: SLA 목표 설정 및 모니터링

