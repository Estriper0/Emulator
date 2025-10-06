package app

import (
	"encoding/csv"
	"os"
)

type Emulator struct {
	user       string
	path       string
	scriptPath string
	vfs        [][]string
	vfsPath    string
}

func NewEmulator(script string, vfsPath string) *Emulator {
	return &Emulator{
		user:       "user",
		path:       "~",
		scriptPath: script,
		vfsPath:    vfsPath,
	}
}

func (e *Emulator) Run() error {
	file, err := os.Open(e.vfsPath)
	if err != nil {
		return err
	}
	reader := csv.NewReader(file)

	vfs, err := reader.ReadAll()
	if err != nil {
		return err
	}
	file.Close()
	e.vfs = vfs

	if e.scriptPath != "" {
		file, err := os.Open(e.scriptPath)
		if err != nil {
			return err
		}
		e.command_reader(file, true)
		file.Close()
	}

	e.command_reader(os.Stdin, false)
	return nil
}

func (e *Emulator) Save() error {
	file, err := os.Create(e.vfsPath)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(file)
	writer.WriteAll(e.vfs)
	writer.Flush()
	return nil
}
