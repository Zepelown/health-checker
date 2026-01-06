# 고급 헬스 체크 기능

## 개요

기본적인 HTTP GET 요청 외에, 더 복잡하고 세밀한 헬스 체크 기능을 제공합니다.

## 기능 요구사항

### 1. HTTP 메서드 커스터마이징

#### 지원 메서드
- GET (기본)
- POST
- PUT
- PATCH
- HEAD
- OPTIONS

#### 사용 예제
```yaml
endpoints:
  - name: "API Health"
    url: https://api.example.com/health
    method: POST
    body: '{"check": true}'
    headers:
      Content-Type: application/json
```

### 2. 커스텀 헤더 및 인증

#### 기본 인증
```yaml
endpoints:
  - name: "Protected API"
    url: https://api.example.com/health
    auth:
      type: basic
      username: admin
      password: secret
```

#### Bearer 토큰
```yaml
endpoints:
  - name: "Authenticated API"
    url: https://api.example.com/health
    headers:
      Authorization: "Bearer YOUR_TOKEN"
```

#### API 키
```yaml
endpoints:
  - name: "API Key Auth"
    url: https://api.example.com/health
    headers:
      X-API-Key: YOUR_API_KEY
```

### 3. 요청 본문 (Request Body)

#### JSON 본문
```yaml
endpoints:
  - name: "POST Health Check"
    url: https://api.example.com/health
    method: POST
    body: |
      {
        "service": "api",
        "check": true
      }
    headers:
      Content-Type: application/json
```

#### Form 데이터
```yaml
endpoints:
  - name: "Form Health Check"
    url: https://api.example.com/health
    method: POST
    body: |
      service=api&check=true
    headers:
      Content-Type: application/x-www-form-urlencoded
```

#### 파일 업로드
```yaml
endpoints:
  - name: "File Upload Check"
    url: https://api.example.com/upload
    method: POST
    body_file: test.json
    headers:
      Content-Type: application/json
```

### 4. 응답 검증

#### 상태 코드 검증
```yaml
endpoints:
  - name: "Custom Status"
    url: https://api.example.com/health
    expected_status: 201  # 200이 아닌 다른 상태 코드도 정상으로 간주
    # 또는
    expected_status_range: [200, 299]  # 2xx 모두 정상
```

#### 응답 본문 검증
```yaml
endpoints:
  - name: "Body Validation"
    url: https://api.example.com/health
    response_validation:
      # JSON 경로 검증
      json_path:
        - path: "$.status"
          expected: "healthy"
        - path: "$.version"
          expected: "1.0.0"
      
      # 정규식 검증
      regex:
        - pattern: '"status"\s*:\s*"healthy"'
      
      # 포함 여부 검증
      contains:
        - "healthy"
        - "ok"
      
      # 포함하지 않아야 함
      not_contains:
        - "error"
        - "down"
```

#### 응답 헤더 검증
```yaml
endpoints:
  - name: "Header Validation"
    url: https://api.example.com/health
    response_headers:
      X-Health-Check: "ok"
      Content-Type: "application/json"
```

### 5. SSL/TLS 검증

#### SSL 인증서 검증
```yaml
endpoints:
  - name: "SSL Check"
    url: https://api.example.com/health
    ssl:
      verify: true
      # 인증서 만료 임계값
      expiry_threshold: 30d  # 30일 이내 만료 시 경고
      # 자체 서명 인증서 허용
      allow_self_signed: false
      # 특정 CA 인증서 사용
      ca_cert: /path/to/ca.crt
```

#### TLS 설정
```yaml
endpoints:
  - name: "TLS Config"
    url: https://api.example.com/health
    tls:
      min_version: "1.2"
      max_version: "1.3"
      cipher_suites:
        - "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
```

### 6. DNS 체크

#### DNS 레코드 확인
```yaml
endpoints:
  - name: "DNS Check"
    url: https://api.example.com/health
    dns:
      # A 레코드 확인
      a_record: "192.168.1.1"
      # 또는
      a_record_any: true  # A 레코드가 존재하기만 하면 됨
      
      # CNAME 확인
      cname: "api.example.com"
      
      # MX 레코드 확인
      mx_record:
        - "mail.example.com"
      
      # DNS 서버 지정
      nameserver: "8.8.8.8"
```

### 7. 연결 체크 (TCP/UDP)

#### TCP 포트 체크
```yaml
endpoints:
  - name: "TCP Port Check"
    type: tcp
    host: api.example.com
    port: 8080
    timeout: 5s
```

#### UDP 포트 체크
```yaml
endpoints:
  - name: "UDP Port Check"
    type: udp
    host: api.example.com
    port: 53
    timeout: 5s
```

### 8. 데이터베이스 헬스 체크

#### PostgreSQL
```yaml
endpoints:
  - name: "PostgreSQL Health"
    type: postgres
    host: db.example.com
    port: 5432
    database: mydb
    username: user
    password: password
    query: "SELECT 1"
    timeout: 5s
```

#### MySQL
```yaml
endpoints:
  - name: "MySQL Health"
    type: mysql
    host: db.example.com
    port: 3306
    database: mydb
    username: user
    password: password
    query: "SELECT 1"
    timeout: 5s
```

#### Redis
```yaml
endpoints:
  - name: "Redis Health"
    type: redis
    host: redis.example.com
    port: 6379
    password: password
    command: "PING"
    expected_response: "PONG"
    timeout: 5s
```

### 9. 커스텀 체크 스크립트

#### 외부 스크립트 실행
```yaml
endpoints:
  - name: "Custom Script Check"
    type: script
    script: /path/to/check.sh
    timeout: 30s
    # 스크립트 종료 코드 0이면 성공
```

