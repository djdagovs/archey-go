package cmd

import (
	"fmt"
	"os"

	archey "github.com/alexdreptu/archey-go/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	version  = "0.1.0"
	author   = "Alexandru Dreptu <alexdreptu@gmail.com>"
	homeLink = "https://github.com/alexdreptu/archey-go"
	bugsLink = "https://github.com/alexdreptu/archey-go/issues"
)

var (
	noOS              bool
	noArch            bool
	noKernel          bool
	noUser            bool
	noHostname        bool
	noUptime          bool
	noUpSince         bool
	noWM              bool
	noDE              bool
	noGTK2Theme       bool
	noGTK2IconTheme   bool
	noGTK2Font        bool
	noGTK2CursorTheme bool
	noGTK3Theme       bool
	noGTK3IconTheme   bool
	noGTK3CursorTheme bool
	noGTK3Font        bool
	noTerminal        bool
	noShell           bool
	noEditor          bool
	noPackages        bool
	noMemory          bool
	noSwap            bool
	noCPU             bool
	noRoot            bool
	noHome            bool
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

var config string

var usageTemplate = `Version: {{version}}
Author: {{author}}
Home: {{homeLink}}

Usage:
      {{.Name}} [flags]

Example:
      {{.Name}} {{.Example}}

Flags:
{{.Flags.FlagUsages}}
Report bugs to {{bugsLink}}
`

var RootCmd = &cobra.Command{
	Use:   "archey-go",
	Short: "Archey-go is a tool to display prettified system info on Arch Linux",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			cmd.Help()
			os.Exit(0)
		}

		opt := archey.New()

		opt.Show.OS = viper.GetBool("show.no_os")
		opt.Show.Arch = viper.GetBool("show.no_arch")
		opt.Show.Kernel = viper.GetBool("show.no_kernel")
		opt.Show.User = viper.GetBool("show.no_user")
		opt.Show.Hostname = viper.GetBool("show.no_hostname")
		opt.Show.Uptime = viper.GetBool("show.no_uptime")
		opt.Show.UpSince = viper.GetBool("show.no_up_since")
		opt.Show.WM = viper.GetBool("show.no_wm")
		opt.Show.DE = viper.GetBool("show.no_de")
		opt.Show.GTK2Theme = viper.GetBool("show.no_gtk2_theme")
		opt.Show.GTK2IconTheme = viper.GetBool("show.no_gtk2_icon_theme")
		opt.Show.GTK2Font = viper.GetBool("show.no_gtk2_font")
		opt.Show.GTK2CursorTheme = viper.GetBool("show.no_gtk2_cursor_theme")
		opt.Show.GTK3Theme = viper.GetBool("show.no_gtk3_theme")
		opt.Show.GTK3IconTheme = viper.GetBool("show.no_gtk3_icon_theme")
		opt.Show.GTK3Font = viper.GetBool("show.no_gtk3_font")
		opt.Show.GTK3CursorTheme = viper.GetBool("show.no_gtk3_cursor_theme")
		opt.Show.Terminal = viper.GetBool("show.no_terminal")
		opt.Show.Shell = viper.GetBool("show.no_shell")
		opt.Show.Editor = viper.GetBool("show.no_editor")
		opt.Show.Packages = viper.GetBool("show.no_packages")
		opt.Show.Memory = viper.GetBool("show.no_memory")
		opt.Show.Swap = viper.GetBool("show.no_swap")
		opt.Show.CPU = viper.GetBool("show.no_cpu")
		opt.Show.Root = viper.GetBool("show.no_root")
		opt.Show.Home = viper.GetBool("show.no_home")

		if viper.GetString("options.sep") != "" {
			opt.Sep = viper.GetString("options.sep")
		}

		if viper.GetString("options.memory_unit") != "" {
			opt.MemoryUnit = viper.GetString("options.memory_unit")
		}

		if viper.GetString("options.swap_unit") != "" {
			opt.SwapUnit = viper.GetString("options.swap_unit")
		}

		if viper.GetString("options.disk_unit") != "" {
			opt.DiskUnit = viper.GetString("options.disk_unit")
		}

		opt.Paths = viper.GetStringSlice("options.paths")
		opt.PathFull = viper.GetBool("options.path_full")
		opt.ShellFull = viper.GetBool("options.shell_full")

		if viper.GetString("options.up_since_format") != "" {
			opt.UpSinceFormat = viper.GetString("options.up_since_format")
		}

		if viper.GetString("colors.name_color") != "" {
			opt.Colors.Name = viper.GetString("colors.name_color")
		}

		if viper.GetString("colors.text_color") != "" {
			opt.Colors.Text = viper.GetString("colors.text_color")
		}

		if viper.GetString("colors.sep_color") != "" {
			opt.Colors.Sep = viper.GetString("colors.sep_color")
		}

		if len(viper.GetStringSlice("colors.body_color")) != 0 {
			opt.Colors.Body = viper.GetStringSlice("colors.body_color")
		}

		if viper.GetBool("options.no_color") {
			archey.NoColor()
		}

		if cmd.Flag("list-colors").Changed {
			archey.ListColors()
			os.Exit(1)
		}

		if cmd.Flag("version").Changed {
			fmt.Println("version " + version)
			os.Exit(0)
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
	cobra.AddTemplateFunc("version", func() string { return version })
	cobra.AddTemplateFunc("author", func() string { return author })
	cobra.AddTemplateFunc("homeLink", func() string { return homeLink })
	cobra.AddTemplateFunc("bugsLink", func() string { return bugsLink })
	RootCmd.Example = `--body-color 111 --name-color 150 --sep ' ->' --sep-color 191 \
	--shell-full --memory-unit mb --no-swap --paths /tmp,/usr --path-full`
	RootCmd.SetUsageTemplate(usageTemplate)
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
	RootCmd.Flags().BoolVar(&noGTK2CursorTheme, "no-gtk2-cursor-theme", false, "don't print GTK2 cursor theme name")
	RootCmd.Flags().BoolVar(&noGTK3Theme, "no-gtk3-theme", false, "don't print GTK3 theme name")
	RootCmd.Flags().BoolVar(&noGTK3IconTheme, "no-gtk3-icon-theme", false, "don't print GTK3 icon theme name")
	RootCmd.Flags().BoolVar(&noGTK3Font, "no-gtk3-font", false, "don't print GTK3 font name")
	RootCmd.Flags().BoolVar(&noGTK3CursorTheme, "no-gtk3-cursor-theme", false, "don't print GTK3 cursor theme name")
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
	RootCmd.Flags().StringVarP(&config, "config", "c", "", "config file")
	RootCmd.Flags().BoolP("version", "v", false, "print version")
	RootCmd.Flags().BoolP("no-color", "n", false, "don't use any colors")
	RootCmd.Flags().BoolP("list-colors", "l", false, "print all colors and styles")

	viper.BindPFlag("show.no_os", RootCmd.Flags().Lookup("no-os"))
	viper.BindPFlag("show.no_arch", RootCmd.Flags().Lookup("no-arch"))
	viper.BindPFlag("show.no_kernel", RootCmd.Flags().Lookup("no-kernel"))
	viper.BindPFlag("show.no_user", RootCmd.Flags().Lookup("no-user"))
	viper.BindPFlag("show.no_hostname", RootCmd.Flags().Lookup("no-hostname"))
	viper.BindPFlag("show.no_uptime", RootCmd.Flags().Lookup("no-uptime"))
	viper.BindPFlag("show.no_up_since", RootCmd.Flags().Lookup("no-up-since"))
	viper.BindPFlag("show.no_wm", RootCmd.Flags().Lookup("no-wm"))
	viper.BindPFlag("show.no_de", RootCmd.Flags().Lookup("no-de"))
	viper.BindPFlag("show.no_gtk2_theme", RootCmd.Flags().Lookup("no-gtk2-theme"))
	viper.BindPFlag("show.no_gtk2_icon_theme", RootCmd.Flags().Lookup("no-gtk2-icon-theme"))
	viper.BindPFlag("show.no_gtk2_font", RootCmd.Flags().Lookup("no-gtk2-font"))
	viper.BindPFlag("show.no_gtk2_cursor_theme", RootCmd.Flags().Lookup("no-gtk2-cursor-theme"))
	viper.BindPFlag("show.no_gtk3_theme", RootCmd.Flags().Lookup("no-gtk3-theme"))
	viper.BindPFlag("show.no_gtk3_icon_theme", RootCmd.Flags().Lookup("no-gtk3-icon-theme"))
	viper.BindPFlag("show.no_gtk3_font", RootCmd.Flags().Lookup("no-gtk3-font"))
	viper.BindPFlag("show.no_gtk3_cursor_theme", RootCmd.Flags().Lookup("no-gtk3-cursor-theme"))
	viper.BindPFlag("show.no_terminal", RootCmd.Flags().Lookup("no-terminal"))
	viper.BindPFlag("show.no_shell", RootCmd.Flags().Lookup("no-shell"))
	viper.BindPFlag("show.no_editor", RootCmd.Flags().Lookup("no-editor"))
	viper.BindPFlag("show.no_packages", RootCmd.Flags().Lookup("no-packages"))
	viper.BindPFlag("show.no_memory", RootCmd.Flags().Lookup("no-memory"))
	viper.BindPFlag("show.no_swap", RootCmd.Flags().Lookup("no-swap"))
	viper.BindPFlag("show.no_cpu", RootCmd.Flags().Lookup("no-cpu"))
	viper.BindPFlag("show.no_Root", RootCmd.Flags().Lookup("no-root"))
	viper.BindPFlag("show.no_Home", RootCmd.Flags().Lookup("no-home"))

	viper.BindPFlag("options.sep", RootCmd.Flags().Lookup("sep"))
	viper.BindPFlag("options.memory_unit", RootCmd.Flags().Lookup("memory-unit"))
	viper.BindPFlag("options.swap_unit", RootCmd.Flags().Lookup("swap-unit"))
	viper.BindPFlag("options.disk_unit", RootCmd.Flags().Lookup("disk-unit"))
	viper.BindPFlag("options.paths", RootCmd.Flags().Lookup("paths"))
	viper.BindPFlag("options.path_full", RootCmd.Flags().Lookup("path-full"))
	viper.BindPFlag("options.shell_full", RootCmd.Flags().Lookup("shell-full"))
	viper.BindPFlag("options.up_since_format", RootCmd.Flags().Lookup("up-since-format"))
	viper.BindPFlag("options.no_color", RootCmd.Flags().Lookup("no-color"))

	viper.BindPFlag("colors.name_color", RootCmd.Flags().Lookup("name-color"))
	viper.BindPFlag("colors.text_color", RootCmd.Flags().Lookup("text-color"))
	viper.BindPFlag("colors.sep_color", RootCmd.Flags().Lookup("sep-color"))
	viper.BindPFlag("colors.body_color", RootCmd.Flags().Lookup("body-color"))
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
