package cmd

import (
	"fmt"
	"os"

	archey "github.com/alectic/archey-go/lib"
	"github.com/spf13/cobra"
)

var confFile string

// Options
var (
	sep        string
	memoryUnit string
	diskUnit   string
	paths      []string
	pathFull   bool
)

// Show
var (
	noUser     bool
	noHostname bool
	noOS       bool
	noArch     bool
	noKernel   bool
	noUptime   bool
	noUpSince  bool
	noWM       bool
	noDE       bool
	noTerminal bool
	noShell    bool
	noEditor   bool
	noPackages bool
	noMemory   bool
	noCPU      bool
	noRoot     bool
	noHome     bool
)

var (
	nameColor  string
	textColor  string
	sepColor   string
	bodyColor  string
	noColor    bool
	listColors bool
)

var RootCmd = &cobra.Command{
	Use:   "archey-go",
	Short: "Archey-go is a tool to display prettified system info on Arch Linux",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			cmd.Help()
			os.Exit(0)
		}

		opt := archey.New()

		opt.Show.OS = noOS
		opt.Show.Arch = noArch
		opt.Show.Kernel = noKernel
		opt.Show.User = noUser
		opt.Show.Hostname = noHostname
		opt.Show.Uptime = noUptime
		opt.Show.UpSince = noUpSince
		opt.Show.WM = noWM
		opt.Show.DE = noDE
		opt.Show.Terminal = noTerminal
		opt.Show.Shell = noShell
		opt.Show.Editor = noEditor
		opt.Show.Packages = noPackages
		opt.Show.Memory = noMemory
		opt.Show.CPU = noCPU
		opt.Show.Root = noRoot
		opt.Show.Home = noHome
		opt.PathFull = pathFull

		if sep != "" {
			opt.Sep = sep
		}

		if memoryUnit != "" {
			opt.MemoryUnit = memoryUnit
		}

		if diskUnit != "" {
			opt.DiskUnit = diskUnit
		}

		if len(paths) != 0 {
			opt.Paths = paths
		}

		if nameColor != "" {
			opt.Colors.Name = nameColor
		}

		if textColor != "" {
			opt.Colors.Text = textColor
		}

		if sepColor != "" {
			opt.Colors.Sep = sepColor
		}

		if bodyColor != "" {
			opt.Colors.Body = bodyColor
		}

		// TODO: implement - alectic (08 Nov 2016)
		if noColor {
			fmt.Println("to be implemented")
			os.Exit(1)
		}

		if listColors {
			archey.ListColors()
			os.Exit(1)
		}

		info, err := opt.Render()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println(info)
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.Flags().BoolVar(&noOS, "no-os", false, "don't print os name")
	RootCmd.Flags().BoolVar(&noArch, "no-arch", false, "don't print architecture")
	RootCmd.Flags().BoolVar(&noKernel, "no-kernel", false, "don't print kernel version")
	RootCmd.Flags().BoolVar(&noUser, "no-user", false, "don't print user")
	RootCmd.Flags().BoolVar(&noHostname, "no-hostname", false, "don't print hostname")
	RootCmd.Flags().BoolVar(&noUptime, "no-uptime", false, "don't print uptime")
	RootCmd.Flags().BoolVar(&noUpSince, "no-up-since", false, "don't print up since")
	RootCmd.Flags().BoolVar(&noWM, "no-wm", false, "don't print Window Manager name")
	RootCmd.Flags().BoolVar(&noDE, "no-de", false, "don't print Desktop Environment name")
	RootCmd.Flags().BoolVar(&noTerminal, "no-terminal", false, "don't print terminal name")
	RootCmd.Flags().BoolVar(&noShell, "no-shell", false, "don't print shell name")
	RootCmd.Flags().BoolVar(&noEditor, "no-editor", false, "don't print editor name")
	RootCmd.Flags().BoolVar(&noPackages, "no-packages", false, "don't print packages count")
	RootCmd.Flags().BoolVar(&noMemory, "no-memory", false, "don't print memory usage")
	RootCmd.Flags().BoolVar(&noCPU, "no-cpu", false, "don't print CPU model")
	RootCmd.Flags().BoolVar(&noRoot, "no-root", false, "don't print root disk usage")
	RootCmd.Flags().BoolVar(&noRoot, "no-home", false, "don't print home disk usage")
	RootCmd.Flags().StringVar(&sep, "sep", "", "separator string")
	RootCmd.Flags().StringVar(&memoryUnit, "memory-unit", "", "unit to use for memory usage")
	RootCmd.Flags().StringVar(&diskUnit, "disk-unit", "", "unit to use for disk usage")
	RootCmd.Flags().StringSliceVar(&paths, "paths", nil, "additional paths to add to disk usage info")
	RootCmd.Flags().BoolVar(&pathFull, "path-full", false, "show full paths")
	RootCmd.Flags().StringVar(&nameColor, "name-color", "", "color of the variable name")
	RootCmd.Flags().StringVar(&textColor, "text-color", "", "color of the text")
	RootCmd.Flags().StringVar(&sepColor, "sep-color", "", "color of the separator")
	RootCmd.Flags().StringVar(&bodyColor, "body-color", "", "color of the logo body")
	RootCmd.Flags().BoolVar(&noColor, "no-color", false, "don't use any colors")
	RootCmd.Flags().BoolVar(&listColors, "list-colors", false, "print all colors and styles")
}

func initConfig() {
	/*  */
}
