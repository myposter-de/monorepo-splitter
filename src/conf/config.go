package conf

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"splitter/version"
	"strings"
)

const (
	configName = "splitter"
)

var extensions = [2]string{"yaml", "yml"}

type AuthFunc func() (http.AuthMethod, error)

type Config struct {
	Root            `yaml:"root"`
	Packages        `yaml:"packages"`
	Actions         []string             `yaml:"actions"`
	VersionValue    version.Version      `yaml:"-"`
	PackageAuthFunc AuthFunc             `yaml:"-"`
	PackageAuth     http.AuthMethod      `yaml:"-"`
	RootAuth        transport.AuthMethod `yaml:"-"`
}

type Root struct {
	Branch string `yaml:"branch"`
	Path   string `yaml:"path"`
	Remote string `yaml:"remote"`
}

type Packages struct {
	Prefix string `yaml:"prefix"`
	Items  []*Pkg `yaml:"items"`
	Branch string `yaml:"branch"`
}

type Pkg struct {
	Remote string `yaml:"remote"`
	Url    string `yaml:"url"`
	Path   string `yaml:"path"`
}

func LoadConfig(name string, authFunc AuthFunc, version version.Version) (*Config, error) {
	if name != "" {
		if strings.HasPrefix(name, "~") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}
			name = filepath.Join(homeDir, name[1:])
		}
		return loadConfig(name, authFunc, version)
	}
	for _, ext := range extensions {
		filename := fmt.Sprintf("%s.%s", configName, ext)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			continue
		}
		return loadConfig(filename, authFunc, version)
	}
	return nil, fmt.Errorf("no suitable conf file found")
}

func loadConfig(filename string, authFunc AuthFunc, version version.Version) (*Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var c Config

	if err = yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	c.PackageAuthFunc = authFunc
	c.VersionValue = version
	current, err := user.Current()
	if err != nil {
		return nil, err
	}
	rootAuth, err := ssh.NewPublicKeysFromFile(
		current.Name,
		filepath.Join(current.HomeDir, ".ssh/id_rsa"),
		"",
	)
	if err != nil {
		return nil, err
	}
	c.RootAuth = rootAuth

	for _, item := range c.Items {
		if item.Path == "" {
			item.Path = item.Remote
		}
	}

	return &c, nil
}
