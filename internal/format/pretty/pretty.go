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
		metric, ok := dMsg.(messages.DogStatsDMetric)
		if dMsg.Type() != messages.MetricMessageType || !ok {
			return nil
		}

		str := h.styledMetricType(metric)
		str += h.styledMetricName(metric, h.NameWidth)
		str += h.styledMetricValue(metric, h.ValueWidth)
		str += h.styledTags(metric, h.ExtraTags)
		fmt.Println(str)
		return nil
	}
}

func (h *Handler) styledMetricType(metric messages.DogStatsDMetric) string {
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
		Foreground(fg)
	return style.Render(strings.ToUpper(metricType.String()))
}

func (h *Handler) styledMetricName(metric messages.DogStatsDMetric, width int) string {
	text := fmt.Sprintf("%s | %s", metric.Namespace, metric.Name)
	if len(text) > width {
		text = text[:width-3] + "..."
	}
	style := lipgloss.NewStyle().
		Width(50).
		MaxWidth(50).
		Foreground(h.Theme.Lavender())
	return style.Render(text)
}

func (h *Handler) styledMetricValue(metric messages.DogStatsDMetric, width int) string {
	value := fmt.Sprintf("%.2f", metric.FloatValue)
	if len(value) > width {
		value = value[:width-3] + "..."
	}
	style := lipgloss.NewStyle().
		Width(15).
		MaxWidth(15).
		Foreground(h.Theme.Sapphire())
	if metric.MetricType == messages.TimerMetricType {
		value += "ms"
	}
	return style.Render(value)
}

func (h *Handler) styledTags(metric messages.DogStatsDMetric, extraTags []string) string {
	style := lipgloss.NewStyle().
		Foreground(h.Theme.Subtext0())
	var tags []string
	for _, tag := range append(extraTags, metric.Tags...) {
		tags = append(tags, strings.TrimSpace(tag))
	}
	prefix := style.Foreground(h.Theme.Overlay0()).Render("TAGS =")
	return prefix + style.SetString(tags...).Render()
}
