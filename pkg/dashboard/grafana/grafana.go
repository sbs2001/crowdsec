package grafana

import (
	"fmt"
	"io/ioutil"

	"github.com/crowdsecurity/crowdsec/pkg/csconfig"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Grafana struct {
	Container *Container
	Config    *Config
}

var (
	metabaseDefaultUser     = "crowdsec@crowdsec.net"
	metabaseDefaultPassword = "!!Cr0wdS3c_M3t4b4s3??"
	containerName           = "/crowdsec-grafana"
	metabaseImage           = "grafana/grafana"

	metabaseSQLiteDBURL = "https://crowdsec-statics-assets.s3-eu-west-1.amazonaws.com/grafana_dashboard.zip"
)

type Config struct {
	Database   *csconfig.DatabaseCfg `yaml:"database"`
	ListenAddr string                `yaml:"listen_addr"`
	ListenPort string                `yaml:"listen_port"`
	ListenURL  string                `yaml:"listen_url"`
	Username   string                `yaml:"username"`
	Password   string                `yaml:"password"`
}

func NewGrafana(configPath string) (*Grafana, error) {
	g := &Grafana{}
	if err := g.LoadConfig(configPath); err != nil {
		return g, err
	}
	if err := g.Init(); err != nil {
		return g, err
	}
}

func (g *Grafana) LoadConfig(configPath string) error {
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	config := &Config{}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return err
	}
	if config.Username == "" {
		return fmt.Errorf("'username' not found in configuration file '%s'", configPath)
	}

	if config.Password == "" {
		return fmt.Errorf("'password' not found in configuration file '%s'", configPath)
	}

	if config.ListenURL == "" {
		return fmt.Errorf("'listen_url' not found in configuration file '%s'", configPath)
	}

	g.Config = config

	if err := g.Init(); err != nil {
		return err
	}

	return nil

}

func (g *Grafana) Init() error {
	var err error
	var DBConnectionURI string
	var remoteDBAddr string

	g.Container, err = NewContainer(g.Config.ListenAddr, g.Config.ListenPort, containerName, metabaseImage)
	if err != nil {
		return errors.Wrap(err, "container init")
	}

	return nil
}