package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SlackMessage represents the payload structure for Slack webhook
type SlackMessage struct {
	Text string `json:"text"`
}

// SendSlack sends a message to Slack using the webhook URL
func SendSlack(webhookURL string, message string) error {
	// Slack 메시지 구조 생성
	slackMsg := SlackMessage{
		Text: message,
	}

	// JSON으로 마샬링
	jsonData, err := json.Marshal(slackMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack message: %w", err)
	}

	// HTTP POST 요청 생성
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Content-Type 헤더 설정
	req.Header.Set("Content-Type", "application/json")

	// HTTP 클라이언트로 요청 전송
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 응답 상태 코드 확인
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Slack API returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}
