package archey

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"text/template"

	"github.com/alectic/sysinfo"
	"github.com/mgutz/ansi"
)

type info map[string]string

type Show struct {
	User     bool
	Hostname bool
	OS       bool
	Kernel   bool
	Uptime   bool
	WM       bool
	DE       bool
	Terminal bool
	Shell    bool
	Editor   bool
	Packages bool
	Memory   bool
	CPU      bool
	Root     bool
	Home     bool
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
	Paths      string
	Show       Show
	Colors     Colors
}

// default options
const (
	defSep        = ":"         // default separator
	defDiskUnit   = "gb"        // default unit for printing disk storage
	defMemoryUnit = defDiskUnit // default unit for printing memory usage
)

const maxPaths = 5 // max number of paths allowed

// Name Sep Info
const infoFormat = "%s%s %s" // eg. OS: Linux

// default info colors
const (
	defNameColor  = "yellow+h" // default color of the variable name
	defTextColor  = "white+h"  // default color of the text
	defSepColor   = "white"    // default color of the separator
	defBodyColor1 = "cyan+h"   // default color of upper body of the logo
	defBodyColor2 = "cyan"     // default color of lower body of the logo
	resetColor    = "reset"    // reset color
)

const archLogo = `
                  {{.bCol1}}##{{.reset}}                    {{.info0}}
                 {{.bCol1}}####{{.reset}}                   {{.info1}}
                {{.bCol1}}######{{.reset}}                  {{.info2}}
               {{.bCol1}}########{{.reset}}                 {{.info3}}
              {{.bCol1}}##########{{.reset}}                {{.info4}}
             {{.bCol1}}############{{.reset}}               {{.info5}}
            {{.bCol1}}##############{{.reset}}              {{.info6}}
           {{.bCol1}}################{{.reset}}             {{.info7}}
          {{.bCol1}}##################{{.reset}}            {{.info8}}
         {{.bCol1}}###########{{.bCol2}}######{{.bCol1}}###{{.reset}}           {{.info9}}
        {{.bCol1}}###{{.bCol2}}#################{{.bCol1}}##{{.reset}}          {{.info10}}
       {{.bCol1}}##{{.bCol2}}#######{{.reset}}      {{.bCol2}}#########{{.reset}}         {{.info11}}
      {{.bCol2}}########;{{.reset}}        {{.bCol2}};########{{.reset}}        {{.info12}}
     {{.bCol2}}########;{{.reset}}          {{.bCol2}};########{{.reset}}       {{.info13}}
    {{.bCol2}}##########.{{.reset}}        {{.bCol2}}.##########{{.reset}}      {{.info14}}
   {{.bCol2}}#######{{.reset}}                  {{.bCol2}}#######{{.reset}}     {{.info15}}
  {{.bCol2}}#####{{.reset}}                        {{.bCol2}}#####{{.reset}}    {{.info16}}
 {{.bCol2}}###{{.reset}}                              {{.bCol2}}###{{.reset}}   {{.info17}}
{{.bCol2}}##{{.reset}}                                  {{.bCol2}}##{{.reset}}  {{.info18}}
`

var (
	ErrInvalidMemUnit = func(u string) error {
		return fmt.Errorf("invalid memory unit '%s'", u)
	}
	ErrInvalidDiskUnit = func(u string) error {
		return fmt.Errorf("invalid disk unit '%s'", u)
	}
	ErrExcessivePaths = errors.New("excessive number of paths")
)

