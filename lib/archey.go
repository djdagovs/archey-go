package archey

import (
	"bytes"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/alexdreptu/sysinfo"
	utils "github.com/alexdreptu/utils-go"
	"github.com/mgutz/ansi"
)

type info map[string]string

type Show struct {
	OS        bool
	Arch      bool
	Kernel    bool
	User      bool
	Hostname  bool
	Uptime    bool
	UpSince   bool
	WM        bool
	DE        bool
	GTK2Theme bool
	GTK2Icons bool
	GTK2Font  bool
	GTK3Theme bool
	GTK3Icons bool
	GTK3Font  bool
	Terminal  bool
	Shell     bool
	Editor    bool
	Packages  bool
	Memory    bool
	Swap      bool
	CPU       bool
	Root      bool
	Home      bool
}

type Colors struct {
	Name string
	Text string
	Sep  string
	Body string
}

type Options struct {
	Sep        string
	DiskUnit   string
	MemoryUnit string
	SwapUnit   string
	Paths      []string
	PathFull   bool
	ShellFull  bool
	Show       Show
	Colors     Colors
}

// pacman's local database of installed packages
const pacmanDir = "/var/lib/pacman/local"

// default options
const (
	defSep        = ":"         // default separator
	defDiskUnit   = "gb"        // default unit for disk usage
	defMemoryUnit = defDiskUnit // default unit for memory usage
	defSwapUnit   = defMemoryUnit
)

// Name Sep Info
const infoFormat = "%s%s %s" // eg. OS: Linux

// default info colors
const (
	defNameColor  = "111"     // default color of the variable name
	defTextColor  = "white+h" // default color of the text
	defSepColor   = "white"   // default color of the separator
	defBodyColor1 = "111"     // default color of upper body of the logo
	defBodyColor2 = "69"      // default color of lower body of the logo
	resetColor    = "reset"   // reset color
)

const archLogo = `
                  {{.bCol1}}##{{.reset}}                      {{.info0}}
                 {{.bCol1}}####{{.reset}}                     {{.info1}}
                {{.bCol1}}######{{.reset}}                    {{.info2}}
               {{.bCol1}}########{{.reset}}                   {{.info3}}
              {{.bCol1}}##########{{.reset}}                  {{.info4}}
             {{.bCol1}}############{{.reset}}                 {{.info5}}
            {{.bCol1}}##############{{.reset}}                {{.info6}}
           {{.bCol1}}################{{.reset}}               {{.info7}}
          {{.bCol1}}##################{{.reset}}              {{.info8}}
         {{.bCol1}}#########{{.bCol2}}########{{.bCol1}}###{{.reset}}             {{.info9}}
        {{.bCol1}}###{{.bCol2}}#################{{.bCol1}}##{{.reset}}            {{.info10}}
       {{.bCol1}}##{{.bCol2}}#######{{.reset}}      {{.bCol2}}#########{{.reset}}           {{.info11}}
      {{.bCol2}}########;{{.reset}}        {{.bCol2}};########{{.reset}}          {{.info12}}
     {{.bCol2}}########;{{.reset}}          {{.bCol2}};########{{.reset}}         {{.info13}}
    {{.bCol2}}##########.{{.reset}}        {{.bCol2}}.##########{{.reset}}        {{.info14}}
   {{.bCol2}}#######{{.reset}}                  {{.bCol2}}#######{{.reset}}       {{.info15}}
  {{.bCol2}}#####{{.reset}}                        {{.bCol2}}#####{{.reset}}      {{.info16}}
 {{.bCol2}}###{{.reset}}                              {{.bCol2}}###{{.reset}}     {{.info17}}
{{.bCol2}}##{{.reset}}                                  {{.bCol2}}##{{.reset}}    {{.info18}}`

var (
	ErrInvalidMemUnit = func(u string) error {
		return fmt.Errorf("invalid memory unit '%s'", u)
	}
	ErrInvalidSwapUnit = func(u string) error {
		return fmt.Errorf("invalid swap unit '%s'", u)
	}
	ErrInvalidDiskUnit = func(u string) error {
		return fmt.Errorf("invalid disk unit '%s'", u)
	}
)

var gtkrc2 = filepath.Join(os.Getenv("HOME"), ".gtkrc-2.0")
var gtkrc3 = filepath.Join(os.Getenv("HOME"), ".config/gtk-3.0/settings.ini")

