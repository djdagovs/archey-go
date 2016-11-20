package cmd

import (
	"fmt"
	"os"

	archey "github.com/alexdreptu/archey-go/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	noOS            bool
	noArch          bool
	noKernel        bool
	noUser          bool
	noHostname      bool
	noUptime        bool
	noUpSince       bool
	noWM            bool
	noDE            bool
	noGTK2Theme     bool
	noGTK2IconTheme bool
	noGTK2Font      bool
	noGTK3Theme     bool
	noGTK3IconTheme bool
	noGTK3Font      bool
	noTerminal      bool
	noShell         bool
	noEditor        bool
	noPackages      bool
	noMemory        bool
	noSwap          bool
	noCPU           bool
	noRoot          bool
	noHome          bool
)

var (
	nameColor string
	textColor string
	sepColor  string
	bodyColor []string
)

// options
var (
	sep           string
	diskUnit      string
	memoryUnit    string
	swapUnit      string
	paths         []string
	pathFull      bool
	shellFull     bool
	upSinceFormat string
	noColor       bool
)

var listColors bool
var config string

var RootCmd = &cobra.Command{
	Use:   "archey-go",
	Short: "Archey-go is a tool to display prettified system info on Arch Linux",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			cmd.Help()
			os.Exit(0)
		}

		opt := archey.New()

		opt.Show.OS = viper.GetBool("show.noOS")
		opt.Show.Arch = viper.GetBool("show.noArch")
		opt.Show.Kernel = viper.GetBool("show.noKernel")
		opt.Show.User = viper.GetBool("show.noUser")
		opt.Show.Hostname = viper.GetBool("show.noHostname")
		opt.Show.Uptime = viper.GetBool("show.noUptime")
		opt.Show.UpSince = viper.GetBool("show.noUpSince")
		opt.Show.WM = viper.GetBool("show.noWM")
		opt.Show.DE = viper.GetBool("show.noDE")
		opt.Show.GTK2Theme = viper.GetBool("show.noGTK2Theme")
		opt.Show.GTK2IconTheme = viper.GetBool("show.noGTK2IconTheme")
		opt.Show.GTK2Font = viper.GetBool("show.noGTK2Font")
		opt.Show.GTK3Theme = viper.GetBool("show.noGTK3Theme")
		opt.Show.GTK3IconTheme = viper.GetBool("show.noGTK3IconTheme")
		opt.Show.GTK3Font = viper.GetBool("show.noGTK3Font")
		opt.Show.Terminal = viper.GetBool("show.noTerminal")
		opt.Show.Shell = viper.GetBool("show.noShell")
		opt.Show.Editor = viper.GetBool("show.noEditor")
		opt.Show.Packages = viper.GetBool("show.noPackages")
		opt.Show.Memory = viper.GetBool("show.noMemory")
		opt.Show.Swap = viper.GetBool("show.noSwap")
		opt.Show.CPU = viper.GetBool("show.noCPU")
		opt.Show.Root = viper.GetBool("show.noRoot")
		opt.Show.Home = viper.GetBool("show.noHome")

		if viper.GetString("options.sep") != "" {
			opt.Sep = viper.GetString("options.sep")
		}

		if viper.GetString("options.memoryUnit") != "" {
			opt.MemoryUnit = viper.GetString("options.memoryUnit")
		}

		if viper.GetString("options.swapUnit") != "" {
			opt.SwapUnit = viper.GetString("options.swapUnit")
		}

		if viper.GetString("options.diskUnit") != "" {
			opt.DiskUnit = viper.GetString("options.diskUnit")
		}

		// NOTE: slice binds to pflag not handled correctly by viper - alexdreptu (10 Nov 2016)
		opt.Paths = viper.GetStringSlice("options.paths")
		opt.PathFull = viper.GetBool("options.pathFull")
		opt.ShellFull = viper.GetBool("options.shellFull")

		if viper.GetString("options.upSinceFormat") != "" {
			opt.UpSinceFormat = viper.GetString("options.upSinceFormat")
		}

		if viper.GetString("colors.nameColor") != "" {
			opt.Colors.Name = viper.GetString("colors.nameColor")
		}

		if viper.GetString("colors.textColor") != "" {
			opt.Colors.Text = viper.GetString("colors.textColor")
		}

		if viper.GetString("colors.sepColor") != "" {
			opt.Colors.Sep = viper.GetString("colors.sepColor")
		}

		if len(viper.GetStringSlice("colors.bodyColor")) != 0 {
			opt.Colors.Body = viper.GetStringSlice("colors.bodyColor")
		}

		// TODO: implement - alexdreptu (08 Nov 2016)
		if viper.GetBool("options.noColor") {
			fmt.Println("to be implemented")
			os.Exit(1)
		}

		if cmd.Flag("list-colors").Changed {
			archey.ListColors()
			os.Exit(1)
		}

		info, err := opt.Render()
		if err != nil {
			return err
		}

		fmt.Println(info)
		return nil
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
	RootCmd.Flags().BoolVar(&noGTK2Theme, "no-gtk2-theme", false, "don't print GTK2 theme name")
	RootCmd.Flags().BoolVar(&noGTK2IconTheme, "no-gtk2-icon-theme", false, "don't print GTK2 icon theme name")
	RootCmd.Flags().BoolVar(&noGTK2Font, "no-gtk2-font", false, "don't print GTK2 font name")
	RootCmd.Flags().BoolVar(&noGTK3Theme, "no-gtk3-theme", false, "don't print GTK3 theme name")
	RootCmd.Flags().BoolVar(&noGTK3IconTheme, "no-gtk3-icon-theme", false, "don't print GTK3 icon theme name")
	RootCmd.Flags().BoolVar(&noGTK3Font, "no-gtk3-font", false, "don't print GTK3 font name")
	RootCmd.Flags().BoolVar(&noTerminal, "no-terminal", false, "don't print terminal name")
	RootCmd.Flags().BoolVar(&noShell, "no-shell", false, "don't print shell name")
	RootCmd.Flags().BoolVar(&noEditor, "no-editor", false, "don't print editor name")
	RootCmd.Flags().BoolVar(&noPackages, "no-packages", false, "don't print packages count")
	RootCmd.Flags().BoolVar(&noMemory, "no-memory", false, "don't print memory usage")
	RootCmd.Flags().BoolVar(&noSwap, "no-swap", false, "don't print swap usage")
	RootCmd.Flags().BoolVar(&noCPU, "no-cpu", false, "don't print CPU model")
	RootCmd.Flags().BoolVar(&noRoot, "no-root", false, "don't print root disk usage")
	RootCmd.Flags().BoolVar(&noHome, "no-home", false, "don't print home disk usage")
	RootCmd.Flags().StringVar(&sep, "sep", "", "separator string")
	RootCmd.Flags().StringVar(&memoryUnit, "memory-unit", "", "unit to use for memory usage")
	RootCmd.Flags().StringVar(&swapUnit, "swap-unit", "", "unit to use for swap usage")
	RootCmd.Flags().StringVar(&diskUnit, "disk-unit", "", "unit to use for disk usage")
	RootCmd.Flags().StringSliceVar(&paths, "paths", nil, "additional paths to add to disk usage info")
	RootCmd.Flags().BoolVar(&pathFull, "path-full", false, "show full paths")
	RootCmd.Flags().BoolVar(&shellFull, "shell-full", false, "print shell's full path instead of its name")
	RootCmd.Flags().StringVar(&upSinceFormat, "up-since-format", "", "strftime format for up since")
	RootCmd.Flags().StringVar(&nameColor, "name-color", "", "color of the variable name")
	RootCmd.Flags().StringVar(&textColor, "text-color", "", "color of the text")
	RootCmd.Flags().StringVar(&sepColor, "sep-color", "", "color of the separator")
	RootCmd.Flags().StringSliceVar(&bodyColor, "body-color", nil, "color of the logo body")
	// TODO: implement - Alexandru Dreptu (20 Nov 2016)
	RootCmd.Flags().BoolVar(&noColor, "no-color", false, "don't use any colors")
	RootCmd.Flags().BoolVar(&listColors, "list-colors", false, "print all colors and styles")
	RootCmd.Flags().StringVar(&config, "config", "", "config file")

	viper.BindPFlag("show.noOS", RootCmd.Flags().Lookup("no-os"))
	viper.BindPFlag("show.noArch", RootCmd.Flags().Lookup("no-arch"))
	viper.BindPFlag("show.noKernel", RootCmd.Flags().Lookup("no-kernel"))
	viper.BindPFlag("show.noUser", RootCmd.Flags().Lookup("no-user"))
	viper.BindPFlag("show.noHostname", RootCmd.Flags().Lookup("no-hostname"))
	viper.BindPFlag("show.noUptime", RootCmd.Flags().Lookup("no-uptime"))
	viper.BindPFlag("show.noUpSince", RootCmd.Flags().Lookup("no-up-since"))
	viper.BindPFlag("show.noWM", RootCmd.Flags().Lookup("no-wm"))
	viper.BindPFlag("show.noDE", RootCmd.Flags().Lookup("no-de"))
	viper.BindPFlag("show.noGTK2Theme", RootCmd.Flags().Lookup("no-gtk2-theme"))
	viper.BindPFlag("show.noGTK2IconTheme", RootCmd.Flags().Lookup("no-gtk2-icon-theme"))
	viper.BindPFlag("show.noGTK2Font", RootCmd.Flags().Lookup("no-gtk2-font"))
	viper.BindPFlag("show.noGTK3Theme", RootCmd.Flags().Lookup("no-gtk3-theme"))
	viper.BindPFlag("show.noGTK3IconTheme", RootCmd.Flags().Lookup("no-gtk3-icon-theme"))
	viper.BindPFlag("show.noGTK3Font", RootCmd.Flags().Lookup("no-gtk3-font"))
	viper.BindPFlag("show.noTerminal", RootCmd.Flags().Lookup("no-terminal"))
	viper.BindPFlag("show.noShell", RootCmd.Flags().Lookup("no-shell"))
	viper.BindPFlag("show.noEditor", RootCmd.Flags().Lookup("no-editor"))
	viper.BindPFlag("show.noPackages", RootCmd.Flags().Lookup("no-packages"))
	viper.BindPFlag("show.noMemory", RootCmd.Flags().Lookup("no-memory"))
	viper.BindPFlag("show.noSwap", RootCmd.Flags().Lookup("no-swap"))
	viper.BindPFlag("show.noCPU", RootCmd.Flags().Lookup("no-cpu"))
	viper.BindPFlag("show.noRoot", RootCmd.Flags().Lookup("no-root"))
	viper.BindPFlag("show.noHome", RootCmd.Flags().Lookup("no-home"))

	viper.BindPFlag("options.sep", RootCmd.Flags().Lookup("sep"))
	viper.BindPFlag("options.memoryUnit", RootCmd.Flags().Lookup("memory-unit"))
	viper.BindPFlag("options.swapUnit", RootCmd.Flags().Lookup("swap-unit"))
	viper.BindPFlag("options.diskUnit", RootCmd.Flags().Lookup("disk-unit"))
	viper.BindPFlag("options.paths", RootCmd.Flags().Lookup("paths"))
	viper.BindPFlag("options.pathFull", RootCmd.Flags().Lookup("path-full"))
	viper.BindPFlag("options.shellFull", RootCmd.Flags().Lookup("shell-full"))
	viper.BindPFlag("options.upSinceFormat", RootCmd.Flags().Lookup("up-since-format"))
	viper.BindPFlag("options.noColor", RootCmd.Flags().Lookup("no-color"))

	viper.BindPFlag("colors.nameColor", RootCmd.Flags().Lookup("name-color"))
	viper.BindPFlag("colors.textColor", RootCmd.Flags().Lookup("text-color"))
	viper.BindPFlag("colors.sepColor", RootCmd.Flags().Lookup("sep-color"))
	viper.BindPFlag("colors.bodyColor", RootCmd.Flags().Lookup("body-color"))
}

func initConfig() {
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/archey-go")
	viper.AddConfigPath("$HOME/.archey-go")
	viper.AddConfigPath("/etc/archey-go")

	if config != "" {
		viper.SetConfigFile(config)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	viper.ReadInConfig()
}
