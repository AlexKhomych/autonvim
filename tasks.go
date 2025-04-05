package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

type BaseTask struct {
	Name   string
	Config any
}

func (t BaseTask) Initialize() {
	if len(t.Name) == 0 {
		t.Name = FReflectName(t)
		slog.Warn("task name is empty, using struct name instead: %s", "task_name", t.Name)
	}
}

type Task interface {
	Validate() error
	Run() error
	Update() error
}

type InstallPackageConfig struct {
	name   string
	path   string
	isSudo bool
}

type InstallPackageTask struct {
	BaseTask
	th TaskHelper
	vh ValidationHelper
}

func (t *InstallPackageTask) Validate() error {
	if err := t.vh.ValidateBaseTask(t.BaseTask, t.Config); err != nil {
		return FPrefixError(t.Name, err.Error())
	}

	cfg, _ := t.Config.(InstallPackageConfig)

	if len(cfg.name) == 0 {
		return fmt.Errorf("package name cannot be empty.")
	}
	if err := t.vh.ValidatePath(cfg.path, false); len(cfg.path) > 0 && err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	return nil
}

func (t *InstallPackageTask) Run() error {
	cfg, _ := t.Config.(InstallPackageConfig)
	isInstalled, err := t.th.IsPackageInstalled(cfg.name, cfg.isSudo)
	if err != nil {
		return FPrefixError(t.Name, "failed to check package installation")
	}
	if isInstalled {
		promptAsk := fmt.Sprintf("Package '%s' is already installed. Would you like to install/update it? (y/n): ", cfg.name)
		if !FPrompt(promptAsk) {
			return nil
		}
	}
	var identifier string
	if len(cfg.path) > 0 {
		identifier = cfg.path
	} else {
		identifier = cfg.name
	}

	cmd := "apt"
	args := []string{"install", "--yes", identifier}

	if _, err := FRunCommand(cmd, args, cfg.isSudo); err != nil {
		return FPrefixError(t.Name, "failed to install the package")
	}

	return nil
}

type NeovimLSPConfig struct {
	path   Path
	url    string
	isSudo bool
}

type NeovimLSPTask struct {
	BaseTask
	th TaskHelper
	vh ValidationHelper
}

func (t *NeovimLSPTask) Validate() error {
	if err := t.vh.ValidateBaseTask(t.BaseTask, t.Config); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	cfg, _ := t.Config.(NeovimLSPConfig)

	if err := t.vh.ValidatePath(cfg.path.path, true); err != nil {
		return fmt.Errorf("%s: validation failed, %s", t.Name, err.Error())
	}
	if err := t.vh.ValidateURL(cfg.url); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	return nil
}

func (t *NeovimLSPTask) Run() error {
	cfg, _ := t.Config.(NeovimLSPConfig)
	dstPath := cfg.path.Join()

	if err := t.th.GitClone(cfg.url, dstPath, cfg.isSudo); err != nil {
		return fmt.Errorf("failed to clone git repository: %v", err)
	}
	return nil
}

type OhMyZshConfig struct {
	tmpDir   string
	path     Path
	username string
	url      string
	isSudo   bool
}

type OhMyZshTask struct {
	BaseTask
	th TaskHelper
	vh ValidationHelper
}

func (t *OhMyZshTask) Validate() error {
	if err := t.vh.ValidateBaseTask(t.BaseTask, t.Config); err != nil {
		return FPrefixError(t.Name, err.Error())
	}

	cfg, _ := t.Config.(OhMyZshConfig)

	if len(cfg.username) == 0 {
		return FPrefixError(t.Name, "empty username value")
	}
	if err := t.vh.ValidatePath(cfg.tmpDir, true); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	if err := t.vh.ValidateURL(cfg.url); err != nil {
		return FPrefixError(t.Name, err.Error())
	}

	return nil
}

func (t *OhMyZshTask) Run() error {
	cfg, _ := t.Config.(OhMyZshConfig)
	scriptPath := filepath.Join(cfg.tmpDir, "install.sh")

	if err := t.th.Download(cfg.url, scriptPath, cfg.isSudo); err != nil {
		return fmt.Errorf("%v", err)
	}
	if _, err := FRunCommand("/bin/sh", []string{scriptPath}, false); err != nil {
		return err
	}

	if _, err := FRunCommand("/usr/bin/chsh", []string{cfg.username, "-s", "/bin/zsh"}, true); err != nil {
		return fmt.Errorf("failed to change shell: %v", err)
	}

	return nil
}