#### 인라인 스크립트
```yaml
endpoints:
  - name: "Inline Script"
    type: script
    script_content: |
      #!/bin/bash
      curl -f https://api.example.com/health || exit 1
    timeout: 30s
```

### 10. 체인 체크 (의존성 체크)

#### 체인 체크 정의
```yaml
endpoints:
  - name: "API Health"
    url: https://api.example.com/health
    depends_on:
      - "Database Health"
      - "Cache Health"
    # 의존성이 실패하면 이 체크도 실패로 간주
```

## 구현 계획

### 1. 체크 타입 인터페이스

```go
// internal/checker/checker.go
type Checker interface {
    Check(ctx context.Context) (CheckResult, error)
    Validate() error
}

type CheckResult struct {
    Success      bool
    StatusCode   int
    Latency      time.Duration
    ResponseBody []byte
    ResponseHeaders map[string]string
    Error        error
    Metadata     map[string]interface{}
}
```

### 2. 각 체크 타입 구현

#### HTTP 체크
```go
// internal/checker/http.go
type HTTPChecker struct {
    URL            string
    Method         string
    Headers        map[string]string
    Body           []byte
    ExpectedStatus int
    Timeout        time.Duration
    SSL            *SSLConfig
    Validation     *ResponseValidation
}

func (h *HTTPChecker) Check(ctx context.Context) (CheckResult, error) {
    // HTTP 요청 수행
    // 응답 검증
    // 결과 반환
}
```

#### TCP 체크
```go
// internal/checker/tcp.go
type TCPChecker struct {
    Host    string
    Port    int
    Timeout time.Duration
}

func (t *TCPChecker) Check(ctx context.Context) (CheckResult, error) {
    conn, err := net.DialTimeout("tcp", 
        fmt.Sprintf("%s:%d", t.Host, t.Port), 
        t.Timeout)
    if err != nil {
        return CheckResult{Success: false, Error: err}, err
    }
    defer conn.Close()
    return CheckResult{Success: true}, nil
}
```

#### 데이터베이스 체크
```go
// internal/checker/database.go
type DatabaseChecker struct {
    Type     string  // postgres, mysql, redis
    Host     string
    Port     int
    Database string
    Username string
    Password string
    Query    string
    Timeout  time.Duration
}

func (d *DatabaseChecker) Check(ctx context.Context) (CheckResult, error) {
    // 데이터베이스 연결 및 쿼리 실행
}
```

### 3. 응답 검증 엔진

```go
// internal/checker/validation.go
type ResponseValidation struct {
    JSONPath    []JSONPathCheck
    Regex       []string
    Contains    []string
    NotContains []string
    Headers     map[string]string
}

type JSONPathCheck struct {
    Path     string
    Expected interface{}
}

func (v *ResponseValidation) Validate(body []byte, headers map[string]string) error {
    // JSON 경로 검증
    // 정규식 검증
    // 포함 여부 검증
    // 헤더 검증
}
```

### 4. SSL 검증

```go
// internal/checker/ssl.go
type SSLConfig struct {
    Verify            bool
    ExpiryThreshold   time.Duration
    AllowSelfSigned   bool
    CACert            string
    MinVersion        string
    MaxVersion        string
    CipherSuites      []string
}

func CheckSSL(url string, config SSLConfig) (SSLInfo, error) {
    // SSL 인증서 정보 조회
    // 만료일 확인
    // TLS 버전 확인
}
```

## 사용 예제

### 예제 1: POST 요청으로 헬스 체크
```bash
health-checker run \
  --url https://api.example.com/health \
  --method POST \
  --body '{"check": true}' \
  --header "Content-Type: application/json" \
  --header "Authorization: Bearer token"
```

### 예제 2: 응답 본문 검증
```yaml
# config.yaml
endpoints:
  - name: "API with Validation"
    url: https://api.example.com/health
    response_validation:
      json_path:
        - path: "$.status"
          expected: "healthy"
      contains:
        - "ok"
```

```bash
health-checker monitor --config config.yaml
```

### 예제 3: SSL 인증서 만료 체크
```yaml
endpoints:
  - name: "SSL Expiry Check"
    url: https://api.example.com
    ssl:
      expiry_threshold: 30d
```

### 예제 4: TCP 포트 체크
```yaml
endpoints:
  - name: "Database Port"
    type: tcp
    host: db.example.com
    port: 5432
    timeout: 5s
```

### 예제 5: 데이터베이스 체크
```yaml
endpoints:
  - name: "PostgreSQL"
    type: postgres
    host: db.example.com
    port: 5432
    database: mydb
    username: user
    password: ${DB_PASSWORD}
    query: "SELECT 1"
```

## 고려사항

### 1. 보안
- 비밀번호/토큰 안전한 저장
- 환경 변수 사용
- 암호화된 설정 파일

### 2. 성능
- 복잡한 검증 로직의 오버헤드
- 대량의 응답 본문 처리
- 타임아웃 관리

### 3. 유연성
- 다양한 체크 타입 지원
- 확장 가능한 검증 시스템
- 플러그인 아키텍처

## 향후 확장 가능성

1. **gRPC 헬스 체크**: gRPC 서비스 헬스 체크
2. **GraphQL 헬스 체크**: GraphQL 쿼리 기반 체크
3. **웹소켓 체크**: 웹소켓 연결 및 메시지 교환
4. **메시지 큐 체크**: RabbitMQ, Kafka 등
5. **컨테이너 체크**: Docker/Kubernetes 컨테이너 상태
6. **클라우드 서비스 체크**: AWS, GCP, Azure 서비스 상태


