# 설정 파일 기반 관리

## 개요

명령줄 플래그 대신 설정 파일을 사용하여 여러 엔드포인트와 복잡한 설정을 관리할 수 있도록 합니다.

## 기능 요구사항

### 1. 설정 파일 형식

#### YAML 형식 (권장)
- 가독성 우수
- 주석 지원
- 중첩 구조 표현 용이

#### JSON 형식
- 기계적 파싱 용이
- 프로그래밍 언어 간 호환성

#### TOML 형식
- 간단한 구조
- 주석 지원

### 2. 설정 파일 구조

#### 기본 구조
```yaml
# health-checker.yaml
version: "1.0"

# 전역 설정 (모든 엔드포인트에 기본값으로 적용)
global:
  interval: 60s
  timeout: 5s
  latency_threshold: 3s
  slack_webhook: https://hooks.slack.com/services/GLOBAL/WEBHOOK
  discord_webhook: https://discord.com/api/webhooks/GLOBAL/WEBHOOK

# 엔드포인트 목록
endpoints:
  - name: "Production API"
    url: https://api.prod.example.com/health
    interval: 30s
    timeout: 5s
    latency_threshold: 1s
    enabled: true
    
    # HTTP 설정
    method: GET
    headers:
      Authorization: "Bearer ${API_TOKEN}"
      X-Custom-Header: "value"
    body: ""
    
    # 응답 검증
    expected_status: 200
    response_validation:
      json_path:
        - path: "$.status"
          expected: "healthy"
    
    # SSL 설정
    ssl:
      verify: true
      expiry_threshold: 30d
    
    # 알림 설정
    notifications:
      slack_webhook: https://hooks.slack.com/services/API/WEBHOOK
      discord_webhook: https://discord.com/api/webhooks/API/WEBHOOK
      rules:
        - name: "Critical"
          condition: status >= 500
          channels: [slack, email]
  
  - name: "Staging API"
    url: https://api.staging.example.com/health
    interval: 60s
    enabled: true
  
  - name: "Database Health"
    type: postgres
    host: db.example.com
    port: 5432
    database: mydb
    username: ${DB_USER}
    password: ${DB_PASSWORD}
    query: "SELECT 1"
    enabled: false  # 일시적으로 비활성화

# 알림 채널 설정
notifications:
  channels:
    slack:
      webhook: https://hooks.slack.com/services/YOUR/WEBHOOK
      enabled: true
    
    discord:
      webhook: https://discord.com/api/webhooks/YOUR/WEBHOOK
      enabled: true
    
    email:
      smtp:
        host: smtp.gmail.com
        port: 587
        username: ${EMAIL_USER}
        password: ${EMAIL_PASSWORD}
      to:
        - admin@example.com
        - ops@example.com
      enabled: true
    
    pagerduty:
      integration_key: ${PAGERDUTY_KEY}
      enabled: false
  
  # 전역 알림 규칙
  rules:
    - name: "All Critical Failures"
      condition: status >= 500 OR consecutive_failures >= 3
      channels: [slack, email, pagerduty]
      enabled: true

# 메트릭 설정
metrics:
  enabled: true
  storage:
    type: file  # file, sqlite, postgres
    path: ./metrics.json
    retention: 30d  # 30일간 보관
  
  prometheus:
    enabled: true
    port: 9090
    path: /metrics

# 대시보드 설정
dashboard:
  enabled: true
  port: 8080
  auth:
    enabled: true
    username: admin
    password: ${DASHBOARD_PASSWORD}
  theme: light  # light, dark

# 로깅 설정
logging:
  level: info  # debug, info, warn, error
  format: json  # text, json
  output: stdout  # stdout, file, both
  file: ./health-checker.log
  max_size: 100MB
  max_backups: 10
  max_age: 30d
```

### 3. 환경 변수 지원

#### 변수 치환
```yaml
endpoints:
  - name: "API"
    url: https://api.example.com
    headers:
      Authorization: "Bearer ${API_TOKEN}"
    # 또는
    headers:
      Authorization: "Bearer ${API_TOKEN:-default_token}"  # 기본값
```

