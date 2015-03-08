package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/parkghost/hiradio"
	"github.com/parkghost/hiradio/cmd/internal/config"

	log "github.com/Sirupsen/logrus"
)

const (
	channelIDKey = "channelID"
	playerKey    = "player"
	proxyPortKey = "proxyPort"
)

func playCmd(args []string) {
	// load config from file
	cfgPath, err := configPath("play.json")
	if err != nil {
		log.Warnf("Failed to load configuration: %s", err)
	}
	cfg, err := loadConfig(cfgPath)
	if err != nil {
		log.Warnf("Failed to load configuration: %s", err)
	}

	// flag settings
	fs := flag.NewFlagSet("play", flag.ExitOnError)
	app := fs.String("player", cfg.GetString(playerKey, ""), "The player which supports HTTP Live Streaming")
	port := fs.Int("port", cfg.GetInt(proxyPortKey, 1077), "Port for the proxy server")
	verbose := fs.Bool("verbose", false, "Print output from the player")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: hiradio play [options] [ChannelID]

Play radio on player

The options are:`)
		fs.PrintDefaults()
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

	// save current config
	if cfgPath != "" {
		cfg.Set(playerKey, *app)
		cfg.Set(proxyPortKey, *port)
		cfg.Set(channelIDKey, channelID)
		if err := config.SaveTo(cfgPath, cfg); err != nil {
			log.Warnf("Failed to save configuration: %s", err)
		}
	}

	// run proxy server
	proxyServer := proxy{
		address: ":" + strconv.Itoa(*port),
	}
	go func() {
		if err := proxyServer.Run(); err != nil {
			log.Fatalf("Failed to start proxy: %s", err)
		}
	}()

	// run audio player
	playlist := fmt.Sprintf("http://localhost:%d/stream/%d.m3u8", *port, channelID)
	if *app == "" {
		fmt.Printf("Open URL with player:\n%s\n", playlist)
		select {}
	} else {
		p := player{
			app:     *app,
			url:     playlist,
			verbose: *verbose,
		}
		if err := p.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

var errNotFound = errors.New("key not found")

func getChannelID(args []string, cfg *config.Config) (int, error) {
	if len(args) > 0 {
		channelID, err := strconv.Atoi(args[0])
		if err != nil {
			return -1, err
		}
		return channelID, nil
	}

	channelID := cfg.GetInt(channelIDKey, -1)
	if channelID == -1 {
		return -1, errNotFound
	}
	return channelID, nil
}

type proxy struct {
	address string
}

var routeRE = regexp.MustCompile(`/stream/(\d+).m3u8`)

func (p *proxy) Run() error {
	return http.ListenAndServe(p.address, p)
}

func (p *proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	matched := routeRE.FindStringSubmatch(req.RequestURI)
	if matched == nil {
		http.NotFound(rw, req)
		return
	}
	channelID, _ := strconv.Atoi(matched[1])

	pl, err := hiradio.GetPlaylist(channelID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, req, pl.URL, http.StatusTemporaryRedirect)
}

type player struct {
	app     string
	url     string
	verbose bool

	cmd *exec.Cmd
}

func (p *player) Run() error {
	p.cmd = exec.Command(p.app, p.url)
	if p.verbose {
		p.cmd.Stdout = os.Stdout
		p.cmd.Stderr = os.Stderr
	}
	if err := p.cmd.Run(); err != nil {
		return err
	}
	return nil
}
