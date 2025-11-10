package foserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)
import "github.com/apex/log"
import "golang.org/x/sync/errgroup"

var client *http.Client
var host string = "localhost:8080"
var upstreams = []string{
	"localhost:8081",
	"localhost:8082",
	"localhost:8083",
	"localhost:8084",
	"localhost:8085",
}

func fetch(ctx context.Context, url string) (map[string]any, error) {
	var out map[string]any

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode >= http.StatusMultipleChoices {
		return nil, err
	}

	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func aggregateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2500*time.Millisecond)
	defer cancel()

	merged := make(map[string]any)
	g, ctx := errgroup.WithContext(ctx)
	results := make([]map[string]any, len(upstreams))
	for idx, upstream := range upstreams {
		i, u := idx, fmt.Sprintf("http://%s/", upstream)
		g.Go(func() error {
			res, err := fetch(ctx, u)
			if err != nil {
				return err
			}
			results[i] = res
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		log.WithError(err).Error("fan out request failed")
	}
	for _, result := range results {
		for k, v := range result {
			merged[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(merged)
}

func Execute(ctx context.Context, c *http.Client) error {
	client = c

	// create upstream servers
	servers := make([]*http.Server, len(upstreams))
	for i, addr := range upstreams {
		servers[i] = listenAndServe(addr, http.StatusOK, i+1)
	}

	// set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// create local listener
	mux := http.NewServeMux()
	mux.HandleFunc("/aggregate", aggregateHandler)
	localServer := &http.Server{Addr: host, Handler: mux}
	go func() {
		if err := localServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Error("local listener failed")
		}
	}()
	log.WithField("addr", host).Info("server listening")

	<-sigChan
	log.Info("received shutdown signal, initiating graceful shutdown...")
	ctxShutdown, cancel := context.WithTimeout(ctx, 2500*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(servers) + 1)
	go func() {
		defer wg.Done()
		if err := localServer.Shutdown(ctxShutdown); err != nil {
			log.WithError(err).Error("graceful shutdown failed")
		}
	}()
	for index, server := range servers {
		go func(i int, s *http.Server) {
			defer wg.Done()
			if err := s.Shutdown(ctxShutdown); err != nil {
				log.WithError(err).WithField("node", i+1).Error("graceful shutdown failed for upstream server")
			}
		}(index, server)
	}
	wg.Wait()
	log.Info("all servers completed shutdown")

	return nil
}
