# 대시보드/UI

## 개요

CLI 기반의 텍스트 로그 대신, 웹 기반 대시보드를 제공하여 여러 엔드포인트의 상태를 한눈에 확인하고 시각화할 수 있도록 합니다.

## 기능 요구사항

### 1. 웹 대시보드

#### 기본 기능
- 실시간 상태 모니터링
- 여러 엔드포인트 동시 표시
- 상태별 색상 코딩 (정상/경고/장애)
- 응답 시간 시각화
- 가용성 통계 표시

#### 고급 기능
- 히스토리 차트 (응답 시간, 가용성)
- 필터링 및 검색
- 엔드포인트 그룹화
- 알림 이력
- 설정 관리

### 2. UI 디자인

#### 메인 대시보드
```
┌─────────────────────────────────────────────────────────────┐
│ Health Checker Dashboard                    [Settings] [Log]│
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │ API Server   │  │ Web Server   │  │ DB Health    │     │
│  │ ✅ UP        │  │ ✅ UP        │  │ ⚠️ DEGRADED  │     │
│  │ 150ms        │  │ 200ms        │  │ 1.2s         │     │
│  │ 99.9% uptime │  │ 99.5% uptime │  │ 98.0% uptime │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
│                                                              │
│  Response Time (Last 24h)                                   │
│  ┌────────────────────────────────────────────────────┐   │
│  │     ╱╲                                              │   │
│  │    ╱  ╲    ╱╲                                       │   │
│  │   ╱    ╲  ╱  ╲                                      │   │
│  │  ╱      ╲╱    ╲                                     │   │
│  └────────────────────────────────────────────────────┘   │
│                                                              │
│  Availability (Last 7 days)                                 │
│  ┌────────────────────────────────────────────────────┐   │
│  │ 100%│ ████████████████████████████████████████     │   │
│  │  99%│ ████████████████████████████████████████     │   │
│  │  98%│ ████████████████████████████████████████     │   │
│  └────────────────────────────────────────────────────┘   │
│                                                              │
│  Recent Events                                              │
│  [10:30:00] ✅ API Server - 200 (150ms)                    │
│  [10:29:30] ⚠️  DB Health - Slow response (1.2s)           │
│  [10:29:00] ✅ Web Server - 200 (200ms)                    │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 3. 기술 스택 옵션

#### 옵션 1: Go 템플릿 + HTML/CSS/JavaScript
- **장점**: 
  - Go 표준 라이브러리만 사용
  - 가볍고 빠름
  - 의존성 최소화
- **단점**: 
  - 프론트엔드 개발 복잡도 증가
  - 현대적인 UI 구현 어려움

#### 옵션 2: Go + React/Vue (권장)
- **장점**: 
  - 현대적인 UI/UX
  - 컴포넌트 기반 개발
  - 풍부한 차트 라이브러리
- **단점**: 
  - 빌드 프로세스 복잡
  - 프론트엔드 번들 크기

#### 옵션 3: Go + HTMX
- **장점**: 
  - 서버 사이드 렌더링
  - 간단한 구현
  - 현대적인 인터랙션
- **단점**: 
  - 복잡한 UI 구현 제한적

### 4. API 엔드포인트

#### RESTful API 설계

```go
// 상태 조회
GET /api/status
Response: {
  "endpoints": [
    {
      "name": "API Server",
      "url": "https://api.example.com",
      "status": "up",
      "status_code": 200,
      "latency_ms": 150,
      "last_check": "2025-01-15T10:30:00Z",
      "availability_24h": 99.9
    }
  ]
}

// 통계 조회
GET /api/stats?url=https://example.com&period=24h
Response: {
  "url": "https://example.com",
  "period": "24h",
  "availability": 99.9,
  "average_latency_ms": 150,
  "total_checks": 1440,
  "successful": 1435,
  "failed": 5
}

// 히스토리 조회
GET /api/history?url=https://example.com&from=2025-01-15T00:00:00Z&to=2025-01-15T23:59:59Z
Response: {
  "url": "https://example.com",
  "checks": [
    {
      "timestamp": "2025-01-15T10:00:00Z",
      "status": 200,
      "latency_ms": 150,
      "success": true
    }
  ]
}

