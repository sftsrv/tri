package command

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Replace `(<` with `(?P<` to support JS style regexps
func createRegexp(re string) (*regexp.Regexp, error) {
	formatted := strings.ReplaceAll(re, `(<`, `(?P<`)

	return regexp.Compile(formatted)
}

func generateCommand(base string, pattern string, input string) (bin string, args []string, err error) {
	re, err := createRegexp(pattern)

	result := base
	names := re.SubexpNames()

	matches := re.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		for groupIndex, submatch := range match {
			groupName := names[groupIndex]

			if groupName == "" {
				groupName = strconv.Itoa(groupIndex)
			}

			replacer := fmt.Sprintf("$%s", groupName)
			result = strings.ReplaceAll(result, replacer, submatch)
		}
	}

	// Replace lone $ only after all others have been done.
	// This is a convenience syntax for $0
	result = strings.ReplaceAll(result, "$", input)
	if len(result) > 0 {
		parts := strings.Split(result, " ")
		bin := parts[0]
		args := parts[1:]

		return bin, args, nil
	}

	return bin, args, fmt.Errorf("could not create command from: \n  base: %s \n  pattern: %s \n  input: %s \n  result: %s", base, pattern, input, result)
}

func useCommand(preview string, input string, width int) (string, []string) {
	if len(preview) > 0 {
		parts := strings.Split(preview, " ")
		bin := parts[0]
		args := append(parts[1:], input)

		return bin, args
	}

	_, err := exec.LookPath("bat")
	if err == nil {
		return "bat", []string{"--color=always", "--number", "--terminal-width", strconv.Itoa(width), input}
	} else {
		return "cat", []string{input}
	}
}

// If `pattern` is provided will use command generation - otherwise will default to
// simple append-based behavior
func CreateCommand(base string, pattern string, input string, width int) (*exec.Cmd, error) {
	if pattern == "" {
		bin, args := useCommand(base, input, width)
		cmd := exec.Command(bin, args...)
		return cmd, nil
	}

	var err error
	bin, args, err := generateCommand(base, pattern, input)

	if err != nil {
		return nil, err
	}

	return exec.Command(bin, args...), nil

}
