package pretty

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mroyme/dogstatsd-local/internal/messages"
)

func (h *Handler) StyledMetricType(metric messages.DogStatsDMetric) string {
	var fg lipgloss.AdaptiveColor
	var label string
	switch metric.MetricType {
	case messages.CountMetricType:
		fg = h.Theme.Green()
		label = "COUNT"
	case messages.HistogramMetricType:
		fg = h.Theme.Blue()
		label = "HIST"
	case messages.DistributionMetricType:
		fg = h.Theme.Sky()
		label = "DISTR"
	case messages.GaugeMetricType:
		fg = h.Theme.Teal()
		label = "GAUGE"
	case messages.TimerMetricType:
		fg = h.Theme.Mauve()
		label = "TIMER"
	case messages.SetMetricType:
		fg = h.Theme.Pink()
		label = "SET"
	default:
		fg = h.Theme.Text()
		label = strings.ToUpper(metric.MetricType.String())
	}
	style := lipgloss.NewStyle().
		Width(11).
		Underline(true).
		Foreground(fg)
	return style.Render(label)
}

func (h *Handler) StyledMetricName(metric messages.DogStatsDMetric, width int) string {
	if width < 50 {
		width = 50
	}
	namespace := metric.Namespace
	name := metric.Name
	lenNamespace := len(namespace)
	lenName := len(name)
	textLen := lenNamespace + lenName

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
