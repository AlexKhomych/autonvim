package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCommand returns command status code and error message.
// Following status codes are possible:
// 0 for successful execution
// -1 for failed execution
// any other non-negative number representing command exit code,
// which is useful for case checks.
func RunCommand(cmd string, args []string, useSudo bool) (int, error) {
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

// IsEmpty checks if path exists and is empty.
// Directory with no files in it is considered to be empty.
// Can return an error.
func IsEmpty(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check path stat: %v", err)
	}
	if info.IsDir() {
		files, err := os.ReadDir(path)
		if err != nil {
			return false, fmt.Errorf("failed to list directory: %v", err)
		}
		return len(files) == 0, nil
	}

	return false, nil
}

// CreateDirectory makes directory with any parent folders that yet to exist.
func CreateDirectory(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	return nil
}

// DeletePath removes specified path, be it a directory or a single file.
func DeletePath(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to delete path: %v", err)
	}
	return nil
}

// Similar to DeletePath, but uses sudo.
func DeletePathPrivileged(path string) error {
	if _, err := RunCommand("rm", []string{"-rf", path}, true); err != nil {
		return fmt.Errorf("failed to delete path: %v", err)
	}
	return nil
}

// Prompt asks a user for an input (y/n)
// and returns boolean representing:
// true if lowercase answer is (y) and false otherwise
func Prompt(ask string) bool {
	var userInput string
	fmt.Println(ask)
	fmt.Scanln(&userInput)
	userInput = strings.ToLower(userInput)
	return userInput == "y"
}

// CloneGitRepo executed git clone command with git url and destination as arguments.
func CloneGitRepo(repoURL, targetDir string, useSudo bool) error {
	if _, err := RunCommand("git", []string{"clone", repoURL, targetDir}, useSudo); err != nil {
		return fmt.Errorf("failed to clone git repository: %v", err)
	}
	return nil
}

// Download makes use of curl with -L and -o flags.
func Download(url, path string, useSudo bool) error {
	if _, err := RunCommand("curl", []string{"-L", url, "-o", path}, useSudo); err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	return nil
}

// Move executes mv command.
func Move(src, dst string, useSudo bool) error {
	if _, err := RunCommand("mv", []string{src, dst}, useSudo); err != nil {
		return fmt.Errorf("failed to move: %v", err)
	}
	return nil
}

// UpdatePermission executed chmod with perm string such as u+g or 0644.
func UpdatePermission(path, perm string, useSudo bool) error {
	if _, err := RunCommand("chmod", []string{perm, path}, useSudo); err != nil {
		return fmt.Errorf("failed to update permission: %v", err)
	}
	return nil
}

// UpdatePermissionRecursively similar to UpdatePermission but with --recursive flag
func UpdatePermissionRecursively(path, perm string, useSudo bool) error {
	if _, err := RunCommand("chmod", []string{"--recursive", perm, path}, useSudo); err != nil {
		return fmt.Errorf("failed to update permission: %v", err)
	}
	return nil
}

// UpdateOwnership uses chown and argument similar to system-like user:group or just user.
func UpdateOwnership(path, owner string, useSudo bool) error {
	if _, err := RunCommand("chown", []string{owner, path}, useSudo); err != nil {
		return fmt.Errorf("failed to update ownership: %v", err)
	}
	return nil
}

// UpdateOwnershipRecursively similar to UpdateOwnership but with --recursive flag.
func UpdateOwnershipRecursively(path, owner string, useSudo bool) error {
	if _, err := RunCommand("chown", []string{"--recursive", owner, path}, useSudo); err != nil {
		return fmt.Errorf("failed to update ownership: %v", err)
	}
	return nil
}

// IsPackageInstalled uses dpkg-query with --status flag to verify if package exists.
// Will print package information if present or unavailability otherwise.
func IsPackageInstalled(pkgName string, useSudo bool) (bool, error) {
	errCode, err := RunCommand("dpkg-query", []string{"--status", pkgName}, useSudo)
	if err != nil && errCode == 1 {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check package presence: %v", err)
	}
	return true, nil
}

// Package struct used to install packages.
// Name - package name like htop, ripgrep
// URL - package url, for https://url/package.deb files(like ripgrep)
// Path - download path for URL-used packages
//
// Use only Name for packages to be installed via apt install name
// Specify URL and Path for packages to be installed from .deb files,
// downloaded from the Internet.
type Package struct {
	Name string
	URL  string
	Path string
}

// InstallPackage uses apt to install packages.
// If pkg.URL is not empty, it will download .deb file 
// to pkg.Path and install it with apt as well.
func InstallPackage(pkg *Package, useSudo bool) error {
	cmd := "apt"
	args := []string{"install", "--yes"}
	if len(pkg.URL) > 0 {
		args = append(args, pkg.Path)
	} else {
		args = append(args, pkg.Name)
	}
	if _, err := RunCommand(cmd, args, useSudo); err != nil {
		return fmt.Errorf("failed to install package: %v", err)
	}
	return nil
}

// ExtractTar uses tar xzf with -C flag for destination.
func ExtractTar(tarFile, dstDir string, useSudo bool) error {
	if _, err := RunCommand("tar", []string{"xzf", tarFile, "-C", dstDir}, useSudo); err != nil {
		return fmt.Errorf("failed to extract tar.gz file: %v", err)
	}
	return nil
}

// AppendContent appends content to a file,
// in APPEND|WRONLY mode.
func AppendContent(file, content string) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open a file %s: %v", file, err)
	}
	defer f.Close()
	if _, err = f.WriteString(content); err != nil {
		return fmt.Errorf("failed to append content: %v", err)
	}
	return nil
}
