# Oido Google Calendar MCP Extension

List, create, search, and read Google Calendar events via Service Account using the Model Context Protocol.

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
4. Configure settings (paste service account JSON key) in the plugin settings panel

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
- Service account with JSON key

## Setup

### 1. Create a Service Account

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing
3. Enable the **Google Calendar API**
4. Go to **IAM & Admin** → **Service Accounts** → **Create Service Account**
5. Name it "Oido Calendar" → **Create and Continue**
6. Skip granting permissions (not needed) → **Done**
7. Click the service account → **Keys** → **Add Key** → **Create New Key**
8. Choose **JSON** → **Create** → the JSON key file downloads automatically

### 2. Share Your Calendar with the Service Account

1. Copy `client_email` from the downloaded JSON key (looks like `oido-calendar@your-project.iam.gserviceaccount.com`)
2. Open [Google Calendar](https://calendar.google.com/)
3. Find the calendar you want to access → click **⋮** → **Settings and sharing**
4. Under **Share with specific people**, add the service account email
5. Set permission to **Make changes to events** (for full access)

### 3. Configure Extension

Open the JSON key file, copy the **entire content**, and paste it as `CALENDAR_SERVICE_ACCOUNT_JSON`.  

**If newlines cause issues**, use base64 instead:
```bash
base64 -w0 /path/to/service-account.json
# Copy the output → paste as CALENDAR_SERVICE_ACCOUNT_B64
```

| Variable | Description | Default |
|----------|-------------|---------|
| `CALENDAR_SERVICE_ACCOUNT_JSON` | Service account JSON key (raw) | *(optional)* |
| `CALENDAR_SERVICE_ACCOUNT_B64` | Base64-encoded JSON key (alternative) | *(optional)* |
| `CALENDAR_ID` | Calendar ID (use your email for your main calendar) | `primary` |
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

### `g_calendar_list_events`
List upcoming events from the calendar. Optionally filter by `max_results` and `time_min`.

### `g_calendar_get_event`
Get full event details by event ID.

### `g_calendar_create_event`
Create a new event. Requires `summary`, `start_time`, and `end_time` (RFC3339 format).

### `g_calendar_search_events`
Search events by query string matching title, description, or location.

## Architecture

```
┌─────────────┐     stdio      ┌─────────────────────────────┐
│  Qwen CLI   │ ◄────────────► │  oido-google-calendar-mcp    │
│             │                │                             │
│             │                │  ┌───────────────────────┐  │
│             │                │  │ Calendar API Client    │──► Google Calendar API
│             │                │  │ (Service Account)              │  │
│             │                │  └───────────────────────┘  │
└─────────────┘                └─────────────────────────────┘
```

## License

MIT
