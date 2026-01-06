# 알림 채널 확장

## 개요

현재 health-checker는 Slack과 Discord만 지원합니다. 다양한 알림 채널을 추가하고, 더 세밀한 알림 규칙을 설정할 수 있도록 확장합니다.

## 기능 요구사항

### 1. 추가 알림 채널

#### 이메일 알림
- SMTP 서버를 통한 이메일 전송
- HTML 형식 지원
- 첨부 파일 (로그, 리포트)

#### PagerDuty 연동
- PagerDuty Events API v2
- 인시던트 자동 생성/해결
- 우선순위 설정

#### 커스텀 웹훅
- 일반적인 HTTP POST 요청
- JSON 페이로드 커스터마이징
- 인증 헤더 지원

#### 텔레그램 봇
- Telegram Bot API
- 실시간 알림
- 인라인 키보드 지원

#### Microsoft Teams
- Teams Incoming Webhook
- Adaptive Cards 지원
- 리치 포맷팅

#### SMS (선택 사항)
- Twilio, AWS SNS 등
- 긴급 알림용

### 2. 알림 규칙 시스템

#### 기본 규칙
```yaml
notifications:
  channels:
    - type: slack
      webhook: https://hooks.slack.com/services/...
      enabled: true
    
    - type: email
      smtp:
        host: smtp.gmail.com
        port: 587
        username: user@example.com
        password: password
      to: admin@example.com
      enabled: true

  rules:
    # 모든 장애에 대해 Slack 알림
    - name: "All Failures"
      condition: status != 200 OR error != nil
      channels: [slack]
      enabled: true
    
    # 심각한 장애에 대해 이메일 + PagerDuty
    - name: "Critical Failures"
      condition: status >= 500 OR consecutive_failures >= 3
      channels: [email, pagerduty]
      enabled: true
    
    # 응답 시간이 느릴 때만 Slack
    - name: "Slow Response"
      condition: latency > 2s AND status == 200
      channels: [slack]
      enabled: true
    
    # 특정 시간대에만 알림 (업무 시간)
    - name: "Business Hours Only"
      condition: status != 200
      time_window:
        start: "09:00"
        end: "18:00"
        timezone: "Asia/Seoul"
        weekdays: [monday, tuesday, wednesday, thursday, friday]
      channels: [slack]
      enabled: true
    
    # 연속 실패 시에만 알림 (알림 스팸 방지)
    - name: "Consecutive Failures"
      condition: consecutive_failures >= 3
      channels: [slack, email]
      cooldown: 5m  # 5분간 중복 알림 방지
      enabled: true
```

### 3. 알림 템플릿

#### 커스터마이징 가능한 메시지 템플릿
```yaml
templates:
  failure:
    slack: |
      🚨 *Site Failure Detected*
      *URL:* {url}
      *Status:* {status}
      *Latency:* {latency}
      *Time:* {timestamp}
      *Error:* {error}
    
    email:
      subject: "🚨 Site Failure: {url}"
      body: |
        <h2>Site Failure Detected</h2>
        <p><strong>URL:</strong> {url}</p>
        <p><strong>Status:</strong> {status}</p>
        <p><strong>Latency:</strong> {latency}</p>
        <p><strong>Time:</strong> {timestamp}</p>
        {#if error}
        <p><strong>Error:</strong> {error}</p>
        {/if}
    
    pagerduty:
      severity: "critical"
      summary: "Site Failure: {url}"
      source: "health-checker"
      custom_details:
        url: "{url}"
        status: "{status}"
        latency: "{latency}"
```

### 4. 알림 그룹핑 및 집계

#### 알림 집계
- 짧은 시간 내 여러 알림을 하나로 묶기
- 다이제스트 형식으로 전송
- 알림 스팸 방지

```yaml
aggregation:
  enabled: true
  window: 5m  # 5분간 알림 집계
  max_alerts: 10  # 최대 10개까지 집계
  format: digest  # digest 또는 individual
```

### 5. 알림 이력 및 추적

#### 알림 로그
- 전송된 알림 기록
- 전송 성공/실패 추적
- 재시도 로직

## 구현 계획

### 1. 알림 인터페이스

```go
// internal/notifier/notifier.go
type Notifier interface {
    Send(message NotificationMessage) error
    Name() string
    Validate() error
}

type NotificationMessage struct {
    Title       string
    Body        string
    Severity    string  // info, warning, error, critical
    URL         string
    Status      int
    Latency     time.Duration
    Error       error
    Timestamp   time.Time
    Metadata    map[string]interface{}
}
```

### 2. 각 알림 채널 구현

#### 이메일 알림
```go
// internal/notifier/email.go
type EmailNotifier struct {
    SMTPHost     string
    SMTPPort     int
    Username     string
    Password     string
    From         string
    To           []string
    Subject      string
    Template     string
}

func (e *EmailNotifier) Send(msg NotificationMessage) error {
    // SMTP를 통한 이메일 전송
}
```

#### PagerDuty 연동
```go
// internal/notifier/pagerduty.go
type PagerDutyNotifier struct {
    IntegrationKey string
    APIURL         string
}

func (p *PagerDutyNotifier) Send(msg NotificationMessage) error {
    // PagerDuty Events API v2 호출
    event := map[string]interface{}{
        "routing_key": p.IntegrationKey,
        "event_action": "trigger",
        "payload": map[string]interface{}{
            "summary": msg.Title,
            "severity": msg.Severity,
            "source": "health-checker",
            "custom_details": msg.Metadata,
        },
    }
    // HTTP POST 요청
}
```

