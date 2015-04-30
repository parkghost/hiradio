package hiradio_test

import (
	"testing"

	"github.com/parkghost/hiradio"
)

func TestIntegrationListChannels(t *testing.T) {
	if testing.Short() {
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

func TestIntegrationGetPlaylist(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping intergration test")
	}

	p, err := hiradio.GetPlaylist(232)
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	t.Logf("Playlist: %+v", p)
}

func TestIntegrationGetChannelInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping intergration test")
	}

	rl, err := hiradio.GetChannelInfo(232)
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	t.Logf("ChannelInfo: %+v", rl)
}

func TestIntegrationListRankings(t *testing.T) {
	if testing.Short() {
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