#### 환경 변수 파일 (.env)
```bash
# .env
API_TOKEN=your_token_here
DB_USER=dbuser
DB_PASSWORD=dbpassword
EMAIL_USER=email@example.com
EMAIL_PASSWORD=emailpassword
PAGERDUTY_KEY=pagerduty_key
DASHBOARD_PASSWORD=admin123
```

### 4. 설정 파일 검증

#### 스키마 검증
- 설정 파일 형식 검증
- 필수 필드 확인
- 타입 검증
- 값 범위 검증

#### 예제
```bash
health-checker validate --config config.yaml
```

출력:
```
✓ Configuration file is valid
✓ Found 3 endpoints
✓ All notification channels are configured
⚠ Warning: Database Health endpoint is disabled
```

### 5. 설정 파일 템플릿

#### 기본 템플릿 생성
```bash
health-checker init --config config.yaml
```

생성된 파일:
```yaml
# health-checker.yaml
# This is a template configuration file
# Fill in your endpoints and settings

version: "1.0"

global:
  interval: 60s
  timeout: 5s

endpoints:
  - name: "Example Endpoint"
    url: https://example.com
    interval: 60s
    enabled: true
```

### 6. 설정 파일 병합

#### 여러 설정 파일 병합
```bash
health-checker monitor \
  --config base.yaml \
  --config overrides.yaml \
  --config local.yaml
```

나중에 로드된 파일이 이전 설정을 덮어씁니다.

### 7. 동적 설정 리로드

#### 핫 리로드
- 설정 파일 변경 감지
- 런타임에 설정 업데이트
- 엔드포인트 추가/제거/수정

```bash
health-checker monitor --config config.yaml --watch
```

시그널 기반 리로드:
```bash
# SIGHUP 시그널로 설정 리로드
kill -HUP <pid>
```

## 구현 계획

### 1. 설정 구조체

```go
// internal/config/config.go
type Config struct {
    Version    string                 `yaml:"version"`
    Global     GlobalConfig           `yaml:"global"`
    Endpoints  []EndpointConfig       `yaml:"endpoints"`
    Notifications NotificationsConfig `yaml:"notifications"`
    Metrics    MetricsConfig          `yaml:"metrics"`
    Dashboard  DashboardConfig        `yaml:"dashboard"`
    Logging    LoggingConfig          `yaml:"logging"`
}

type GlobalConfig struct {
    Interval        string            `yaml:"interval"`
    Timeout         string            `yaml:"timeout"`
    LatencyThreshold string          `yaml:"latency_threshold"`
    SlackWebhook    string            `yaml:"slack_webhook"`
    DiscordWebhook  string            `yaml:"discord_webhook"`
}

type EndpointConfig struct {
    Name             string            `yaml:"name"`
    URL              string            `yaml:"url"`
    Type             string            `yaml:"type"`  // http, tcp, postgres, etc.
    Interval         string            `yaml:"interval"`
    Timeout          string            `yaml:"timeout"`
    LatencyThreshold string            `yaml:"latency_threshold"`
    Enabled          bool              `yaml:"enabled"`
    Method           string            `yaml:"method"`
    Headers          map[string]string `yaml:"headers"`
    Body             string            `yaml:"body"`
    ExpectedStatus   int               `yaml:"expected_status"`
    ResponseValidation *ResponseValidation `yaml:"response_validation"`
    SSL              *SSLConfig        `yaml:"ssl"`
    Notifications    *EndpointNotifications `yaml:"notifications"`
}
```

### 2. 설정 로더

```go
// internal/config/loader.go
type Loader struct {
    paths []string
    env   map[string]string
}

func (l *Loader) Load(path string) (*Config, error) {
    // 파일 읽기
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    
    // 환경 변수 치환
    data = l.substituteEnvVars(data)
    
    // YAML 파싱
    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, err
    }
    
    // 검증
    if err := l.Validate(&config); err != nil {
        return nil, err
    }
    
    // 기본값 채우기
    l.setDefaults(&config)
    
    return &config, nil
}

func (l *Loader) Merge(configs ...*Config) *Config {
    // 여러 설정 파일 병합
}
```

