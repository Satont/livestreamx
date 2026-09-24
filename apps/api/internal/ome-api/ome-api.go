package ome_api

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/imroc/req/v3"
	"github.com/satont/stream/apps/api/internal/config"
)

var ErrNotFound = errors.New("stream not found")

func New(cfg config.Config) *OmeApi {
	client := req.C()
	if cfg.OmeApiAccessToken != "" {
		client.SetCommonHeader(
			"Authorization",
			"Basic "+base64.StdEncoding.EncodeToString([]byte(cfg.OmeApiAccessToken)),
		)
	}

	return &OmeApi{
		config: cfg,
		client: client,
	}
}

type OmeApi struct {
	config config.Config
	client *req.Client
}

type StreamInfo struct {
	Name        string
	CreatedTime *time.Time
}

type apiResponse[T any] struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Response   T      `json:"response"`
}

type streamDetail struct {
	Name  string `json:"name"`
	Input struct {
		CreatedTime *time.Time `json:"createdTime"`
	} `json:"input"`
}

// ListStreams returns names of all active (input) streams.
func (c *OmeApi) ListStreams(ctx context.Context) ([]string, error) {
	data := &apiResponse[[]string]{}

	if err := c.get(
		ctx,
		fmt.Sprintf("/v1/vhosts/%s/apps/%s/streams", c.config.OmeVHost, c.config.OmeApp),
		data,
	); err != nil {
		return nil, err
	}

	return data.Response, nil
}

// GetStream returns info about a single stream. ErrNotFound is returned
// when the stream is offline.
func (c *OmeApi) GetStream(ctx context.Context, name string) (*StreamInfo, error) {
	data := &apiResponse[streamDetail]{}

	err := c.get(
		ctx,
		fmt.Sprintf(
			"/v1/vhosts/%s/apps/%s/streams/%s",
			c.config.OmeVHost,
			c.config.OmeApp,
			name,
		),
		data,
	)
	if err != nil {
		return nil, err
	}

	return &StreamInfo{
		Name:        data.Response.Name,
		CreatedTime: data.Response.Input.CreatedTime,
	}, nil
}

func (c *OmeApi) get(ctx context.Context, path string, result any) error {
	resp, err := c.client.R().SetContext(ctx).SetSuccessResult(result).Get(c.config.OmeApiAddr + path)
	if err != nil {
		return err
	}

	if resp.StatusCode == 404 {
		return ErrNotFound
	}

	if !resp.IsSuccessState() {
		return fmt.Errorf("unexpected response from ome api: %s %s", resp.Status, resp.String())
	}

	return nil
}
