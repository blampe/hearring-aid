package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"go.uber.org/zap/buffer"
	"golang.org/x/time/rate"
)

type cli struct {
	Token string `default:"2" help:"tadb token, defaults to the free API token"`
	RPS   int64  `default:"5" help:"maximum requests per second"`
	Port  int    `default:"8080" help:"the port to listen on"`
}

func proxy(w http.ResponseWriter, r *http.Request) {

}

// Reduce allocations by pooling buffers.
var _buffers = buffer.NewPool()

func (c *cli) Run() error {
	// There's probably a simpler way to compute this duration but I haven't had my coffee.
	every := time.Duration(float64(time.Second) / float64(time.Duration(2)*time.Second) * float64(time.Second))
	throttle := rate.NewLimiter(rate.Every(every), 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		err := throttle.Wait(ctx)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if !strings.HasPrefix(r.URL.Path, "/v1/caa/") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		url := fixCAAPath(r.URL)

		req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		defer resp.Body.Close()

		w.WriteHeader(resp.StatusCode)
		if resp.StatusCode == http.StatusOK {
			w.Header().Add("Cache-Control", "s-maxage=31536000, max-age=86400") // 1 year in CDN, 1 day locally.
		}

		buf := _buffers.Get()
		defer buf.Free()
		if buf.Len() == 0 {
			buf.AppendBytes(make([]byte, 1024*1024))
		}
		io.CopyBuffer(w, resp.Body, buf.Bytes())
	})

	addr := fmt.Sprintf(":%d", c.Port)
	server := &http.Server{Handler: mux, Addr: addr}

	return server.ListenAndServe()
}

// fixCAAPath reconstructs the upstream CoverArtArchive URL.
//
// https://images.musicinfo.pro/v1/caa/foo/bar
//
// becomes
//
// https://coverartarchive.org/release/foo/bar

func fixCAAPath(u *url.URL) *url.URL {
	url := u
	url.Path = strings.Replace(url.Path, "/v1/caa/", "/release/", 1)
	url.Host = "coverartarchive.org"
	url.Scheme = "https"
	return url
}

func main() {
	kctx := kong.Parse(&cli{})
	err := kctx.Run()
	if err != nil {
		log.Fatal(err.Error())
	}

}
