// MIT License
//
// Copyright (c) 2016 Alexandru Dreptu
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package archey

import (
	"bufio"
	"fmt"
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

var ErrFileEmpty = func(f string) error {
	return fmt.Errorf("file '%s' is empty", f)
}

// GetWM returns Window Manager name
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
		"mutter":        "Mutter",
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

// GetDE returns the Desktop Environment name
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

// GetGTKInfo reads gtkrc and returns a GTK type
// containing theme name, icon theme name, font name and cursor theme name
func GetGTKInfo(f string) (GTK, error) {
	var gtk GTK

	file, err := os.Open(f)
	if err != nil {
		return gtk, err
	}

	fstat, err := file.Stat()
	if err != nil {
		return gtk, err
	}

	if fstat.Size() == 0 {
		return gtk, ErrFileEmpty(file.Name())
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
		return gtk, err
	}

	return gtk, nil
}

func ListColors() {
	ansi.PrintStyles()
}

func NoColor() {
	ansi.DisableColors(true)
}
