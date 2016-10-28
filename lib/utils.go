package archey

import (
	"io/ioutil"

	"github.com/mgutz/ansi"
)

// pacman's local database of installed packages
const pacmanDir = "/var/lib/pacman/local"

// TODO: implement - alectic (29 Oct 2016)
func getWM() string {
	//wms := map[string]string{
	//"awesome":       "Awesome",
	//"blackbox":      "Blackbox",
	//"bspwm":         "bspwm",
	//"dwm":           "DWM",
	//"enlightenment": "Enlightenment",
	//"fluxbox":       "Fluxbox",
	//"fvwm":          "FVWM",
	//"herbstluftwm":  "herbstluftwm",
	//"i3":            "i3",
	//"icewm":         "IceWM",
	//"kwin":          "KWin",
	//"metacity":      "Metacity",
	//"musca":         "Musca",
	//"openbox":       "Openbox",
	//"pekwm":         "PekWM",
	//"ratpoison":     "ratpoison",
	//"scrotwm":       "ScrotWM",
	//"subtle":        "subtle",
	//"monsterwm":     "MonsterWM",
	//"wmaker":        "Window Maker",
	//"wmfs":          "Wmfs",
	//"wmii":          "wmii",
	//"xfwm4":         "Xfwm",
	//"qtile":         "QTile",
	//"wingo":         "Wingo",
	//}

	return "to be implemented"
}

// TODO: implement - alectic (29 Oct 2016)
func getDE() string {
	return "to be implemented"
}

func countPackages() (int, error) {
	var count int
	files, err := ioutil.ReadDir(pacmanDir)
	if err != nil {
		return 0, err
	}

	for _, file := range files {
		// count only directories
		if file.IsDir() {
			count++
		}
	}
	return count, nil
}

func ListColors() {
	ansi.PrintStyles()
}
