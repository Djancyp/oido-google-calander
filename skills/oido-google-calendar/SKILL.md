---
name: oido-google-calendar
description: List, create, search, and read Google Calendar events via Service Account
---

# Oido Google Calendar Extension

## Overview

The Oido Google Calendar extension provides tools to list, create, search, and read events on Google Calendar via the Calendar API v3 with a service account. Use these tools when users ask about their schedule, want to create meetings, or need to find specific events.

## Available Tools

### `g_calendar_list_events`

List upcoming events from the calendar.

**Parameters:**
- `max_results` (number, optional): Maximum number of events to return (default: 20)
- `time_min` (string, optional): Start of time range (RFC3339 format, default: now)

**When to use:**
- User asks "what's on my calendar?"
- User wants to see upcoming events
- User wants to check their schedule

**Example Usage:**
```
User: "Show me my next 10 events"
→ Call g_calendar_list_events with max_results: 10
```

**Response Format:**
Returns formatted table with ID, Summary, Start, End.

### `g_calendar_get_event`

Get full details of a specific event by its ID.

**Parameters:**
- `event_id` (string, required): ID of the event to retrieve

**When to use:**
- User asks for details of a specific event
- User wants to see full event information
- User references an event ID from a previous g_calendar_list_events result

**Example Usage:**
```
User: "Show me details for event abc123"
→ Call g_calendar_get_event with event_id: "abc123"
```

**Response Format:**
Returns Summary, Description, Start, End, Location, Status, Creator, Attendees.

### `g_calendar_create_event`

Create a new event on the calendar.

**Parameters:**
- `summary` (string, required): Title of the event
- `start_time` (string, required): Start time in RFC3339 format (e.g. 2026-05-15T09:00:00Z)
- `end_time` (string, required): End time in RFC3339 format (e.g. 2026-05-15T10:00:00Z)
- `description` (string, optional): Description or notes for the event
- `location` (string, optional): Location of the event

**When to use:**
- User asks to create an event or meeting
- User wants to add something to their calendar
- User wants to schedule a new appointment

**Example Usage:**
```
User: "Create a meeting tomorrow at 2pm for 1 hour called 'Sprint Review'"
→ Call g_calendar_create_event with summary: "Sprint Review", start_time: "2026-05-03T14:00:00Z", end_time: "2026-05-03T15:00:00Z"
```

**Response Format:**
Returns confirmation with event ID, summary, start, and end times.

### `g_calendar_search_events`

Search events by query string matching title, description, or location.

**Parameters:**
- `query` (string, required): Search term to match in event title, description, or location
- `max_results` (number, optional): Maximum number of results to return (default: 20)

**When to use:**
- User asks to search for events
- User wants to find events about a specific topic
- User wants to filter calendar by keyword

**Example Usage:**
```
User: "Find events about 'meeting'"
→ Call g_calendar_search_events with query: "meeting"
```

**Response Format:**
Returns formatted table with ID, Summary, Start, End for matching events.

## Best Practices

1. **List before reading**: Use g_calendar_list_events first to show available events, then g_calendar_get_event for details
2. **ID-based access**: Always use the event ID from g_calendar_list_events/g_calendar_search_events when calling g_calendar_get_event
3. **RFC3339 times**: Always use RFC3339 format for time parameters
4. **Default limits**: Respect max_results limits (default 20 for list/search)
5. **Context**: Mention results are from the user's primary calendar unless otherwise specified

## Example Interactions

### Listing Events

```
User: "What's on my calendar?"

Assistant: Let me check your upcoming events.
[Calls g_calendar_list_events with max_results: 20]

Upcoming Events (3):

ID                                     | Summary                      | Start                | End
---------------------------------------+------------------------------+----------------------+----------------------
abc123def456                           | Q2 Planning Meeting          | 2026-05-15T09:00:00Z | 2026-05-15T10:00:00Z
def789abc012                           | Lunch with Alice             | 2026-05-15T12:00:00Z | 2026-05-15T13:00:00Z
...

Would you like details on any of these?
```

### Reading an Event

```
User: "Show me details for event abc123def456"

Assistant: Let me fetch that event's details.
[Calls g_calendar_get_event with event_id: "abc123def456"]

Summary:     Q2 Planning Meeting
ID:          abc123def456
Status:      confirmed
Start:       2026-05-15T09:00:00Z
End:         2026-05-15T10:00:00Z
Location:    Conference Room A
Creator:     user@gmail.com
Attendees:   [alice@example.com bob@example.com]
-------
Quarterly planning session. Please review the agenda beforehand.
```

### Creating an Event

```
User: "Schedule a lunch meeting with Alice for tomorrow at 12pm lasting 1 hour"

Assistant: I'll create that event for you.
[Calls g_calendar_create_event with summary: "Lunch with Alice", start_time: "...T12:00:00Z", end_time: "...T13:00:00Z"]

Event created successfully:

Summary:     Lunch with Alice
ID:          xyz789abc123
Start:       2026-05-03T12:00:00Z
End:         2026-05-03T13:00:00Z
```

### Searching Events

```
User: "Find events about 'sprint'"

Assistant: Let me search your calendar for "sprint".
[Calls g_calendar_search_events with query: "sprint"]

Search Results for "sprint" (2):

ID                                     | Summary                      | Start                | End
---------------------------------------+------------------------------+----------------------+----------------------
abc789def012                           | Sprint Planning              | 2026-05-10T10:00:00Z | 2026-05-10T11:00:00Z
ghi345jkl678                           | Sprint Retrospective         | 2026-05-14T15:00:00Z | 2026-05-14T16:00:00Z
```

## Limitations

- **Service Account only**: Requires a Google Cloud service account JSON key (no browser OAuth flow)
- **Must share calendar**: Add the service account email to your calendar's share settings
- **Single calendar**: Reads from one calendar at a time (defaults to primary)
- **No recurrence expansion**: Recurring events may appear as single instances
- **No attendee management**: Creates events without adding attendees
- **No reminders**: Does not set event reminders
- **No color/label customization**: Creates events with default settings
- Requires environment variable: CALENDAR_SERVICE_ACCOUNT_JSON (paste the entire JSON key content)
- **Google Cloud setup required**: Requires enabling Calendar API and creating a service account
- **RFC3339 time format**: All times must be in RFC3339 format

## Related Commands

- `/list-events` - List upcoming events (custom command)
- `/get-event` - Get full event details (custom command)
- `/create-event` - Create a new event (custom command)
- `/search-events` - Search events by query (custom command)

## Triggers

Use these tools when you see:
- "calendar" or "schedule" or "event"
- "meeting" or "appointment"
- "what's on my calendar" or "upcoming events"
- "create event" or "add to calendar"
- "find event" or "search calendar"
- Google Calendar, schedule, planning
