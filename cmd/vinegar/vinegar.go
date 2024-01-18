// vinegar.go

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/vinegarhq/vinegar/config"
	"github.com/vinegarhq/vinegar/config/editor"
	"github.com/vinegarhq/vinegar/internal/dirs"
	"github.com/vinegarhq/vinegar/roblox"
	"github.com/vinegarhq/vinegar/splash"
	"github.com/vinegarhq/vinegar/sysinfo"
	"github.com/vinegarhq/vinegar/wine"
)

var (
	BinPrefix string
	Version   string
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: vinegar [-config filepath] player|studio [args...]")
	fmt.Fprintln(os.Stderr, "       vinegar [-config filepath] exec prog args...")
	fmt.Fprintln(os.Stderr, "       vinegar [-config filepath] kill|winetricks|sysinfo")
	fmt.Fprintln(os.Stderr, "       vinegar delete|edit|submit|version")
	os.Exit(1)
}

func main() {
	configPath := flag.String("config", filepath.Join(dirs.Config, "config.toml"), "config.toml file which should be used")
	flag.Parse()

	cmd := flag.Arg(0)
	args := flag.Args()

	wine.Wine = "wine64"

	switch cmd {
	case "delete", "edit", "submit", "version":
		switch cmd {
		case "delete":
			Delete()
		case "edit":
			if err := editor.Edit(*configPath); err != nil {
				log.Fatal(err)
			}
		case "submit":
			if err := SubmitMerlin(); err != nil {
				log.Fatal(err)
			}
		case "version":
			fmt.Println("Vinegar", Version)
		}
	case "sysinfo", "player", "studio", "exec", "kill", "winetricks":
		cfg, err := config.Load(*configPath)
		if err != nil {
			log.Fatal(err)
		}

		pfx := wine.New(dirs.Prefix, os.Stderr)
		if err := os.MkdirAll(dirs.Prefix, 0o755); err != nil {
			log.Fatal(err)
		}

		switch cmd {
		case "sysinfo":
			Sysinfo(&pfx)
		case "exec":
			if len(args) < 2 {
				usage()
			}

			if err := pfx.Wine(args[1], args[2:]...).Run(); err != nil {
				log.Fatal(err)
			}
		case "kill":
			pfx.Kill()
		case "winetricks":
			if err := pfx.Winetricks(); err != nil {
				log.Fatal(err)
			}
		case "player":
			NewBinary(roblox.Player, &cfg, &pfx).Main(args[1:]...)
		case "studio":
			NewBinary(roblox.Studio, &cfg, &pfx).Main(args[1:]...)
		}
	default:
		usage()
	}
}

func Delete() {
	log.Println("Deleting Wineprefix")
	if err := os.RemoveAll(dirs.Prefix); err != nil {
		log.Fatal(err)
	}
}

func Sysinfo(pfx *wine.Prefix) {
	cmd := pfx.Wine("--version")
	cmd.Stdout = nil
	ver, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	var revision string
	bi, _ := debug.ReadBuildInfo()
	for _, bs := range bi.Settings {
		if bs.Key == "vcs.revision" {
			revision = fmt.Sprintf("(%s)", bs.Value)
		}
	}

	info := `* Vinegar: %s %s
* Distro: %s
* Processor: %s
* Kernel: %s
* Wine: %s`

	fmt.Printf(info, Version, revision, sysinfo.Distro, sysinfo.CPU.Name, sysinfo.Kernel, ver)

	fmt.Println("* Cards:")
	for i, c := range sysinfo.Cards {
		fmt.Printf("  * Card %d: %s %s %s\n", i, c.Driver, path.Base(c.Device), c.Path)
	}
}

func LogFile(name string) (*os.File, error) {
	if err := dirs.Mkdirs(dirs.Logs); err != nil {
		return nil, err
	}

	path := filepath.Join(dirs.Logs, name+"-"+time.Now().Format(time.RFC3339)+".log")

	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create %s log file: %w", name, err)
	}

	log.Printf("Logging to file: %s", path)

	return file, nil
}

func (b *Binary) Main(args ...string) {
	b.Config.Env.Setenv()

	logFile, err := LogFile(b.Type.String())
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logOutput := io.MultiWriter(logFile, os.Stderr)
	b.Prefix.Output = logOutput
	log.SetOutput(logOutput)

	firstRun := false
	if _, err := os.Stat(filepath.Join(b.Prefix.Dir(), "drive_c", "windows")); err != nil {
		firstRun = true
	}

	if !wine.WineLook() {
		b.Splash.Dialog(DialogNoWine, false)
		log.Fatalf("%s is required to run roblox", wine.Wine)
	}

	go func() {
		err := b.Splash.Run()
		if errors.Is(splash.ErrClosed, err) {
			log.Printf("Splash window closed!")
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			return
		}

		if err != nil {
			log.Fatalf("splash: %s", err)
		}
	}()

	errHandler := func(err error) {
		if !b.GlobalConfig.Splash.Enabled || b.Splash.IsClosed() {
			log.Fatal(err)
		}

		log.Println(err)
		b.Splash.LogPath = logFile.Name()
		b.Splash.Invalidate()
		b.Splash.Dialog(fmt.Sprintf(DialogFailure, err), false)
		os.Exit(1)
	}

	if firstRun {
		log.Printf("Initializing wineprefix at %s", b.Prefix.Dir())
		b.Splash.SetMessage("Initializing wineprefix")

		if err := b.Prefix.SetDPI(97); err != nil {
			b.Splash.SetMessage(err.Error())
			errHandler(err)
		}
	}

	func() {
		if len(args) < 1 {
			return
		}

		c := regexp.MustCompile(`channel:([^+]*)`).FindStringSubmatch(args[0])
		if len(c) < 1 {
			return
		}

		if c[1] != "" && c[1] != b.Config.Channel {
			r := b.Splash.Dialog(
				fmt.Sprintf(DialogReqChannel, c[1], b.Config.Channel),
				true,
			)
			if r {
				log.Println("Switching user channel temporarily to", c[1])
				b.Config.Channel = c[1]
			}
		}
	}()

	b.Splash.SetDesc(b.Config.Channel)

	if err := b.Setup(); err != nil {
		b.Splash.SetMessage("Failed to setup Roblox")
		errHandler(err)
	}

	if err := b.Run(args...); err != nil {
		b.Splash.SetMessage("Failed to run Roblox")
		errHandler(err)
	}
}
