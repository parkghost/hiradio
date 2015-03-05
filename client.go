// Package hiradio provides a client for using the Hichannel API.
package hiradio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	libraryVersion   = "0.1"
	defaultEndpoint  = "http://hichannel.hinet.net/radio/"
	defaultUserAgent = "hiradio/" + libraryVersion
)

// A Client manages communicating with Hichannel API.
type Client struct {
	client *http.Client

	// Endpoint for API requests.
	Endpoint string

	// User agent used when communicating with the Hichannel API.
	UserAgent string
}

// Channel represents a Hichannel channel.
//
// Source: http://hichannel.hinet.net/radio/channelList.do?radioType=&freqType=&freq=&area=&pN=%d
type Channel struct {
	ID    int
	Title string
	Image string
	Type  RadioType

	ProgramName string
}

// RadioType represents a Hichannel radio type.
//
// Mappings:
//	1	音樂
//	2	生活資訊
//	3	新聞
//	4	綜合
//	5	外語
//	6	多元文化
//	7	交通
type RadioType int

func (rt RadioType) String() string {
	switch rt {
	case 1:
		return "音樂"
	case 2:
		return "生活資訊"
	case 3:
		return "新聞"
	case 4:
		return "綜合"
	case 5:
		return "外語"
	case 6:
		return "多元文化"
	case 7:
		return "交通"
	}
	return "Unknown"
}

// ListChannels list all channels.
func (c *Client) ListChannels() ([]Channel, error) {
	// resolve page size of ChannelLists
	initPage, err := c.fetchChannelList(1)
	if err != nil {
		return nil, err
	}

	pageSize := initPage.PageSize
	var list []Channel
	list, err = appendChannel(list, initPage.List)
	if err != nil {
		return nil, err
	}

	// fetch rest of ChannelLists
	pages := make([]int, 0, pageSize-1)
	for i := 2; i <= pageSize; i++ {
		pages = append(pages, i)
	}
	restOfPages, err := c.fetchChannelLists(pages)
	if err != nil {
		return nil, err
	}

	for _, page := range restOfPages {
		list, err = appendChannel(list, page.List)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}

type channelList struct {
	PageNo   int       `json:"pageNo"`
	PageSize int       `json:"pageSize"`
	List     []channel `json:"list"`
}

type channel struct {
	IsChannel    bool   `json:"isChannel"`
	ChannelImage string `json:"channel_image,omitempty"`
	ChannelTitle string `json:"channel_title,omitempty"`
	RadioType    int    `json:"radio_type,omitempty,string"`
	ProgramName  string `json:"program_name,omitempty"`
	ChannelID    int    `json:"channel_id,omitempty,string"`
}

func (c *Client) fetchChannelList(page int) (*channelList, error) {
	url := fmt.Sprintf("%schannelList.do?pN=%d", c.Endpoint, page)
	cl := new(channelList)
	if err := c.fetchObject(url, cl); err != nil {
		return nil, err
	}
	return cl, nil
}

func (c *Client) fetchObject(url string, v interface{}) error {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", c.UserAgent)
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if err = checkResponse(res); err != nil {
		return err
	}

	dec := json.NewDecoder(res.Body)
	return dec.Decode(v)
}

type responseError struct {
	Response *http.Response
	Message  string
}

func (r *responseError) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message)
}

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	respErr := &responseError{Response: r}
	if data, err := ioutil.ReadAll(r.Body); err == nil && data != nil {
		respErr.Message = string(data)
	}
	return respErr
}

func appendChannel(dst []Channel, src []channel) ([]Channel, error) {
	for _, e := range src {
		if e.IsChannel {
			dst = append(dst, Channel{
				ID:          e.ChannelID,
				Title:       e.ChannelTitle,
				Image:       e.ChannelImage,
				Type:        RadioType(e.RadioType),
				ProgramName: e.ProgramName,
			})
		}
	}
	return dst, nil
}

func (c *Client) fetchChannelLists(pages []int) ([]channelList, error) {
	n := len(pages)
	aggregate := make(chan channelList, n)
	errCh := make(chan error, n)
	for _, page := range pages {
		go func(i int) {
			cl, err := c.fetchChannelList(i)
			if err != nil {
				errCh <- err
				return
			}
			aggregate <- *cl
		}(page)
	}

	cls := make([]channelList, 0, n)
	nr := 0
	for {
		select {
		case result := <-aggregate:
			cls = append(cls, result)
			nr++
			if n == nr {
				close(aggregate)
				return cls, nil
			}
		case err := <-errCh:
			return nil, err
		}
	}
}

