# TraceHub CLI

Command-line interface for [TraceHub](https://github.com/tapiaw38/tracehub-server) - Monitor and analyze traces from your terminal.

## Features

- 🔗 **Connect to Projects** - Easy project selection and switching
- 📊 **Stream Traces** - Real-time trace viewing with filtering
- ❌ **Error Monitoring** - Quickly view errors and issues
- 🎨 **Colored Output** - Beautiful terminal UI with syntax highlighting
- 💾 **Persistent Config** - Saves server and project settings locally
- 🔑 **Secure Authentication** - API key-based authentication

## Installation

### From Source

```bash
git clone https://github.com/tapiaw38/tracehub-cli.git
cd tracehub-cli
go build -o tracehub ./cmd/tracehub
sudo mv tracehub /usr/local/bin/
```

### From Binary (Coming Soon)

```bash
# Download from releases
curl -L https://github.com/tapiaw38/tracehub-cli/releases/latest/download/tracehub-$(uname -s)-$(uname -m) -o tracehub
chmod +x tracehub
sudo mv tracehub /usr/local/bin/
```

## Quick Start

### 1. Configure Connection

First, configure the CLI to connect to your TraceHub server:

```bash
tracehub configure \
  --server http://localhost:8080 \
  --api-key th_your_api_key_here
```

This saves the configuration to `~/.tracehub/config.json`.

### 2. List Projects

See all available projects:

```bash
tracehub projects
```

Output:
```
Found 2 projects:

  • backend-api
    ID: a1b2c3d4-...
    Language: javascript
    Description: Production API server

  • payment-service
    ID: e5f6g7h8-...
    Language: go
```

### 3. Connect to a Project

```bash
tracehub connect a1b2c3d4-...
```

Output:
```
✓ Connected to project: backend-api

Available commands:
  tracehub stream    - Stream traces in real-time
  tracehub errors    - View detected errors
```

### 4. Stream Traces

View traces in real-time:

```bash
tracehub stream
```

With filters:

```bash
# Only errors
tracehub stream --level error

# Limit results
tracehub stream --limit 50
```

Output:
```
Streaming traces from project: backend-api

[15:34:22] INFO    Application started
[15:34:23] WARN    High memory usage detected
          Source: server.js:145
[15:34:24] ERROR   Database connection failed
          Source: db.js:23
          Stack: Error: Connection timeout...
```

### 5. View Errors

See recent errors:

```bash
tracehub errors
```

Output:
```
Found 3 recent errors:

  ❌ Database connection failed
     Time: 2026-01-10 15:34:24
     Source: db.js:23

  ❌ Nil pointer dereference
     Time: 2026-01-10 15:35:10
     Source: handler.go:89
```

## Commands

### `configure`

Set up server connection:

```bash
tracehub configure --server <url> --api-key <key>
```

**Flags:**
- `--server` - TraceHub server URL (required)
- `--api-key` - API key for authentication (required)

### `projects`

List all available projects:

```bash
tracehub projects
```

### `connect`

Connect to a specific project:

```bash
tracehub connect <project-id>
```

Sets the current project for subsequent commands.

### `stream`

Stream traces in real-time:

```bash
tracehub stream [flags]
```

**Flags:**
- `--level <level>` - Filter by level (info, warn, error, fatal)
- `--limit <n>` - Number of traces to display (default: 100)

**Examples:**

```bash
# Stream all traces
tracehub stream

# Only errors and fatals
tracehub stream --level error

# Last 50 traces
tracehub stream --limit 50
```

### `errors`

View recent errors:

```bash
tracehub errors
```

Shows the last 20 error-level traces.

### `version`

Show CLI version:

```bash
tracehub version
```

## Configuration

Configuration is stored in `~/.tracehub/config.json`:

```json
{
  "server_url": "http://localhost:8080",
  "api_key": "th_your_api_key",
  "current_project": {
    "id": "a1b2c3d4-...",
    "name": "backend-api"
  }
}
```

### Manual Configuration

You can also edit the configuration file directly:

```bash
vim ~/.tracehub/config.json
```

## Environment Variables

Override configuration with environment variables:

```bash
export TRACEHUB_SERVER=http://localhost:8080
export TRACEHUB_API_KEY=th_your_api_key

tracehub stream
```

## Output Formatting

The CLI uses colored output for better readability:

- 🔴 **Red** - Errors and fatal logs
- 🟡 **Yellow** - Warnings
- 🟢 **Green** - Info logs
- 🔵 **Cyan** - Debug logs

## Troubleshooting

### Connection Refused

```
Error: Failed to connect to server
```

**Solutions:**
1. Verify server is running: `curl http://localhost:8080/api/v1/health`
2. Check server URL in config
3. Ensure firewall allows connection

### Authentication Error

```
Error: server returned status 401
```

**Solutions:**
1. Verify API key is correct
2. Check API key hasn't expired
3. Reconfigure: `tracehub configure --server ... --api-key ...`

### No Project Connected

```
Error: No project connected
```

**Solution:**
```bash
tracehub projects           # List available projects
tracehub connect <id>       # Connect to a project
```

## Development

### Build from Source

```bash
go mod download
go build -o tracehub ./cmd/tracehub
```

### Run Tests

```bash
go test ./...
```

## Architecture

```
tracehub-cli/
├── cmd/tracehub/          # Main CLI entry point
├── internal/
│   ├── client/            # HTTP client for TraceHub API
│   ├── config/            # Configuration management
│   └── tui/               # Terminal UI components
└── pkg/models/            # Data models
```

## Roadmap

- [ ] Interactive TUI mode with real-time updates
- [ ] Export traces to file (JSON, CSV)
- [ ] Search and filtering with advanced queries
- [ ] AI-powered error analysis
- [ ] Automatic PR creation from CLI
- [ ] Multi-project monitoring
- [ ] Alerts and notifications
- [ ] Performance metrics visualization

## Integration

Works seamlessly with:
- [TraceHub Server](https://github.com/tapiaw38/tracehub-server)
- [TraceHub JavaScript SDK](https://github.com/tapiaw38/tracehub-javascript-sdk)

## License

MIT License

## Links

- [Documentation](https://github.com/tapiaw38/tracehub-server#readme)
- [Issues](https://github.com/tapiaw38/tracehub-cli/issues)
- [Releases](https://github.com/tapiaw38/tracehub-cli/releases)

---

Made with ❤️ for developers who love the terminal
