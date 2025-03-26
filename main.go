package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func main() {
	// Load config
	// Validate
	// Execute

	clean, tempDir := MCreateTempDir()
	defer clean()

	packages := map[string]string{
		"htop":            "",
		"curl":            "",
		"zsh":             "",
		"git":             "",
		"vim":             "",
		"build-essential": "",
	}

	packagesExtra := map[string]string{
		"ripgrep": "https://github.com/BurntSushi/ripgrep/releases/download/14.1.0/ripgrep_14.1.0-1_amd64.deb",
	}

	MInstallPackages(packages, tempDir)
	MInstallPackages(packagesExtra, tempDir)

	MInstallLSP()

	MSetupOhMyZsh(tempDir)

	MInstallNeovim(tempDir)

	MConfigNeovim(tempDir)

	MSetupGolang(tempDir)

	MSetupTypescript(tempDir)

	MNote()
}

func MCreateTempDir() (func(), string) {
	tempDir, err := os.MkdirTemp("", "nvim-auto")
	if err != nil {
		os.Exit(1)
	}
	return func() {
		os.RemoveAll(tempDir)
	}, tempDir
}

func MInstallPackages(packages map[string]string, tempDir string) {
	for pkgName, pkgURL := range packages {
		isInstalled, err := IsPackageInstalled(pkgName, false)
		if err != nil {
			log.Fatalf("%v", err)
		}
		pkg := &Package{
			Name: pkgName,
			URL:  pkgURL,
			Path: filepath.Join(tempDir, fmt.Sprintf("%s.deb", pkgName)),
		}
		if isInstalled {
			promptAsk := fmt.Sprintf("Package '%s' is already installed. Would you like to install/update it? (y/n): ", pkgName)
			if Prompt(promptAsk) {
				MInstallPackage(pkg)
			}
		} else {
			MInstallPackage(pkg)
		}
	}
}

func MInstallPackage(pkg *Package) {
	if len(pkg.URL) > 0 {
		if err := Download(pkg.URL, pkg.Path, false); err != nil {
			log.Fatalf("%v", err)
		}
	}
	if err := InstallPackage(pkg, true); err != nil {
		log.Fatalf("%v", err)
	}
}

func AHandleDirectory(dstPath string, isPrivileged bool) bool {
	isEscape := false
	isEmpty, err := IsEmpty(dstPath)
	if err != nil {
		log.Fatalf("failed to check if directory is empty: %v", err)
	}
	if !isEmpty {
		isOverride := Prompt(fmt.Sprintf("Directory '%s' is not empty. Do you want to overwrite it? (y/n): ", dstPath))
		isEscape = !isOverride
		if isOverride && isPrivileged {
			if err := DeletePathPrivileged(dstPath); err != nil {
				log.Fatalf("failed to delete path: %v", err)
			}
		} else if isOverride {
			if err := DeletePath(dstPath); err != nil {
				log.Fatalf("failed to delete path: %v", err)
			}
		}
	}
	return isEscape
}

func MInstallLSP() {
	dstPath := GetConfig().Destinations["nvim-lspconfig"]
	if isEscape := AHandleDirectory(dstPath, false); isEscape {
		return
	}

	if err := CloneGitRepo("https://github.com/neovim/nvim-lspconfig", dstPath, false); err != nil {
		log.Fatalf("failed to clone git repository: %v", err)
	}
}

func MSetupOhMyZsh(tmpDir string) {
	dstPath := GetConfig().Destinations["oh-my-zsh"]
	if isEscape := AHandleDirectory(dstPath, false); isEscape {
		return
	}

	installScript := filepath.Join(tmpDir, "install.sh")
	installScriptURL := GetConfig().URLs["oh-my-zsh/install.sh"]
	shell := GetConfig().Envs["SHELL"]
	user, err := user.Current()
	if err != nil {
		log.Fatalf("failed to get current user: %v", err)
	}

	if err := Download(installScriptURL, installScript, false); err != nil {
		log.Fatalf("%v", err)
	}
	if _, err := RunCommand("/bin/sh", []string{installScript}, false); err != nil {
		log.Fatalf("%v", err)
	}

	if shell != "/bin/zsh" {
		if _, err := RunCommand("/usr/bin/chsh", []string{user.Username, "-s", "/bin/zsh"}, true); err != nil {
			log.Fatalf("failed to change shell: %v", err)
		}
	}
}

