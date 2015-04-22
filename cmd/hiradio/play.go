package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"

	"github.com/parkghost/hiradio"
	"github.com/parkghost/hiradio/cmd/internal/config"
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
		Warnf("Failed to load configuration: %s", err)
	}
	cfg, err := loadConfig(cfgPath)
	if err != nil {
		Warnf("Failed to load configuration: %s", err)
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
		Fatalf("Failed to parse ChannelID: %s", args)
	}

	// save current config
	if cfgPath != "" {
		cfg.Set(playerKey, *app)
		cfg.Set(proxyPortKey, *port)
		cfg.Set(channelIDKey, channelID)
		if err := config.SaveTo(cfgPath, cfg); err != nil {
			Warnf("Failed to save configuration: %s", err)
		}
	}

	// run proxy server
	proxyServer := proxy{
		address: ":" + strconv.Itoa(*port),
	}
	go func() {
		if err := proxyServer.Run(); err != nil {
			Fatalf("Failed to start proxy: %s", err)
		}
	}()

	// run audio player
	quit := make(chan os.Signal)
	playlist := fmt.Sprintf("http://localhost:%d/stream/%d.m3u8", *port, channelID)
	if *app == "" {
		fmt.Printf("Open URL with player: %s\n", playlist)
	} else {
		p := player{
			app:     *app,
			url:     playlist,
			verbose: *verbose,
		}
		go func() {
			if err := p.Run(); err != nil {
				Fatal(err)
			}
			close(quit)
		}()
	}

	fmt.Print("Press ctrl-c to exit")
	signal.Notify(quit, os.Kill, os.Interrupt)
	<-quit
	println()
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

func (p *proxy) Run() error {
	return http.ListenAndServe(p.address, p)
}

var routeRE = regexp.MustCompile(`/stream/(\d+).m3u8`)

func (p *proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	matched := routeRE.FindStringSubmatch(req.RequestURI)
	if matched == nil {
		http.NotFound(rw, req)
		return
	}
	channelID, _ := strconv.Atoi(matched[1])

	pl, err := hiradio.GetPlaylist(channelID)
	if err != nil {
		Warnf("Failed to get playlist: %s", err)
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
		fmt.Println()
	}
	if err := p.cmd.Run(); err != nil {
		return err
	}
	return nil
}
