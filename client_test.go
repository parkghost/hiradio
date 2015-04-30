package hiradio

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	url, _ := url.Parse(server.URL)
	client = NewClient(&http.Client{})
	client.Endpoint = url.String() + "/radio/"
}

func teardown() {
	server.Close()
}

func TestFetchChannelList(t *testing.T) {
	setup()
	defer teardown()
	test := `{
    "pageNo": 1,
    "pageSize": 4,
    "list":[
        {
            "channel_id": "1471",
            "channel_image": "14a7b76cf9c00000340a.jpg",
            "channel_title": "NER教育電臺 臺北總臺AM",
            "isChannel": true,
            "program_name": "校園健康筆記",
            "radio_type": "4"
        }
    ]
}`
	want := channelList{1, 4, []channel{{true, "14a7b76cf9c00000340a.jpg", "NER教育電臺 臺北總臺AM", 4, "校園健康筆記", 1471}}}
	mux.HandleFunc("/radio/channelList.do", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(test))
	})

	got, err := client.fetchChannelList(1)
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	if !reflect.DeepEqual(*got, want) {
		t.Fatalf("got %+v, want %+v", *got, want)
	}
}

func TestFetchChannelLists(t *testing.T) {
	setup()
	defer teardown()
	test := `{
    "pageNo": 1,
    "pageSize": 4,
    "list":[
        {
            "channel_id": "1471",
            "channel_image": "14a7b76cf9c00000340a.jpg",
            "channel_title": "NER教育電臺 臺北總臺AM",
            "isChannel": true,
            "program_name": "校園健康筆記",
            "radio_type": "4"
        }
    ]
}`
	want := []channelList{
		{1, 4, []channel{{true, "14a7b76cf9c00000340a.jpg", "NER教育電臺 臺北總臺AM", 4, "校園健康筆記", 1471}}},
		{1, 4, []channel{{true, "14a7b76cf9c00000340a.jpg", "NER教育電臺 臺北總臺AM", 4, "校園健康筆記", 1471}}},
	}
	mux.HandleFunc("/radio/channelList.do", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(test))
	})
	pages := []int{2, 3}
	got, err := client.fetchChannelLists(pages)
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestAppendChannel(t *testing.T) {
	test := []channel{
		{true, "14a7b212625000003d2d.jpg", "大千電台", 2, "Super Live Show", 109},
		{IsChannel: false},
		{true, "14ab932e39800000b250.jpg", "大漢之音", 6, "客家恁靚！-主持人Rita、Vera", 300},
	}
	want := []Channel{
		{109, "大千電台", "14a7b212625000003d2d.jpg", 2, "Super Live Show"},
		{300, "大漢之音", "14ab932e39800000b250.jpg", 6, "客家恁靚！-主持人Rita、Vera"},
	}

	var dst []Channel
	got, err := appendChannel(dst, test)
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestListChannels(t *testing.T) {
	setup()
	defer teardown()
	test := `{
    "pageNo": 1,
    "pageSize": 3,
    "list":[
        {
            "channel_id": "1471",
            "channel_image": "14a7b76cf9c00000340a.jpg",
            "channel_title": "NER教育電臺 臺北總臺AM",
            "isChannel": true,
            "program_name": "校園健康筆記",
            "radio_type": "4"
        }
    ]
}`
	want := []Channel{
		{1471, "NER教育電臺 臺北總臺AM", "14a7b76cf9c00000340a.jpg", 4, "校園健康筆記"},
		{1471, "NER教育電臺 臺北總臺AM", "14a7b76cf9c00000340a.jpg", 4, "校園健康筆記"},
		{1471, "NER教育電臺 臺北總臺AM", "14a7b76cf9c00000340a.jpg", 4, "校園健康筆記"},
	}
	mux.HandleFunc("/radio/channelList.do", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(test))
	})
	got, err := client.ListChannels()
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestGetPlaylist(t *testing.T) {
	setup()
	defer teardown()
	test := `{
    "channel_collect": false,
    "channel_title": "飛碟電台",
    "playRadio": "http://radio-hichannel.cdn.hinet.net/live/pool/hich-ra000072/ra-hls/index.m3u8?token1=1h9kB2PACU_DrzMNIgRC9Q&token2=0hMOrYO0s-RLwkZS7PMH0g&expire1=1425121978&expire2=1425143578",
    "programName": "亞洲美食王（陳鴻）",
    "vaParamter": "&2&id=232"
}`
	want := Playlist{"http://radio-hichannel.cdn.hinet.net/live/pool/hich-ra000072/ra-hls/index.m3u8?token1=1h9kB2PACU_DrzMNIgRC9Q&token2=0hMOrYO0s-RLwkZS7PMH0g&expire1=1425121978&expire2=1425143578"}
	mux.HandleFunc("/radio/play.do", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(test))
	})
	got, err := client.GetPlaylist(232)
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	if !reflect.DeepEqual(*got, want) {
		t.Fatalf("got %+v, want %+v", *got, want)
	}
}