type InstallNeovimConfig struct {
	path    Path
	shrc    ShrcConfig
	tarPath string
	isSudo  bool
}

type InstallNeovimTask struct {
	BaseTask
	th TaskHelper
	vh ValidationHelper
}

func (t *InstallNeovimTask) Validate() error {
	if err := t.vh.ValidateBaseTask(t.BaseTask, t.Config); err != nil {
		return FPrefixError(t.Name, err.Error())
	}

	cfg, _ := t.Config.(InstallNeovimConfig)
	if err := t.vh.ValidatePath(cfg.path.path, true); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	if err := t.vh.ValidatePath(cfg.tarPath, false); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	if err := t.vh.ValidatePath(cfg.shrc.path, false); err != nil {
		return FPrefixError(t.Name, err.Error())
	}

	return nil
}

func (t *InstallNeovimTask) Run() error {
	cfg, _ := t.Config.(InstallNeovimConfig)
	dstPath := cfg.path.path

	if err := t.th.ExtractTar(cfg.tarPath, dstPath, cfg.isSudo); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	if err := t.th.AppendContent(cfg.shrc.path, cfg.shrc.content); err != nil {
		return FPrefixError(t.Name, err.Error())
	}

	return nil
}

type NeovimDotConfig struct {
	path     Path
	url      string
	tmpDir   string
	subpaths []string
	isSudo   bool
}

type NeovimDotTask struct {
	BaseTask
	th TaskHelper
	vh ValidationHelper
}

func (t *NeovimDotTask) Validate() error {
	if err := t.vh.ValidateBaseTask(t.BaseTask, t.Config); err != nil {
		return FPrefixError(t.Name, err.Error())
	}

	cfg, _ := t.Config.(NeovimDotConfig)

	if err := t.vh.ValidatePath(cfg.path.path, true); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	if err := t.vh.ValidatePath(cfg.tmpDir, true); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	if err := t.vh.ValidateURL(cfg.url); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	return nil
}

func (t *NeovimDotTask) Run() error {
	cfg, _ := t.Config.(NeovimDotConfig)
	dstPath := cfg.path.Join()

	if err := t.th.GitClone(cfg.url, cfg.tmpDir, cfg.isSudo); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
  if err := FCreateDir(dstPath); err != nil {
    return FPrefixError(t.Name, err.Error())
  }

	for _, subpath := range cfg.subpaths {
		src := filepath.Join(cfg.tmpDir, "nvim", subpath)
		dst := filepath.Join(dstPath, subpath)
		if err := t.th.Move(src, dst, cfg.isSudo); err != nil {
			return FPrefixError(t.Name, err.Error())
		}
	}
	return nil
}

type DownloadConfig struct {
	path   Path
	url    string
	isSudo bool
}

type DownloadTask struct {
	BaseTask
	th TaskHelper
	vh ValidationHelper
}

func (t *DownloadTask) Validate() error {
	if err := t.vh.ValidateBaseTask(t.BaseTask, t.Config); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	cfg, _ := t.Config.(DownloadConfig)

	if err := t.vh.ValidateURL(cfg.url); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	if err := t.vh.ValidatePath(cfg.path.path, true); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	if len(cfg.path.subpath) == 0 {
		return FPrefixError(t.Name, "download filename(subpath) is missing")
	}
	return nil
}

func (t *DownloadTask) Run() error {
	cfg, _ := t.Config.(DownloadConfig)

	if err := t.th.Download(cfg.url, cfg.path.Join(), cfg.isSudo); err != nil {
		return err
	}

	return nil
}

type InstallGolangConfig struct {
	path    Path
	shrc    ShrcConfig
	tarPath string
	isSudo  bool
}

type Path struct {
	path    string
	subpath string
}

type ShrcConfig struct {
	path    string
	content string
}

func (p *Path) Join() string {
	return filepath.Join(p.path, p.subpath)
}

type InstallGolangTask struct {
	BaseTask
	th TaskHelper
	vh ValidationHelper
}

func (t InstallGolangTask) Validate() error {
	if err := t.vh.ValidateBaseTask(t.BaseTask, t.Config); err != nil {
		return FPrefixError(t.Name, err.Error())
	}

	cfg, _ := t.Config.(InstallGolangConfig)

	if err := t.vh.ValidatePath(cfg.path.path, true); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	if err := t.vh.ValidatePath(cfg.tarPath, false); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	if err := t.vh.ValidatePath(cfg.shrc.path, false); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	return nil
}