func MInstallNeovim(tmpDir string) {
	dstPath := GetConfig().Destinations["nvim"]
	nvimURL := GetConfig().URLs["nvim/tar.gz"]
	nvimTmp := filepath.Join(tmpDir, "nvim-linux-x86_64.tar.gz")
	nvimTar := filepath.Join(tmpDir, "nvim-linux-x86_64")
	zshrc := GetConfig().Destinations["zshrc"]

	if isEscape := AHandleDirectory(dstPath, true); isEscape {
		return
	}

	if err := Download(nvimURL, nvimTmp, false); err != nil {
		log.Fatalf("%v", err)
	}
	if err := ExtractTar(nvimTmp, tmpDir, false); err != nil {
		log.Fatalf("%v", err)
	}
	if err := Move(nvimTar, dstPath, true); err != nil {
		log.Fatalf("%v", err)
	}
	pathContent := "PATH=$PATH:/opt/nvim/bin\n"
	if err := AppendContent(zshrc, pathContent); err != nil {
		log.Fatalf("%v", err)
	}
}

func MConfigNeovim(tmpDir string) {
	dstPath := GetConfig().Destinations["nvim/config"]
	configURL := GetConfig().URLs["nvim/config"]
	configTmp := filepath.Join(tmpDir, "neovim-dot")
	configSubpaths := strings.Split(GetConfig().Destinations["nvim/config/subpaths"], "|")

	if err := CloneGitRepo(configURL, configTmp, false); err != nil {
		log.Fatalf("%v", err)
	}

	for _, subpath := range configSubpaths {
		src := filepath.Join(configTmp, subpath)
		dst := filepath.Join(dstPath, subpath)
		if err := Move(src, dst, false); err != nil {
			log.Fatalf("%v", err)
		}
	}
}

func MSetupGolang(tmpDir string) {
	dstPath := GetConfig().Destinations["go"]
	goURL := GetConfig().URLs["go/tar.gz"]
	goTmp := filepath.Join(tmpDir, "go1.24.1.linux-amd64.tar.gz")
	goTar := filepath.Join(tmpDir, "go")
	zshrc := GetConfig().Destinations["zshrc"]

	if isEscape := AHandleDirectory(dstPath, true); isEscape {
		return
	}

	if err := Download(goURL, goTmp, false); err != nil {
		log.Fatalf("%v", err)
	}
	if err := ExtractTar(goTmp, tmpDir, false); err != nil {
		log.Fatalf("%v", err)
	}
	if err := Move(goTar, dstPath, true); err != nil {
		log.Fatalf("%v", err)
	}
	if _, err := RunCommand(fmt.Sprintf("%s/bin/go", dstPath), []string{"install", "golang.org/x/tools/gopls@latest"}, false); err != nil {
		log.Fatalf("%v", err)
	}
	pathContent := fmt.Sprintf("PATH=$PATH:/opt/go/bin:%s/go/bin\n", GetConfig().Envs["HOME"])
	if err := AppendContent(zshrc, pathContent); err != nil {
		log.Fatalf("%v", err)
	}
}

func MSetupTypescript(tmpDir string) {
	tsURL := GetConfig().URLs["ts/install.sh"]
	tsVER := GetConfig().Opts["TS_VERSION"]
	envHOME := GetConfig().Envs["HOME"]
	installScript := filepath.Join(tmpDir, "ts_install.sh")
	dstPath := fmt.Sprintf("%s/.nvm/versions/node/v%s", envHOME, tsVER)
	zshrc := GetConfig().Destinations["zshrc"]

	if isEscape := AHandleDirectory(dstPath, false); isEscape {
		return
	}

	if err := Download(tsURL, installScript, false); err != nil {
		log.Fatalf("%v", err)
	}
	if err := UpdatePermission(installScript, "u+x", false); err != nil {
		log.Fatalf("%v", err)
	}
	if _, err := RunCommand("/bin/zsh", []string{"-c", installScript}, false); err != nil {
		log.Fatalf("%v", err)
	}
	if _, err := RunCommand("/bin/zsh", []string{"-c", fmt.Sprintf("source %s/.nvm/nvm.sh && nvm install %s", envHOME, tsVER)}, false); err != nil {
		log.Fatalf("%v", err)
	}
	if _, err := RunCommand("/bin/zsh", []string{"-c", fmt.Sprintf("source %s/.nvm/nvm.sh && %s/.nvm/versions/node/v%s/bin/npm install -g typescript-language-server typescript", envHOME, envHOME, tsVER)}, false); err != nil {
		log.Fatalf("%v", err)
	}
	pathContent := "export NVM_DIR=\"$HOME/.nvm\"\n[ -s \"$NVM_DIR/nvm.sh\" ] && \\. \"$NVM_DIR/nvm.sh\"  # This loads nvm\n[ -s \"$NVM_DIR/bash_completion\" ] && \\. \"$NVM_DIR/bash_completion\"  # This loads nvm bash_completion\n"
	if err := AppendContent(zshrc, pathContent); err != nil {
		log.Fatalf("%v", err)
	}
}

func MNote() {
	if _, err := RunCommand("echo", []string{"\n into 'zsh'\nby running 'zsh' command in your terminal."}, false); err != nil {
		log.Fatalf("failed to print note message: %v", err)
	}
}
