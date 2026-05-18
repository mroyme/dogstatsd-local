package messages

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type EventPriority string

const (
	EventPriorityNormal EventPriority = "normal"
	EventPriorityLow    EventPriority = "low"
)

type EventAlertType string

const (
	EventAlertTypeInfo    EventAlertType = "info"
	EventAlertTypeWarning EventAlertType = "warning"
	EventAlertTypeError   EventAlertType = "error"
	EventAlertTypeSuccess EventAlertType = "success"
)

type DogStatsDEvent struct {
	Title          string
	Text           string
	Timestamp      time.Time
	Hostname       string
	Priority       EventPriority
	AlertType      EventAlertType
	AggregationKey string
	SourceType     string
	Tags           []string
	RawData        []byte
}

func (d DogStatsDEvent) Data() []byte {
	return d.RawData
}

func (DogStatsDEvent) Type() DogStatsDMessageType {
	return EventMessageType
}

func parseDogStatsDEventMessage(buf []byte) (DogStatsDMessage, error) {
	// DogStatsD event format:
	//   _e{<TITLE_LEN>,<TEXT_LEN>}:<TITLE>|<TEXT>|d:<TIMESTAMP>|h:<HOSTNAME>|p:<PRIORITY>|t:<ALERT_TYPE>|k:<AGG_KEY>|s:<SRC_TYPE>|#<TAGS>
	ev := DogStatsDEvent{
		RawData:   buf,
		Timestamp: time.Now(),
		Priority:  EventPriorityNormal,
		AlertType: EventAlertTypeInfo,
		Tags:      make([]string, 0),
	}

	body := string(buf)

	// Parse the _e{titleLen,textLen}: header
	if !strings.HasPrefix(body, "_e{") {
		return nil, errors.New("invalid event: missing _e{ prefix")
	}

	// Find the closing brace
	braceEnd := strings.Index(body, "}")
	if braceEnd < 0 {
		return nil, errors.New("invalid event: missing closing brace")
	}

	// Parse title and text lengths
	lenPart := body[3:braceEnd]
	lenPieces := strings.SplitN(lenPart, ",", 2)
	if len(lenPieces) != 2 {
		return nil, errors.New("invalid event: missing title or text length")
	}
	titleLen, err := strconv.Atoi(lenPieces[0])
	if err != nil {
		return nil, errors.New("invalid event: invalid title length")
	}
	textLen, err := strconv.Atoi(lenPieces[1])
	if err != nil {
		return nil, errors.New("invalid event: invalid text length")
	}

	// After the closing brace there must be a ':'
	headerEnd := braceEnd + 1
	if headerEnd >= len(body) || body[headerEnd] != ':' {
		return nil, errors.New("invalid event: missing ':' after brace")
	}
	payload := body[headerEnd+1:]

	// Extract title using the declared length (UTF-8 byte length)
	if titleLen > len(payload) {
		return nil, errors.New("invalid event: title length exceeds payload")
	}
	ev.Title = payload[:titleLen]
	payload = payload[titleLen:]

	// After the title there must be a '|'
	if len(payload) == 0 || payload[0] != '|' {
		return nil, errors.New("invalid event: missing '|' after title")
	}
	payload = payload[1:]

	// Extract text using the declared length
	if textLen > len(payload) {
		return nil, errors.New("invalid event: text length exceeds payload")
	}
	ev.Text = payload[:textLen]
	payload = payload[textLen:]

	// The remaining payload contains optional metadata fields separated by '|'
	if len(payload) > 0 && payload[0] == '|' {
		payload = payload[1:]
	}

	if payload == "" {
		return ev, nil
	}

	pieces := strings.Split(payload, "|")
	for _, piece := range pieces {
		if strings.HasPrefix(piece, "d:") {
			ts, err := strconv.ParseInt(piece[2:], 10, 64)
			if err == nil {
				ev.Timestamp = time.Unix(ts, 0)
			}
			continue
		}
		if strings.HasPrefix(piece, "h:") {
			ev.Hostname = piece[2:]
			continue
		}
		if strings.HasPrefix(piece, "p:") {
			ev.Priority = EventPriority(piece[2:])
			continue
		}
		if strings.HasPrefix(piece, "t:") {
			ev.AlertType = EventAlertType(piece[2:])
			continue
		}
		if strings.HasPrefix(piece, "k:") {
			ev.AggregationKey = piece[2:]
			continue
		}
		if strings.HasPrefix(piece, "s:") {
			ev.SourceType = piece[2:]
			continue
		}
		if strings.HasPrefix(piece, "#") {
			tags := strings.Split(piece[1:], ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
			ev.Tags = append(ev.Tags, tags...)
			continue
		}
	}

	return ev, nil
}