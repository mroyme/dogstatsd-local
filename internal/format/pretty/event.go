package pretty

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mroyme/dogstatsd-local/internal/messages"
)

func (h *Handler) StyledEventAlertType(ev messages.DogStatsDEvent) string {
	var fg lipgloss.AdaptiveColor
	var label string
	switch ev.AlertType {
	case messages.EventAlertTypeError:
		fg = h.Theme.Red()
		label = "ERR"
	case messages.EventAlertTypeWarning:
		fg = h.Theme.Yellow()
		label = "WARN"
	case messages.EventAlertTypeSuccess:
		fg = h.Theme.Green()
		label = "OK"
	default:
		fg = h.Theme.Blue()
		label = "INFO"
	}
	style := lipgloss.NewStyle().
		Width(11).
		Underline(true).
		Foreground(fg)
	return style.Render(label)
}

func (h *Handler) StyledEventTitle(ev messages.DogStatsDEvent, width int) string {
	if width < 50 {
		width = 50
	}
	text := lipgloss.NewStyle().Bold(true).Foreground(h.Theme.Pink()).Render(ev.Title)
	return lipgloss.NewStyle().
		Width(width).
		MaxWidth(width).
		Render(text)
}

func (h *Handler) StyledEventText(ev messages.DogStatsDEvent) string {
	if ev.Text == "" {
		return ""
	}
	return lipgloss.NewStyle().Foreground(h.Theme.Subtext1()).Render(ev.Text)
}
