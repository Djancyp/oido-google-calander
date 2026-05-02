package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type CalendarSettings struct {
	ServiceAccountJSON string
	CalendarID         string
	MaxResults         int64
	Timezone           string
}

type EventSummary struct {
	ID        string
	Summary   string
	StartTime string
	EndTime   string
	Location  string
}

type EventDetail struct {
	ID          string
	Summary     string
	Description string
	StartTime   string
	EndTime     string
	Location    string
	Status      string
	Creator     string
	Attendees   []string
}

func DefaultCalendarSettings() *CalendarSettings {
	return &CalendarSettings{
		CalendarID: "primary",
		MaxResults: 20,
	}
}

func parseCalendarSettings() (*CalendarSettings, error) {
	settings := DefaultCalendarSettings()

	settings.ServiceAccountJSON = os.Getenv("CALENDAR_SERVICE_ACCOUNT_JSON")

	if settings.ServiceAccountJSON == "" {
		b64 := os.Getenv("CALENDAR_SERVICE_ACCOUNT_B64")
		if b64 != "" {
			decoded, err := base64.StdEncoding.DecodeString(b64)
			if err != nil {
				return nil, fmt.Errorf("failed to decode CALENDAR_SERVICE_ACCOUNT_B64: %w (use standard base64 encoding)", err)
			}
			settings.ServiceAccountJSON = string(decoded)
		}
	}

	if id := os.Getenv("CALENDAR_ID"); id != "" {
		settings.CalendarID = id
	}

	if mr := os.Getenv("CALENDAR_MAX_RESULTS"); mr != "" {
		if n, err := strconv.ParseInt(mr, 10, 64); err == nil {
			settings.MaxResults = n
		}
	}

	if tz := os.Getenv("CALENDAR_TIMEZONE"); tz != "" {
		settings.Timezone = tz
	}

	if settings.ServiceAccountJSON == "" {
		return nil, fmt.Errorf("missing required env var: set CALENDAR_SERVICE_ACCOUNT_JSON (raw JSON) or CALENDAR_SERVICE_ACCOUNT_B64 (base64-encoded JSON)")
	}

	return settings, nil
}

type CalendarClient struct {
	settings *CalendarSettings
}

func NewCalendarClient() (*CalendarClient, error) {
	settings, err := parseCalendarSettings()
	if err != nil {
		return nil, err
	}
	return &CalendarClient{settings: settings}, nil
}

func (c *CalendarClient) getService(ctx context.Context) (*calendar.Service, error) {
	input := strings.TrimSpace(c.settings.ServiceAccountJSON)
	raw := []byte(input)

	conf, err := google.JWTConfigFromJSON(raw, calendar.CalendarScope)
	if err != nil {
		preview := string(raw)
		if len(preview) > 120 {
			preview = preview[:120] + "..."
		}
		hexPreview := ""
		for i, b := range raw {
			if i >= 30 {
				break
			}
			hexPreview += fmt.Sprintf("%02x ", b)
		}
		detail := fmt.Sprintf(
			"error: %v\nfirst 30 bytes (hex): %s\nstring preview: %s\n\nTip: use CALENDAR_SERVICE_ACCOUNT_B64 instead of CALENDAR_SERVICE_ACCOUNT_JSON to avoid env var escaping issues.\n  base64 -w0 /path/to/key.json  → paste result into CALENDAR_SERVICE_ACCOUNT_B64",
			err, hexPreview, preview)
		return nil, fmt.Errorf("%s", detail)
	}

	client := conf.Client(ctx)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}

	return srv, nil
}

