package pretty

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mroyme/dogstatsd-local/internal/messages"
)

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
	if sc.Hostname == "" {
		return ""
	}
	return lipgloss.NewStyle().Bold(true).Foreground(h.Theme.Sapphire()).Render(sc.Hostname)
}

func (h *Handler) StyledServiceCheckMessage(sc messages.DogStatsDServiceCheck) string {
	if sc.Message == "" {
		return ""
	}
	return lipgloss.NewStyle().Foreground(h.Theme.Subtext1()).Render(sc.Message)
}