func (o *Options) Render() (string, error) {
	data, err := getData(o)
	if err != nil {
		return "", err
	}

	t, err := template.New("Arch Logo").Parse(archLogo)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err = t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func getData(opt *Options) (map[string]string, error) {
	info, err := getFormattedInfo(opt)
	if err != nil {
		return nil, err
	}

	var bCol1 string
	var bCol2 string
	if opt.Colors.Body != "" {
		bCol1 = ansi.ColorCode(opt.Colors.Body)
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

	logoLength := len(strings.Split(archLogo, "\n")) - 1
	for i := 0; i < logoLength; i++ {
		key := "info" + strconv.Itoa(i)
		data[key] = ""
	}

	for i, v := range info {
		key := "info" + strconv.Itoa(i)
		data[key] = v
	}

	return data, nil
}

func getFormattedInfo(opt *Options) ([]string, error) {
	nameColor := ansi.ColorFunc(opt.Colors.Name)
	textColor := ansi.ColorFunc(opt.Colors.Text)
	sepColor := ansi.ColorFunc(opt.Colors.Sep)

	// hold the info format lines
	info := []string{}

	node := sysinfo.NewNode()
	if err := node.Get(); err != nil {
		return nil, err
	}

	if !opt.Show.OS {
		distroFormat := fmt.Sprintf(infoFormat,
			nameColor("OS"), sepColor(opt.Sep), textColor(node.OSName))
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

	// TODO: implement - alectic (27 Oct 2016)
	if !opt.Show.Uptime {
		uptimeFormat := fmt.Sprintf(infoFormat,
			nameColor("Uptime"), sepColor(opt.Sep), textColor("to be implemented"))
		info = append(info, uptimeFormat)
	}

	if !opt.Show.WM {
		wmFormat := fmt.Sprintf(infoFormat,
			nameColor("Window Manager"), sepColor(opt.Sep), textColor(getWM()))
		info = append(info, wmFormat)
	}

	if !opt.Show.DE {
		deFormat := fmt.Sprintf(infoFormat,
			nameColor("Desktop Environment"), sepColor(opt.Sep), textColor(getDE()))
		info = append(info, deFormat)
	}

	if !opt.Show.Terminal {
		terminalFormat := fmt.Sprintf(infoFormat,
			nameColor("Terminal"), sepColor(opt.Sep), textColor(os.Getenv("TERM")))
		info = append(info, terminalFormat)
	}

	if !opt.Show.Shell {
		shellFormat := fmt.Sprintf(infoFormat,
			nameColor("Shell"), sepColor(opt.Sep), textColor(os.Getenv("SHELL")))
		info = append(info, shellFormat)
	}

	if !opt.Show.Editor {
		editorFormat := fmt.Sprintf(infoFormat,
			nameColor("Editor"), sepColor(opt.Sep), textColor(os.Getenv("EDITOR")))
		info = append(info, editorFormat)
	}

	if !opt.Show.Packages {
		n, err := countPkgs()
		if err != nil {
			return nil, err
		}

		packagesFormat := fmt.Sprintf(infoFormat,
			nameColor("Packages"), sepColor(opt.Sep), textColor(strconv.Itoa(n)))
		info = append(info, packagesFormat)
	}

	if !opt.Show.Memory {
		mem := sysinfo.NewMem()
		if err := mem.Get(); err != nil {
			return nil, err
		}

		var memUsage string
		switch opt.MemoryUnit {
		case "mb":
			memUsage = fmt.Sprintf("%.1f MB / %.1f MB",
				mem.UsedMemInMB(), mem.TotalMemInMB())
		case "gb":
			memUsage = fmt.Sprintf("%.1f GB / %.1f GB",
				mem.UsedMemInGB(), mem.TotalMemInGB())
		}

		memoryFormat := fmt.Sprintf(infoFormat,
			nameColor("Memory"), sepColor(opt.Sep), textColor(memUsage))
		info = append(info, memoryFormat)
	}

	if !opt.Show.CPU {
		cpu := sysinfo.NewCPU()
		if err := cpu.Get(); err != nil {
			return nil, err
		}

		cpuFormat := fmt.Sprintf(infoFormat,
			nameColor("CPU"), sepColor(opt.Sep), textColor(cpu.Name))
		info = append(info, cpuFormat)
	}

	if !opt.Show.Root {
		rootfs := sysinfo.NewFS()
		if err := rootfs.Get("/"); err != nil {
			return nil, err
		}
		var rootfsUsage string
		switch opt.DiskUnit {
		case "mb":
			rootfsUsage = fmt.Sprintf("%.1f MB / %.1f MB",
				rootfs.UsedSpaceInMB(), rootfs.TotalSizeInMB())
		case "gb":
			rootfsUsage = fmt.Sprintf("%.1f GB / %.1f GB",
				rootfs.UsedSpaceInGB(), rootfs.TotalSizeInGB())
		default:
			return nil, ErrInvalidDiskUnit(opt.DiskUnit)
		}
		rootFormat := fmt.Sprintf(infoFormat,
			nameColor("Root"), sepColor(opt.Sep), textColor(rootfsUsage))
		info = append(info, rootFormat)
	}

	if !opt.Show.Home {
		homefs := sysinfo.NewFS()
		if err := homefs.Get("/home"); err != nil {
			return nil, err
		}

		var homefsUsage string
		switch opt.DiskUnit {
		case "mb":
			homefsUsage = fmt.Sprintf("%.1f MB / %.1f MB",
				homefs.UsedSpaceInMB(), homefs.TotalSizeInMB())
		case "gb":
			homefsUsage = fmt.Sprintf("%.1f GB / %.1f GB",
				homefs.UsedSpaceInGB(), homefs.TotalSizeInGB())
		default:
			return nil, ErrInvalidDiskUnit(opt.DiskUnit)
		}

		homeFormat := fmt.Sprintf(infoFormat,
			nameColor("Home"), sepColor(opt.Sep), textColor(homefsUsage))
		info = append(info, homeFormat)
	}

	if len(opt.Paths) != 0 {
		paths := strings.Split(opt.Paths, ",")
		// if passed paths exceeds maxPaths
		if len(paths) > maxPaths {
			return nil, ErrExcessivePaths
		}

		for _, path := range paths {
			pathfs := sysinfo.NewFS()
			if err := pathfs.Get(path); err != nil {
				return nil, err
			}
			var pathfsUsage string
			switch opt.DiskUnit {
			case "mb":
				pathfsUsage = fmt.Sprintf("%.1f MB / %.1f MB",
					pathfs.UsedSpaceInMB(), pathfs.TotalSizeInMB())
			case "gb":
				pathfsUsage = fmt.Sprintf("%.1f GB / %.1f GB", pathfs.UsedSpaceInGB(), pathfs.TotalSizeInGB())
			default:
				return nil, ErrInvalidDiskUnit(opt.DiskUnit)
			}
			pathFormat := fmt.Sprintf(infoFormat,
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
		DiskUnit:   defDiskUnit,
		Show: Show{
			User:     true,
			Hostname: true,
			OS:       true,
			Kernel:   true,
			Uptime:   true,
			WM:       true,
			DE:       true,
			Terminal: true,
			Shell:    true,
			Editor:   true,
			Packages: true,
			Memory:   true,
			CPU:      true,
			Root:     true,
			Home:     true,
		},
		Colors: Colors{
			Name: defNameColor,
			Sep:  defSepColor,
			Text: defTextColor,
		},
	}
}
