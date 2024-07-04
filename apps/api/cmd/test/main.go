package main

import (
	"fmt"
	"strings"

	"github.com/grafov/m3u8"
	"github.com/imroc/req/v3"
	"github.com/kr/pretty"
	"golang.org/x/sync/errgroup"
)

func main() {
	resolutions := []string{"", "720p_", "480p_", "360p_"}
	var masterPlaylist m3u8.MasterPlaylist
	masterPlaylist.SetVersion(9)

	var errwg errgroup.Group
	for _, resolution := range resolutions {
		errwg.Go(
			func() error {
				resp, err := req.Get(
					fmt.Sprintf(
						"http://localhost:5173/api/streams/read/%s780eac34-5b06-4647-9ae5-38b4fae26033/index.m3u8",
						resolution,
					),
				)
				if err != nil {
					panic(err)
				}
				p, t, err := m3u8.DecodeFrom(resp.Body, true)
				if err != nil {
					panic(err)
				}

				switch t {
				case m3u8.MASTER:
					masterpl := p.(*m3u8.MasterPlaylist)
					for _, variant := range masterpl.Variants {
						res := strings.Split(variant.Resolution, "x")

						variant.URI = fmt.Sprintf(
							"http://localhost:5173/api/streams/read/%sp_780eac34-5b06-4647-9ae5-38b4fae26033/index.m3u8",
							res[1],
						)
						masterPlaylist.Variants = append(masterPlaylist.Variants, variant)
					}
				}

				return nil
			},
		)
	}

	if err := errwg.Wait(); err != nil {
		panic(err)
	}

	pretty.Println(masterPlaylist.String())
}
