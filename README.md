# Oido Google Calendar MCP Extension

List, create, search, and read Google Calendar events via OAuth2 using the Model Context Protocol.

## Features

- **List Events**: View upcoming calendar events
- **Read Events**: Fetch full event details by ID
- **Create Events**: Add new events to your calendar
- **Search Events**: Find events by title, description, or location

## Installation

### Option 1: Upload via Plugins UI (Recommended)

1. Download the latest release zip for your platform from [GitHub Releases](../../releases)
   - Linux: `oido-google-calendar-linux-amd64.zip`
   - macOS (Apple Silicon): `oido-google-calendar-darwin-arm64.zip`
2. Open Qwen CLI → Plugins UI
3. Upload the zip file
4. Configure settings (client ID, secret, refresh token) in the plugin settings panel

### Option 2: Build from Source

```bash
git clone <repo-url>
cd oido-google-calendar
make build
```

Then point your plugin configuration to the built `oido-google-calendar-mcp` binary.

### Option 3: Manual Install from Release Artifacts

```bash
curl -LO https://github.com/<owner>/<repo>/releases/latest/download/oido-google-calendar-linux-amd64.zip
unzip oido-google-calendar-linux-amd64.zip -d oido-google-calendar
./oido-google-calendar/oido-google-calendar-mcp
```

## Requirements

- Go 1.26+
- Google Cloud project with Calendar API enabled
- OAuth2 credentials (client ID + secret) and refresh token

## Setup

### 1. Enable Google Calendar API

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing
3. Enable the **Google Calendar API**
4. Go to **Credentials** → **Create Credentials** → **OAuth client ID**
5. Choose **Desktop app** as application type
6. Download the JSON credentials file

### 2. Get a Refresh Token

Use the downloaded credentials to obtain a refresh token (requires a one-time OAuth2 flow):

```bash
# Install Google OAuth2 CLI tool or use a script
# The refresh token grants long-lived access to the Calendar API
```

### 3. Configure Extension

Set the following environment variables (or configure via plugin settings):

| Variable | Description | Default |
|----------|-------------|---------|
| `CALENDAR_CLIENT_ID` | OAuth2 client ID | *(required)* |
| `CALENDAR_CLIENT_SECRET` | OAuth2 client secret | *(required)* |
| `CALENDAR_REFRESH_TOKEN` | OAuth2 refresh token | *(required)* |
| `CALENDAR_ID` | Calendar ID to use | `primary` |
| `CALENDAR_MAX_RESULTS` | Max events per list/search | `20` |
| `CALENDAR_TIMEZONE` | Timezone for events | *(system timezone)* |

## Build

```bash
make build
```

## Package for Distribution

```bash
make dist
```

This creates `dist/oido-google-calendar.zip` for upload via the Plugins UI.

## Tools

### `list_events`
List upcoming events from the calendar. Optionally filter by `max_results` and `time_min`.

### `get_event`
Get full event details by event ID.

### `create_event`
Create a new event. Requires `summary`, `start_time`, and `end_time` (RFC3339 format).

### `search_events`
Search events by query string matching title, description, or location.

## Architecture

```
┌─────────────┐     stdio      ┌─────────────────────────────┐
│  Qwen CLI   │ ◄────────────► │  oido-google-calendar-mcp    │
│             │                │                             │
│             │                │  ┌───────────────────────┐  │
│             │                │  │ Calendar API Client    │──► Google Calendar API
│             │                │  │ (OAuth2)              │  │
│             │                │  └───────────────────────┘  │
└─────────────┘                └─────────────────────────────┘
```

## License

MIT
