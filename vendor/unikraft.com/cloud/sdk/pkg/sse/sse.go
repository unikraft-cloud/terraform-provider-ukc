// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

// Package sse implements Server-Sent Events (SSE) parsing.
package sse

import (
	"bufio"
	"bytes"
	"io"
)

// Event represents a Server-Sent Event.
type Event struct {
	ID    string
	Event string
	Data  []byte
	Retry string
}

// Reader reads Server-Sent Events from a stream.
type Reader struct {
	scanner *bufio.Scanner
}

// NewReader creates a new SSE reader for the given io.Reader.
func NewReader(r io.Reader) *Reader {
	scanner := bufio.NewScanner(r)
	// Use a custom split function to handle newlines
	scanner.Split(scanLines)

	return &Reader{
		scanner: scanner,
	}
}

// ReadEvent reads a single event from the stream.
// Returns io.EOF when the stream is closed.
func (r *Reader) ReadEvent() (*Event, error) {
	event := &Event{}
	inEvent := false

	// Process lines until we complete an event
	for r.scanner.Scan() {
		line := r.scanner.Bytes()

		// Empty line marks the end of an event
		if len(line) == 0 {
			if inEvent {
				// End of event
				return event, nil
			}
			continue
		}

		inEvent = true

		// Check for comment
		if len(line) > 0 && line[0] == ':' {
			// Comment line, ignore
			continue
		}

		// Process field
		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) < 2 {
			// Invalid format, skip this line
			continue
		}

		fieldName := parts[0]

		var fieldValue []byte
		if len(parts) > 1 {
			fieldValue = bytes.TrimPrefix(parts[1], []byte(" "))
		}

		// Handle the field based on its name
		switch string(fieldName) {
		case "event":
			event.Event = string(fieldValue)
		case "data":
			if event.Data == nil {
				event.Data = fieldValue
			} else {
				// Data fields are concatenated with a newline
				event.Data = append(append(event.Data, '\n'), fieldValue...)
			}
		case "id":
			event.ID = string(fieldValue)
		case "retry":
			event.Retry = string(fieldValue)
		}
	}

	// Check for scanner error
	if err := r.scanner.Err(); err != nil {
		return nil, err
	}

	// If we've processed some data but hit EOF
	if inEvent {
		return event, nil
	}

	// End of stream
	return nil, io.EOF
}

// scanLines is a custom split function for bufio.Scanner that handles different
// line endings.
func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// Try to find LF (\n) or CRLF (\r\n)
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a line
		advance = i + 1
		if i > 0 && data[i-1] == '\r' {
			// CRLF
			token = data[:i-1]
		} else {
			// LF
			token = data[:i]
		}
		return advance, token, nil
	}

	// If we're at EOF, return the remaining data
	if atEOF {
		return len(data), data, nil
	}

	// Request more data
	return 0, nil, nil
}
