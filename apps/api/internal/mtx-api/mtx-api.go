package mtx_api

import (
	"context"
	"fmt"
	"time"

	"github.com/imroc/req/v3"
	"github.com/satont/stream/apps/api/internal/config"
)

func New(cfg config.Config) *MtxApi {
	return &MtxApi{
		config: cfg,
	}
}

type MtxApi struct {
	config config.Config
}

type PathInfo struct {
	Name     string `json:"name"`
	ConfName string `json:"confName"`
	Source   *struct {
		Type string `json:"type"`
		Id   string `json:"id"`
	} `json:"source"`
	Ready         bool          `json:"ready"`
	ReadyTime     *time.Time    `json:"readyTime"`
	Tracks        []string      `json:"tracks"`
	BytesReceived int           `json:"bytesReceived"`
	BytesSent     int           `json:"bytesSent"`
	Readers       []interface{} `json:"readers"`
}

func (c *MtxApi) GetPathInfoByApyKey(ctx context.Context, apiKey string) (*PathInfo, error) {
	data := &PathInfo{}

	addr := fmt.Sprintf("%s%s/%s", c.config.MediaMtxAddr, ":9997/v3/paths/get", apiKey)

	_, err := req.SetSuccessResult(&data).SetContext(ctx).Get(addr)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *MtxApi) GetPaths(ctx context.Context) ([]PathInfo, error) {
	data := &PathsResponse{}

	addr := fmt.Sprintf("%s%s", c.config.MediaMtxAddr, ":9997/v3/paths/list")

	resp, err := req.SetSuccessResult(&data).SetContext(ctx).Get(addr)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("failed to get paths: %s", resp.String())
	}

	return data.Items, nil
}

type PathsResponse struct {
	PageCount int        `json:"pageCount"`
	Items     []PathInfo `json:"items"`
}