// Playlist represents a url of m3u8 file.
//
// Source: http://hichannel.hinet.net/radio/play.do?id=%d
type Playlist struct {
	URL string `json:"playRadio"`
}

// GetPlaylist fetches a playlist for specified channel.
func (c *Client) GetPlaylist(channelID int) (*Playlist, error) {
	return c.fetchPlaylist(channelID)
}

func (c *Client) fetchPlaylist(channelID int) (*Playlist, error) {
	url := fmt.Sprintf("%splay.do?id=%d", c.Endpoint, channelID)
	p := new(Playlist)
	if err := c.fetchObject(url, p); err != nil {
		return nil, err
	}

	if p.URL == "" {
		return nil, fmt.Errorf("playlist not found, channelID: %d", channelID)
	}
	return p, nil
}

// ChannelInfo represents a Hichannel channel information.
//
// Source: http://hichannel.hinet.net/radio/getProgramList.do?channelId=%d
type ChannelInfo struct {
	Area     string    `json:"channel_area"`
	Desc     string    `json:"channel_desc"`
	Image    string    `json:"channel_image"`
	Title    string    `json:"channel_title"`
	TypeText string    `json:"channel_type"`
	List     []Program `json:"list"`
}

// Program represents a Hichannel program.
//
// Source: http://hichannel.hinet.net/radio/getProgramList.do?channelId=%d
type Program struct {
	EndTime   string `json:"end_time"`
	Name      string `json:"name"`
	On        bool   `json:"on"`
	StartTime string `json:"start_time"`
}

// GetChannelInfo fetches a channel information for specified channel.
func (c *Client) GetChannelInfo(channelID int) (*ChannelInfo, error) {
	return c.fetchChannelInfo(channelID)
}

func (c *Client) fetchChannelInfo(channelID int) (*ChannelInfo, error) {
	url := fmt.Sprintf("%sgetProgramList.do?channelId=%d", c.Endpoint, channelID)
	ci := new(ChannelInfo)
	if err := c.fetchObject(url, ci); err != nil {
		return nil, err
	}

	if ci.Title == "" {
		return nil, fmt.Errorf("channel not found, channelID: %d", channelID)
	}
	return ci, nil
}

type rankingList struct {
	List []Ranking `json:"list"`
}

// Ranking represents a Hichannel ranking.
//
// Source: http://hichannel.hinet.net/radio/getRanking.do
type Ranking struct {
	ID    int `json:"channel_id,string"`
	Value int `json:"channel_rank,string"`
}

// ListRankings list all rankings.
func (c *Client) ListRankings() ([]Ranking, error) {
	rl, err := c.fetchRankingList()
	if err != nil {
		return nil, err
	}
	return rl.List, nil
}

func (c *Client) fetchRankingList() (*rankingList, error) {
	url := c.Endpoint + "getRanking.do"
	rl := new(rankingList)
	if err := c.fetchObject(url, rl); err != nil {
		return nil, err
	}
	return rl, nil
}

// NewClient returns a new Client.
func NewClient(client *http.Client) *Client {
	c := new(Client)
	c.client = client
	c.Endpoint = defaultEndpoint
	c.UserAgent = defaultUserAgent
	return c
}

// DefaultClient is the default Client.
//
// Http request will be canceled if read/write time more than 1 minus.
var DefaultClient = &Client{
	client: &http.Client{
		Timeout: 1 * time.Minute,
	},
	Endpoint:  defaultEndpoint,
	UserAgent: defaultUserAgent,
}

// ListChannels list all channels.
func ListChannels() ([]Channel, error) {
	return DefaultClient.ListChannels()
}

// GetPlaylist fetches a playlist for specified channel.
func GetPlaylist(channelID int) (*Playlist, error) {
	return DefaultClient.GetPlaylist(channelID)
}

// GetChannelInfo fetches a channel information for specified channel.
func GetChannelInfo(channelID int) (*ChannelInfo, error) {
	return DefaultClient.GetChannelInfo(channelID)
}

// ListRankings list all rankings.
func ListRankings() ([]Ranking, error) {
	return DefaultClient.ListRankings()
}
