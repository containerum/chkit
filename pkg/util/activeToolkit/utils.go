package activeToolkit

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func IsStop(word string) bool {
	word = strings.TrimSpace(strings.ToLower(word))
	switch word {
	case "n", "no", "nope", "stop", "q", "quit", "exit", "e", "-":
		return true
	default:
		return false
	}
}

func AskLine(promt string) (string, bool) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s", promt)
	var input string
	for scanner.Scan() {
		input = scanner.Text()
		break
	}
	return input, scanner.Err() == io.EOF
}

func Yes(message string) (bool, string) {
	fmt.Printf("%s [Y/N]: ", message)
	scanner := bufio.NewScanner(os.Stdin)
	answer := "N"
	for scanner.Scan() {
		answer = strings.TrimSpace(scanner.Text())
		break
	}
	return strings.ToLower(answer) == "y", answer
}

func AskWord(message string) (string, bool) {
	fmt.Printf("%s", message)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	var input string
	for scanner.Scan() {
		input = scanner.Text()
		break
	}
	return input, scanner.Err() == io.EOF
}

func AskFieldToChange(fields []string) (int, bool) {
	_, n, exit := Options("Which field do you want to change? (print no to stop): ",
		true,
		fields...)
	return n, !exit
}

func Options(msg string, withStop bool, options ...string) (string, int, bool) {
	fmt.Printf("%s\n", msg)
	for i, opt := range options {
		fmt.Printf("%d) %s\n", i+1, opt)
	}
	for {
		nStr, exit := AskLine("Choose wisely: ")
		if exit || (withStop && IsStop(nStr)) {
			return "", -1, true
		}
		if n, err := strconv.Atoi(nStr); err != nil {
			for i, opt := range options {
				if opt == strings.TrimSpace(nStr) {
					return opt, i, false
				}
			}
			fmt.Printf("Option %q not found :( Try again.\n", nStr)
			continue
		} else if n > 0 && n <= len(options) {
			return options[n-1], n - 1, false
		} else {
			fmt.Printf("Option %d not found :( Try again.\n", n)
			continue
		}
	}
}

func OrString(str, def string) string {
	if strings.TrimSpace(str) == "" {
		return def
	}
	return str
}

func OrStringer(str fmt.Stringer, def string) string {
	if str == nil {
		return def
	}
	return str.String()
}

func OrValue(val interface{}, def string) string {
	if val == nil {
		return def
	}
	return fmt.Sprintf("%v", val)
}
