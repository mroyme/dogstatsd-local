package messages

import (
	"testing"
	"time"
)

func TestParseDogStatsDMetricMessage(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantName     string
		wantNs       string
		wantType     DogStatsDMetricType
		wantValue    float64
		wantTags     []string
		wantRate     float64
		wantExtras   []string
		wantDuration time.Duration
		wantErr      bool
	}{
		{
			name:      "counter with namespace",
			input:     "page.views:1|c",
			wantNs:    "page",
			wantName:  "views",
			wantType:  CounterMetricType,
			wantValue: 1.0,
			wantRate:  1.0,
		},
		{
			name:      "gauge without namespace",
			input:     "fuel:0.5|g",
			wantName:  "fuel",
			wantType:  GaugeMetricType,
			wantValue: 0.5,
			wantRate:  1.0,
		},
		{
			name:      "histogram with sample rate",
			input:     "song.length:240|h|@0.5",
			wantNs:    "song",
			wantName:  "length",
			wantType:  HistogramMetricType,
			wantValue: 240.0,
			wantRate:  0.5,
		},
		{
			name:      "distribution metric",
			input:     "page.views:1:2:32|d",
			wantNs:    "page",
			wantName:  "views",
			wantType:  DistributionMetricType,
			wantValue: 1.0,
			wantRate:  1.0,
		},
		{
			name:      "set metric",
			input:     "users.uniques:1234|s",
			wantNs:    "users",
			wantName:  "uniques",
			wantType:  SetMetricType,
			wantValue: 1234.0,
			wantRate:  1.0,
		},
		{
			name:         "timer with duration",
			input:        "request.time:150|ms",
			wantNs:       "request",
			wantName:     "time",
			wantType:     TimerMetricType,
			wantValue:    150.0,
			wantRate:     1.0,
			wantDuration: 150 * time.Millisecond,
		},
		{
			name:      "counter with tags",
			input:     "users.online:1|c|#country:china",
			wantNs:    "users",
			wantName:  "online",
			wantType:  CounterMetricType,
			wantValue: 1.0,
			wantTags:  []string{"country:china"},
			wantRate:  1.0,
		},
		{
			name:      "counter with tags and sample rate",
			input:     "users.online:1|c|@0.5|#country:china",
			wantNs:    "users",
			wantName:  "online",
			wantType:  CounterMetricType,
			wantValue: 1.0,
			wantTags:  []string{"country:china"},
			wantRate:  0.5,
		},
		{
			name:      "counter with multiple tags",
			input:     "users.online:1|c|#country:china,env:prod",
			wantNs:    "users",
			wantName:  "online",
			wantType:  CounterMetricType,
			wantValue: 1.0,
			wantTags:  []string{"country:china", "env:prod"},
			wantRate:  1.0,
		},
		{
			name:       "counter with extras",
			input:      "metric.name:1|c|#tag1:value1|extra",
			wantNs:     "metric",
			wantName:   "name",
			wantType:   CounterMetricType,
			wantValue:  1.0,
			wantTags:   []string{"tag1:value1"},
			wantExtras: []string{"extra"},
			wantRate:   1.0,
		},
		{
			name:    "missing value",
			input:   "metric.name|c",
			wantErr: true,
		},
		{
			name:    "missing type",
			input:   "metric.name:1",
			wantErr: true,
		},
		{
			name:    "unknown metric type",
			input:   "metric.name:1|x",
			wantErr: true,
		},
		{
			name:    "invalid value",
			input:   "metric.name:abc|c",
			wantErr: true,
		},
		{
			name:    "invalid sample rate",
			input:   "metric.name:1|c|@abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := ParseDogStatsDMessage([]byte(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			metric, ok := msg.(DogStatsDMetric)
			if !ok {
				t.Fatalf("expected DogStatsDMetric, got %T", msg)
			}
			if metric.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", metric.Name, tt.wantName)
			}
			if metric.Namespace != tt.wantNs {
				t.Errorf("Namespace = %q, want %q", metric.Namespace, tt.wantNs)
			}
			if metric.MetricType != tt.wantType {
				t.Errorf("MetricType = %v, want %v", metric.MetricType, tt.wantType)
			}
			if metric.FloatValue != tt.wantValue {
				t.Errorf("FloatValue = %v, want %v", metric.FloatValue, tt.wantValue)
			}
			if metric.SampleRate != tt.wantRate {
				t.Errorf("SampleRate = %v, want %v", metric.SampleRate, tt.wantRate)
			}
			if tt.wantDuration != 0 && metric.DurationValue != tt.wantDuration {
				t.Errorf("DurationValue = %v, want %v", metric.DurationValue, tt.wantDuration)
			}
			if len(metric.Tags) != len(tt.wantTags) {
				t.Errorf("Tags = %v, want %v", metric.Tags, tt.wantTags)
			} else {
				for i, tag := range metric.Tags {
					if tag != tt.wantTags[i] {
						t.Errorf("Tags[%d] = %q, want %q", i, tag, tt.wantTags[i])
					}
				}
			}
			if len(metric.Extras) != len(tt.wantExtras) {
				t.Errorf("Extras = %v, want %v", metric.Extras, tt.wantExtras)
			} else {
				for i, extra := range metric.Extras {
					if extra != tt.wantExtras[i] {
						t.Errorf("Extras[%d] = %q, want %q", i, extra, tt.wantExtras[i])
					}
				}
			}
		})
	}
}

