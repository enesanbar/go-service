package messages

import (
	"bytes"
	"encoding/json"
	"time"
)

// WIP
type Metadata struct {
	PublisherName string    `json:"publisherName"`
	MessageName   string    `json:"messageName"`
	PublishDate   time.Time `json:"publishDate"`
	Traceparent   string    `json:"traceparent"`
	Tracestate    string    `json:"tracestate"`
	SpanID        string    `json:"spanId"`
}

type Message[T any] struct {
	Metadata Metadata `json:"metadata"`
	Payload  T        `json:"payload"`
}

func (m *Message[T]) UnmarshalPayload(o any) error {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(m.Payload)
	return json.NewDecoder(buf).Decode(o)
}

func (m *Metadata) GetPublisherName() string {
	return m.PublisherName
}

func (m *Metadata) GetMessageName() string {
	return m.MessageName
}

func (m *Metadata) GetPublishDate() string {
	return m.PublishDate.Format(time.RFC3339)
}

func (m *Metadata) GetTraceparent() string {
	return m.Traceparent
}

func (m *Metadata) GetTracestate() string {
	return m.Tracestate
}
