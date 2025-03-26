package main

import (
	"os"
	"path/filepath"
	"strings"
)

var globalConfig *Config

type Config struct {
	Envs         map[string]string `yaml:"envs"`
	Destinations map[string]string `yaml:"destinations"`
	URLs         map[string]string `yaml:"urls"`
	Opts         map[string]string `yaml:"opts"`
}

func NewConfig() *Config {
	globalConfig = &Config{
		Envs:         make(map[string]string),
		Destinations: make(map[string]string),
		URLs:         make(map[string]string),
		Opts:         make(map[string]string),
	}
	return globalConfig
}

func GetConfig() *Config {
	if globalConfig == nil {
		config := NewConfig()
		config.FillConfig()
		globalConfig = config
	}
	return globalConfig
}

func (c *Config) FillConfig() {
	c.Envs["HOME"] = os.Getenv("HOME")
	c.Envs["SHELL"] = os.Getenv("SHELL")

	c.Destinations["oh-my-zsh"] = filepath.Join(c.Envs["HOME"], ".oh-my-zsh")
	c.Destinations["zshrc"] = filepath.Join(c.Envs["HOME"], ".zshrc")
	c.URLs["oh-my-zsh/install.sh"] = "https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh"

	c.Destinations["nvim-lspconfig"] = filepath.Join(c.Envs["HOME"], ".config/nvim/pack/nvim/start/nvim-lspconfig")

	c.Destinations["nvim"] = "/opt/nvim"
	c.URLs["nvim/tar.gz"] = "https://github.com/neovim/neovim/releases/download/v0.10.4/nvim-linux-x86_64.tar.gz"

	c.Destinations["nvim/config"] = filepath.Join(c.Envs["HOME"], ".config")                          // Usually under /home/user/.config
	c.Destinations["nvim/config/subpaths"] = strings.Join([]string{"nvim/init.lua", "nvim/lua"}, "|") // git/*nvim/init.lua* will be moved to /home/user/.config/*nvim/init.lua* and so on
	c.URLs["nvim/config"] = "https://github.com/AlexKhomych/neovim-dot.git"

	c.Destinations["go"] = "/opt/go"
	c.URLs["go/tar.gz"] = "https://go.dev/dl/go1.24.1.linux-amd64.tar.gz"

	c.URLs["ts/install.sh"] = "https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.2/install.sh"
	c.Opts["TS_VERSION"] = "22.14.0"
}
