package ui

import "strings"

var ColorMap = map[string]Attribute{
	"fg":           ColorWhite,
	"bg":           ColorDefault,
	"border.fg":    ColorWhite,
	"label.fg":     ColorGreen,
	"par.fg":       ColorYellow,
	"par.label.bg": ColorWhite,
}

func ThemeAttr(name string) Attribute {
	return lookUpAttr(ColorMap, name)
}

func lookUpAttr(clrmap map[string]Attribute, name string) Attribute {

	a, ok := clrmap[name]
	if ok {
		return a
	}

	ns := strings.Split(name, ".")
	for i := range ns {
		nn := strings.Join(ns[i:len(ns)], ".")
		a, ok = ColorMap[nn]
		if ok {
			break
		}
	}

	return a
}