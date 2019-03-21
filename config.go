package chart

import (
	"os"
	"path/filepath"

	"github.com/jinzhu/configor"

	"github.com/ecletus-pkg/admin"

	"github.com/moisespsena/go-error-wrap"
)

const (
	DEFAULT_CONFIG_DIR = "config/charts"
	URI                = "/reports.json"
)

type ChartConfig struct {
	ResourceName string
	AdminName    string
	Label        string
	Description  string
	Site         bool
	Factory      string
	Options      map[string]interface{}
	// JSON Value
	DataSet string
	ID      string
}

type Config struct {
	Uri struct {
		SiteUri  string
		AdminUri string
	}
	Charts map[string]*ChartConfig
}

func (c *Config) LoadDefaults() {
	if c.Uri.SiteUri == "" {
		c.Uri.SiteUri = URI
	}
	if c.Uri.AdminUri == "" {
		c.Uri.AdminUri = URI
	}
}

func (c *Config) Merge(cfg *Config) {
	if c.Charts == nil {
		c.Charts = map[string]*ChartConfig{}
	}
	for k, chart := range cfg.Charts {
		c.Charts[k] = chart
	}
}

func LoadConfigFile(pth string) (config *Config, err error) {
	config = &Config{}
	if err = configor.Load(config, pth); err != nil {
		return nil, errwrap.Wrap(err, "Load config file %q", pth)
	}
	for id, cfg := range config.Charts {
		cfg.ID = id
		if cfg.Options == nil {
			cfg.Options = map[string]interface{}{}
		}
		if cfg.ResourceName != "" && cfg.AdminName == "" {
			cfg.AdminName = admin_plugin.DEFAULT_ADMIN
		}
	}
	return
}

func LoadConfigDir(configDir string) (config *Config, err error) {
	config = &Config{}

	if configDir == "" {
		configDir = DEFAULT_CONFIG_DIR
	}

	err = filepath.Walk(configDir, func(pth string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name()[0] == '.' {
				return filepath.SkipDir
			}
			return nil
		}
		if info.Name()[0] == '.' {
			return nil
		}
		if filepath.Ext(info.Name()) == ".yaml" {
			cfg, err := LoadConfigFile(pth)
			if err != nil {
				return err
			}
			config.Merge(cfg)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return config, nil
}
