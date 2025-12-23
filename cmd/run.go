/*
Copyright © 2025 Zepelown
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"health-checker/internal/checker"
	"health-checker/internal/notifier"

	"github.com/spf13/cobra"
)

var (
	urlFlag              string
	intervalFlag         string
	timeoutFlag          string
	slackWebhookFlag     string
	discordWebhookFlag   string
	latencyThresholdFlag string
	testModeFlag         bool
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start health checking for a URL",
	Long: `Start periodic health checking for a URL and send notifications (Slack/Discord) on failure.

Example:
  health-checker run --url https://example.com --interval 60s --timeout 5s
  health-checker run --url https://example.com --slack-webhook <url> --discord-webhook <url>`,
	Run: func(cmd *cobra.Command, args []string) {
		// 플래그 값 검증
		if urlFlag == "" {
			fmt.Println("Error: --url flag is required")
			os.Exit(1)
		}

		// duration 파싱
		interval, err := time.ParseDuration(intervalFlag)
		if err != nil {
			fmt.Printf("Error: invalid interval format: %v\n", err)
			os.Exit(1)
		}

		timeout, err := time.ParseDuration(timeoutFlag)
		if err != nil {
			fmt.Printf("Error: invalid timeout format: %v\n", err)
			os.Exit(1)
		}

		// duration 파싱: latency threshold (선택 사항)
		var latencyThreshold time.Duration
		latencyThresholdValue := latencyThresholdFlag
		if latencyThresholdValue == "" {
			latencyThresholdValue = os.Getenv("LATENCY_THRESHOLD")
		}
		if latencyThresholdValue != "" {
			latencyThreshold, err = time.ParseDuration(latencyThresholdValue)
			if err != nil {
				fmt.Printf("Error: invalid latency-threshold format: %v\n", err)
				os.Exit(1)
			}
			if latencyThreshold <= 0 {
				fmt.Println("Error: latency-threshold must be greater than 0")
				os.Exit(1)
			}
		}

		// Notification 설정 구성 (환경변수 또는 플래그)
		notifConfig := notifier.NotificationConfig{
			SlackWebhook:   slackWebhookFlag,
			DiscordWebhook: discordWebhookFlag,
		}

		// 환경변수에서 webhook URL 가져오기 (플래그가 없을 경우)
		if notifConfig.SlackWebhook == "" {
			notifConfig.SlackWebhook = os.Getenv("SLACK_WEBHOOK_URL")
		}
		if notifConfig.DiscordWebhook == "" {
			notifConfig.DiscordWebhook = os.Getenv("DISCORD_WEBHOOK_URL")
		}

		// 시그널 핸들링 (Ctrl+C로 깔끔하게 종료)
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		// Ticker로 주기적 체크
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		fmt.Printf("Starting health check for %s (interval: %s, timeout: %s)\n", urlFlag, interval, timeout)
		if status := notifier.GetNotificationStatus(notifConfig); status != "" {
			fmt.Println(status)
		}
		if latencyThreshold > 0 {
			fmt.Printf("Latency threshold enabled: %s\n", latencyThreshold)
		}
		if testModeFlag {
			fmt.Println("Test mode: notifications will be sent for all status codes (including 200)")
		}
		fmt.Println("Press Ctrl+C to stop")

		// 첫 체크 즉시 실행
		performCheck(urlFlag, timeout, notifConfig, testModeFlag, latencyThreshold)

		// 주기적 체크 루프
		for {
			select {
			case <-ticker.C:
				performCheck(urlFlag, timeout, notifConfig, testModeFlag, latencyThreshold)
			case <-sigChan:
				fmt.Println("\nShutting down...")
				return
			}
		}
	},
}

func performCheck(url string, timeout time.Duration, config notifier.NotificationConfig, testMode bool, latencyThreshold time.Duration) {
	status, latency, err := checker.CheckURL(url, timeout)

	if err != nil {
		log.Printf("❌ [%s] Error: %v (latency: %v)\n", url, err, latency)

		// 알림 전송 (Slack, Discord 모두)
		if notifier.HasAnyNotification(config) {
			message := fmt.Sprintf("🚨 사이트 장애 감지: %s\n에러: %v\n응답 시간: %v", url, err, latency)
			notifier.SendToAll(config, message)
		}
		return
	}

	if status != 200 {
		log.Printf("⚠️  [%s] Status: %d (latency: %v)\n", url, status, latency)

		// 알림 전송 (Slack, Discord 모두)
		if notifier.HasAnyNotification(config) {
			message := fmt.Sprintf("🚨 사이트 장애 감지: %s\n상태 코드: %d\n응답 시간: %v", url, status, latency)
			notifier.SendToAll(config, message)
		}
		return
	}

	// Latency threshold check (for successful 200 responses)
	if latencyThreshold > 0 && latency > latencyThreshold {
		log.Printf("⏱️  [%s] Slow response: %v (threshold: %v, status: %d)\n", url, latency, latencyThreshold, status)

		if notifier.HasAnyNotification(config) {
			message := fmt.Sprintf("🚨 응답 지연 임계값 초과: %s\n응답 시간: %v\n임계값: %v\n상태 코드: %d", url, latency, latencyThreshold, status)
			notifier.SendToAll(config, message)
		}
		return
	}

	log.Printf("✅ [%s] Status: %d (latency: %v)\n", url, status, latency)

	// 테스트 모드일 때는 정상 상태(200)에서도 알림 전송
	if testMode && notifier.HasAnyNotification(config) {
		message := fmt.Sprintf("✅ 사이트 정상: %s\n상태 코드: %d\n응답 시간: %v", url, status, latency)
		notifier.SendToAll(config, message)
	}
}

func init() {
	rootCmd.AddCommand(runCmd)

	// 플래그 정의
	runCmd.Flags().StringVarP(&urlFlag, "url", "u", "", "URL to check (required)")
	runCmd.Flags().StringVarP(&intervalFlag, "interval", "i", "60s", "Check interval (e.g., 60s, 1m)")
	runCmd.Flags().StringVarP(&timeoutFlag, "timeout", "t", "5s", "Request timeout (e.g., 5s, 10s)")
	runCmd.Flags().StringVarP(&slackWebhookFlag, "slack-webhook", "s", "", "Slack webhook URL (or use SLACK_WEBHOOK_URL env var)")
	runCmd.Flags().StringVarP(&discordWebhookFlag, "discord-webhook", "d", "", "Discord webhook URL (or use DISCORD_WEBHOOK_URL env var)")
	runCmd.Flags().StringVar(&latencyThresholdFlag, "latency-threshold", "", "Latency threshold for considering slow responses as failures (e.g., 3s, 500ms). Can also be set via LATENCY_THRESHOLD env var")
	runCmd.Flags().BoolVar(&testModeFlag, "test", false, "Test mode: send notifications for all status codes (including 200)")

	// url 플래그를 필수로 설정
	runCmd.MarkFlagRequired("url")
}
