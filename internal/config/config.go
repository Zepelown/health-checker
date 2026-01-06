package config

type Config struct {
	Version   string           `yaml:"version"`
	Global    GlobalConfig     `yaml:"global"`
	Endpoints []EndpointConfig `yaml:"endpoints"`
}

type GlobalConfig struct {
	Version          string `yaml:"version"`
	Timeout          string `yaml:"timeout"`
	LatencyThreshold string `yaml:"latency_threshold"`
	SlackWebhook     string `yaml:"slack_webhook"`
	DiscordWebhook   string `yaml:"discord_webhook"`
}

type EndpointConfig struct {
	Name             string `yaml:"name"`
	URL              string `yaml:"url"`
	Type             string `yaml:"type"` //http, tcp, postgres, etc.
	Interval         string `yaml:"interval"`
	Timeout          string `yaml:"timeout"`
	LatencyThreshold string `yaml:"latency_threshold"`
	Enabled          bool   `yaml:"enabled"`
	Method           string `yaml:"method"` //GET, POST, PUT, DELETE, etc.
	Headers          string `yaml:"headers"`
	Body             string `yaml:"body"`
	ExpectedStatus   int    `yaml:"expected_status"`
	// ResponseValidation *ResponseValidation `yaml:"response_validation"`
	// SSL *SSLConfig `yaml:"ssl"`
	// Notifications *EndpointNotifications `yaml:"notifications"`
}
