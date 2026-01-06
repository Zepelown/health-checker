# Health Checker

[한국어](README.ko.md) | [English](../README.md)

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

- `--latency-threshold`: 응답 지연 임계값 (지정 시, 이를 초과하는 응답도 장애로 간주)
  ```bash
  --latency-threshold 3s     # 3초 이상 걸리면 장애 알림
  --latency-threshold 500ms  # 500ms 이상 걸리면 장애 알림
  ```

- `--test`: 테스트 모드 (정상 상태에서도 알림 전송)
  ```bash
  --test
  ```

모든 duration 관련 옵션(`--interval`, `--timeout`, `--latency-threshold` 및 관련 환경 변수)은 Go의 `time.ParseDuration` 형식을 따르며, 다음 단위들을 지원합니다:

- `ns` (나노초), `us`/`µs` (마이크로초), `ms` (밀리초)
- `s` (초), `m` (분), `h` (시간)

예: `500ms`, `2s`, `1.5s`, `3m`, `1h30m`

### 환경 변수

플래그 대신 환경 변수를 사용할 수 있습니다:

- `SLACK_WEBHOOK_URL`: Slack Webhook URL
- `DISCORD_WEBHOOK_URL`: Discord Webhook URL
- `LATENCY_THRESHOLD`: 응답 지연 임계값 (예: `3s`, `500ms`)

```bash
# Windows PowerShell
$env:SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
$env:DISCORD_WEBHOOK_URL="https://discord.com/api/webhooks/YOUR/WEBHOOK/URL"
$env:LATENCY_THRESHOLD="3s"
.\health-checker.exe run --url https://example.com

# Linux/Mac
export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
export DISCORD_WEBHOOK_URL="https://discord.com/api/webhooks/YOUR/WEBHOOK/URL"
export LATENCY_THRESHOLD="3s"
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
  --latency-threshold 2s \
  --discord-webhook https://discord.com/api/webhooks/YOUR/WEBHOOK/URL
```

## 알림 조건

### 일반 모드 (기본)

다음 경우에 알림이 전송됩니다:

- HTTP 요청 실패 (연결 오류, 타임아웃 등)
- HTTP 상태 코드가 200이 아닌 경우 (500, 404, 503 등)
- (선택) `--latency-threshold` 또는 `LATENCY_THRESHOLD`가 설정된 경우, 응답 시간이 임계값을 초과할 때

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


