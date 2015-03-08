package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/parkghost/hiradio"
	"github.com/parkghost/hiradio/cmd/internal/config"

	log "github.com/Sirupsen/logrus"
)

func infoCmd(args []string) {
	// load config from file
	cfgPath, err := configPath("info.json")
	if err != nil {
		log.Warnf("Failed to load configuration: %s", err)
	}
	cfg, err := loadConfig(cfgPath)
	if err != nil {
		log.Warnf("Failed to load configuration: %s", err)
	}

	// flag settings
	fs := flag.NewFlagSet("info", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: hiradio info [ChannelID]

Display radio information and program list`)
		os.Exit(1)
	}

	// parse arguments
	fs.Parse(args)
	if fs.NArg() > 1 {
		fs.Usage()
		return
	}
	channelID, err := getChannelID(fs.Args(), cfg)
	if err != nil {
		if err == errNotFound {
			fs.Usage()
			return
		}
		log.Fatalf("Failed to parse ChannelID: %s", args)
	}

	// fetch channel info
	info, err := hiradio.GetChannelInfo(channelID)
	if err != nil {
		log.Fatal(err)
	}

	printChannelInfo(channelID, info)

	// save current config
	if cfgPath != "" {
		cfg.Set(channelIDKey, channelID)
		if err := config.SaveTo(cfgPath, cfg); err != nil {
			log.Warnf("Failed to save configuration: %s", err)
		}
	}
}

func printChannelInfo(id int, info *hiradio.ChannelInfo) {
	fmt.Printf("編號: %d\n頻道: %s\n類型: %s\n地點: %s\n簡介: %s\n",
		id,
		info.Title,
		info.TypeText,
		info.Area,
		info.Desc)
	fmt.Println("節目表:")
	var playing bool
	for _, p := range info.List {
		cursor := "  "
		if p.On && !playing {
			cursor = ">>"
			playing = true
		}
		fmt.Printf("%s %s ~ %s  %s\n", cursor, p.StartTime, p.EndTime, p.Name)
	}
}
