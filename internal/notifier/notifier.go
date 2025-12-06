package notifier

import (
	"fmt"
	"log"
)

// Notifier defines the interface for notification services
type Notifier interface {
	Send(webhookURL string, message string) error
}

// SlackNotifier implements Notifier for Slack
type SlackNotifier struct{}

// Send sends a message to Slack
func (n *SlackNotifier) Send(webhookURL string, message string) error {
	return SendSlack(webhookURL, message)
}

// DiscordNotifier implements Notifier for Discord
type DiscordNotifier struct{}

// Send sends a message to Discord
func (d *DiscordNotifier) Send(webhookURL string, message string) error {
	return SendDiscord(webhookURL, message)
}

// NotificationConfig holds webhook URLs for different notification services
type NotificationConfig struct {
	SlackWebhook   string
	DiscordWebhook string
}

// SendToAll sends notifications to all configured services
// Each notification failure is logged but doesn't stop other notifications
func SendToAll(config NotificationConfig, message string) {
	if config.SlackWebhook != "" {
		slackNotifier := &SlackNotifier{}
		if err := slackNotifier.Send(config.SlackWebhook, message); err != nil {
			log.Printf("Failed to send Slack notification: %v\n", err)
		}
	}

	if config.DiscordWebhook != "" {
		discordNotifier := &DiscordNotifier{}
		if err := discordNotifier.Send(config.DiscordWebhook, message); err != nil {
			log.Printf("Failed to send Discord notification: %v\n", err)
		}
	}
}

// HasAnyNotification checks if at least one notification service is configured
func HasAnyNotification(config NotificationConfig) bool {
	return config.SlackWebhook != "" || config.DiscordWebhook != ""
}

// GetNotificationStatus returns a string describing which notifications are enabled
func GetNotificationStatus(config NotificationConfig) string {
	var services []string
	if config.SlackWebhook != "" {
		services = append(services, "Slack")
	}
	if config.DiscordWebhook != "" {
		services = append(services, "Discord")
	}
	if len(services) == 0 {
		return ""
	}

	// 서비스 목록을 읽기 쉬운 형식으로 변환
	serviceList := ""
	for i, service := range services {
		if i > 0 {
			if i == len(services)-1 {
				serviceList += " and "
			} else {
				serviceList += ", "
			}
		}
		serviceList += service
	}

	return fmt.Sprintf("%s notifications enabled", serviceList)
}
