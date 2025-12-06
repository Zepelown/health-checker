# Health Checker

웹사이트 가용성을 주기적으로 체크하고, 장애 발생 시 Slack 또는 Discord로 알림을 보내는 Go 기반 헬스 체커입니다.

## 기능

- 주기적인 웹사이트 헬스 체크
- HTTP 상태 코드 및 응답 시간 모니터링
- Slack Webhook을 통한 알림 지원
- Discord Webhook을 통한 알림 지원
- Slack과 Discord 동시 사용 가능
- 테스트 모드 지원 (정상 상태에서도 알림 전송)

## 설치

### Go가 설치된 경우

```bash
go install github.com/your-username/health-checker@latest
```

또는 소스에서 빌드:

```bash
git clone https://github.com/your-username/health-checker.git
cd health-checker
go build
```

## 사용 방법

### 기본 사용법

```bash
# Windows PowerShell
.\health-checker.exe run --url https://example.com

# Linux/Mac
./health-checker run --url https://example.com
```

### 명령어 옵션

#### 필수 옵션

- `--url`, `-u`: 체크할 URL (필수)
  ```bash
  --url https://example.com
  ```

#### 선택 옵션

- `--interval`, `-i`: 체크 주기 (기본값: `60s`)
  ```bash
  --interval 30s    # 30초마다 체크
  --interval 5m     # 5분마다 체크
  ```

- `--timeout`, `-t`: 요청 타임아웃 (기본값: `5s`)
  ```bash
  --timeout 10s     # 10초 타임아웃
  ```

- `--slack-webhook`, `-s`: Slack Webhook URL
  ```bash
  --slack-webhook https://hooks.slack.com/services/YOUR/WEBHOOK/URL
  ```

- `--discord-webhook`, `-d`: Discord Webhook URL
  ```bash
  --discord-webhook https://discord.com/api/webhooks/YOUR/WEBHOOK/URL
  ```

- `--test`: 테스트 모드 (정상 상태에서도 알림 전송)
  ```bash
  --test
  ```

### 환경 변수

플래그 대신 환경 변수를 사용할 수 있습니다:

- `SLACK_WEBHOOK_URL`: Slack Webhook URL
- `DISCORD_WEBHOOK_URL`: Discord Webhook URL

```bash
# Windows PowerShell
$env:SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
$env:DISCORD_WEBHOOK_URL="https://discord.com/api/webhooks/YOUR/WEBHOOK/URL"
.\health-checker.exe run --url https://example.com

# Linux/Mac
export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
export DISCORD_WEBHOOK_URL="https://discord.com/api/webhooks/YOUR/WEBHOOK/URL"
./health-checker run --url https://example.com
```

## 사용 예제

### 1. 기본 헬스 체크 (알림 없음)

```bash
.\health-checker.exe run --url https://example.com --interval 60s
```

### 2. Slack 알림만 사용

```bash
.\health-checker.exe run \
  --url https://example.com \
  --interval 30s \
  --slack-webhook https://hooks.slack.com/services/YOUR/WEBHOOK/URL
```

### 3. Discord 알림만 사용

```bash
.\health-checker.exe run \
  --url https://example.com \
  --interval 30s \
  --discord-webhook https://discord.com/api/webhooks/YOUR/WEBHOOK/URL
```

### 4. Slack과 Discord 동시 사용

```bash
.\health-checker.exe run \
  --url https://example.com \
  --interval 60s \
  --slack-webhook https://hooks.slack.com/services/YOUR/WEBHOOK/URL \
  --discord-webhook https://discord.com/api/webhooks/YOUR/WEBHOOK/URL
```

### 5. 테스트 모드 (정상 상태에서도 알림)

```bash
.\health-checker.exe run \
  --url https://example.com \
  --discord-webhook https://discord.com/api/webhooks/YOUR/WEBHOOK/URL \
  --test
```

### 6. 짧은 주기로 빠른 모니터링

```bash
.\health-checker.exe run \
  --url https://example.com \
  --interval 10s \
  --timeout 3s \
  --discord-webhook https://discord.com/api/webhooks/YOUR/WEBHOOK/URL
```

## 알림 조건

### 일반 모드 (기본)

다음 경우에 알림이 전송됩니다:

- HTTP 요청 실패 (연결 오류, 타임아웃 등)
- HTTP 상태 코드가 200이 아닌 경우 (500, 404, 503 등)

### 테스트 모드 (`--test` 플래그)

모든 상태에서 알림이 전송됩니다:

- 에러 발생 시: 장애 알림
- 200이 아닌 상태 코드: 장애 알림
- 200 OK: 정상 알림 (테스트 모드에서만)

## 알림 메시지 형식

### 장애 알림

```
🚨 사이트 장애 감지: https://example.com
상태 코드: 500
응답 시간: 2.5s
```

또는

```
🚨 사이트 장애 감지: https://example.com
에러: connection timeout
응답 시간: 5s
```

### 정상 알림 (테스트 모드)

```
✅ 사이트 정상: https://example.com
상태 코드: 200
응답 시간: 150ms
```

## Webhook 설정 방법

### Slack Webhook 설정

1. [Slack API](https://api.slack.com/apps)에서 새 앱 생성
2. Incoming Webhooks 활성화
3. Webhook URL 복사
4. `--slack-webhook` 플래그 또는 `SLACK_WEBHOOK_URL` 환경 변수에 설정

### Discord Webhook 설정

1. Discord 채널 설정 → 연동 → 웹후크
2. 새 웹후크 생성
3. 웹후크 URL 복사
4. `--discord-webhook` 플래그 또는 `DISCORD_WEBHOOK_URL` 환경 변수에 설정

## 종료 방법

프로그램을 종료하려면 `Ctrl+C`를 누르세요.

```
Press Ctrl+C to stop
^C
Shutting down...
```

## 프로젝트 구조

```
health-checker/
├── cmd/
│   ├── root.go      # 루트 커맨드
│   └── run.go       # run 서브커맨드
├── internal/
│   ├── checker/     # 헬스 체크 로직
│   │   └── checker.go
│   └── notifier/    # 알림 로직
│       ├── slack.go
│       ├── discord.go
│       └── notifier.go
├── main.go
├── go.mod
└── README.md
```

## 기술 스택

- **언어**: Go
- **CLI 프레임워크**: Cobra
- **HTTP 클라이언트**: net/http (표준 라이브러리)

## 라이선스

MIT License

## 기여

이슈나 Pull Request를 환영합니다!

