# Go 웹사이트 헬스 체커 + 슬랙 알림 봇 개발 계획

## 1. 목표 정의

### 1.1 최종 목표
- 특정 웹사이트(또는 여러 개)의 **가용성(Uptime)** 을 주기적으로 체크하고,
- 장애(접속 불가, 응답 코드 200이 아님, 응답 지연 등)를 감지하면
- **Slack 채널로 알림을 보내는 작은 SRE/DevOps 도구**를 Go로 만든다.

### 1.2 1차 버전 범위 (MVP)
- 단일 URL을 대상으로:
  - 일정 주기(예: 1분)에 한 번 `HTTP GET` 요청
  - 상태 코드, 응답 시간 측정
  - 실패 시 슬랙 Webhook으로 간단한 텍스트 알림 전송
- CLI 옵션으로:
  - `--url`, `--interval`, `--timeout` 정도만 지원

---

## 2. 기술 스택 및 구조

### 2.1 기술 스택
- 언어: **Go**
- CLI 프레임워크: **Cobra**
- HTTP 클라이언트: 표준 라이브러리 `net/http`
- 슬랙 연동: **Incoming Webhook** (표준 HTTP POST)
- 배포(선택): Docker + Cloudtype 또는 Fly.io

### 2.2 구성 요소
1. **CLI 진입점 (`main.go` + `cmd/root.go`)**
   - `health-checker` 실행 시 기본 명령과 플래그 처리
2. **체크 로직 패키지 (`internal/checker`)**
   - URL 체크, 응답 시간/상태 코드 반환
3. **알림 패키지 (`internal/notifier`)**
   - Slack Webhook 호출
4. **루프/스케줄러 (`cmd/run.go` or 내부 goroutine)**
   - `time.Ticker`로 주기적인 체크 수행

---

## 3. 기능 설계

### 3.1 기본 기능

#### 3.1.1 단일 URL 체크 기능
- 입력: URL 문자열
- 처리:
  - `http.Client`로 `GET` 요청
  - `timeout` 설정
  - 시작 시간 기록 후 요청 → 응답까지 경과 시간 계산
- 출력:
  - 상태 코드 (int)
  - 응답 시간 (`time.Duration`)
  - 에러 (`error`)

#### 3.1.2 CLI 인터페이스
명령 예시:
health-checker run
--url https://example.com
--interval 60s
--timeout 5s

text

지원 플래그:
- `--url` (필수)
- `--interval` (기본값: `60s`)
- `--timeout` (기본값: `5s`)
- `--slack-webhook` (선택, 기본은 환경변수로 처리)

### 3.2 슬랙 알림 기능

#### 3.2.1 알림 조건
- HTTP 요청이 실패 (`error != nil`)
- 상태 코드가 200이 아님 (예: 500, 404 등)
- 응답 시간이 threshold보다 길 때 (예: 3초 이상) → 추후 옵션화

#### 3.2.2 메시지 포맷
- 간단 버전:
  - `"🚨 사이트 장애 감지: https://example.com (status=500, latency=2.5s)"`
- 확장 가능 버전:
  - JSON payload에 `attachments`를 포함해 색상, 필드 등 추가

#### 3.2.3 Webhook 설계
- Webhook URL은 가능하면 **직접 커밋하지 않고**:
  - 환경 변수 `SLACK_WEBHOOK_URL` 에서 읽기
  - 또는 CLI 플래그 `--slack-webhook` 로 전달

---

## 4. 단계별 구현 로드맵

### 4.1 1단계: 프로젝트 뼈대 만들기

1. Go 모듈 초기화
mkdir go-health-checker
cd go-health-checker
go mod init github.com/<username>/go-health-checker

text
2. Cobra 초기화
go install github.com/spf13/cobra-cli@latest
cobra-cli init

text
3. `run` 서브커맨드 추가
cobra-cli add run

text

