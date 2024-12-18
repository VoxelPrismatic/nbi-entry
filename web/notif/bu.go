package notif

import (
	"fmt"
	"nbientry/web"
	"strings"

	"github.com/a-h/templ"
	templruntime "github.com/a-h/templ/runtime"
)

var _ = web.Migrate(ApplicationSegment{})

type ApplicationSegment struct {
	Id    int `gorm:"primaryKey"`
	Name  string
	Color string
}

func AllApplicationSegments() []ApplicationSegment {
	return web.GetSorted(ApplicationSegment{}, "name ASC")
}

func (app ApplicationSegment) color() templ.CSSClass {
	return templ_color(app.Color)
}

func templ_color(color string) templ.CSSClass {
	css := templruntime.GetBuilder()
	if strings.HasPrefix(color, "#") {
		if len(color) == 4 || len(color) == 7 {
			css.WriteString(string(templ.SanitizeCSS(`--color`, templ.SafeCSSProperty(color))))
		}
	} else {
		if len(color) == 4 {
			css.WriteString(string(templ.SanitizeCSS(`--color`, templ.SafeCSSProperty(fmt.Sprintf("var(--sakura-paint-%s)", color)))))
		}
	}

	css_id := templ.CSSID(`color`, css.String())
	return templ.ComponentCSSClass{
		ID:    css_id,
		Class: templ.SafeCSS(`.` + css_id + `{` + css.String() + `}`),
	}
}
