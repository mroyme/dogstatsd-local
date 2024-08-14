package pretty

import (
	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
)

type LipglossAdaptiveTheme interface {
	Rosewater() lipgloss.AdaptiveColor
	Flamingo() lipgloss.AdaptiveColor
	Pink() lipgloss.AdaptiveColor
	Mauve() lipgloss.AdaptiveColor
	Red() lipgloss.AdaptiveColor
	Maroon() lipgloss.AdaptiveColor
	Peach() lipgloss.AdaptiveColor
	Yellow() lipgloss.AdaptiveColor
	Green() lipgloss.AdaptiveColor
	Teal() lipgloss.AdaptiveColor
	Sky() lipgloss.AdaptiveColor
	Sapphire() lipgloss.AdaptiveColor
	Blue() lipgloss.AdaptiveColor
	Lavender() lipgloss.AdaptiveColor
	Text() lipgloss.AdaptiveColor
	Subtext1() lipgloss.AdaptiveColor
	Subtext0() lipgloss.AdaptiveColor
	Overlay2() lipgloss.AdaptiveColor
	Overlay1() lipgloss.AdaptiveColor
	Overlay0() lipgloss.AdaptiveColor
	Surface2() lipgloss.AdaptiveColor
	Surface1() lipgloss.AdaptiveColor
	Surface0() lipgloss.AdaptiveColor
	Crust() lipgloss.AdaptiveColor
	Mantle() lipgloss.AdaptiveColor
	Base() lipgloss.AdaptiveColor
	Name() string
}

type CatppuccinAdaptiveTheme struct {
	light catppuccin.Flavour
	dark  catppuccin.Flavour
}

func (t *CatppuccinAdaptiveTheme) Rosewater() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Rosewater().Hex,
		Dark:  t.dark.Rosewater().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Flamingo() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Flamingo().Hex,
		Dark:  t.dark.Flamingo().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Pink() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Pink().Hex,
		Dark:  t.dark.Pink().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Mauve() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Mauve().Hex,
		Dark:  t.dark.Mauve().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Red() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Red().Hex,
		Dark:  t.dark.Red().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Maroon() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Maroon().Hex,
		Dark:  t.dark.Maroon().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Peach() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Peach().Hex,
		Dark:  t.dark.Peach().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Yellow() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Yellow().Hex,
		Dark:  t.dark.Yellow().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Green() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Green().Hex,
		Dark:  t.dark.Green().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Teal() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Teal().Hex,
		Dark:  t.dark.Teal().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Sky() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Sky().Hex,
		Dark:  t.dark.Sky().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Sapphire() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Sapphire().Hex,
		Dark:  t.dark.Sapphire().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Blue() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Blue().Hex,
		Dark:  t.dark.Blue().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Lavender() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Lavender().Hex,
		Dark:  t.dark.Lavender().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Text() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Text().Hex,
		Dark:  t.dark.Text().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Subtext1() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Subtext1().Hex,
		Dark:  t.dark.Subtext1().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Subtext0() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Subtext0().Hex,
		Dark:  t.dark.Subtext0().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Overlay2() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Overlay2().Hex,
		Dark:  t.dark.Overlay2().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Overlay1() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Overlay1().Hex,
		Dark:  t.dark.Overlay1().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Overlay0() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Overlay0().Hex,
		Dark:  t.dark.Overlay0().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Surface2() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Surface2().Hex,
		Dark:  t.dark.Surface2().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Surface1() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Surface1().Hex,
		Dark:  t.dark.Surface1().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Surface0() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Surface0().Hex,
		Dark:  t.dark.Surface0().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Crust() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Crust().Hex,
		Dark:  t.dark.Crust().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Mantle() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Mantle().Hex,
		Dark:  t.dark.Mantle().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Base() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.light.Base().Hex,
		Dark:  t.dark.Base().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Name() string {
	if lipgloss.HasDarkBackground() {
		return t.dark.Name()
	}
	return t.light.Name()
}
