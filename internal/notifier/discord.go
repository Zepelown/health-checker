package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// DiscordMessage represents the payload structure for Discord webhook
type DiscordMessage struct {
	Content string `json:"content"`
}

// SendDiscord sends a message to Discord using the webhook URL
func SendDiscord(webhookURL string, message string) error {
	// Discord 메시지 구조 생성
	discordMsg := DiscordMessage{
		Content: message,
	}

	// JSON으로 마샬링
	jsonData, err := json.Marshal(discordMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal Discord message: %w", err)
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

	// 응답 상태 코드 확인 (Discord는 204 No Content도 성공으로 처리)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Discord API returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}
