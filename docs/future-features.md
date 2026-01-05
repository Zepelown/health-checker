## 1. 다중 엔드포인트 모니터링 (가장 자연스러운 확장)
`health-checker monitor --config endpoints.yaml`
- 여러 URL을 동시에 모니터링
- 각 엔드포인트별 독립적인 헬스 체크
- 헬스 체커의 핵심 기능 확장
## 2. 상세 메트릭 및 통계
- 응답 시간 히스토리
- 가용성 통계 (99.9% SLA 등)
- 에러율 추적
- Prometheus 메트릭 노출
## 3. 대시보드/UI
- 웹 대시보드로 모든 엔드포인트 상태 시각화
- 실시간 상태 모니터링
- 히스토리 차트
## 4. 알림 채널 확장
- 이메일 알림
- PagerDuty 연동
- 커스텀 웹훅
- 알림 규칙 세분화 (중요도별, 시간대별)
## 5. 고급 헬스 체크 기능
- HTTP 메서드 커스터마이징 (POST, PUT 등)
- 커스텀 헤더/인증
- 응답 본문 검증
- SSL 인증서 만료 체크
- DNS 체크
## 6. 설정 파일 기반 관리
```
endpoints:
  - name: "API Server"
    url: https://api.example.com
    interval: 30s
    timeout: 5s
    expected_status: 200
    latency_threshold: 1s
  
  - name: "Database Health"
    url: https://db.example.com/health
    interval: 60s
    method: GET
    headers:
      Authorization: "Bearer token"
```
