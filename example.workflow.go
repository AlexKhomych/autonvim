package main

import (
	"os"
	"path/filepath"
)

const (
	HomePath   string = "/home/alex/"
	ShrcPath   string = "/home/alex/.zshrc"

	OhMyZshURL string = "https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh"
	GolangURL  string = "https://go.dev/dl/go1.24.1.linux-amd64.tar.gz"
	NvmURL     string = "https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.2/install.sh"
	NvimLSPURL string = "https://github.com/neovim/nvim-lspconfig"
	NvimDotURL string = "https://github.com/AlexKhomych/neovim-dot.git"
)

func ExampleRun() {
	clear, tmpDir := CreateTempDir()
	defer clear()

	check(Package(tmpDir))
	check(OhMyZsh(tmpDir))
	check(Neovim(tmpDir))
	check(NeovimLSP())
	check(Golang(tmpDir))
	check(Typescript(tmpDir))
	check(DotConfig(tmpDir))
}

func Package(tmpDir string) error {
	downloadConfig := DownloadConfig{
		url: "https://github.com/BurntSushi/ripgrep/releases/download/14.1.0/ripgrep_14.1.0-1_amd64.deb",
		path: Path{
			path:    tmpDir,
			subpath: "ripgrep_14.1.0-1_amd64.deb",
		},
	}
	downloadTask := DownloadTask{
		BaseTask: BaseTask{
			Name:   "DownloadTask",
			Config: downloadConfig,
		},
	}

	check(downloadTask.Validate())
	check(downloadTask.Run())

	config := InstallPackageConfig{
		name:   "ripgrep",
		path:   downloadConfig.path.Join(),
		isSudo: true,
	}
	task := InstallPackageTask{
		BaseTask: BaseTask{
			Config: config,
			Name:   "InstallPackage",
		},
	}
	check(task.Validate())
	check(task.Run())

	return nil
}

func Neovim(tmpDir string) error {
	downloadConfig := DownloadConfig{
		path: Path{
			path:    tmpDir,
			subpath: "nvim-linux-x86_64.tar.gz",
		},
		url:    "https://github.com/neovim/neovim/releases/download/v0.10.4/nvim-linux-x86_64.tar.gz",
		isSudo: false,
	}

	downloadTask := DownloadTask{
		BaseTask: BaseTask{
			Name:   "DownloadTask",
			Config: downloadConfig,
		},
	}

	check(downloadTask.Validate())
	check(downloadTask.Run())

	installConfig := InstallNeovimConfig{
		path: Path{
			path:    filepath.Join(HomePath, ".local/share"),
			subpath: "nvim-linux-x86_64",
		},
		shrc: ShrcConfig{
			path:    ShrcPath,
			content: "export PATH=$PATH:/home/alex/.local/share/nvim-linux-x86_64/bin\n",
		},
		tarPath: downloadConfig.path.Join(),
		isSudo:  false,
	}

	installTask := InstallNeovimTask{
		BaseTask: BaseTask{
			Name:   "InstallNeovimTask",
			Config: installConfig,
		},
	}

	opts := OverwriteOptions{
		path:   installConfig.path,
		isSudo: false,
	}

	isSkip := HandleOverwrite(opts)
	if isSkip {
		return nil
	}

	check(installTask.Validate())
	check(installTask.Run())

	return nil
}

func DotConfig(tmpDir string) error {
	tmpDir = filepath.Join(tmpDir, "neovim-dot")
	if err := FCreateDir(tmpDir); err != nil {
		return err
	}
	config := NeovimDotConfig{
		path: Path{
			path:    filepath.Join(HomePath, ".config"),
			subpath: "nvim",
		},
		tmpDir:   tmpDir,
		url:      NvimDotURL,
		subpaths: []string{"init.lua", "lua"},
		isSudo:   false,
	}

	task := NeovimDotTask{
		BaseTask: BaseTask{
			Name:   "NeovimDotTask",
			Config: config,
		},
	}

	for _, subpath := range config.subpaths {
		path := Path{
			path:    config.path.Join(),
			subpath: subpath,
		}
		opts := OverwriteOptions{
			path:   path,
			isSudo: false,
		}

		isSkip := HandleOverwrite(opts)
		if isSkip {
			continue
		}
	}

	check(task.Validate())
	check(task.Run())

	return nil
}

func NeovimLSP() error {
	config := NeovimLSPConfig{
		path: Path{
			path:    filepath.Join(HomePath, ".config/nvim/pack/nvim/start"),
			subpath: "nvim-lspconfig",
		},
		url:    NvimLSPURL,
		isSudo: false,
	}

	task := NeovimLSPTask{
		BaseTask: BaseTask{
			Name:   "NeovimLSPTask",
			Config: config,
		},
	}

	opts := OverwriteOptions{
		path:   config.path,
		isSudo: false,
	}

	isSkip := HandleOverwrite(opts)
	if isSkip {
		return nil
	}

	check(task.Validate())
	check(task.Run())

	return nil
}

