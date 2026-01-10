package config

type Config struct {
	Version       string              `yaml:"version"`
	Global        GlobalConfig        `yaml:"global"`
	Endpoints     []EndpointConfig    `yaml:"endpoints"`
	Notifications NotificationsConfig `yaml:"notifications"`
	Metrics       MetricsConfig       `yaml:"metrics"`
	Dashboard     DashboardConfig     `yaml:"dashboard"`
	Logging       LoggingConfig       `yaml:"logging"`
}

type GlobalConfig struct {
	Version          string `yaml:"version"`
	Timeout          string `yaml:"timeout"`
	LatencyThreshold string `yaml:"latency_threshold"`
	SlackWebhook     string `yaml:"slack_webhook"`
	DiscordWebhook   string `yaml:"discord_webhook"`
}

type EndpointConfig struct {
	Name               string                 `yaml:"name"`
	URL                string                 `yaml:"url"`
	Type               string                 `yaml:"type"` //http, tcp, postgres, etc.
	Interval           string                 `yaml:"interval"`
	Timeout            string                 `yaml:"timeout"`
	LatencyThreshold   string                 `yaml:"latency_threshold"`
	Enabled            bool                   `yaml:"enabled"`
	Method             string                 `yaml:"method"` //GET, POST, PUT, DELETE, etc.
	Headers            string                 `yaml:"headers"`
	Body               string                 `yaml:"body"`
	ExpectedStatus     int                    `yaml:"expected_status"`
	ResponseValidation *ResponseValidation    `yaml:"response_validation"`
	SSL                *SSLConfig             `yaml:"ssl"`
	Notifications      *EndpointNotifications `yaml:"notifications"`
}
type ResponseValidation struct {
	JSONPath []JSONPathRule `yaml:"json_path"`
}

type JSONPathRule struct {
	Path     string `yaml:"path"`
	Expected string `yaml:"expected"`
}

type SSLConfig struct {
	Verify          bool   `yaml:"verify"`
	ExpiryThreshold string `yaml:"expiry_threshold"`
}

type EndpointNotifications struct {
	SlackWebhook   string             `yaml:"slack_webhook"`
	DiscordWebhook string             `yaml:"discord_webhook"`
	Rules          []NotificationRule `yaml:"rules"`
}

type NotificationRule struct {
	Name      string   `yaml:"name"`
	Condition string   `yaml:"condition"`
	Channels  []string `yaml:"channels"`
}

type NotificationsConfig struct {
	Channels map[string]ChannelConfig `yaml:"channels"`
	Rules    []NotificationRule       `yaml:"rules"`
}

type ChannelConfig struct {
	Webhook        string      `yaml:"webhook,omitempty"`
	Enabled        bool        `yaml:enabled"`
	SMTP           *SMTPConfig `yaml:"smtp,omitempty"`
	To             []string    `yaml:"to,omitempty"`
	IntegrationKey string      `yaml:"integration_key,omitempty"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type MetricsConfig struct {
	Enabled    bool             `yaml:"enabled"`
	Storage    StorageConfig    `yaml:"storage"`
	Prometheus PrometheusConfig `yaml:"prometheus"`
}

type StorageConfig struct {
	Type      string `yaml:"type"`
	Path      string `yaml:"path"`
	Retention string `yaml:"retention"`
}

type PrometheusConfig struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port"`
	Path    string `yaml:"path"`
}

type DashboardConfig struct {
	Enabled bool       `yaml:"enabled"`
	Port    int        `yaml:"port"`
	Auth    AuthConfig `yaml:"auth"`
	Theme   string     `yaml:"theme"`
}

type AuthConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type LoggingConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	File       string `yaml:"file"`
	MaxSize    string `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     string `yaml:"max_age"`
}