func TestParseDogStatsDServiceCheckMessage(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantName   string
		wantStatus ServiceCheckStatus
		wantTags   []string
		wantHost   string
		wantMsg    string
		wantErr    bool
	}{
		{
			name:       "critical with tags and hostname",
			input:      "_sc|my_service|2|#tag1:value1,tag2|h:host1",
			wantName:   "my_service",
			wantStatus: ServiceCheckCritical,
			wantTags:   []string{"tag1:value1", "tag2"},
			wantHost:   "host1",
		},
		{
			name:       "ok with tags and hostname",
			input:      "_sc|db_check|0|#env:prod|h:db1.example.com",
			wantName:   "db_check",
			wantStatus: ServiceCheckOK,
			wantTags:   []string{"env:prod"},
			wantHost:   "db1.example.com",
		},
		{
			name:       "warning with hostname only",
			input:      "_sc|cache_check|1|h:cache1",
			wantName:   "cache_check",
			wantStatus: ServiceCheckWarning,
			wantTags:   []string{},
			wantHost:   "cache1",
		},
		{
			name:       "unknown with no extras",
			input:      "_sc|unknown_check|3",
			wantName:   "unknown_check",
			wantStatus: ServiceCheckUnknown,
			wantTags:   []string{},
		},
		{
			name:       "with timestamp",
			input:      "_sc|timed_check|0|#env:prod|d:1234567890|h:host1",
			wantName:   "timed_check",
			wantStatus: ServiceCheckOK,
			wantTags:   []string{"env:prod"},
			wantHost:   "host1",
		},
		{
			name:       "with message",
			input:      "_sc|redis.can_connect|1|#env:prod,instance:cache-01|m:High latency detected (250ms)",
			wantName:   "redis.can_connect",
			wantStatus: ServiceCheckWarning,
			wantTags:   []string{"env:prod", "instance:cache-01"},
			wantMsg:    "High latency detected (250ms)",
		},
		{
			name:       "with all fields",
			input:      "_sc|db_check|2|#env:prod|d:1234567890|h:db1.example.com|m:Connection refused",
			wantName:   "db_check",
			wantStatus: ServiceCheckCritical,
			wantTags:   []string{"env:prod"},
			wantHost:   "db1.example.com",
			wantMsg:    "Connection refused",
		},
		{
			name:       "message with pipe character",
			input:      "_sc|Redis connection|2|#env:dev|m:Error: timeout|retrying",
			wantName:   "Redis connection",
			wantStatus: ServiceCheckCritical,
			wantTags:   []string{"env:dev"},
			wantMsg:    "Error: timeout|retrying",
		},
		{
			name:    "missing status",
			input:   "_sc|missing_status",
			wantErr: true,
		},
		{
			name:    "invalid status",
			input:   "_sc|bad_status|abc",
			wantErr: true,
		},
		{
			name:    "status out of range",
			input:   "_sc|out_of_range|5",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := ParseDogStatsDMessage([]byte(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			sc, ok := msg.(DogStatsDServiceCheck)
			if !ok {
				t.Fatalf("expected DogStatsDServiceCheck, got %T", msg)
			}
			if sc.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", sc.Name, tt.wantName)
			}
			if sc.Status != tt.wantStatus {
				t.Errorf("Status = %v, want %v", sc.Status, tt.wantStatus)
			}
			if len(sc.Tags) != len(tt.wantTags) {
				t.Errorf("Tags = %v, want %v", sc.Tags, tt.wantTags)
			} else {
				for i, tag := range sc.Tags {
					if tag != tt.wantTags[i] {
						t.Errorf("Tags[%d] = %q, want %q", i, tag, tt.wantTags[i])
					}
				}
			}
			if sc.Hostname != tt.wantHost {
				t.Errorf("Hostname = %q, want %q", sc.Hostname, tt.wantHost)
			}
			if sc.Message != tt.wantMsg {
				t.Errorf("Message = %q, want %q", sc.Message, tt.wantMsg)
			}
		})
	}
}

