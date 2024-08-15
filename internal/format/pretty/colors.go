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
	Light catppuccin.Flavour
	Dark  catppuccin.Flavour
}

func (t *CatppuccinAdaptiveTheme) Rosewater() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Rosewater().Hex,
		Dark:  t.Dark.Rosewater().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Flamingo() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Flamingo().Hex,
		Dark:  t.Dark.Flamingo().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Pink() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Pink().Hex,
		Dark:  t.Dark.Pink().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Mauve() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Mauve().Hex,
		Dark:  t.Dark.Mauve().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Red() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Red().Hex,
		Dark:  t.Dark.Red().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Maroon() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Maroon().Hex,
		Dark:  t.Dark.Maroon().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Peach() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Peach().Hex,
		Dark:  t.Dark.Peach().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Yellow() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Yellow().Hex,
		Dark:  t.Dark.Yellow().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Green() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Green().Hex,
		Dark:  t.Dark.Green().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Teal() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Teal().Hex,
		Dark:  t.Dark.Teal().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Sky() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Sky().Hex,
		Dark:  t.Dark.Sky().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Sapphire() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Sapphire().Hex,
		Dark:  t.Dark.Sapphire().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Blue() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Blue().Hex,
		Dark:  t.Dark.Blue().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Lavender() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Lavender().Hex,
		Dark:  t.Dark.Lavender().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Text() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Text().Hex,
		Dark:  t.Dark.Text().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Subtext1() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Subtext1().Hex,
		Dark:  t.Dark.Subtext1().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Subtext0() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Subtext0().Hex,
		Dark:  t.Dark.Subtext0().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Overlay2() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Overlay2().Hex,
		Dark:  t.Dark.Overlay2().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Overlay1() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Overlay1().Hex,
		Dark:  t.Dark.Overlay1().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Overlay0() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Overlay0().Hex,
		Dark:  t.Dark.Overlay0().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Surface2() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Surface2().Hex,
		Dark:  t.Dark.Surface2().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Surface1() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Surface1().Hex,
		Dark:  t.Dark.Surface1().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Surface0() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Surface0().Hex,
		Dark:  t.Dark.Surface0().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Crust() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Crust().Hex,
		Dark:  t.Dark.Crust().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Mantle() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Mantle().Hex,
		Dark:  t.Dark.Mantle().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Base() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: t.Light.Base().Hex,
		Dark:  t.Dark.Base().Hex,
	}
}

func (t *CatppuccinAdaptiveTheme) Name() string {
	if lipgloss.HasDarkBackground() {
		return t.Dark.Name()
	}
	return t.Light.Name()
}