// WebSocket (실시간 업데이트)
WS /ws/status
Message: {
  "type": "status_update",
  "endpoint": "API Server",
  "status": "up",
  "latency_ms": 150,
  "timestamp": "2025-01-15T10:30:00Z"
}
```

### 5. 페이지 구성

#### 메인 페이지 (`/`)
- 전체 엔드포인트 상태 개요
- 주요 메트릭 요약
- 최근 이벤트

#### 엔드포인트 상세 페이지 (`/endpoint/:name`)
- 개별 엔드포인트 상세 정보
- 응답 시간 차트
- 가용성 차트
- 상태 코드 분포
- 이벤트 히스토리

#### 통계 페이지 (`/stats`)
- 전체 통계 요약
- 비교 분석
- 트렌드 분석

#### 설정 페이지 (`/settings`)
- 엔드포인트 추가/수정/삭제
- 알림 설정
- 대시보드 설정

### 6. 실시간 업데이트

#### 폴링 방식
- 클라이언트가 주기적으로 API 호출
- 구현 간단
- 서버 부하 증가

#### WebSocket 방식 (권장)
- 서버에서 클라이언트로 푸시
- 실시간성 우수
- 서버 부하 감소

## 구현 계획

### 1. 서버 구조

```go
// cmd/dashboard.go
type DashboardServer struct {
    port     int
    checker  *checker.MultiChecker
    metrics  *metrics.Collector
    router   *http.ServeMux
}

func (s *DashboardServer) Start() error {
    // 정적 파일 서빙
    s.router.Handle("/", http.FileServer(http.Dir("./web/dist")))
    
    // API 엔드포인트
    s.router.HandleFunc("/api/status", s.handleStatus)
    s.router.HandleFunc("/api/stats", s.handleStats)
    s.router.HandleFunc("/api/history", s.handleHistory)
    
    // WebSocket
    s.router.HandleFunc("/ws/status", s.handleWebSocket)
    
    return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router)
}
```

### 2. 프론트엔드 구조 (React 예시)

```
web/
├── public/
│   └── index.html
├── src/
│   ├── components/
│   │   ├── Dashboard.tsx
│   │   ├── EndpointCard.tsx
│   │   ├── StatusChart.tsx
│   │   └── EventList.tsx
│   ├── services/
│   │   └── api.ts
│   ├── hooks/
│   │   └── useWebSocket.ts
│   └── App.tsx
└── package.json
```

### 3. 빌드 프로세스

```makefile
# Makefile
.PHONY: build-web
build-web:
	cd web && npm install && npm run build

.PHONY: embed-web
embed-web:
	go run tools/embed.go

.PHONY: build
build: build-web embed-web
	go build -o health-checker
```

또는 Go 1.16+ embed 사용:

```go
//go:embed web/dist/*
var webFiles embed.FS
```

## 사용 예제

### 예제 1: 기본 대시보드 시작
```bash
health-checker dashboard \
  --port 8080 \
  --config endpoints.yaml
```

브라우저에서 `http://localhost:8080` 접속

### 예제 2: 대시보드와 모니터링 동시 실행
```bash
health-checker monitor \
  --config endpoints.yaml \
  --dashboard \
  --dashboard-port 8080
```

### 예제 3: 인증 추가
```bash
health-checker dashboard \
  --port 8080 \
  --auth \
  --username admin \
  --password secret
```

### 예제 4: 커스텀 테마
```bash
health-checker dashboard \
  --port 8080 \
  --theme dark
```

## 보안 고려사항

### 1. 인증 및 권한
- 기본 인증 (Basic Auth)
- JWT 토큰 기반 인증
- 역할 기반 접근 제어 (RBAC)

### 2. HTTPS 지원
```bash
health-checker dashboard \
  --port 8080 \
  --tls \
  --tls-cert cert.pem \
  --tls-key key.pem
```

### 3. CORS 설정
- 개발 환경: 모든 origin 허용
- 프로덕션: 특정 origin만 허용

### 4. Rate Limiting
- API 엔드포인트별 요청 제한
- DDoS 방어

## 고려사항

### 1. 성능
- 대량의 엔드포인트 처리
- 실시간 업데이트 성능
- 메모리 사용량

### 2. 확장성
- 수평 확장 지원
- 로드 밸런싱
- 세션 관리

### 3. 사용자 경험
- 반응형 디자인 (모바일 지원)
- 빠른 로딩 시간
- 직관적인 UI/UX

## 향후 확장 가능성

1. **다크 모드**: 테마 전환 기능
2. **알림 설정 UI**: 웹에서 알림 규칙 설정
3. **엔드포인트 관리 UI**: 웹에서 엔드포인트 추가/수정/삭제
4. **리포트 생성**: PDF/CSV 리포트 다운로드
5. **모바일 앱**: React Native 기반 모바일 앱
6. **플러그인 시스템**: 커스텀 위젯 및 차트

