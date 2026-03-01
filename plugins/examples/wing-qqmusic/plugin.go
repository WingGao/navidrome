package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/navidrome/navidrome/plugins/pdk/go/metadata"
	"github.com/navidrome/navidrome/plugins/pdk/go/pdk"
)

const (
	requestTimeoutMs = 10000
)

var (
	ErrNotFound        = errors.New("not found")
	baseURL     string = "http://192.168.3.50:12000"
)

type AlbumInfoWithImage struct {
	metadata.AlbumInfoResponse
	Images []metadata.ImageInfo `json:"images,omitempty"`
}

// 安装 https://github.com/WebAssembly/binaryen/releases
// 安装 https://tinygo.org/
// 打包 tinygo build -o plugin.wasm -target wasip1 -buildmode=c-shared .
type WingQQMusicAgent struct{}

//func (WingQQMusicAgent) OnInit(ctx context.Context, req *api.InitRequest) (*api.InitResponse, error) {
//	log.Printf("WingQQMusicAgent Plugin initializing...")
//	//if _baseURL, ok := req.Config["base_url"]; !ok || _baseURL == "" {
//	//	return &api.InitResponse{Error: "baseurl configuration is required"}, nil
//	//} else {
//	//	baseURL = _baseURL
//	//	log.Printf("Using baseURL: %s", baseURL)
//	//}
//	baseURL = "http://192.168.3.50:12000"
//	return &api.InitResponse{}, nil
//}

// GetArtistURL is not implemented for Wing QQ Music
func (*WingQQMusicAgent) GetArtistURL(input metadata.ArtistRequest) (*metadata.ArtistURLResponse, error) {
	if strings.HasPrefix(input.MBID, "qq-") {
		return &metadata.ArtistURLResponse{URL: "https://y.qq.com/n/ryqq/singer/" + input.MBID[3:]}, nil
	}
	return nil, ErrNotFound
}

// GetArtistImages fetches artist images from Wing QQ Music API
func (*WingQQMusicAgent) GetArtistImages(input metadata.ArtistRequest) (*metadata.ArtistImagesResponse, error) {
	params := map[string]string{
		"id":   input.ID,
		"name": input.Name,
		"mbid": input.MBID,
	}

	resp, err := get("/nv_artist_images", params)
	if err != nil {
		return nil, err
	}

	var images []metadata.ImageInfo
	if err := json.Unmarshal(resp.Body(), &images); err != nil {
		return nil, fmt.Errorf("failed to decode artist images response: %v", err)
	}

	return &metadata.ArtistImagesResponse{Images: images}, nil
}

func getAlbumInfo(mbid string) (*AlbumInfoWithImage, error) {
	params := map[string]string{
		"mbid": mbid,
	}

	resp, err := get("/nv_album", params)
	if err != nil {
		return nil, err
	}

	var albumInfo AlbumInfoWithImage
	if err := json.Unmarshal(resp.Body(), &albumInfo); err != nil {
		return nil, fmt.Errorf("failed to decode album info response: %v", err)
	}
	return &albumInfo, nil
}

func (*WingQQMusicAgent) GetAlbumInfo(input metadata.AlbumRequest) (*metadata.AlbumInfoResponse, error) {
	resp, err := getAlbumInfo(input.MBID)
	if err != nil {
		return nil, err
	}
	return &resp.AlbumInfoResponse, nil
}

func (*WingQQMusicAgent) GetAlbumImages(input metadata.AlbumRequest) (*metadata.AlbumImagesResponse, error) {
	resp, err := getAlbumInfo(input.MBID)
	if err != nil {
		return nil, err
	}
	return &metadata.AlbumImagesResponse{Images: resp.Images}, nil
}

// Helper method to make HTTP GET requests
func get(path string, params map[string]string) (*pdk.HTTPResponse, error) {
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
	req := pdk.NewHTTPRequest(pdk.MethodGet, u.String())

	resp := req.Send()

	if resp.Status() != 200 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.Status(), string(resp.Body()))
	}

	return &resp, nil
}

func main() {}

func init() {
	metadata.Register(&WingQQMusicAgent{})
}