### 3. 환경 변수 치환

```go
// internal/config/env.go
func substituteEnvVars(data []byte) []byte {
    // ${VAR} 또는 ${VAR:-default} 형식 치환
    re := regexp.MustCompile(`\$\{([^}]+)\}`)
    return re.ReplaceAllFunc(data, func(match []byte) []byte {
        // 환경 변수 추출 및 치환
    })
}
```

### 4. 설정 검증

```go
// internal/config/validator.go
type Validator struct{}

func (v *Validator) Validate(config *Config) error {
    // 버전 확인
    if config.Version == "" {
        return errors.New("version is required")
    }
    
    // 엔드포인트 검증
    for i, endpoint := range config.Endpoints {
        if err := v.validateEndpoint(&endpoint); err != nil {
            return fmt.Errorf("endpoint[%d]: %w", i, err)
        }
    }
    
    // 알림 채널 검증
    if err := v.validateNotifications(&config.Notifications); err != nil {
        return err
    }
    
    return nil
}
```

### 5. 동적 리로드

```go
// internal/config/watcher.go
type Watcher struct {
    path   string
    config *Config
    mutex  sync.RWMutex
    reload chan struct{}
}

func (w *Watcher) Watch() error {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return err
    }
    defer watcher.Close()
    
    if err := watcher.Add(w.path); err != nil {
        return err
    }
    
    for {
        select {
        case event := <-watcher.Events:
            if event.Op&fsnotify.Write == fsnotify.Write {
                w.reloadConfig()
            }
        case err := <-watcher.Errors:
            log.Printf("watcher error: %v", err)
        }
    }
}
```

## 사용 예제

### 예제 1: 기본 설정 파일 사용
```bash
# config.yaml 생성 후
health-checker monitor --config config.yaml
```

### 예제 2: 환경 변수와 함께 사용
```bash
export API_TOKEN=your_token
export DB_PASSWORD=secret
health-checker monitor --config config.yaml
```

### 예제 3: .env 파일 사용
```bash
# .env 파일 생성
cat > .env << EOF
API_TOKEN=your_token
DB_PASSWORD=secret
EOF

health-checker monitor --config config.yaml --env .env
```

### 예제 4: 설정 파일 검증
```bash
health-checker validate --config config.yaml
```

### 예제 5: 설정 파일 템플릿 생성
```bash
health-checker init --config config.yaml
```

### 예제 6: 동적 리로드
```bash
health-checker monitor --config config.yaml --watch
# 다른 터미널에서 config.yaml 수정 시 자동 리로드
```

### 예제 7: 여러 설정 파일 병합
```bash
health-checker monitor \
  --config base.yaml \
  --config env-specific.yaml \
  --config local-overrides.yaml
```

## 고려사항

### 1. 보안
- 비밀번호/토큰 평문 저장 방지
- 환경 변수 사용 권장
- 설정 파일 권한 관리
- 암호화된 설정 파일 지원 (선택 사항)

### 2. 유효성 검증
- 스키마 검증
- 타입 검증
- 값 범위 검증
- 의존성 검증

### 3. 호환성
- 버전 관리
- 하위 호환성
- 마이그레이션 도구

### 4. 성능
- 대용량 설정 파일 처리
- 파싱 성능
- 메모리 사용량

## 향후 확장 가능성

1. **원격 설정**: HTTP/HTTPS에서 설정 파일 로드
2. **설정 버전 관리**: Git과 통합하여 설정 변경 이력 추적
3. **설정 템플릿**: 재사용 가능한 설정 템플릿
4. **설정 UI**: 웹 대시보드에서 설정 편집
5. **설정 검증 API**: 설정 파일 검증을 위한 별도 API
6. **설정 백업/복원**: 설정 파일 자동 백업


