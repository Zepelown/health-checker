# 배포 가이드

외부 서버에서 Health Checker를 실행하는 방법을 안내합니다.

## 목차

1. [환경 변수 설정](#환경-변수-설정)
2. [Docker를 사용한 배포](#docker를-사용한-배포)
3. [Linux 서버에 직접 배포](#linux-서버에-직접-배포)
4. [systemd 서비스로 실행](#systemd-서비스로-실행)

---

## 환경 변수 설정

프로젝트 루트에 `.env.example` 파일을 참고하여 `.env` 파일을 생성하세요.

```bash
# .env.example을 복사
cp .env.example .env

# .env 파일을 편집하여 실제 값 입력
nano .env
```

필요한 환경 변수:
- `URL`: 체크할 URL (필수)
- `SLACK_WEBHOOK_URL`: Slack Webhook URL (선택)
- `DISCORD_WEBHOOK_URL`: Discord Webhook URL (선택)
- `LATENCY_THRESHOLD`: 응답 지연 임계값 (선택, 예: `3s`, `500ms`)

---

## Docker를 사용한 배포

### 방법 1: docker-compose 사용 (권장)

1. `.env` 파일 생성 및 설정:
```bash
cp .env.example .env
# .env 파일 편집
nano .env
```

2. docker-compose로 실행:
```bash
docker-compose up -d
```

3. 로그 확인:
```bash
docker-compose logs -f
```

4. 중지:
```bash
docker-compose down
```

### 방법 2: Docker 직접 사용

1. 이미지 빌드:
```bash
docker build -t health-checker .
```

2. 컨테이너 실행:
```bash
docker run -d \
  --name health-checker \
  --restart unless-stopped \
  -e URL=https://example.com \
  -e INTERVAL=60s \
  -e TIMEOUT=5s \
  -e SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL \
  -e DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/YOUR/WEBHOOK/URL \
  -e LATENCY_THRESHOLD=3s \
  health-checker
```

3. 로그 확인:
```bash
docker logs -f health-checker
```

4. 중지:
```bash
docker stop health-checker
docker rm health-checker
```

### 환경 변수 파일 사용

`.env` 파일을 사용하여 환경 변수를 관리할 수 있습니다:

```bash
docker run -d \
  --name health-checker \
  --restart unless-stopped \
  --env-file .env \
  health-checker
```

---

## Linux 서버에 직접 배포

### 1. 바이너리 빌드

로컬에서 빌드:
```bash
# Linux용 바이너리 빌드
GOOS=linux GOARCH=amd64 go build -o health-checker-linux

# 또는 서버에서 직접 빌드
go build -o health-checker
```

### 2. 서버에 업로드

```bash
# SCP를 사용하여 서버에 업로드
scp health-checker-linux user@server:/opt/health-checker/health-checker

# 또는 rsync 사용
rsync -avz health-checker-linux user@server:/opt/health-checker/health-checker
```

### 3. 실행 권한 부여

```bash
ssh user@server
chmod +x /opt/health-checker/health-checker
```

### 4. 환경 변수 설정

서버에서 환경 변수를 설정합니다:

```bash
# ~/.bashrc 또는 ~/.profile에 추가
export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
export DISCORD_WEBHOOK_URL="https://discord.com/api/webhooks/YOUR/WEBHOOK/URL"
export LATENCY_THRESHOLD="3s"

# 또는 /etc/health-checker/env 파일 생성
sudo mkdir -p /etc/health-checker
sudo nano /etc/health-checker/env
# 위 환경 변수들을 입력
```

### 5. 수동 실행

```bash
/opt/health-checker/health-checker run \
  --url https://example.com \
  --interval 60s \
  --timeout 5s
```

---

## systemd 서비스로 실행

### 1. 서비스 파일 설치

```bash
# 서비스 파일 복사
sudo cp health-checker.service /etc/systemd/system/

# 환경 변수 파일 생성 (선택 사항)
sudo mkdir -p /etc/health-checker
sudo nano /etc/health-checker/env
# SLACK_WEBHOOK_URL=...
# DISCORD_WEBHOOK_URL=...
# LATENCY_THRESHOLD=...
```

### 2. 서비스 파일 수정

`/etc/systemd/system/health-checker.service` 파일을 편집하여 URL과 환경 변수를 설정합니다:

```bash
sudo nano /etc/systemd/system/health-checker.service
```

주요 수정 사항:
- `ExecStart`의 `%i` 부분을 실제 URL로 변경하거나, `--url https://example.com` 형식으로 수정
- `Environment` 섹션에 실제 webhook URL 설정
- 또는 `EnvironmentFile=/etc/health-checker/env` 주석 해제

### 3. 전용 사용자 생성 (선택 사항, 권장)

```bash
sudo useradd -r -s /bin/false health-checker
sudo chown -R health-checker:health-checker /opt/health-checker
```

### 4. 서비스 시작

```bash
# systemd 재로드
sudo systemctl daemon-reload

# 서비스 시작
sudo systemctl start health-checker

# 부팅 시 자동 시작 설정
sudo systemctl enable health-checker

# 상태 확인
sudo systemctl status health-checker

# 로그 확인
sudo journalctl -u health-checker -f
```

### 5. 서비스 관리

```bash
# 서비스 중지
sudo systemctl stop health-checker

# 서비스 재시작
sudo systemctl restart health-checker

# 서비스 상태 확인
sudo systemctl status health-checker

# 로그 확인
sudo journalctl -u health-checker -f
```

---

## 여러 URL 모니터링

여러 URL을 모니터링하려면 각각 별도의 서비스 인스턴스를 실행하세요:

### Docker Compose 예제

```yaml
version: '3.8'

services:
  health-checker-site1:
    build: .
    container_name: health-checker-site1
    restart: unless-stopped
    environment:
      - URL=https://site1.com
      - SLACK_WEBHOOK_URL=${SLACK_WEBHOOK_URL}
    command: ./health-checker run --url https://site1.com --interval 60s

  health-checker-site2:
    build: .
    container_name: health-checker-site2
    restart: unless-stopped
    environment:
      - URL=https://site2.com
      - SLACK_WEBHOOK_URL=${SLACK_WEBHOOK_URL}
    command: ./health-checker run --url https://site2.com --interval 120s
```

### systemd 예제

각 URL마다 별도의 서비스 파일을 생성:

```bash
sudo cp health-checker.service /etc/systemd/system/health-checker-site1.service
sudo cp health-checker.service /etc/systemd/system/health-checker-site2.service

# 각 서비스 파일에서 URL 수정
sudo nano /etc/systemd/system/health-checker-site1.service
sudo nano /etc/systemd/system/health-checker-site2.service
```

---

## 문제 해결

### 로그 확인

**Docker:**
```bash
docker logs -f health-checker
```

**systemd:**
```bash
sudo journalctl -u health-checker -f
```

### 환경 변수 확인

**Docker:**
```bash
docker exec health-checker env
```

**systemd:**
```bash
sudo systemctl show health-checker --property=Environment
```

### 네트워크 연결 확인

서버에서 대상 URL에 접근 가능한지 확인:
```bash
curl -I https://example.com
```

---

## 보안 고려사항

1. **환경 변수 보호**: `.env` 파일이나 webhook URL을 Git에 커밋하지 마세요.
2. **파일 권한**: 바이너리와 설정 파일의 권한을 적절히 설정하세요.
3. **방화벽**: 필요한 경우 방화벽 규칙을 설정하세요.

---

## 참고 자료

- [README.md](../README.md) - 기본 사용법
- [.env.example](../.env.example) - 환경 변수 예제