### 4.2 2단계: 단일 체크 함수 구현

- 파일: `internal/checker/checker.go`

구현 항목:
- `CheckURL(url string, timeout time.Duration) (status int, latency time.Duration, err error)`

테스트:
- `go test ./internal/checker` (간단한 유닛 테스트 작성 가능)

### 4.3 3단계: run 커맨드에 연결

- 파일: `cmd/run.go`
- 할 일:
- 플래그 정의:
 - `url`, `interval`, `timeout`
- `Run` 함수에서:
 - 플래그 값 읽기
 - `time.ParseDuration`으로 문자열을 `time.Duration`으로 변환
 - 단 한 번 `CheckURL` 호출 → 결과를 콘솔 출력

### 4.4 4단계: 주기적 실행 루프 추가

- `Run` 함수 수정:
- `time.NewTicker(interval)` 사용
- `for` 루프 안에서:
 - 각 tick마다 `CheckURL` 호출
 - 결과를 콘솔에 로깅

### 4.5 5단계: Slack Notifier 추가

1. 패키지 추가:
- `internal/notifier/slack.go`
2. 함수 설계:
- `SendSlack(webhookURL, message string) error`
3. `run` 커맨드에서:
- 슬랙 Webhook URL을 환경변수나 플래그로 읽기
- 체크 결과가 장애 조건일 때 `SendSlack` 호출

### 4.6 6단계: 설정 정리 및 사용성 개선

- 여러 URL 지원 (추후):
- `--url` 플래그를 여러 번 받거나, 설정 파일 읽기
- 로그 포맷 개선:
- `log` 패키지 또는 `zap` 같은 로거 사용
- 종료 처리:
- `os.Signal` (Ctrl+C) 핸들링해서 깔끔하게 종료

---

## 5. 배포 계획 (선택 사항)

### 5.1 Docker 이미지 만들기

- `Dockerfile` 작성:
- multi-stage build로 최종 이미지를 작게 유지
- 예:
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o health-checker .

FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/health-checker .
CMD ["./health-checker", "run", "--url", "https://example.com", "--interval", "60s"]

text

### 5.2 Cloudtype 또는 Fly.io에 배포

- Docker 이미지 푸시 (GitHub Container Registry 또는 Docker Hub)
- Cloudtype/Fly.io에서 컨테이너 기반 앱으로 배포
- 슬랙 Webhook URL은 **환경 변수**로 설정

---

## 6. 포트폴리오용 문서화 계획

### 6.1 README에 꼭 적을 내용
- 프로젝트 개요: "Go로 만든 웹사이트 Uptime 모니터링 + Slack 알림 봇"
- 사용 방법:
- 설치 방법 (`go install` or `docker run`)
- 예제 명령어
- 내부 구조:
- 패키지 구조, 핵심 흐름 다이어그램 간단히 설명
- 기술 포인트:
- `goroutine`, `time.Ticker`, `net/http`, Slack Webhook 활용
- 한계와 개선 여지:
- 분산 환경, 여러 인스턴스, 상태 저장, 대시보드 등 확장 아이디어

### 6.2 추가로 시도해 볼 확장 아이디어
- 여러 URL을 JSON 설정 파일로 관리
- 응답 코드, 응답 시간 등을 Prometheus Metrics로 노출 (`/metrics`)
- 간단한 웹 UI로 상태 보여주기 (Go + 템플릿)

---

## 7. 오늘 할 일 체크리스트

- [ ] Go 모듈 & Cobra 프로젝트 초기화
- [ ] `run` 커맨드 생성 및 `--url`, `--interval`, `--timeout` 플래그 추가
- [ ] `CheckURL` 기본 함수 구현 및 한 번 실행해서 상태/응답속도 출력
- [ ] Slack Webhook 발급 받아 `.env` 또는 환경 변수로 설정
- [ ] 장애 시 Slack으로 한 번이라도 알림 보내보기