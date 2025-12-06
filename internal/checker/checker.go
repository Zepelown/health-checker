package checker

import (
	"net/http"
	"time"
)

// CheckURL performs an HTTP GET request to the given URL and returns
// the status code, latency, and any error that occurred.
func CheckURL(url string, timeout time.Duration) (status int, latency time.Duration, err error) {
	// HTTP 클라이언트 생성 (타임아웃 설정)
	client := &http.Client{
		Timeout: timeout,
	}

	// 시작 시간 기록
	start := time.Now()

	// HTTP GET 요청
	resp, err := client.Get(url)
	if err != nil {
		// 에러 발생 시 latency는 측정된 시간만큼 반환
		return 0, time.Since(start), err
	}
	defer resp.Body.Close()

	// 응답 시간 계산
	latency = time.Since(start)

	// 상태 코드 반환
	return resp.StatusCode, latency, nil
}
