package app

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	rightsUsers = map[string]string{
		"0": "---",
		"1": "--x",
		"2": "-w-",
		"3": "-wx",
		"4": "r--",
		"5": "r-x",
		"6": "rw-",
		"7": "rwx",
	}
)

func (e *Emulator) command_handler(comm string, flags map[string]string, args []string) string {
	switch comm {
	case "":
		return ""
	case "ls":
		res := strings.Builder{}
		c := 0
		s, ok := flags["-l"]
		if len(args) > 0 || s != "" {
			return "unknown args"
		}
		for _, record := range e.vfs {
			if record[0] == e.path {
				if ok {
					if len(flags) > 1 {
						return "unknown flag"
					}
					if c != 0 {
						res.WriteString("\n")
					}
					for _, s := range record[5] {
						res.WriteString(rightsUsers[string(s)])
					}
					res.WriteString(" ")
					res.WriteString(record[4] + " " + record[6] + " " + record[1])
					c++
				} else {
					if len(flags) > 0 {
						return "unknown flag"
					}
					res.WriteString(record[1] + " ")
				}

			}
		}
		return res.String()
	case "cd":
		if len(args) > 1 {
			return "unknown args"
		}
		if len(flags) > 0 {
			return "unknown flag"
		}
		for _, record := range e.vfs {
			if len(args) == 0 {
				return "empty path"
			} else if args[0] == "." {
				e.path = "~"
				return ""
			} else if (e.path + "/" + args[0]) == (record[0] + "/" + record[1]) {
				if record[2] == "true" {
					e.path = e.path + "/" + record[1]
					return ""
				}
				return "can not open"
			} else if ("~" + args[0]) == (record[0] + "/" + record[1]) {
				if record[2] == "true" {
					e.path = "~/" + record[1]
					return ""
				}
				return "can not open"
			}
		}
		return "not found"
	case "rev":
		if len(args) > 1 {
			return "unknown args"
		}
		if len(flags) > 0 {
			return "unknown flag"
		}
		for _, record := range e.vfs {
			if (e.path+"/"+args[0] == record[0]+"/"+record[1]) || ("~"+args[0] == record[0]+"/"+record[1]) {
				n, _ := strconv.Atoi(string(record[5]))
				if (record[4] == e.user && n/100 >= 4) || (record[4] != e.user && n%10 >= 4) || e.user == "root" {
					res := DecodeFromBase64(record[3])
					return reverse(res)
				}
				return "no rights"
			}
		}
		return "not found"
	case "head":
		if len(args) == 0 {
			return "invalid command"
		}
		if len(args) > 1 {
			return "unknown args"
		}
		nStr, ok := flags["-n"]
		for _, record := range e.vfs {
			if (e.path+"/"+args[0] == record[0]+"/"+record[1]) || ("~"+args[0] == record[0]+"/"+record[1]) {
				n, _ := strconv.Atoi(string(record[5]))
				if !((record[4] == e.user && n/100 >= 4) || (record[4] != e.user && n%10 >= 4) || e.user == "root") {
					return "no rights"
				}
				res := DecodeFromBase64(record[3])
				if ok {
					if len(flags) > 1 {
						return "unknown flag"
					}
					n, err := strconv.Atoi(nStr)
					if err != nil {
						return "invalid command"
					}
					return firstNString(res, n)
				}
				return firstNString(res, 10)

			}
		}
		return "not found"
	case "chown":
		if len(args) < 2 {
			return "invalid command"
		}
		for i, record := range e.vfs {
			if record[0] == e.path && args[1] == record[1] {
				if record[4] != e.user && e.user != "root" {
					return "no rights"
				}
				e.vfs[i][4] = args[0]
				tm := time.Now()
				e.vfs[i][6] = fmt.Sprintf("%d-%d-%d %d:%d:%d", tm.Day(), tm.Month(), tm.Year(), tm.Hour(), tm.Minute(), tm.Second())
				return ""
			}
		}
		return "not found"
	case "chmod":
		if len(args) < 2 {
			return "invalid command"
		}
		if !correctRights(args[0]) {
			return "invalid command"
		}
		for i, record := range e.vfs {
			if record[0] == e.path && args[1] == record[1] {
				if record[4] != e.user && e.user != "root" {
					return "no rights"
				}
				e.vfs[i][5] = args[0]
				tm := time.Now()
				e.vfs[i][6] = fmt.Sprintf("%d-%d-%d %d:%d:%d", tm.Day(), tm.Month(), tm.Year(), tm.Hour(), tm.Minute(), tm.Second())
				return ""
			}
		}
		return "not found"
	case "sudo":
		usr := e.user
		e.user = "root"
		resp := e.command_handler(args[0], map[string]string{}, args[1:])
		e.user = usr
		return resp
	case "exit":
		return ""
	default:
		return "unknown command"
	}
}

func reverse(s string) string {
	runes := []rune(strings.ReplaceAll(s, "\\n", "\n"))
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func firstNString(s string, n int) string {
	s = strings.ReplaceAll(s, "\\n", "\n")
	arrString := strings.Split(s, "\n")
	if n > len(arrString) {
		n = len(arrString)
	}
	return strings.Join(arrString[0:n], "\n")
}

func correctRights(s string) bool {
	sRune := []rune(s)
	if len(sRune) != 3 {
		return false
	}
	for _, el := range sRune {
		elStr := string(el)
		if elStr != "0" && elStr != "1" && elStr != "2" && elStr != "3" && elStr != "4" && elStr != "5" && elStr != "6" && elStr != "7" {
			return false
		}
	}
	return true
}

func EncodeToBase64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}
func DecodeFromBase64(input string) string {
	decoded, _ := base64.StdEncoding.DecodeString(input)
	return string(decoded)
}