func TestParseDogStatsDEventMessage(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		wantTitle      string
		wantText       string
		wantPriority   EventPriority
		wantAlertType  EventAlertType
		wantTags       []string
		wantHostname   string
		wantAggKey     string
		wantSourceType string
		wantErr        bool
	}{
		{
			name:          "simple event",
			input:         "_e{6,15}:title1|text with pipes",
			wantTitle:     "title1",
			wantText:      "text with pipes",
			wantPriority:  EventPriorityNormal,
			wantAlertType: EventAlertTypeInfo,
			wantTags:      []string{},
		},
		{
			name:          "event with tags and alert type",
			input:         "_e{21,21}:An exception occurred|Cannot parse CSV file|t:warning|#err_type:bad_file",
			wantTitle:     "An exception occurred",
			wantText:      "Cannot parse CSV file",
			wantPriority:  EventPriorityNormal,
			wantAlertType: EventAlertTypeWarning,
			wantTags:      []string{"err_type:bad_file"},
		},
		{
			name:           "event with all fields",
			input:          "_e{5,17}:title|Cannot parse JSON|h:host1|p:low|t:error|k:aggkey1|s:source1|#env:prod,region:us",
			wantTitle:      "title",
			wantText:       "Cannot parse JSON",
			wantPriority:   EventPriorityLow,
			wantAlertType:  EventAlertTypeError,
			wantHostname:   "host1",
			wantAggKey:     "aggkey1",
			wantSourceType: "source1",
			wantTags:       []string{"env:prod", "region:us"},
		},
		{
			name:          "event with text containing escaped newline",
			input:         "_e{5,12}:title|line1\\nline2|t:success",
			wantTitle:     "title",
			wantText:      "line1\\nline2",
			wantAlertType: EventAlertTypeSuccess,
			wantPriority:  EventPriorityNormal,
			wantTags:      []string{},
		},
		{
			name:    "missing closing brace",
			input:   "_e{5,12title|text",
			wantErr: true,
		},
		{
			name:    "missing comma in lengths",
			input:   "_e{5}:title|text",
			wantErr: true,
		},
		{
			name:    "non-numeric title length",
			input:   "_e{abc,3}:title|text",
			wantErr: true,
		},
		{
			name:    "title length exceeds payload",
			input:   "_e{100,3}:title|text",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := ParseDogStatsDMessage([]byte(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			ev, ok := msg.(DogStatsDEvent)
			if !ok {
				t.Fatalf("expected DogStatsDEvent, got %T", msg)
			}
			if ev.Title != tt.wantTitle {
				t.Errorf("Title = %q, want %q", ev.Title, tt.wantTitle)
			}
			if ev.Text != tt.wantText {
				t.Errorf("Text = %q, want %q", ev.Text, tt.wantText)
			}
			if ev.Priority != tt.wantPriority {
				t.Errorf("Priority = %q, want %q", ev.Priority, tt.wantPriority)
			}
			if ev.AlertType != tt.wantAlertType {
				t.Errorf("AlertType = %q, want %q", ev.AlertType, tt.wantAlertType)
			}
			if ev.Hostname != tt.wantHostname {
				t.Errorf("Hostname = %q, want %q", ev.Hostname, tt.wantHostname)
			}
			if ev.AggregationKey != tt.wantAggKey {
				t.Errorf("AggregationKey = %q, want %q", ev.AggregationKey, tt.wantAggKey)
			}
			if ev.SourceType != tt.wantSourceType {
				t.Errorf("SourceType = %q, want %q", ev.SourceType, tt.wantSourceType)
			}
			if len(ev.Tags) != len(tt.wantTags) {
				t.Errorf("Tags = %v, want %v", ev.Tags, tt.wantTags)
			} else {
				for i, tag := range ev.Tags {
					if tag != tt.wantTags[i] {
						t.Errorf("Tags[%d] = %q, want %q", i, tag, tt.wantTags[i])
					}
				}
			}
		})
	}
}

