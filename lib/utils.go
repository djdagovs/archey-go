package archey

import (
	"bufio"
	"os"
	"strings"

	utils "github.com/alexdreptu/utils-go"
	"github.com/mgutz/ansi"
)

type GTK struct {
	Theme  string
	Icons  string
	Font   string
	Cursor string
}

func GetWM() string {
	wm := "None"
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

func GetDE() string {
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

func GetGTKInfo(f string) GTK {
	var gtk GTK

	file, err := os.Open(f)
	if err != nil {
		return gtk
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if !strings.HasPrefix(line, "gtk") {
			continue
		}

		fields := strings.Split(line, "=")
		key := fields[0]
		value := strings.Trim(fields[1], "\"")

		switch key {
		case "gtk-theme-name":
			gtk.Theme = value
		case "gtk-icon-theme-name":
			gtk.Icons = value
		case "gtk-font-name":
			gtk.Font = value
		case "gtk-cursor-theme-name":
			gtk.Cursor = value
		}
	}

	if err := scanner.Err(); err != nil {
		return gtk
	}

	return gtk
}

func ListColors() {
	ansi.PrintStyles()
}

func NoColor() {
	ansi.DisableColors(true)
}