func OhMyZsh(tmpDir string) error {
	config := OhMyZshConfig{
		tmpDir: tmpDir,
		path: Path{
			path:    HomePath,
			subpath: ".oh-my-zsh",
		},
		username: "alex",
		url:      OhMyZshURL,
		isSudo:   false,
	}

	task := OhMyZshTask{
		BaseTask: BaseTask{
			Name:   "OhMyZshTask",
			Config: config,
		},
	}

	opts := OverwriteOptions{
		path:   config.path,
		isSudo: false,
	}

	isSkip := HandleOverwrite(opts)
	if isSkip {
		return nil
	}

	check(task.Validate())
	check(task.Run())

	return nil
}

func Golang(tmpDir string) error {
	downloadConfig := DownloadConfig{
		path: Path{
			path:    tmpDir,
			subpath: "go1.24.1.linux-amd64.tar.gz",
		},
		url:    GolangURL,
		isSudo: false,
	}

	downloadTask := DownloadTask{
		BaseTask: BaseTask{
			Name:   "DownloadTask",
			Config: downloadConfig,
		},
	}

	check(downloadTask.Validate())
	check(downloadTask.Run())

	installConfig := InstallGolangConfig{
		path: Path{
			path:    filepath.Join(HomePath, ".local/share"),
			subpath: "go",
		},
		tarPath: filepath.Join(tmpDir, "go1.24.1.linux-amd64.tar.gz"),
		shrc: ShrcConfig{
			path:    ShrcPath,
			content: "export PATH=$PATH:/home/alex/.local/share/go/bin:/home/alex/go/bin\n",
		},
		isSudo: false,
	}

	installTask := InstallGolangTask{
		BaseTask: BaseTask{
			Name:   "InstallGolangTask",
			Config: installConfig,
		},
	}

	opts := OverwriteOptions{
		path:   installConfig.path,
		isSudo: false,
	}

	isSkip := HandleOverwrite(opts)
	if isSkip {
		return nil
	}

	check(installTask.Validate())
	check(installTask.Run())

	return nil
}

func Typescript(tmpDir string) error {
	downloadConfig := DownloadConfig{
		path: Path{
			path:    tmpDir,
			subpath: "nvm_install.sh",
		},
		url:    NvmURL,
		isSudo: false,
	}

	downloadTask := DownloadTask{
		BaseTask: BaseTask{
			Name:   "DownloadTask",
			Config: downloadConfig,
		},
	}

	check(downloadTask.Validate())
	check(downloadTask.Run())

	installConfig := InstallTypescriptConfig{
		version:        "22.14.0",
		installNVMPath: downloadConfig.path.Join(),
		homePath:       HomePath,
		shrc: ShrcConfig{
			path:    ShrcPath,
			content: "export NVM_DIR=\"$HOME/.nvm\"\n[ -s \"$NVM_DIR/nvm.sh\" ] && \\. \"$NVM_DIR/nvm.sh\"  # This loads nvm\n[ -s \"$NVM_DIR/bash_completion\" ] && \\. \"$NVM_DIR/bash_completion\"  # This loads nvm bash_completion\n",
		},
		isSudo: false,
	}

	installTask := InstallTypescriptTask{
		BaseTask: BaseTask{
			Name:   "InstallTypescriptTask",
			Config: installConfig,
		},
	}

	opts := OverwriteOptions{
		path: Path{
			path:    HomePath,
			subpath: filepath.Join(".nvm/versions/node/v" + installConfig.version),
		},
		isSudo: false,
	}
	isSkip := HandleOverwrite(opts)
	if isSkip {
		return nil
	}

	check(installTask.Validate())
	check(installTask.Run())

	return nil
}

type OverwriteOptions struct {
	path   Path
	isSudo bool
}

func HandleOverwrite(o OverwriteOptions) bool {
	var isSkip bool
	action := func(isYes bool) error {
		isSkip = !isYes
		if isSkip {
			return nil
		}
		deleteConfig := DeletePathConfig{
			path:   o.path.Join(),
			isSudo: o.isSudo,
		}
		deleteTask := DeletePathTask{
			BaseTask: BaseTask{
				Name:   "DeletePathTask",
				Config: deleteConfig,
			},
		}
		check(deleteTask.Validate())
		check(deleteTask.Run())

		return nil
	}

	config := DirectoryPromptConfig{
		path:   o.path.Join(),
		isSudo: o.isSudo,
		action: action,
	}
	task := DirectoryPromptTask{
		BaseTask: BaseTask{
			Name:   "DirectoryPromptTask",
			Config: config,
		},
	}

	if err := FCreateDir(o.path.path); err != nil {
		return false
	}

	check(task.Validate())
	check(task.Run())

	return isSkip
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func CreateTempDir() (func(), string) {
	tmpDir, err := os.MkdirTemp("", "autonvim")
	if err != nil {
		panic(err)
	}
	return func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			panic(err)
		}
	}, tmpDir
}