func (t InstallGolangTask) Run() error {
	cfg, _ := t.Config.(InstallGolangConfig)
	dstPath := cfg.path.path

	if err := t.th.ExtractTar(cfg.tarPath, dstPath, cfg.isSudo); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	if _, err := FRunCommand(filepath.Join(dstPath, "go/bin/go"), []string{"install", "golang.org/x/tools/gopls@latest"}, false); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	if err := t.th.AppendContent(cfg.shrc.path, cfg.shrc.content); err != nil {
		return FPrefixError(t.Name, err.Error())
	}

	return nil
}

type InstallTypescriptConfig struct {
	version        string
	installNVMPath string
	homePath       string
	shrc           ShrcConfig
	isSudo         bool
}

type InstallTypescriptTask struct {
	BaseTask
	th TaskHelper
	vh ValidationHelper
}

func (t InstallTypescriptTask) Validate() error {
	if err := t.vh.ValidateBaseTask(t.BaseTask, t.Config); err != nil {
		return FPrefixError(t.Name, err.Error())
	}

	cfg, _ := t.Config.(InstallTypescriptConfig)

	if err := t.vh.ValidatePath(cfg.installNVMPath, false); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	if err := t.vh.ValidatePath(cfg.homePath, false); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	if err := t.vh.ValidatePath(cfg.shrc.path, false); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	if len(cfg.version) == 0 {
		return FPrefixError(t.Name, "no version is specified")
	}

	return nil
}

func (t InstallTypescriptTask) Run() error {
	cfg, _ := t.Config.(InstallTypescriptConfig)

	if err := t.th.UpdatePermission(cfg.installNVMPath, "u+x", cfg.isSudo); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	if _, err := FRunCommand("/bin/zsh", []string{"-c", cfg.installNVMPath}, cfg.isSudo); err != nil {
		return FPrefixError(t.Name, "failed to install nvm")
	}
	if _, err := FRunCommand("/bin/zsh", []string{"-c", fmt.Sprintf("source %s/.nvm/nvm.sh && nvm install %s", cfg.homePath, cfg.version)}, false); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	if _, err := FRunCommand("/bin/zsh", []string{"-c", fmt.Sprintf("source %s/.nvm/nvm.sh && %s/.nvm/versions/node/v%s/bin/npm install -g typescript-language-server typescript", cfg.homePath, cfg.homePath, cfg.version)}, false); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	if err := t.th.AppendContent(cfg.shrc.path, cfg.shrc.content); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	return nil
}

type DeletePathConfig struct {
	path   string
	isSudo bool
}

type DeletePathTask struct {
	BaseTask
	th TaskHelper
	vh ValidationHelper
}

func (t *DeletePathTask) Validate() error {
	if err := t.vh.ValidateBaseTask(t.BaseTask, t.Config); err != nil {
		return FPrefixError(t.Name, err.Error())
	}

	cfg, _ := t.Config.(DeletePathConfig)

	if err := t.vh.ValidatePath(cfg.path, false); err != nil {
		return FPrefixError(t.Name, err.Error())
	}

	return nil
}

func (t DeletePathTask) Run() error {
	cfg, _ := t.Config.(DeletePathConfig)
	if err := t.th.DeletePath(cfg.path, cfg.isSudo); err != nil {
		return FPrefixError(t.Name, err.Error())
	}
	return nil
}

type DirectoryPromptConfig struct {
	path   string
	isSudo bool
	action func(bool) error
}
type DirectoryPromptTask struct {
	BaseTask
	vh ValidationHelper
	th TaskHelper
}

func (t *DirectoryPromptTask) Validate() error {
	if err := t.vh.ValidateBaseTask(t.BaseTask, t.Config); err != nil {
		return FPrefixError(t.Name, err.Error())
	}

	return nil
}

func (t *DirectoryPromptTask) Run() error {
	cfg, _ := t.Config.(DirectoryPromptConfig)
	_, err := os.Stat(cfg.path)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		slog.Error(err.Error(), "task_name", t.Name)
		return err
	}

	ask := fmt.Sprintf("Directory '%s' is not empty. Do you want to overwrite it? (y/n): ", cfg.path)

	t.th.PromptAction(ask, cfg.action)

	return nil
}
