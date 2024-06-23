package mtx_api

import (
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

type StreamInfo struct {
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

func (c *MtxApi) GetStreamInfo() (*StreamInfo, error) {
	data := &StreamInfo{}

	_, err := req.SetSuccessResult(&data).Get(c.config.StreamPathInfoAddr)
	if err != nil {
		return nil, err
	}

	return data, nil
}
