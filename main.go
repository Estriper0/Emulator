package main

import (
	"flag"
	"fmt"

	"github.com/Estriper0/emulator/app"
)

func main() {
	vfsPath := flag.String("vfs", "vfs.csv", "path to vfs")
	startPath := flag.String("start_script", "", "path to start script")
	flag.Parse()

	emulator := app.NewEmulator(*startPath, *vfsPath)
	if err := emulator.Run(); err != nil {
		fmt.Println(err)
		return
	}

	if err := emulator.Save(); err != nil {
		fmt.Println(err)
		return
	}
}
