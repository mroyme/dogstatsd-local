package pretty

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/mroyme/dogstatsd-local/internal/messages"
	"strings"
)

type Handler struct {
	Logger     *log.Logger
	Theme      *CatppuccinAdaptiveTheme
	ExtraTags  []string
	NameWidth  int
	ValueWidth int
}

func (h *Handler) New() messages.OutputHandler {
	return func(msg []byte) error {
		dMsg, err := messages.ParseDogStatsDMessage(msg)
		if err != nil {
			h.Logger.Error(err)
			return nil
		}

		switch dMsg.Type() {
		case messages.MetricMessageType:
			metric, ok := dMsg.(messages.DogStatsDMetric)
			if !ok {
				return nil
			}
			str := h.StyledMetricType(metric)
			str += h.StyledMetricName(metric, h.NameWidth)
			str += h.StyledMetricValue(metric, h.ValueWidth)
			str += h.StyledTags(metric, h.ExtraTags)
			fmt.Println(str)

		case messages.ServiceCheckMessageType:
			sc, ok := dMsg.(messages.DogStatsDServiceCheck)
			if !ok {
				return nil
			}
			str := h.StyledServiceCheckStatus(sc)
			str += h.StyledServiceCheckName(sc, h.NameWidth)
			str += h.StyledServiceCheckHostname(sc)
			str += h.StyledServiceCheckMessage(sc)
			str += h.StyledServiceCheckTags(sc, h.ExtraTags)
			fmt.Println(str)
		}

		return nil
	}
}

func (h *Handler) StyledMetricType(metric messages.DogStatsDMetric) string {
	var fg lipgloss.AdaptiveColor
	metricType := metric.MetricType
	switch metricType {
	case messages.CounterMetricType:
		fg = h.Theme.Green()
	case messages.HistogramMetricType:
		fg = h.Theme.Blue()
	case messages.GaugeMetricType:
		fg = h.Theme.Teal()
	case messages.TimerMetricType:
		fg = h.Theme.Mauve()
	case messages.SetMetricType:
		fg = h.Theme.Pink()
	default:
		fg = h.Theme.Text()
	}
	style := lipgloss.NewStyle().
		Width(11).
		Underline(true).
		Foreground(fg)
	return style.Render(strings.ToUpper(metricType.String()))
}

func (h *Handler) StyledMetricName(metric messages.DogStatsDMetric, width int) string {
	// Minimum supported width is 50
	if width < 50 {
		width = 50
	}
	namespace := metric.Namespace
	name := metric.Name
	lenNamespace := len(namespace)
	lenName := len(name)
	textLen := lenNamespace + lenName

	// 3 for the separator " | " + 1 for gap with the next field
	if textLen > width-4 {
		diff := textLen - (width - 4)
		if lenName-diff > 20 {
			name = name[:lenName-diff-1] + "~"
		} else if lenNamespace-diff > 20 {
			namespace = namespace[:lenNamespace-diff-1] + "~"
		} else {
			sub := diff / 2
			name = name[:lenName-sub-1] + "~"
			if diff%2 != 0 {
				sub++
			}
			namespace = namespace[:lenNamespace-sub-1] + "~"
		}
	}
	var text string
	if lenNamespace != 0 {
		text += lipgloss.NewStyle().Foreground(h.Theme.Lavender()).Render(namespace) + " | "
	}
	text += lipgloss.NewStyle().Bold(true).Foreground(h.Theme.Pink()).Render(name)
	return lipgloss.NewStyle().
		Width(width).
		MaxWidth(width).
		Render(text)
}

func (h *Handler) StyledMetricValue(metric messages.DogStatsDMetric, width int) string {
	value := fmt.Sprintf("%.2f", metric.FloatValue)
	if len(value) > width {
		value = value[:width-3] + "..."
	}
	style := lipgloss.NewStyle().
		Width(15).
		MaxWidth(15).
		Bold(true).
		Foreground(h.Theme.Sapphire())
	if metric.MetricType == messages.TimerMetricType {
		value += "ms"
	}
	return style.Render(value)
}

func (h *Handler) StyledTags(metric messages.DogStatsDMetric, extraTags []string) string {
	style := lipgloss.NewStyle().
		Foreground(h.Theme.Overlay0()).
		Italic(true)
	return style.SetString(append(extraTags, metric.Tags...)...).Render()
}

func (h *Handler) StyledServiceCheckStatus(sc messages.DogStatsDServiceCheck) string {
	var fg lipgloss.AdaptiveColor
	switch sc.Status {
	case messages.ServiceCheckOK:
		fg = h.Theme.Green()
	case messages.ServiceCheckWarning:
		fg = h.Theme.Yellow()
	case messages.ServiceCheckCritical:
		fg = h.Theme.Red()
	default:
		fg = h.Theme.Overlay1()
	}
	style := lipgloss.NewStyle().
		Width(11).
		Underline(true).
		Foreground(fg)
	return style.Render(sc.Status.String())
}

func (h *Handler) StyledServiceCheckName(sc messages.DogStatsDServiceCheck, width int) string {
	if width < 50 {
		width = 50
	}
	text := lipgloss.NewStyle().Bold(true).Foreground(h.Theme.Pink()).Render(sc.Name)
	return lipgloss.NewStyle().
		Width(width).
		MaxWidth(width).
		Render(text)
}

func (h *Handler) StyledServiceCheckHostname(sc messages.DogStatsDServiceCheck) string {
	style := lipgloss.NewStyle().
		Width(15).
		MaxWidth(15).
		Bold(true).
		Foreground(h.Theme.Sapphire())
	return style.Render(sc.Hostname)
}

func (h *Handler) StyledServiceCheckMessage(sc messages.DogStatsDServiceCheck) string {
	if sc.Message == "" {
		return ""
	}
	style := lipgloss.NewStyle().
		Width(30).
		MaxWidth(30).
		Foreground(h.Theme.Subtext1())
	return style.Render(sc.Message)
}

func (h *Handler) StyledServiceCheckTags(sc messages.DogStatsDServiceCheck, extraTags []string) string {
	style := lipgloss.NewStyle().
		Foreground(h.Theme.Overlay0()).
		Italic(true)
	return style.SetString(append(extraTags, sc.Tags...)...).Render()
}