#### 커스텀 웹훅
```go
// internal/notifier/webhook.go
type WebhookNotifier struct {
    URL     string
    Method  string  // POST, PUT
    Headers map[string]string
    Body    string  // JSON 템플릿
}

func (w *WebhookNotifier) Send(msg NotificationMessage) error {
    // 커스텀 HTTP 요청
}
```

### 3. 알림 규칙 엔진

```go
// internal/notifier/rules.go
type Rule struct {
    Name        string
    Condition   string  // 표현식 또는 스크립트
    Channels    []string
    TimeWindow  *TimeWindow
    Cooldown    time.Duration
    Enabled     bool
}

type TimeWindow struct {
    Start    string
    End      string
    Timezone string
    Weekdays []string
}

func (r *Rule) Evaluate(checkResult CheckResult) bool {
    // 조건 평가
    // 시간 윈도우 확인
    // 쿨다운 확인
}
```

### 4. 알림 집계

```go
// internal/notifier/aggregator.go
type Aggregator struct {
    window    time.Duration
    maxAlerts int
    alerts    []NotificationMessage
    mutex     sync.Mutex
}

func (a *Aggregator) Add(msg NotificationMessage) {
    // 알림 추가
    // 윈도우 초과 시 전송
}

func (a *Aggregator) Flush() []NotificationMessage {
    // 집계된 알림 반환
}
```

## 사용 예제

### 예제 1: 이메일 알림 설정
```bash
health-checker run \
  --url https://example.com \
  --email-smtp-host smtp.gmail.com \
  --email-smtp-port 587 \
  --email-username user@example.com \
  --email-password password \
  --email-to admin@example.com
```

### 예제 2: PagerDuty 연동
```bash
health-checker run \
  --url https://example.com \
  --pagerduty-key YOUR_INTEGRATION_KEY
```

### 예제 3: 설정 파일로 알림 규칙 정의
```yaml
# config.yaml
endpoints:
  - url: https://api.example.com
    notifications:
      rules:
        - name: "Critical"
          condition: status >= 500
          channels: [slack, email, pagerduty]
        - name: "Warning"
          condition: latency > 2s
          channels: [slack]
```

```bash
health-checker monitor --config config.yaml
```

### 예제 4: 커스텀 웹훅
```bash
health-checker run \
  --url https://example.com \
  --webhook-url https://api.example.com/alerts \
  --webhook-method POST \
  --webhook-header "Authorization: Bearer token" \
  --webhook-body '{"alert": "{title}", "url": "{url}"}'
```

## 알림 메시지 예시

### Slack (Rich Format)
```json
{
  "text": "🚨 Site Failure Detected",
  "blocks": [
    {
      "type": "section",
      "text": {
        "type": "mrkdwn",
        "text": "*URL:* https://example.com\n*Status:* 500\n*Latency:* 2.5s"
      }
    },
    {
      "type": "actions",
      "elements": [
        {
          "type": "button",
          "text": {
            "type": "plain_text",
            "text": "View Details"
          },
          "url": "https://dashboard.example.com/endpoint/api"
        }
      ]
    }
  ]
}
```

### PagerDuty
```json
{
  "routing_key": "YOUR_INTEGRATION_KEY",
  "event_action": "trigger",
  "payload": {
    "summary": "Site Failure: https://example.com",
    "severity": "critical",
    "source": "health-checker",
    "custom_details": {
      "url": "https://example.com",
      "status": 500,
      "latency_ms": 2500
    }
  }
}
```

### 이메일 (HTML)
```html
<!DOCTYPE html>
<html>
<head>
  <style>
    .alert { background-color: #f44336; color: white; padding: 20px; }
    .details { background-color: #f5f5f5; padding: 15px; margin-top: 10px; }
  </style>
</head>
<body>
  <div class="alert">
    <h2>🚨 Site Failure Detected</h2>
  </div>
  <div class="details">
    <p><strong>URL:</strong> https://example.com</p>
    <p><strong>Status:</strong> 500</p>
    <p><strong>Latency:</strong> 2.5s</p>
    <p><strong>Time:</strong> 2025-01-15 10:30:00</p>
  </div>
</body>
</html>
```

## 고려사항

### 1. 보안
- 비밀번호/토큰 안전한 저장
- 환경 변수 사용 권장
- 암호화된 설정 파일

### 2. 신뢰성
- 알림 전송 실패 시 재시도
- 백오프 전략
- 데드 레터 큐

### 3. 성능
- 비동기 알림 전송
- 알림 큐 관리
- 배치 전송

### 4. 비용
- SMS, PagerDuty 등 유료 서비스 비용 고려
- 알림 빈도 제한

## 향후 확장 가능성

1. **알림 템플릿 에디터**: 웹 UI에서 템플릿 편집
2. **알림 테스트**: 알림 전송 테스트 기능
3. **알림 통계**: 알림 전송 성공률, 응답 시간 등
4. **알림 피드백**: 사용자 응답 추적 (예: Slack 버튼 클릭)
5. **알림 에스컬레이션**: 일정 시간 후 자동 에스컬레이션
6. **알림 음소거**: 특정 시간대 또는 조건에서 알림 비활성화