func (c *CalendarClient) ListEvents(ctx context.Context, maxResults int64, timeMin string) ([]EventSummary, error) {
	if maxResults <= 0 {
		maxResults = c.settings.MaxResults
	}

	srv, err := c.getService(ctx)
	if err != nil {
		return nil, err
	}

	req := srv.Events.List(c.settings.CalendarID).
		MaxResults(maxResults).
		OrderBy("startTime").
		SingleEvents(true)

	if timeMin != "" {
		req.TimeMin(timeMin)
	} else {
		req.TimeMin(time.Now().Format(time.RFC3339))
	}

	events, err := req.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	if len(events.Items) == 0 {
		return []EventSummary{}, nil
	}

	var summaries []EventSummary
	for _, event := range events.Items {
		start := event.Start.DateTime
		if start == "" {
			start = event.Start.Date
		}
		end := event.End.DateTime
		if end == "" {
			end = event.End.Date
		}

		summaries = append(summaries, EventSummary{
			ID:        event.Id,
			Summary:   event.Summary,
			StartTime: start,
			EndTime:   end,
			Location:  event.Location,
		})
	}

	return summaries, nil
}

func (c *CalendarClient) GetEvent(ctx context.Context, eventID string) (*EventDetail, error) {
	srv, err := c.getService(ctx)
	if err != nil {
		return nil, err
	}

	event, err := srv.Events.Get(c.settings.CalendarID, eventID).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	start := event.Start.DateTime
	if start == "" {
		start = event.Start.Date
	}
	end := event.End.DateTime
	if end == "" {
		end = event.End.Date
	}

	var attendees []string
	for _, a := range event.Attendees {
		attendees = append(attendees, a.Email)
	}

	creator := ""
	if event.Creator != nil {
		creator = event.Creator.Email
	}

	return &EventDetail{
		ID:          event.Id,
		Summary:     event.Summary,
		Description: event.Description,
		StartTime:   start,
		EndTime:     end,
		Location:    event.Location,
		Status:      event.Status,
		Creator:     creator,
		Attendees:   attendees,
	}, nil
}

func (c *CalendarClient) CreateEvent(ctx context.Context, summary, startTime, endTime, description, location string) (*EventDetail, error) {
	srv, err := c.getService(ctx)
	if err != nil {
		return nil, err
	}

	event := &calendar.Event{
		Summary:     summary,
		Description: description,
		Location:    location,
		Start: &calendar.EventDateTime{
			DateTime: startTime,
			TimeZone: c.settings.Timezone,
		},
		End: &calendar.EventDateTime{
			DateTime: endTime,
			TimeZone: c.settings.Timezone,
		},
	}

	created, err := srv.Events.Insert(c.settings.CalendarID, event).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	start := created.Start.DateTime
	if start == "" {
		start = created.Start.Date
	}
	end := created.End.DateTime
	if end == "" {
		end = created.End.Date
	}

	var attendees []string
	for _, a := range created.Attendees {
		attendees = append(attendees, a.Email)
	}

	creator := ""
	if created.Creator != nil {
		creator = created.Creator.Email
	}

	return &EventDetail{
		ID:          created.Id,
		Summary:     created.Summary,
		Description: created.Description,
		StartTime:   start,
		EndTime:     end,
		Location:    created.Location,
		Status:      created.Status,
		Creator:     creator,
		Attendees:   attendees,
	}, nil
}

func (c *CalendarClient) SearchEvents(ctx context.Context, query string, maxResults int64) ([]EventSummary, error) {
	if query == "" {
		return nil, fmt.Errorf("query is required")
	}

	if maxResults <= 0 {
		maxResults = c.settings.MaxResults
	}

	srv, err := c.getService(ctx)
	if err != nil {
		return nil, err
	}

	events, err := srv.Events.List(c.settings.CalendarID).
		Q(query).
		MaxResults(maxResults).
		OrderBy("startTime").
		SingleEvents(true).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search events: %w", err)
	}

	if len(events.Items) == 0 {
		return []EventSummary{}, nil
	}

	var summaries []EventSummary
	for _, event := range events.Items {
		start := event.Start.DateTime
		if start == "" {
			start = event.Start.Date
		}
		end := event.End.DateTime
		if end == "" {
			end = event.End.Date
		}

		summaries = append(summaries, EventSummary{
			ID:        event.Id,
			Summary:   event.Summary,
			StartTime: start,
			EndTime:   end,
			Location:  event.Location,
		})
	}

	return summaries, nil
}
