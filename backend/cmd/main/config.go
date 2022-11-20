package main

import (
	"time"

	"github.com/BurntSushi/toml"
)

type (
	// Duration is a wrapper for time.Duration variable
	Duration struct {
		time.Duration
	}

	// CfgService contains all service settings
	CfgService struct {
		Host     string `toml:"host"`
		Port     int    `toml:"port"`
		Secret   string `toml:"secret"`
		Timeouts struct {
			Write *Duration `toml:"write"`
			Read  *Duration `toml:"read"`
			Idle  *Duration `toml:"idle"`
		} `toml:"timeouts"`
	}

	// CfgDb contains settings for DB connection
	CfgDb struct {
		Host       string    `toml:"host"`
		Port       int       `toml:"port"`
		User       string    `toml:"user"`
		Password   string    `toml:"password"`
		DbName     string    `toml:"dbname"`
		Collection string    `toml:"collection"`
		Timeout    *Duration `toml:"timeout"`
	}

	// CfgFormController - form controller settings
	CfgFormController struct {
		AllowedFormats []string `toml:"allowed_formats"`
		MaxFileSize    int      `toml:"max_file_size"`
		JpegQuality    uint     `toml:"jpeg_quality"`
		ShrinkPhotoTo  uint     `toml:"shrink_photo_to"`
		ThumbnailSize  uint     `toml:"thumbnail_size"`
		AddWatermark   bool     `toml:"add_watermark"`
	}

	// CfgStorage - path and owner for files saving
	CfgStorage struct {
		Path string `toml:"path"`
		Uid  int    `toml:"uid"`
		Gid  int    `toml:"gid"`
	}

	// Config unites all application settings
	Config struct {
		Service        CfgService        `toml:"service"`
		Db             CfgDb             `toml:"db"`
		FormController CfgFormController `toml:"form"`
		Storage        CfgStorage        `toml:"storage"`
		Path           string
	}
)

// UnmarshalText parses time.Duration from a byte-string
func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

// LoadConfig parses a fileName file to Config structure
func LoadConfig(fileName string) *Config {
	var conf = new(Config)
	if _, err := toml.DecodeFile(fileName, &conf); err != nil {
		panic(err)
	}
	conf.Path = fileName
	return conf
}
