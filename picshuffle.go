// go build -ldflags -H=windowsgui
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jumbleview/decor/screen"
)

func main() {
	var isSilent bool
	flag.BoolVar(&isSilent, "s", false, "silent mode")
	var mustPrintHistory bool
	flag.BoolVar(&mustPrintHistory, "l", false, "print history (log)")
	flag.Parse()
	cmd := flag.Args()

	if !isSilent || mustPrintHistory {
		screen.MakeConsole()
	}
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	execPath := filepath.Dir(ex)
	if len(cmd) != 1 {
		fmt.Printf("...")
		os.Exit(1)
	}
	if mustPrintHistory {
		PrintLog(execPath)
	} else {
		if name, err := ChooseFile(execPath, cmd[0]); err == nil {
			err = screen.SetWallpaper(name, execPath)
			fmt.Printf("%s %s\n", err, name)
		}
	}
	fmt.Printf("Hit \"Enter\" to exit... ")
	s := ""
	fmt.Scanln(&s)
}
