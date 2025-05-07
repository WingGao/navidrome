package wingqq

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/navidrome/navidrome/conf"
	"github.com/navidrome/navidrome/consts"
	"github.com/navidrome/navidrome/core/agents"
	"github.com/navidrome/navidrome/model"
	"github.com/navidrome/navidrome/utils/cache"
	"net/http"
	"net/url"
)

const wingqqAgentName = "wingqq"

type wingAgent struct {
	ds      model.DataStore
	baseURL string
	client  *cache.HTTPClient
}

func wingConstructor(ds model.DataStore) agents.Interface {
	if conf.Server.Wing.BaseURL == "" {
		return nil
	}
	l := &wingAgent{
		ds:      ds,
		baseURL: conf.Server.Wing.BaseURL,
	}
	hc := &http.Client{
		Timeout: consts.DefaultHttpClientTimeOut,
	}
	chc := cache.NewHTTPClient(hc, consts.DefaultHttpClientTimeOut)
	l.client = chc
	return l
}

func (s *wingAgent) AgentName() string {
	return wingqqAgentName
}

func (s *wingAgent) GetAlbumInfo(ctx context.Context, name, artist, mbid string) (*agents.AlbumInfo, error) {
	data := agents.AlbumInfo{}
	form := url.Values{}
	form.Add("name", name)
	form.Add("artist", artist)
	form.Add("mbid", mbid)
	err := s.get(ctx, "/nv_album", form, &data)
	return &data, err
}

func (s *wingAgent) GetArtistImages(ctx context.Context, id, name, mbid string) ([]agents.ExternalImage, error) {
	data := make([]agents.ExternalImage, 0)
	form := url.Values{}
	form.Add("id", id)
	form.Add("name", name)
	form.Add("mbid", mbid)
	err := s.get(ctx, "/nv_artist_images", form, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *wingAgent) get(ctx context.Context, path string, q url.Values, out interface{}) error {
	req, _ := http.NewRequestWithContext(ctx, "GET", s.baseURL+path, nil)
	req.URL.RawQuery = q.Encode()
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("Invalid response from WingQQ: " + resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func init() {
	conf.AddHook(func() {
		agents.Register(wingqqAgentName, wingConstructor)
	})
}
