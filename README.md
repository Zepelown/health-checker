# Health Checker

[한국어](docs/README.ko.md) | English

A Go-based health checker that periodically monitors website availability and sends notifications via Slack or Discord when issues are detected.

## Features

- Periodic website health checks
- HTTP status code and response time monitoring
- Slack Webhook notification support
- Discord Webhook notification support
- Support for using both Slack and Discord simultaneously
- Test mode support (sends notifications even when status is healthy)

## Installation

### With Go installed

```bash
go install github.com/your-username/health-checker@latest
```

Or build from source:

```bash
git clone https://github.com/your-username/health-checker.git
cd health-checker
go build
```

## Usage

### Basic Usage

```bash
# Windows PowerShell
.\health-checker.exe run --url https://example.com

# Linux/Mac
./health-checker run --url https://example.com
```

### Command Options

#### Required Options

- `--url`, `-u`: URL to check (required)
  ```bash
  --url https://example.com
  ```

#### Optional Options

- `--interval`, `-i`: Check interval (default: `60s`)
  ```bash
  --interval 30s    # Check every 30 seconds
  --interval 5m     # Check every 5 minutes
  ```

- `--timeout`, `-t`: Request timeout (default: `5s`)
  ```bash
  --timeout 10s     # 10 second timeout
  ```

- `--slack-webhook`, `-s`: Slack Webhook URL
  ```bash
  --slack-webhook https://hooks.slack.com/services/YOUR/WEBHOOK/URL
  ```

- `--discord-webhook`, `-d`: Discord Webhook URL
  ```bash
  --discord-webhook https://discord.com/api/webhooks/YOUR/WEBHOOK/URL
  ```

- `--latency-threshold`: Response latency threshold (if specified, responses exceeding this are considered failures)
  ```bash
  --latency-threshold 3s     # Alert if response takes 3 seconds or more
  --latency-threshold 500ms  # Alert if response takes 500ms or more
  ```

- `--test`: Test mode (sends notifications even when status is healthy)
  ```bash
  --test
  ```

All duration-related options (`--interval`, `--timeout`, `--latency-threshold` and related environment variables) follow Go's `time.ParseDuration` format and support the following units:

- `ns` (nanoseconds), `us`/`µs` (microseconds), `ms` (milliseconds)
- `s` (seconds), `m` (minutes), `h` (hours)

Examples: `500ms`, `2s`, `1.5s`, `3m`, `1h30m`

### Environment Variables

You can use environment variables instead of flags:

- `SLACK_WEBHOOK_URL`: Slack Webhook URL
- `DISCORD_WEBHOOK_URL`: Discord Webhook URL
- `LATENCY_THRESHOLD`: Response latency threshold (e.g., `3s`, `500ms`)

```bash
# Windows PowerShell
$env:SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
$env:DISCORD_WEBHOOK_URL="https://discord.com/api/webhooks/YOUR/WEBHOOK/URL"
$env:LATENCY_THRESHOLD="3s"
.\health-checker.exe run --url https://example.com

# Linux/Mac
export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
export DISCORD_WEBHOOK_URL="https://discord.com/api/webhooks/YOUR/WEBHOOK/URL"
export LATENCY_THRESHOLD="3s"
./health-checker run --url https://example.com
```

## Examples

### 1. Basic Health Check (No Notifications)

```bash
.\health-checker.exe run --url https://example.com --interval 60s
```

### 2. Slack Notifications Only

```bash
.\health-checker.exe run \
  --url https://example.com \
  --interval 30s \
  --slack-webhook https://hooks.slack.com/services/YOUR/WEBHOOK/URL
```

### 3. Discord Notifications Only

```bash
.\health-checker.exe run \
  --url https://example.com \
  --interval 30s \
  --discord-webhook https://discord.com/api/webhooks/YOUR/WEBHOOK/URL
```

### 4. Using Both Slack and Discord

```bash
.\health-checker.exe run \
  --url https://example.com \
  --interval 60s \
  --slack-webhook https://hooks.slack.com/services/YOUR/WEBHOOK/URL \
  --discord-webhook https://discord.com/api/webhooks/YOUR/WEBHOOK/URL
```

### 5. Test Mode (Notifications Even When Healthy)

```bash
.\health-checker.exe run \
  --url https://example.com \
  --discord-webhook https://discord.com/api/webhooks/YOUR/WEBHOOK/URL \
  --test
```

### 6. Fast Monitoring with Short Intervals

```bash
.\health-checker.exe run \
  --url https://example.com \
  --interval 10s \
  --timeout 3s \
  --latency-threshold 2s \
  --discord-webhook https://discord.com/api/webhooks/YOUR/WEBHOOK/URL
```

## Notification Conditions

### Normal Mode (Default)

Notifications are sent in the following cases:

- HTTP request failure (connection errors, timeouts, etc.)
- HTTP status code is not 200 (500, 404, 503, etc.)
- (Optional) When `--latency-threshold` or `LATENCY_THRESHOLD` is set, if response time exceeds the threshold

### Test Mode (`--test` flag)

Notifications are sent in all states:

- On error: Failure notification
- Non-200 status code: Failure notification
- 200 OK: Healthy notification (only in test mode)

## Notification Message Format

### Failure Notification

```
🚨 Site Failure Detected: https://example.com
Status Code: 500
Response Time: 2.5s
```

Or

```
🚨 Site Failure Detected: https://example.com
Error: connection timeout
Response Time: 5s
```

### Healthy Notification (Test Mode)

```
✅ Site Healthy: https://example.com
Status Code: 200
Response Time: 150ms
```

## Webhook Setup

### Slack Webhook Setup

1. Create a new app at [Slack API](https://api.slack.com/apps)
2. Enable Incoming Webhooks
3. Copy the Webhook URL
4. Set it in the `--slack-webhook` flag or `SLACK_WEBHOOK_URL` environment variable

### Discord Webhook Setup

1. Discord channel settings → Integrations → Webhooks
2. Create a new webhook
3. Copy the webhook URL
4. Set it in the `--discord-webhook` flag or `DISCORD_WEBHOOK_URL` environment variable

## Stopping the Program

Press `Ctrl+C` to stop the program.

```
Press Ctrl+C to stop
^C
Shutting down...
```

## Project Structure

```
health-checker/
├── cmd/
│   ├── root.go      # Root command
│   └── run.go       # Run subcommand
├── internal/
│   ├── checker/     # Health check logic
│   │   └── checker.go
│   └── notifier/    # Notification logic
│       ├── slack.go
│       ├── discord.go
│       └── notifier.go
├── main.go
├── go.mod
└── README.md
```

## Tech Stack

- **Language**: Go
- **CLI Framework**: Cobra
- **HTTP Client**: net/http (standard library)

## License

MIT License

## Contributing

Issues and Pull Requests are welcome!
