package hiradio_test

import (
	"flag"
	"testing"

	"github.com/parkghost/hiradio"
)

var testIntegration = flag.Bool("integration", false, "Perform integration tests")

func TestListChannels(t *testing.T) {
	if !*testIntegration {
		t.Skip("skipping intergration test")
	}

	cl, err := hiradio.ListChannels()
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	if len(cl) == 0 {
		t.Fatalf("empty channel list")
	}
	t.Logf("Number of channels: %d", len(cl))
}

func TestGetPlaylist(t *testing.T) {
	if !*testIntegration {
		t.Skip("skipping intergration test")
	}

	p, err := hiradio.GetPlaylist(232)
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	t.Logf("Playlist: %+v", p)
}

func TestGetChannelInfo(t *testing.T) {
	if !*testIntegration {
		t.Skip("skipping intergration test")
	}

	rl, err := hiradio.GetChannelInfo(232)
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	t.Logf("ChannelInfo: %+v", rl)
}

func TestListRankings(t *testing.T) {
	if !*testIntegration {
		t.Skip("skipping intergration test")
	}

	rl, err := hiradio.ListRankings()
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	if len(rl) == 0 {
		t.Fatalf("empty ranking list")
	}
	t.Logf("Number of ranking: %d", len(rl))
}
