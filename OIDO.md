# Oido Google Calendar Extension

List, create, search, and read Google Calendar events via Service Account. This extension provides MCP tools for calendar management.

## Available Tools

### `g_calendar_list_events`
List upcoming events from the calendar.

**Parameters:**
- `max_results` (number, optional): Maximum number of events to return (default: 20)
- `time_min` (string, optional): Start of time range (RFC3339 format, default: now)

**Returns:** Formatted table with ID, Summary, Start, End.

### `g_calendar_get_event`
Get full details of a specific event by its ID.

**Parameters:**
- `event_id` (string, required): ID of the event to retrieve

**Returns:** Summary, Description, Start, End, Location, Status, Creator, Attendees.

### `g_calendar_create_event`
Create a new event on the calendar.

**Parameters:**
- `summary` (string, required): Title of the event
- `start_time` (string, required): Start time in RFC3339 format (e.g. 2026-05-15T09:00:00Z)
- `end_time` (string, required): End time in RFC3339 format (e.g. 2026-05-15T10:00:00Z)
- `description` (string, optional): Description or notes for the event
- `location` (string, optional): Location of the event

**Returns:** Confirmation with event ID, summary, and times.

### `g_calendar_search_events`
Search events by query string matching title, description, or location.

**Parameters:**
- `query` (string, required): Search term to match in event title, description, or location
- `max_results` (number, optional): Maximum number of results to return (default: 20)

**Returns:** Formatted table with matching events (ID, Summary, Start, End).

## Example Usage

```
User: What's on my calendar this week?

Assistant: Let me check your upcoming events.

[Uses g_calendar_list_events tool]

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

[Uses g_calendar_create_event tool]

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

- **Service Account authentication**: Uses Google Calendar API v3 with a service account JWT key (no browser OAuth flow needed)
- **Calendar ID**: Defaults to "primary" (service account's own calendar); use a shared calendar email or ID to access other calendars
- **Must share your calendar**: For the service account to access **your** Google Calendar, share your calendar with the service account email (found in the JSON key under `client_email`)
- **Time format**: All times are in RFC3339 format
- **Event IDs**: Use the event ID from g_calendar_list_events/g_calendar_search_events when calling g_calendar_get_event
- **Environment variable for authentication** (set one of):
  - `CALENDAR_SERVICE_ACCOUNT_JSON`: Paste the entire service account JSON key (raw)
  - `CALENDAR_SERVICE_ACCOUNT_B64`: Base64-encoded version of the JSON (avoids newline issues in env vars)
- **Limits**: Default 20 events for list/search
- **Creating events**: Requires a Google Cloud project with Calendar API enabled and a service account created