func TestGetChannelInfo(t *testing.T) {
	setup()
	defer teardown()
	test := `{
    "channel_area": "北區(基北桃竹苗)",
    "channel_desc": "熱情Play 只想聽音樂，為全方位的音樂電台。節目內容豐富多元，熱情、專業的DJ全天候放送流行音樂，讓您零時差的迅速抓住最新、最流行的中西方娛樂消息。",
    "channel_image": "14abcde694d00000b2fc.jpg",
    "channel_title": "HitFm聯播網 Taipei 北部",
    "channel_type": "音樂",
    "isToday": true,
    "list": [
        {
            "end_time": "02:00",
            "name": "LOVE DJ",
            "on": true,
            "start_time": "00:00"
        },
        {
            "end_time": "09:00",
            "name": "只想聽音樂",
            "on": true,
            "start_time": "02:00"
        }
    ]
}`
	want := ChannelInfo{
		"北區(基北桃竹苗)",
		"熱情Play 只想聽音樂，為全方位的音樂電台。節目內容豐富多元，熱情、專業的DJ全天候放送流行音樂，讓您零時差的迅速抓住最新、最流行的中西方娛樂消息。",
		"14abcde694d00000b2fc.jpg",
		"HitFm聯播網 Taipei 北部",
		"音樂",
		[]Program{
			{"02:00", "LOVE DJ", true, "00:00"},
			{"09:00", "只想聽音樂", true, "02:00"},
		},
	}

	mux.HandleFunc("/radio/getProgramList.do", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(test))
	})
	got, err := client.GetChannelInfo(222)
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	if !reflect.DeepEqual(*got, want) {
		t.Fatalf("got %+v, want %+v", *got, want)
	}
}

func TestListRankings(t *testing.T) {
	setup()
	defer teardown()
	test := `{
    "list": [
        {
            "channel_id": "222",
            "channel_image": "14abcde694d00000b2fc.jpg",
            "channel_rank": "1",
            "channel_title": "HitFm聯播網 Taipei 北部",
            "program_name": "HITO唱片行",
            "radio_type": "1"
        },
        {
            "channel_id": "156",
            "channel_image": "14a7b23f952000001965.jpg",
            "channel_rank": "2",
            "channel_title": "KISS RADIO 大眾廣播電台",
            "program_name": "音樂玩家",
            "radio_type": "1"
        }
    ]
}`
	want := []Ranking{
		{222, 1},
		{156, 2},
	}
	mux.HandleFunc("/radio/getRanking.do", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(test))
	})
	got, err := client.ListRankings()
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestIntegrationListChannels(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping intergration test")
	}

	cl, err := ListChannels()
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

	p, err := GetPlaylist(232)
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	t.Logf("Playlist: %+v", p)
}

func TestIntegrationGetChannelInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping intergration test")
	}

	rl, err := GetChannelInfo(232)
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	t.Logf("ChannelInfo: %+v", rl)
}

func TestIntegrationListRankings(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping intergration test")
	}

	rl, err := ListRankings()
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	if len(rl) == 0 {
		t.Fatalf("empty ranking list")
	}
	t.Logf("Number of ranking: %d", len(rl))
}
