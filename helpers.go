package main

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
)

type TaskHelper struct{}

// IsPackageInstalled uses dpkg-query with --status flag to verify if package exists.
// Will print package information if present or unavailability otherwise.
func (t TaskHelper) IsPackageInstalled(pkgName string, isSudo bool) (bool, error) {
	cmd := "dpkg-query"
	args := []string{"--status", pkgName}
	errCode, err := FRunCommand(cmd, args, isSudo)
	if err != nil && errCode == 1 {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check package presence: %v", err)
	}
	return true, nil
}

// IsPathEmpty checks if path exists and is empty.
func (t TaskHelper) IsPathEmpty(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return true, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to check path stat: %v", err)
	}

	return false, nil
}

// CloneGitRepo executed git clone command with git url and destination as arguments.
func (t TaskHelper) GitClone(repoURL, targetDir string, isSudo bool) error {
	cmd := "git"
	args := []string{"clone", repoURL, targetDir}
	if _, err := FRunCommand(cmd, args, isSudo); err != nil {
		return fmt.Errorf("failed to clone git repository: %v", err)
	}
	return nil
}

// ExtractTar uses tar xzf with -C flag for destination.
func (t TaskHelper) ExtractTar(file, path string, isSudo bool) error {
	cmd := "tar"
	args := []string{"xzf", file, "-C", path}
	if _, err := FRunCommand(cmd, args, isSudo); err != nil {
		return fmt.Errorf("failed to extract tar.gz file: %v", err)
	}
	return nil
}

// Move executes mv command.
func (t TaskHelper) Move(src, dst string, isSudo bool) error {
	cmd := "mv"
	args := []string{src, dst}
	if _, err := FRunCommand(cmd, args, isSudo); err != nil {
		return fmt.Errorf("failed to move: %v", err)
	}
	return nil
}

// Download makes use of curl with -L and -o flags.
func (t TaskHelper) Download(url, path string, isSudo bool) error {
	cmd := "curl"
	args := []string{"-L", url, "-o", path}
	if _, err := FRunCommand(cmd, args, isSudo); err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	return nil
}

// PromptAction as name implplies, prompts for user input (y/n)
// and passes the result into action function.
// The result is evaluated as y = true, n = false
func (t TaskHelper) PromptAction(msg string, action func(bool) error) error {
	isYes := FPrompt(msg)
	if err := action(isYes); err != nil {
		return err
	}
	return nil
}

func (t TaskHelper) DeletePath(path string, isSudo bool) error {
	cmd := "rm"
	args := []string{"-rf", path}
	if _, err := FRunCommand(cmd, args, isSudo); err != nil {
		return fmt.Errorf("failed to delete the path")
	}

	return nil
}

// UpdatePermission executed chmod with perm string such as u+g or 0644.
func (t TaskHelper) UpdatePermission(path, perm string, useSudo bool) error {
	cmd := "chmod"
	args := []string{perm, path}
	if _, err := FRunCommand(cmd, args, useSudo); err != nil {
		return fmt.Errorf("failed to update permission: %v", err)
	}
	return nil
}

// UpdatePermissionRecursively similar to UpdatePermission but with --recursive flag
func (t TaskHelper) UpdatePermissionRecursively(path, perm string, useSudo bool) error {
	cmd := "chmod"
	args := []string{"--recursive", perm, path}
	if _, err := FRunCommand(cmd, args, useSudo); err != nil {
		return fmt.Errorf("failed to update permission: %v", err)
	}
	return nil
}

// UpdateOwnership uses chown and argument similar to system-like user:group or just user.
func (t TaskHelper) UpdateOwnership(path, owner string, useSudo bool) error {
	cmd := "chown"
	args := []string{owner, path}
	if _, err := FRunCommand(cmd, args, useSudo); err != nil {
		return fmt.Errorf("failed to update ownership: %v", err)
	}
	return nil
}

// UpdateOwnershipRecursively similar to UpdateOwnership but with --recursive flag.
func (t TaskHelper) UpdateOwnershipRecursively(path, owner string, useSudo bool) error {
	cmd := "chown"
	args := []string{"--recursive", owner, path}
	if _, err := FRunCommand(cmd, args, useSudo); err != nil {
		return fmt.Errorf("failed to update ownership: %v", err)
	}
	return nil
}

// AppendContent appends content to a file,
// in CREATE/APPEND|WRONLY mode.
func (t TaskHelper) AppendContent(file, content string) error {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open a file %s: %v", file, err)
	}
	defer f.Close()
	if _, err = f.WriteString(content); err != nil {
		return fmt.Errorf("failed to append content: %v", err)
	}
	return nil
}

type ValidationHelper struct{}

func (v ValidationHelper) ValidateBaseTask(t BaseTask, config any) error {
	if config == nil {
		return fmt.Errorf("base validation failed, empty config value")
	}
	if t.Config == nil {
		return fmt.Errorf("validation failed, config is empty")
	}

	if !FReflectEq(t.Config, config) {
		return fmt.Errorf("validation failed, invalid config")
	}

	return nil
}

func (v ValidationHelper) ValidatePath(path string, checkIfDir bool) error {
	if len(path) == 0 {
		return fmt.Errorf("validation failed, empty path value")
	}

	info, err := os.Stat(path)
	if err != nil {
		slog.Error(err.Error())
		return fmt.Errorf("validation failed, encountered an error")
	}
	if checkIfDir && !info.IsDir() {
		return fmt.Errorf("validation failed, not a directory")
	}
	return nil
}

func (v ValidationHelper) ValidateURL(input string) error {
	_, err := url.ParseRequestURI(input)
	return err
}