func TestDogStatsDMetricTypeString(t *testing.T) {
	tests := []struct {
		metricType DogStatsDMetricType
		want       string
	}{
		{GaugeMetricType, "gauge"},
		{CounterMetricType, "counter"},
		{SetMetricType, "set"},
		{TimerMetricType, "timer"},
		{HistogramMetricType, "histogram"},
		{DistributionMetricType, "distribution"},
		{DogStatsDMetricType(99), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.metricType.String(); got != tt.want {
			t.Errorf("DogStatsDMetricType(%d).String() = %q, want %q", tt.metricType, got, tt.want)
		}
	}
}

func TestDogStatsDMessageTypeString(t *testing.T) {
	tests := []struct {
		msgType DogStatsDMessageType
		want    string
	}{
		{MetricMessageType, "metric"},
		{ServiceCheckMessageType, "service_check"},
		{EventMessageType, "event"},
		{DogStatsDMessageType(99), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.msgType.String(); got != tt.want {
			t.Errorf("DogStatsDMessageType(%d).String() = %q, want %q", tt.msgType, got, tt.want)
		}
	}
}

func TestServiceCheckStatusString(t *testing.T) {
	tests := []struct {
		status ServiceCheckStatus
		want   string
	}{
		{ServiceCheckOK, "OK"},
		{ServiceCheckWarning, "WARNING"},
		{ServiceCheckCritical, "CRITICAL"},
		{ServiceCheckUnknown, "UNKNOWN"},
		{ServiceCheckStatus(99), "UNKNOWN"},
	}
	for _, tt := range tests {
		if got := tt.status.String(); got != tt.want {
			t.Errorf("ServiceCheckStatus(%d).String() = %q, want %q", tt.status, got, tt.want)
		}
	}
}

func TestDogStatsDMessageType(t *testing.T) {
	metricMsg, err := ParseDogStatsDMessage([]byte("metric.name:1|c"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if metricMsg.Type() != MetricMessageType {
		t.Errorf("metric Type() = %v, want %v", metricMsg.Type(), MetricMessageType)
	}

	scMsg, err := ParseDogStatsDMessage([]byte("_sc|my_service|0"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if scMsg.Type() != ServiceCheckMessageType {
		t.Errorf("service check Type() = %v, want %v", scMsg.Type(), ServiceCheckMessageType)
	}

	evMsg, err := ParseDogStatsDMessage([]byte("_e{5,4}:title|text"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if evMsg.Type() != EventMessageType {
		t.Errorf("event Type() = %v, want %v", evMsg.Type(), EventMessageType)
	}
}

func TestDogStatsDMessageData(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"service check", "_sc|my_service|0|#env:prod"},
		{"event", "_e{5,4}:title|text"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := []byte(tt.input)
			msg, err := ParseDogStatsDMessage(raw)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(msg.Data()) != string(raw) {
				t.Errorf("Data() = %q, want %q", msg.Data(), raw)
			}
		})
	}
}
