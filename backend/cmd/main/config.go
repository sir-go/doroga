package main

import (
	"flag"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type (
	Duration struct {
		time.Duration
	}

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

	CfgDb struct {
		Host       string    `toml:"host"`
		Port       int       `toml:"port"`
		User       string    `toml:"user"`
		Password   string    `toml:"password"`
		DbName     string    `toml:"dbname"`
		Collection string    `toml:"collection"`
		Timeout    *Duration `toml:"timeout"`
	}

	CfgReqWebForm struct {
		AllowedFormats []string `toml:"allowed_formats"`
		MaxFileSize    int64    `toml:"max_file_size"`
		JpegQuality    uint     `toml:"jpeg_quality"`
		ShrinkPhotoTo  uint     `toml:"shrink_photo_to"`
		ThumbnailSize  uint     `toml:"thumbnail_size"`
		AddWatermark   bool     `toml:"add_watermark"`
	}

	CfgStorage struct {
		Path string `toml:"path"`
		Uid  int    `toml:"uid"`
		Gid  int    `toml:"gid"`
	}

	Config struct {
		Service    CfgService    `toml:"service"`
		Db         CfgDb         `toml:"db"`
		ReqWebForm CfgReqWebForm `toml:"req_web_form"`
		Storage    CfgStorage    `toml:"storage"`
		Path       string
	}
)

func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func ConfigInit() *Config {
	fCfgPath := flag.String("c", DefaultConfFile, "path to conf file")
	flag.Parse()

	conf := new(Config)
	file, err := os.Open(*fCfgPath)
	if err != nil {
		panic(err)
	}

	defer func() {
		if file == nil {
			return
		}
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err = toml.DecodeFile(*fCfgPath, &conf); err != nil {
		panic(err)
	}
	conf.Path = *fCfgPath
	return conf
}
