package main

import (
	"fmt"
	"os"

	archey "github.com/alectic/archey-go/lib"
	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "Archey"
	app.Usage = "a tool to display system info on Arch Linux in a pretty way"
	app.HideVersion = true
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Alexandru Dreptu",
			Email: "alecticwp@gmail.com",
		},
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "no-user",
			Usage: "don't print user",
		},
		cli.BoolFlag{
			Name:  "no-hostname",
			Usage: "don't print hostname",
		},
		cli.BoolFlag{
			Name:  "no-os",
			Usage: "don't print os name",
		},
		cli.BoolFlag{
			Name:  "no-kernel",
			Usage: "don't print kernel version",
		},
		cli.BoolFlag{
			Name:  "no-uptime",
			Usage: "don't print uptime",
		},
		cli.BoolFlag{
			Name:  "no-wm",
			Usage: "don't print Window Manager name",
		},
		cli.BoolFlag{
			Name:  "no-de",
			Usage: "don't print Desktop Environment name",
		},
		cli.BoolFlag{
			Name:  "no-terminal",
			Usage: "don't print terminal name",
		},
		cli.BoolFlag{
			Name:  "no-shell",
			Usage: "don't print shell name",
		},
		cli.BoolFlag{
			Name:  "no-editor",
			Usage: "don't print editor",
		},
		cli.BoolFlag{
			Name:  "no-packages",
			Usage: "don't print packages count",
		},
		cli.BoolFlag{
			Name:  "no-memory",
			Usage: "don't print memory",
		},
		cli.BoolFlag{
			Name:  "no-cpu",
			Usage: "don't print CPU",
		},
		cli.BoolFlag{
			Name:  "no-root",
			Usage: "don't print root info",
		},
		cli.BoolFlag{
			Name:  "no-home",
			Usage: "don't print home info",
		},
		cli.StringFlag{
			Name:  "sep",
			Usage: "separator string",
		},
		cli.StringFlag{
			Name:  "memory-unit",
			Usage: "unit to use for memory usage",
		},
		cli.StringFlag{
			Name:  "disk-unit",
			Usage: "unit to use for disk usage",
		},
		cli.StringFlag{
			Name:  "paths",
			Usage: "additional paths to add to disk usage info",
		},
		cli.StringFlag{
			Name:  "name-color",
			Usage: "color of the variable name",
		},
		cli.StringFlag{
			Name:  "text-color",
			Usage: "color of the text",
		},
		cli.StringFlag{
			Name:  "sep-color",
			Usage: "color of the separator",
		},
		cli.StringFlag{
			Name:  "body-color",
			Usage: "color of the logo body",
		},
		cli.BoolFlag{
			Name:  "no-color",
			Usage: "don't use any colors",
		},
		cli.BoolFlag{
			Name:  "list-colors",
			Usage: "print all colors and styles",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.NArg() > 0 {
			cli.ShowAppHelp(c)
			os.Exit(1)
		}

		opt := archey.New()
		opt.Show.User = c.Bool("no-user")
		opt.Show.Hostname = c.Bool("no-hostname")
		opt.Show.OS = c.Bool("no-os")
		opt.Show.Kernel = c.Bool("no-kernel")
		opt.Show.Uptime = c.Bool("no-uptime")
		opt.Show.WM = c.Bool("no-wm")
		opt.Show.DE = c.Bool("no-de")
		opt.Show.Terminal = c.Bool("no-terminal")
		opt.Show.Shell = c.Bool("no-shell")
		opt.Show.Editor = c.Bool("no-editor")
		opt.Show.Packages = c.Bool("no-packages")
		opt.Show.Memory = c.Bool("no-memory")
		opt.Show.CPU = c.Bool("no-cpu")
		opt.Show.Root = c.Bool("no-root")
		opt.Show.Home = c.Bool("no-home")

		if c.String("sep") != "" {
			opt.Sep = c.String("sep")
		}

		if c.String("memory-unit") != "" {
			opt.MemoryUnit = c.String("memory-unit")
		}

		if c.String("disk-unit") != "" {
			opt.DiskUnit = c.String("disk-unit")
		}

		if len(c.String("paths")) != 0 {
			opt.Paths = c.String("paths")
		}

		if c.String("name-color") != "" {
			opt.Colors.Name = c.String("name-color")
		}

		if c.String("text-color") != "" {
			opt.Colors.Text = c.String("text-color")
		}

		if c.String("sep-color") != "" {
			opt.Colors.Sep = c.String("sep-color")
		}

		if c.String("body-color") != "" {
			opt.Colors.Body = c.String("body-color")
		}

		// TODO: implement - alectic (27 Oct 2016)
		if c.Bool("no-color") {
			fmt.Println("to be implemented")
			os.Exit(0)
		}

		if c.Bool("list-colors") {
			archey.ListColors()
			os.Exit(0)
		}

		info, err := opt.Render()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println(info)
		return nil
	}
	app.Run(os.Args)
}
