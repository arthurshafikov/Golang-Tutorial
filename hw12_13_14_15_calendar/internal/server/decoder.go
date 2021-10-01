package server

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/thewolf27/hw12_13_14_15_calendar/internal/storage"
)

const (
	requestDateKey = "date"
)

func DecodeJSONEvent(body io.ReadCloser, strict bool) (storage.Event, error) {
	var eventReq EventRequest
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&eventReq); err != nil {
		return storage.Event{}, fmt.Errorf(ErrCantParseEventFormat, err)
	}
	event, err := ParseEvent(eventReq, strict)
	if err != nil {
		return storage.Event{}, fmt.Errorf("cant parseEvent %w", err)
	}
	return event, nil
}

func DecodeJSONDate(body io.ReadCloser) (time.Time, error) {
	var dateInput interface{}
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&dateInput); err != nil {
		return time.Time{}, fmt.Errorf(ErrCantParseDateFormat, err)
	}

	decodedMap, ok := dateInput.(map[string]interface{})
	if !ok {
		return time.Time{}, fmt.Errorf("decodeJSONDate cant assign input to map")
	}

	decodedDateStr, ok := decodedMap[requestDateKey].(string)
	if !ok {
		return time.Time{}, fmt.Errorf("decodeJSONDate cant assign date to string")
	}

	t, err := time.Parse(storage.RequestDateFormat, decodedDateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf(ErrCantParseDateFormat, err)
	}

	return t, nil
}
