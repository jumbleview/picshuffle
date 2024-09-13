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
	var isLogPrintMode bool
	flag.BoolVar(&isLogPrintMode, "l", false, "print history (log)")
	flag.Parse()
	cmd := flag.Args()

	if !isSilent || isLogPrintMode {
		screen.MakeConsole()
	}
	greeting := "picshuffle.exe [-s] [-l] path_to_jpeg_folder_or_file"
	var Usage = func() {
		fmt.Println(greeting)
	}
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	execPath := filepath.Dir(ex)
	if len(cmd) != 1 {
		fmt.Printf("Wrong number of arguments %d\n", len(cmd))
		Usage()
		fmt.Printf("Hit \"Enter\" to exit... ")
		s := ""
		fmt.Scanln(&s)
		os.Exit(1)
	}
	if isLogPrintMode {
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
