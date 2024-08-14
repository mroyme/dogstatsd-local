package handler

import (
	"fmt"
	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
	"github.com/mroyme/dogstatsd-local/internal/messages"
	"log"
	"strings"
)

var theme = catppuccin.Mocha

func Color(catppuccinColor catppuccin.Color) lipgloss.Color {
	return lipgloss.Color(catppuccinColor.Hex)
}

func NewPrettyMessageHandler(extraTags []string, nameWidth int, valueWidth int) MessageHandler {
	return func(msg []byte) error {
		dMsg, err := messages.ParseDogStatsDMessage(msg)
		if err != nil {
			log.Println(err.Error())
			return nil
		}
		metric, ok := dMsg.(messages.DogStatsDMetric)
		if dMsg.Type() != messages.MetricMessageType || !ok {
			return nil
		}

		str := styledMetricType(metric)
		str += styledMetricName(metric, nameWidth)
		str += styledMetricValue(metric, valueWidth)
		str += styledTags(metric, extraTags)
		fmt.Println(str)
		return nil
	}
}

func styledMetricType(metric messages.DogStatsDMetric) string {
	var fg catppuccin.Color
	metricType := metric.MetricType
	switch metricType {
	case messages.CounterMetricType:
		fg = theme.Green()
	case messages.HistogramMetricType:
		fg = theme.Blue()
	case messages.GaugeMetricType:
		fg = theme.Teal()
	case messages.TimerMetricType:
		fg = theme.Mauve()
	case messages.SetMetricType:
		fg = theme.Pink()
	default:
		fg = theme.Text()
	}
	style := lipgloss.NewStyle().
		Width(11).
		Background(Color(theme.Base())).
		Foreground(Color(fg))
	return style.Render(strings.ToUpper(metricType.String()))
}

func styledMetricName(metric messages.DogStatsDMetric, width int) string {
	text := fmt.Sprintf("%s | %s", metric.Namespace, metric.Name)
	if len(text) > width {
		text = text[:width-3] + "..."
	}
	style := lipgloss.NewStyle().
		Width(50).
		MaxWidth(50).
		Background(Color(theme.Base())).
		Foreground(Color(theme.Lavender()))
	return style.Render(text)
}

func styledMetricValue(metric messages.DogStatsDMetric, width int) string {
	value := fmt.Sprintf("%.2f", metric.FloatValue)
	if len(value) > width {
		value = value[:width-3] + "..."
	}
	style := lipgloss.NewStyle().
		Width(15).
		MaxWidth(15).
		Background(Color(theme.Base())).
		Foreground(Color(theme.Sapphire()))
	if metric.MetricType == messages.TimerMetricType {
		value += "ms"
	}
	return style.Render(value)
}

func styledTags(metric messages.DogStatsDMetric, extraTags []string) string {
	style := lipgloss.NewStyle().
		Background(Color(theme.Base())).
		Foreground(Color(theme.Subtext0()))
	var tags []string
	for _, tag := range append(extraTags, metric.Tags...) {
		tags = append(tags, strings.TrimSpace(tag))
	}
	prefix := style.Foreground(Color(theme.Overlay0())).Render("TAGS =")
	return prefix + style.SetString(tags...).Render()
}
