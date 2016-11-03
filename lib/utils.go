package archey

import (
	utils "github.com/alectic/utils-go"
	"github.com/mgutz/ansi"
)

func getWM() string {
	var wm = "None"
	wmList := map[string]string{
		"awesome":       "Awesome",
		"blackbox":      "Blackbox",
		"bspwm":         "bspwm",
		"dwm":           "DWM",
		"enlightenment": "Enlightenment",
		"fluxbox":       "Fluxbox",
		"fvwm":          "FVWM",
		"herbstluftwm":  "herbstluftwm",
		"i3":            "i3",
		"icewm":         "IceWM",
		"kwin":          "KWin",
		"metacity":      "Metacity",
		"musca":         "Musca",
		"openbox":       "Openbox",
		"pekwm":         "PekWM",
		"ratpoison":     "ratpoison",
		"scrotwm":       "ScrotWM",
		"subtle":        "subtle",
		"monsterwm":     "MonsterWM",
		"wmaker":        "Window Maker",
		"wmfs":          "Wmfs",
		"wmii":          "wmii",
		"xfwm4":         "Xfwm",
		"qtile":         "QTile",
		"wingo":         "Wingo",
	}

	for k, v := range wmList {
		if utils.IsExistProcName(k) {
			wm = v
			break
		}
	}

	return wm
}

func getDE() string {
	de := "None"
	deList := map[string]string{
		"cinnamon":      "Cinnamon",
		"gnome-session": "GNOME",
		"ksmserver":     "KDE",
		"mate-session":  "MATE",
		"xfce4-session": "Xfce",
		"lxsession":     "LXDE",
	}

	for k, v := range deList {
		if utils.IsExistProcName(k) {
			de = v
			break
		}
	}

	return de
}

func ListColors() {
	ansi.PrintStyles()
}
