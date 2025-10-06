package utils

import (
	"strings"
)

func splitQuote(cmd string) []string {
	cmdRune := []rune(cmd)
	var res []string
	var current strings.Builder
	inQuote := false
	i := 0

	for i < len(cmdRune) {
		c := cmdRune[i]

		if c == '"' {
			inQuote = !inQuote
			i++
			continue
		}

		if !inQuote && c == ' ' {
			if current.Len() > 0 {
				res = append(res, strings.TrimSpace(current.String()))
				current.Reset()
			}
		} else {
			current.WriteRune(c)
		}
		i++
	}

	if current.Len() > 0 {
		res = append(res, strings.TrimSpace(current.String()))
	}

	return res
}

func ParseCommand(cmd string) (string, map[string]string, []string) {
	cmdArr := splitQuote(cmd)
	if len(cmdArr) == 0 {
		return "", make(map[string]string), []string{}
	}

	command := cmdArr[0]
	flags := make(map[string]string)
	args := []string{}

	i := 1
	for i < len(cmdArr) {
		token := cmdArr[i]
		if strings.HasPrefix(token, "-") {
			flag := token
			value := ""
			if i+1 < len(cmdArr) {
				next := cmdArr[i+1]
				if !strings.HasPrefix(next, "-") {
					value = next
					i++
				}
			}
			flags[flag] = value
		} else {
			args = append(args, token)
		}
		i++
	}

	return command, flags, args
}
