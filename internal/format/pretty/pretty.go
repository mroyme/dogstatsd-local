package pretty

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/mroyme/dogstatsd-local/internal/messages"
)

// joinTrailing joins non-empty styled fields with spaces.
func joinTrailing(fields ...string) string {
	var parts []string
	for _, f := range fields {
		if f != "" {
			parts = append(parts, f)
		}
	}
	return strings.Join(parts, " ")
}

type Handler struct {
	Logger     *log.Logger
	Theme      *CatppuccinAdaptiveTheme
	ExtraTags  []string
	NameWidth  int
	ValueWidth int
}

func (h *Handler) New() messages.OutputHandler {
	return func(msg messages.DogStatsDMessage) error {
		switch msg.Type() {
		case messages.MetricMessageType:
			metric, ok := msg.(messages.DogStatsDMetric)
			if !ok {
				return nil
			}
			str := h.StyledMetricType(metric)
			str += h.StyledMetricName(metric, h.NameWidth)
			str += h.StyledMetricValue(metric, h.ValueWidth)
			str += h.StyledTags(metric.Tags, h.ExtraTags)
			fmt.Println(str)

		case messages.ServiceCheckMessageType:
			sc, ok := msg.(messages.DogStatsDServiceCheck)
			if !ok {
				return nil
			}
			str := h.StyledServiceCheckStatus(sc)
			str += h.StyledServiceCheckName(sc, h.NameWidth)
			str += joinTrailing(
				h.StyledServiceCheckHostname(sc),
				h.StyledServiceCheckMessage(sc),
				h.StyledTags(sc.Tags, h.ExtraTags),
			)
			fmt.Println(str)

		case messages.EventMessageType:
			ev, ok := msg.(messages.DogStatsDEvent)
			if !ok {
				return nil
			}
			str := h.StyledEventAlertType(ev)
			str += h.StyledEventTitle(ev, h.NameWidth)
			str += joinTrailing(
				h.StyledEventText(ev),
				h.StyledTags(ev.Tags, h.ExtraTags),
			)
			fmt.Println(str)
		}

		return nil
	}
}

// StyledTags renders extra tags + metric tags in a muted italic style with automatic spacing.
func (h *Handler) StyledTags(tags []string, extraTags []string) string {
	style := lipgloss.NewStyle().
		Foreground(h.Theme.Overlay0()).
		Italic(true)
	return style.SetString(append(extraTags, tags...)...).Render()
}
