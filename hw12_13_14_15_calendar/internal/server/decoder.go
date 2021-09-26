package server

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage"
)

func DecodeJSONEvent(body io.ReadCloser, strict bool) (storage.Event, error) {
	var eventReq storage.EventRequest
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&eventReq); err != nil {
		return storage.Event{}, fmt.Errorf("cant parseEvent %w", err)
	}
	event, err := storage.ParseEvent(eventReq, strict)
	if err != nil {
		return storage.Event{}, fmt.Errorf("cant parseEvent %w", err)
	}
	return event, nil
}

func DecodeJSONDate(body io.ReadCloser) (time.Time, error) {
	var dateInput interface{}
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&dateInput); err != nil {
		return time.Time{}, fmt.Errorf("decodeJSONDate cant decode json date %w", err)
	}

	decodedMap, ok := dateInput.(map[string]interface{})
	if !ok {
		return time.Time{}, fmt.Errorf("decodeJSONDate cant assign input to map")
	}

	decodedDateStr, ok := decodedMap["date"].(string)
	if !ok {
		return time.Time{}, fmt.Errorf("decodeJSONDate cant assign date to strng")
	}

	t, err := time.Parse(storage.RequestDateFormat, decodedDateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("decodeJSONDate cant parse time %w", err)
	}

	return t, nil
}
