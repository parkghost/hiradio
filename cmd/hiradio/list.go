package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/parkghost/hiradio"

	"github.com/moznion/go-unicode-east-asian-width"
)

func listCmd(args []string) {
	// flag settings
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: hiradio list

List channels information`)
		os.Exit(1)
	}
	fs.Parse(args)

	// fetch rankings
	type lrs struct {
		rankings []hiradio.Ranking
		err      error
	}
	rankingsCh := make(chan lrs)
	go func() {
		result, err := hiradio.ListRankings()
		rankingsCh <- lrs{result, err}
	}()

	// fetch channels
	channels, err := hiradio.ListChannels()
	if err != nil {
		Fatal(err)
	}
	result := <-rankingsCh
	if result.err != nil {
		Fatal(result.err)
	}

	// mix channels and rankings
	rc := newRankedChannels(channels, result.rankings)
	sort.Sort(rc)
	printChannelList(rc)
}

type rankedChannel struct {
	hiradio.Channel
	Ranking int
}

type rankedChannels []rankedChannel

func (rcs rankedChannels) Len() int {
	return len(rcs)
}
func (rcs rankedChannels) Swap(i, j int) {
	rcs[i], rcs[j] = rcs[j], rcs[i]
}
func (rcs rankedChannels) Less(i, j int) bool {
	// order by Type > Ranking > ID
	if rcs[i].Type != rcs[j].Type {
		return rcs[i].Type < rcs[j].Type
	}
	if rcs[i].Ranking != rcs[j].Ranking {
		if rcs[i].Ranking == 0 {
			return false
		}
		if rcs[j].Ranking == 0 {
			return true
		}
		return rcs[i].Ranking < rcs[j].Ranking
	}
	return rcs[i].ID < rcs[j].ID
}

func newRankedChannels(channels []hiradio.Channel, rankings []hiradio.Ranking) rankedChannels {
	list := make(rankedChannels, 0, len(channels))
	for _, c := range channels {
		rc := rankedChannel{c, 0}
		for _, r := range rankings {
			if r.ID == c.ID {
				rc.Ranking = r.Value
			}
		}
		list = append(list, rc)
	}
	return list
}

func printChannelList(rc rankedChannels) {
	//coloumn size:4  8  4  30  *
	fmt.Printf("%-2s  %-6s  %-2s  %-28s  %s\n", "編號", "類型", "排行", "頻道", "現在播放節目")
	for _, c := range rc {
		wType := 8 - stringWidth(c.Type.String()) + len([]rune(c.Type.String()))
		ranking := ""
		if c.Ranking > 0 {
			ranking = strconv.Itoa(c.Ranking)
		}
		wTitle := 30 - stringWidth(c.Title) + len([]rune(c.Title))
		fmt.Printf("%4d  %-*s  %4s  %-*s  %s\n",
			c.ID,
			wType, c.Type,
			ranking,
			wTitle, c.Title,
			c.ProgramName)
	}
}

// stringWidth return width of s
func stringWidth(s string) int {
	var n int
	for _, r := range s {
		if eastasianwidth.IsFullwidth(r) {
			n = n + 2
		} else {
			n = n + 1
		}
	}
	return n
}
