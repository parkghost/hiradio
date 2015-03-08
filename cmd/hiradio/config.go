package main

import (
	"path/filepath"

	"github.com/parkghost/hiradio/cmd/internal/config"

	"github.com/surma-dump/goappdata"
)

func configPath(name string) (string, error) {
	dir, err := goappdata.CreatePath("hiradio")
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, name), nil
}

func loadConfig(path string) (*config.Config, error) {
	if path == "" {
		return config.New(), nil
	}
	cfg, err := config.From(path)
	if err != nil {
		if err == config.ErrEmptyFile {
			return config.New(), nil
		}
		return config.New(), err
	}
	return cfg, nil
}
