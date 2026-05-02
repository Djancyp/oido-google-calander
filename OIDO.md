# Oido Google Calendar Extension

List, create, search, and read Google Calendar events via OAuth2. This extension provides MCP tools for calendar management.

## Available Tools

### `list_events`
List upcoming events from the calendar.

**Parameters:**
- `max_results` (number, optional): Maximum number of events to return (default: 20)
- `time_min` (string, optional): Start of time range (RFC3339 format, default: now)

**Returns:** Formatted table with ID, Summary, Start, End.

### `get_event`
Get full details of a specific event by its ID.

**Parameters:**
- `event_id` (string, required): ID of the event to retrieve

**Returns:** Summary, Description, Start, End, Location, Status, Creator, Attendees.

### `create_event`
Create a new event on the calendar.

**Parameters:**
- `summary` (string, required): Title of the event
- `start_time` (string, required): Start time in RFC3339 format (e.g. 2026-05-15T09:00:00Z)
- `end_time` (string, required): End time in RFC3339 format (e.g. 2026-05-15T10:00:00Z)
- `description` (string, optional): Description or notes for the event
- `location` (string, optional): Location of the event

**Returns:** Confirmation with event ID, summary, and times.

### `search_events`
Search events by query string matching title, description, or location.

**Parameters:**
- `query` (string, required): Search term to match in event title, description, or location
- `max_results` (number, optional): Maximum number of results to return (default: 20)

**Returns:** Formatted table with matching events (ID, Summary, Start, End).

## Example Usage

```
User: What's on my calendar this week?

Assistant: Let me check your upcoming events.

[Uses list_events tool]

Upcoming Events (3):

ID                                     | Summary                      | Start                | End
---------------------------------------+------------------------------+----------------------+----------------------
abc123def456                           | Q2 Planning Meeting          | 2026-05-15T09:00:00Z | 2026-05-15T10:00:00Z
def789abc012                           | Lunch with Alice             | 2026-05-15T12:00:00Z | 2026-05-15T13:00:00Z
...
```

```
User: Create a meeting tomorrow at 2pm for 1 hour called 'Sprint Review'

Assistant: I'll create that event for you.

[Uses create_event tool]

Event created successfully:

Summary:  Sprint Review
ID:       xyz123abc456
Start:    2026-05-03T14:00:00Z
End:      2026-05-03T15:00:00Z
```

## When to Use

- User asks to see calendar events or schedule
- User wants to create a new event or meeting
- User asks for event details
- User wants to search for specific events

## Notes

- **OAuth2 authentication**: Uses Google Calendar API v3 with OAuth2 refresh token
- **Calendar ID**: Defaults to "primary" (user's main calendar); can use any accessible calendar ID
- **Time format**: All times are in RFC3339 format
- **Event IDs**: Use the event ID from list_events/search_events when calling get_event
- **Environment variables for authentication**:
  - `CALENDAR_CLIENT_ID` (required): OAuth2 client ID from Google Cloud Console
  - `CALENDAR_CLIENT_SECRET` (required): OAuth2 client secret from Google Cloud Console
  - `CALENDAR_REFRESH_TOKEN` (required): OAuth2 refresh token for long-lived access
- **Limits**: Default 20 events for list/search
- **Creating events**: Requires a Google Cloud project with Calendar API enabled and OAuth2 credentials configured
