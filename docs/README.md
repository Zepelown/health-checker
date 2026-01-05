# Health Checker 문서

이 디렉토리에는 health-checker 프로젝트의 확장 기능 및 개발 계획 문서가 포함되어 있습니다.

## 문서 목록

### 확장 기능 문서

1. **[다중 엔드포인트 모니터링](./multi-endpoint-monitoring.md)**
   - 여러 URL을 동시에 모니터링하는 기능
   - 설정 파일 기반 엔드포인트 관리
   - 병렬 모니터링 구현

2. **[상세 메트릭 및 통계](./metrics-and-statistics.md)**
   - 응답 시간, 가용성 등 상세 메트릭 수집
   - Prometheus 메트릭 노출
   - 통계 계산 및 리포트

3. **[대시보드/UI](./dashboard-ui.md)**
   - 웹 기반 실시간 대시보드
   - RESTful API 및 WebSocket
   - 프론트엔드 구현

4. **[알림 채널 확장](./advanced-notifications.md)**
   - 이메일, PagerDuty, 커스텀 웹훅 등 다양한 알림 채널
   - 알림 규칙 시스템
   - 알림 템플릿 및 집계

5. **[고급 헬스 체크 기능](./advanced-health-checks.md)**
   - HTTP 메서드 커스터마이징
   - 응답 검증 (JSON 경로, 정규식 등)
   - SSL 인증서 체크
   - TCP/UDP 포트 체크
   - 데이터베이스 헬스 체크

6. **[설정 파일 기반 관리](./config-file-management.md)**
   - YAML/JSON 설정 파일 지원
   - 환경 변수 치환
   - 설정 검증 및 동적 리로드

### 계획 문서

- **[개발 계획서](./first-plan.md)**: 초기 프로젝트 계획
- **[로드맵](./ROADMAP.md)**: 확장 기능 개발 로드맵 및 우선순위
- **[향후 기능 요약](./future-features.md)**: 확장 기능 간단 요약

## 빠른 시작

### 현재 기능
현재 health-checker는 다음 기능을 지원합니다:
- 단일 URL 헬스 체크
- Slack/Discord 알림
- 응답 시간 모니터링

자세한 내용은 [README.ko.md](./README.ko.md)를 참조하세요.

### 확장 기능 계획
확장 기능에 대한 상세한 계획은 각 문서를 참조하세요. 개발 우선순위는 [ROADMAP.md](./ROADMAP.md)를 확인하세요.

## 문서 구조

```
docs/
├── README.md                          # 이 파일
├── README.ko.md                       # 한국어 사용자 가이드
├── first-plan.md                      # 초기 개발 계획
├── ROADMAP.md                         # 확장 기능 로드맵
├── future-features.md                 # 향후 기능 요약
├── multi-endpoint-monitoring.md        # 다중 엔드포인트 모니터링
├── metrics-and-statistics.md          # 메트릭 및 통계
├── dashboard-ui.md                    # 대시보드/UI
├── advanced-notifications.md          # 고급 알림 기능
├── advanced-health-checks.md          # 고급 헬스 체크
└── config-file-management.md          # 설정 파일 관리
```

## 기여하기

새로운 기능이나 개선 사항을 제안하려면:
1. GitHub Issues에서 관련 이슈를 확인하거나 생성
2. 해당 기능에 대한 문서를 작성 (선택 사항)
3. Pull Request 제출

## 문의

문서에 대한 질문이나 제안이 있으시면 GitHub Issues를 통해 문의해주세요.

