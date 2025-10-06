package app

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Estriper0/emulator/utils"
)

func (e *Emulator) command_reader(r io.Reader, fromFile bool) {
	scanner := bufio.NewScanner(r)
	var input string
	for {
		if fromFile {
			if !scanner.Scan() {
				fmt.Println("script complete")
				break
			}
			fmt.Printf("%s@%s:%s$ ", e.user, e.vfsPath, e.path)
			input = scanner.Text()
			fmt.Println(input)
		} else {
			fmt.Printf("%s@%s:%s$ ", e.user, e.vfsPath, e.path)
			if !scanner.Scan() {
				break
			}
			input = scanner.Text()
		}

		comm, flags, args := utils.ParseCommand(input)
		if comm == "exit" {
			fmt.Println("exit")
			break
		}

		resp := e.command_handler(comm, flags, args)
		if resp != "" {
			fmt.Println(resp)
		}
	}
	if fromFile {
		fmt.Println("script complete")
	}
}
