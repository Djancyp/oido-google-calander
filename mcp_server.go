package main

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func NewMCPHandler(cc *CalendarClient) *MCPHandler {
	return &MCPHandler{cc: cc}
}

type MCPHandler struct {
	cc *CalendarClient
}

type ListEventsArgs struct {
	MaxResults int64  `json:"max_results" jsonschema:"Maximum number of events to return (default: 20)"`
	TimeMin    string `json:"time_min,omitempty" jsonschema:"Start of the time range to fetch events (RFC3339 format, default: now)"`
}

type GetEventArgs struct {
	EventID string `json:"event_id" jsonschema:"ID of the event to retrieve"`
}

type CreateEventArgs struct {
	Summary     string `json:"summary" jsonschema:"Title of the event"`
	StartTime   string `json:"start_time" jsonschema:"Start time in RFC3339 format (e.g. 2026-05-15T09:00:00Z)"`
	EndTime     string `json:"end_time" jsonschema:"End time in RFC3339 format (e.g. 2026-05-15T10:00:00Z)"`
	Description string `json:"description,omitempty" jsonschema:"Description or notes for the event"`
	Location    string `json:"location,omitempty" jsonschema:"Location of the event"`
}

type SearchEventsArgs struct {
	Query      string `json:"query" jsonschema:"Search term to match in event title, description, or location"`
	MaxResults int64  `json:"max_results" jsonschema:"Maximum number of results to return (default: 20)"`
}

func RunMCPServer() {
	calClient, err := NewCalendarClient()
	if err != nil {
		log.Fatalf("Failed to create Calendar client: %v", err)
	}

	handler := NewMCPHandler(calClient)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "oido-google-calendar",
		Version: "1.0.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_events",
		Description: "List upcoming events from the calendar. Returns event ID, summary, start time, end time, and location.",
	}, handler.HandleListEvents)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_event",
		Description: "Get full details of a specific event by its ID. Returns summary, description, time, location, status, creator, and attendees.",
	}, handler.HandleGetEvent)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_event",
		Description: "Create a new event on the calendar. Requires summary, start_time, and end_time.",
	}, handler.HandleCreateEvent)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_events",
		Description: "Search events by query string matching title, description, or location. Returns matching events with ID, summary, and time.",
	}, handler.HandleSearchEvents)

	ctx := context.Background()
	log.Println("Oido Google Calendar MCP Server starting on stdio...")

	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatalf("MCP server error: %v", err)
	}
}

func formatTimeDisplay(t string) string {
	if len(t) >= 19 {
		return t[:19]
	}
	return t
}

func (h *MCPHandler) HandleListEvents(ctx context.Context, req *mcp.CallToolRequest, args ListEventsArgs) (*mcp.CallToolResult, any, error) {
	events, err := h.cc.ListEvents(ctx, args.MaxResults, args.TimeMin)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	if len(events) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "No upcoming events found."},
			},
		}, nil, nil
	}

	var result string
	result += fmt.Sprintf("Upcoming Events (%d):\n\n", len(events))
	result += "ID                                     | Summary                      | Start                | End\n"
	result += "---------------------------------------+------------------------------+----------------------+----------------------\n"

	for _, e := range events {
		id := truncate(e.ID, 37)
		summary := truncate(e.Summary, 28)
		start := formatTimeDisplay(e.StartTime)
		end := formatTimeDisplay(e.EndTime)
		result += fmt.Sprintf("%-37s | %-28s | %-20s | %s\n", id, summary, start, end)
	}

	result += "\nUse get_event with an event ID to view full details."

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}

func (h *MCPHandler) HandleGetEvent(ctx context.Context, req *mcp.CallToolRequest, args GetEventArgs) (*mcp.CallToolResult, any, error) {
	if args.EventID == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: event_id parameter is required"},
			},
			IsError: true,
		}, nil, nil
	}

	event, err := h.cc.GetEvent(ctx, args.EventID)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	var result string
	result += fmt.Sprintf("Summary:     %s\n", event.Summary)
	result += fmt.Sprintf("ID:          %s\n", event.ID)
	result += fmt.Sprintf("Status:      %s\n", event.Status)
	result += fmt.Sprintf("Start:       %s\n", formatTimeDisplay(event.StartTime))
	result += fmt.Sprintf("End:         %s\n", formatTimeDisplay(event.EndTime))

	if event.Location != "" {
		result += fmt.Sprintf("Location:    %s\n", event.Location)
	}
	if event.Creator != "" {
		result += fmt.Sprintf("Creator:     %s\n", event.Creator)
	}
	if len(event.Attendees) > 0 {
		result += fmt.Sprintf("Attendees:   %v\n", event.Attendees)
	}
	if event.Description != "" {
		result += "-------\n"
		result += event.Description
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}

func (h *MCPHandler) HandleCreateEvent(ctx context.Context, req *mcp.CallToolRequest, args CreateEventArgs) (*mcp.CallToolResult, any, error) {
	if args.Summary == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: summary parameter is required (event title)"},
			},
			IsError: true,
		}, nil, nil
	}

	if args.StartTime == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: start_time parameter is required (RFC3339 format)"},
			},
			IsError: true,
		}, nil, nil
	}

	if args.EndTime == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: end_time parameter is required (RFC3339 format)"},
			},
			IsError: true,
		}, nil, nil
	}

	event, err := h.cc.CreateEvent(ctx, args.Summary, args.StartTime, args.EndTime, args.Description, args.Location)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	var result string
	result += fmt.Sprintf("Event created successfully:\n\n")
	result += fmt.Sprintf("Summary:     %s\n", event.Summary)
	result += fmt.Sprintf("ID:          %s\n", event.ID)
	result += fmt.Sprintf("Start:       %s\n", formatTimeDisplay(event.StartTime))
	result += fmt.Sprintf("End:         %s\n", formatTimeDisplay(event.EndTime))

	if event.Location != "" {
		result += fmt.Sprintf("Location:    %s\n", event.Location)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}

func (h *MCPHandler) HandleSearchEvents(ctx context.Context, req *mcp.CallToolRequest, args SearchEventsArgs) (*mcp.CallToolResult, any, error) {
	if args.Query == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: query parameter is required"},
			},
			IsError: true,
		}, nil, nil
	}

	events, err := h.cc.SearchEvents(ctx, args.Query, args.MaxResults)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	if len(events) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("No events found matching query: %s", args.Query)},
			},
		}, nil, nil
	}

	var result string
	result += fmt.Sprintf("Search Results for \"%s\" (%d):\n\n", args.Query, len(events))
	result += "ID                                     | Summary                      | Start                | End\n"
	result += "---------------------------------------+------------------------------+----------------------+----------------------\n"

	for _, e := range events {
		id := truncate(e.ID, 37)
		summary := truncate(e.Summary, 28)
		start := formatTimeDisplay(e.StartTime)
		end := formatTimeDisplay(e.EndTime)
		result += fmt.Sprintf("%-37s | %-28s | %-20s | %s\n", id, summary, start, end)
	}

	result += "\nUse get_event with an event ID to view full details."

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
