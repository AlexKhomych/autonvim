package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

// RunCommand returns command status code and error message.
// Following status codes are possible:
// 0 for successful execution
// -1 for failed execution
// any other non-negative number representing command exit code,
// which is useful for case checks.
func FRunCommand(cmd string, args []string, useSudo bool) (int, error) {
	if useSudo {
		args = append([]string{cmd}, args...)
		cmd = "sudo"
	}

	command := exec.Command(cmd, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			statusCode := exitError.ExitCode()
			return statusCode, fmt.Errorf("command failed to execute with status code %d: %v", statusCode, err)
		}
		return -1, fmt.Errorf("failed to run %s: %v", cmd, err)
	}
	return 0, nil
}

// Prompt asks a user for an input (y/n)
// and returns boolean representing:
// true if lowercase answer is (y) and false otherwise
func FPrompt(ask string) bool {
	var userInput string
	fmt.Println(ask)
	fmt.Scanln(&userInput)
	userInput = strings.ToLower(userInput)
	return userInput == "y"
}

func FReflectName(i any) string {
	return reflect.TypeOf(i).Elem().Name()
}

func FReflectEq(a any, b any) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}

func FPrefixError(p, msg string) error {
	return fmt.Errorf("%s: %s", p, msg)
}

func FCreateDir(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		slog.Error(err.Error())
		return fmt.Errorf("failed to create directory, %s", path)
	}
	return nil
}