func (o *Options) Render() (string, error) {
	info, err := getFormattedInfo(o)
	if err != nil {
		return "", err
	}

	var bCol1 string
	var bCol2 string
	if o.Colors.Body != "" {
		bCol1 = ansi.ColorCode(o.Colors.Body)
		bCol2 = bCol1
	} else {
		bCol1 = ansi.ColorCode(defBodyColor1)
		bCol2 = ansi.ColorCode(defBodyColor2)
	}

	data := map[string]string{
		"bCol1": bCol1,
		"bCol2": bCol2,
		"reset": ansi.ColorCode(resetColor),
	}

	logoSize := len(strings.Split(archLogo, "\n")) - 1

	// create the info spots for logo
	// and set each to empty string
	for i := 0; i < logoSize; i++ {
		key := "info" + strconv.Itoa(i)
		data[key] = ""
	}

	// fill each info spot with formatted info
	for i, v := range info {
		key := "info" + strconv.Itoa(i)
		data[key] = v
	}

	var dataCount int
	for k, _ := range data {
		if strings.HasPrefix(k, "info") {
			dataCount++
		}
	}

	// create a slice and add the content of archLogo
	// to be able to extend it as necessary
	logo := []string{archLogo}

	if dataCount > logoSize {
		for i := logoSize; i < dataCount; i++ {
			infon := fmt.Sprintf("{{.info%s}}", strconv.Itoa(i))
			logo = append(logo, "\t\t\t\t\t  "+infon)
		}
	}

	// always append one empty line at the end of the info
	logo = append(logo, "")

	t, err := template.New("Arch Logo").Parse(strings.Join(logo, "\n"))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err = t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func getFormattedInfo(opt *Options) ([]string, error) {
	nameColor := ansi.ColorFunc(opt.Colors.Name)
	textColor := ansi.ColorFunc(opt.Colors.Text)
	sepColor := ansi.ColorFunc(opt.Colors.Sep)
	diskUnit := strings.ToLower(opt.DiskUnit)

	// hold the info format lines
	info := []string{}

	node := sysinfo.Node{}
	if err := node.Get(); err != nil {
		return nil, err
	}

	if !opt.Show.OS {
		osName := node.OSName + " " + node.Machine
		if opt.Show.Arch {
			osName = node.OSName
		}

		distroFormat := fmt.Sprintf(infoFormat,
			nameColor("OS"), sepColor(opt.Sep), textColor(osName))
		info = append(info, distroFormat)
	}

	if !opt.Show.Kernel {
		kernelFormat := fmt.Sprintf(infoFormat,
			nameColor("Kernel"), sepColor(opt.Sep), textColor(node.Release))
		info = append(info, kernelFormat)
	}

	if !opt.Show.User {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}
		userFormat := fmt.Sprintf(infoFormat,
			nameColor("User"), sepColor(opt.Sep), textColor(usr.Username))
		info = append(info, userFormat)
	}

	if !opt.Show.Hostname {
		hostnameFormat := fmt.Sprintf(infoFormat,
			nameColor("Hostname"), sepColor(opt.Sep), textColor(node.NodeName))
		info = append(info, hostnameFormat)
	}

	up := sysinfo.Uptime{}
	if err := up.Get(); err != nil {
		return nil, err
	}

	if !opt.Show.Uptime {
		uptimeFormat := fmt.Sprintf(infoFormat,
			nameColor("Uptime"), sepColor(opt.Sep), textColor(up.String()))
		info = append(info, uptimeFormat)
	}

	if !opt.Show.UpSince {
		upSinceFormat := fmt.Sprintf(infoFormat,
			nameColor("Up since"), sepColor(opt.Sep), textColor(up.UpSince()))
		info = append(info, upSinceFormat)
	}

	if !opt.Show.WM {
		wmFormat := fmt.Sprintf(infoFormat,
			nameColor("Window Manager"), sepColor(opt.Sep), textColor(GetWM()))
		info = append(info, wmFormat)
	}

	if !opt.Show.DE {
		deFormat := fmt.Sprintf(infoFormat,
			nameColor("Desktop Environment"), sepColor(opt.Sep), textColor(GetDE()))
		info = append(info, deFormat)
	}

	gtk2 := GetGTKInfo(gtkrc2)

	if !opt.Show.GTK2Theme {
		gtkThemeFormat := fmt.Sprintf(infoFormat,
			nameColor("GTK2 Theme"), sepColor(opt.Sep), textColor(gtk2.Theme))
		info = append(info, gtkThemeFormat)
	}

	if !opt.Show.GTK2Icons {
		gtkIconsFormat := fmt.Sprintf(infoFormat,
			nameColor("GTK2 Icon Theme"), sepColor(opt.Sep), textColor(gtk2.Icons))
		info = append(info, gtkIconsFormat)
	}

	if !opt.Show.GTK2Font {
		gtkFontFormat := fmt.Sprintf(infoFormat,
			nameColor("GTK2 Font"), sepColor(opt.Sep), textColor(gtk2.Font))
		info = append(info, gtkFontFormat)
	}

	gtk3 := GetGTKInfo(gtkrc3)

	if !opt.Show.GTK3Theme {
		gtkThemeFormat := fmt.Sprintf(infoFormat,
			nameColor("GTK3 Theme"), sepColor(opt.Sep), textColor(gtk3.Theme))
		info = append(info, gtkThemeFormat)
	}

	if !opt.Show.GTK3Icons {
		gtkIconsFormat := fmt.Sprintf(infoFormat,
			nameColor("GTK3 Icon Theme"), sepColor(opt.Sep), textColor(gtk3.Icons))
		info = append(info, gtkIconsFormat)
	}

	if !opt.Show.GTK3Font {
		gtkFontFormat := fmt.Sprintf(infoFormat,
			nameColor("GTK3 Font"), sepColor(opt.Sep), textColor(gtk3.Font))
		info = append(info, gtkFontFormat)
	}

	if !opt.Show.Terminal {
		terminalFormat := fmt.Sprintf(infoFormat,
			nameColor("Terminal"), sepColor(opt.Sep), textColor(os.Getenv("TERM")))
		info = append(info, terminalFormat)
	}

	if !opt.Show.Shell {
		var shell string

		if opt.ShellFull {
			shell = os.Getenv("SHELL")
		} else {
			shell = strings.Title(filepath.Base(os.Getenv("SHELL")))
		}

		shellFormat := fmt.Sprintf(infoFormat,
			nameColor("Shell"), sepColor(opt.Sep), textColor(shell))
		info = append(info, shellFormat)
	}

	if !opt.Show.Editor {
		editorFormat := fmt.Sprintf(infoFormat,
			nameColor("Editor"), sepColor(opt.Sep), textColor(os.Getenv("EDITOR")))
		info = append(info, editorFormat)
	}

	if !opt.Show.Packages {
		count, err := utils.CountDir(pacmanDir)
		if err != nil {
			return nil, err
		}

		n := strconv.Itoa(count.Dirs)
		packagesFormat := fmt.Sprintf(infoFormat,
			nameColor("Packages"), sepColor(opt.Sep), textColor(n))
		info = append(info, packagesFormat)
	}

	mem := sysinfo.Mem{}
	if err := mem.Get(); err != nil {
		return nil, err
	}

	if !opt.Show.Memory {
		var memUsage string
		switch strings.ToLower(opt.MemoryUnit) {
		case "mb":
			memUsage = fmt.Sprintf("%.1f MB / %.1f MB",
				mem.UsedMemInMB(), mem.TotalMemInMB())
		case "gb":
			memUsage = fmt.Sprintf("%.1f GB / %.1f GB",
				mem.UsedMemInGB(), mem.TotalMemInGB())
		default:
			return nil, ErrInvalidMemUnit(opt.MemoryUnit)
		}

		memoryFormat := fmt.Sprintf(infoFormat,
			nameColor("Memory"), sepColor(opt.Sep), textColor(memUsage))
		info = append(info, memoryFormat)
	}

	if !opt.Show.Swap {
		var swapUsage string
		switch strings.ToLower(opt.SwapUnit) {
		case "mb":
			swapUsage = fmt.Sprintf("%.1f MB / %.1f MB",
				mem.UsedSwapInMB(), mem.TotalSwapInMB())
		case "gb":
			swapUsage = fmt.Sprintf("%.1f GB / %.1f GB",
				mem.UsedSwapInGB(), mem.TotalSwapInGB())
		default:
			return nil, ErrInvalidSwapUnit(opt.SwapUnit)
		}

		swapFormat := fmt.Sprintf(infoFormat,
			nameColor("Swap"), sepColor(opt.Sep), textColor(swapUsage))
		info = append(info, swapFormat)
	}

	if !opt.Show.CPU {
		cpu := sysinfo.CPU{}
		if err := cpu.Get(); err != nil {
			return nil, err
		}

		cpuFormat := fmt.Sprintf(infoFormat,
			nameColor("CPU"), sepColor(opt.Sep), textColor(cpu.Name))
		info = append(info, cpuFormat)
	}

	if !opt.Show.Root {
		rootfs := sysinfo.FS{}
		if err := rootfs.Get("/"); err != nil {
			return nil, err
		}

		var rootfsUsage string
		switch diskUnit {
		case "mb":
			rootfsUsage = fmt.Sprintf("%.1f MB / %.1f MB",
				rootfs.UsedSpaceInMB(), rootfs.TotalSizeInMB())
		case "gb":
			rootfsUsage = fmt.Sprintf("%.1f GB / %.1f GB",
				rootfs.UsedSpaceInGB(), rootfs.TotalSizeInGB())
		default:
			return nil, ErrInvalidDiskUnit(opt.DiskUnit)
		}

		pathName := "Root"
		if opt.PathFull {
			pathName = "/root"
		}

		rootFormat := fmt.Sprintf(infoFormat,
			nameColor(pathName), sepColor(opt.Sep), textColor(rootfsUsage))
		info = append(info, rootFormat)
	}

	if !opt.Show.Home {
		homefs := sysinfo.FS{}
		if err := homefs.Get("/home"); err != nil {
			return nil, err
		}

		var homefsUsage string
		switch diskUnit {
		case "mb":
			homefsUsage = fmt.Sprintf("%.1f MB / %.1f MB",
				homefs.UsedSpaceInMB(), homefs.TotalSizeInMB())
		case "gb":
			homefsUsage = fmt.Sprintf("%.1f GB / %.1f GB",
				homefs.UsedSpaceInGB(), homefs.TotalSizeInGB())
		default:
			return nil, ErrInvalidDiskUnit(opt.DiskUnit)
		}

		pathName := "Home"
		if opt.PathFull {
			pathName = "/home"
		}

		homeFormat := fmt.Sprintf(infoFormat,
			nameColor(pathName), sepColor(opt.Sep), textColor(homefsUsage))
		info = append(info, homeFormat)
	}

	if len(opt.Paths) != 0 {
		// NOTE: fix to viper's slice bind handling problem - alexdreptu (10 Nov 2016)
		paths := func() []string {
			var sl []string
			// if theres more than one path string in the slice
			if len(opt.Paths) > 1 {
				// iterate over each path
				for _, path := range opt.Paths {
					// append it
					sl = append(sl, path)
				}
			} else {
				// otherwhise split the first and only path string
				sl = strings.Split(opt.Paths[0], ",")
			}
			return sl
		}()

		for _, path := range paths {
			pathfs := sysinfo.FS{}
			if err := pathfs.Get(path); err != nil {
				return nil, err
			}

			var pathfsUsage string
			switch diskUnit {
			case "mb":
				pathfsUsage = fmt.Sprintf("%.1f MB / %.1f MB",
					pathfs.UsedSpaceInMB(), pathfs.TotalSizeInMB())
			case "gb":
				pathfsUsage = fmt.Sprintf("%.1f GB / %.1f GB",
					pathfs.UsedSpaceInGB(), pathfs.TotalSizeInGB())
			default:
				return nil, ErrInvalidDiskUnit(opt.DiskUnit)
			}

			var pathFormat string
			if !opt.PathFull {
				path = strings.Title(strings.ToLower(filepath.Base(path)))
			}

			pathFormat = fmt.Sprintf(infoFormat,
				nameColor(path), sepColor(opt.Sep), textColor(pathfsUsage))
			info = append(info, pathFormat)
		}
	}

	return info, nil
}

func New() *Options {
	return &Options{
		Sep:        defSep,
		MemoryUnit: defMemoryUnit,
		SwapUnit:   defSwapUnit,
		DiskUnit:   defDiskUnit,
		PathFull:   false,
		ShellFull:  false,
		Show: Show{
			OS:        true,
			Arch:      true,
			Kernel:    true,
			User:      true,
			Hostname:  true,
			Uptime:    true,
			UpSince:   true,
			WM:        true,
			DE:        true,
			GTK2Theme: true,
			GTK2Icons: true,
			GTK2Font:  true,
			GTK3Theme: true,
			GTK3Icons: true,
			GTK3Font:  true,
			Terminal:  true,
			Shell:     true,
			Editor:    true,
			Packages:  true,
			Memory:    true,
			CPU:       true,
			Root:      true,
			Home:      true,
		},
		Colors: Colors{
			Name: defNameColor,
			Sep:  defSepColor,
			Text: defTextColor,
		},
	}
}
