//go:build wasip1

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/navidrome/navidrome/plugins/api"
	"github.com/navidrome/navidrome/plugins/host/http"
)

const (
	requestTimeoutMs = 5000
)

var (
	ErrNotFound       = api.ErrNotFound
	ErrNotImplemented = api.ErrNotImplemented

	client         = http.NewHttpService()
	baseURL string = ""
)

type WingQQMusicAgent struct{}

func (WingQQMusicAgent) OnInit(ctx context.Context, req *api.InitRequest) (*api.InitResponse, error) {
	log.Printf("WingQQMusicAgent Plugin initializing...")
	//if _baseURL, ok := req.Config["base_url"]; !ok || _baseURL == "" {
	//	return &api.InitResponse{Error: "baseurl configuration is required"}, nil
	//} else {
	//	baseURL = _baseURL
	//	log.Printf("Using baseURL: %s", baseURL)
	//}
	baseURL = "http://192.168.3.50:12000"
	return &api.InitResponse{}, nil
}

// GetArtistURL is not implemented for Wing QQ Music
func (WingQQMusicAgent) GetArtistURL(ctx context.Context, req *api.ArtistURLRequest) (*api.ArtistURLResponse, error) {
	if strings.HasPrefix(req.GetMbid(), "qq-") {
		return &api.ArtistURLResponse{Url: "https://y.qq.com/n/ryqq/singer/" + req.GetMbid()[3:]}, nil
	}
	return nil, ErrNotFound
}

// GetArtistBiography is not implemented for Wing QQ Music
func (WingQQMusicAgent) GetArtistBiography(context.Context, *api.ArtistBiographyRequest) (*api.ArtistBiographyResponse, error) {
	return nil, ErrNotImplemented
}

// GetArtistImages fetches artist images from Wing QQ Music API
func (WingQQMusicAgent) GetArtistImages(ctx context.Context, req *api.ArtistImageRequest) (*api.ArtistImageResponse, error) {
	params := map[string]string{
		"id":   req.Id,
		"name": req.Name,
		"mbid": req.Mbid,
	}

	resp, err := get(ctx, "/nv_artist_images", params)
	if err != nil {
		return nil, err
	}

	var images []*api.ExternalImage
	if err := json.Unmarshal(resp.Body, &images); err != nil {
		return nil, fmt.Errorf("failed to decode artist images response: %v", err)
	}

	return &api.ArtistImageResponse{Images: images}, nil
}

// Not implemented methods
func (WingQQMusicAgent) GetArtistMBID(context.Context, *api.ArtistMBIDRequest) (*api.ArtistMBIDResponse, error) {
	return nil, ErrNotImplemented
}
func (WingQQMusicAgent) GetSimilarArtists(context.Context, *api.ArtistSimilarRequest) (*api.ArtistSimilarResponse, error) {
	return nil, ErrNotImplemented
}
func (WingQQMusicAgent) GetArtistTopSongs(context.Context, *api.ArtistTopSongsRequest) (*api.ArtistTopSongsResponse, error) {
	return nil, ErrNotImplemented
}
func (WingQQMusicAgent) GetAlbumInfo(ctx context.Context, req *api.AlbumInfoRequest) (*api.AlbumInfoResponse, error) {
	params := map[string]string{
		"name":   req.Name,
		"artist": req.Artist,
		"mbid":   req.Mbid,
	}

	resp, err := get(ctx, "/nv_album", params)
	if err != nil {
		return nil, err
	}

	var albumInfo api.AlbumInfoResponse
	if err := json.Unmarshal(resp.Body, &albumInfo); err != nil {
		return nil, fmt.Errorf("failed to decode album info response: %v", err)
	}

	return &albumInfo, nil
}

func (WingQQMusicAgent) GetAlbumImages(context.Context, *api.AlbumImagesRequest) (*api.AlbumImagesResponse, error) {
	return nil, ErrNotImplemented
}

// Helper method to make HTTP GET requests
func get(ctx context.Context, path string, params map[string]string) (*http.HttpResponse, error) {
	if baseURL == "" {
		return nil, errors.New("baseURL not configured")
	}

	// Build URL with query parameters
	u, err := url.Parse(baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	if len(params) > 0 {
		q := u.Query()
		for key, value := range params {
			q.Add(key, value)
		}
		u.RawQuery = q.Encode()
	}

	// Make HTTP request
	req := &http.HttpRequest{
		Url:       u.String(),
		TimeoutMs: requestTimeoutMs,
	}

	resp, err := client.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.Status != 200 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.Status, string(resp.Body))
	}

	return resp, nil
}

func main() {}

func init() {
	// Configure logging: No timestamps, no source file/line
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.SetPrefix("[WingQQMusic] ")

	api.RegisterMetadataAgent(WingQQMusicAgent{})
}
