package foserver

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/apex/log"
)

// listenAndServe sets up a mini web server that serves a predetermined response.
func listenAndServe(addr string, statusCode int, node int) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if statusCode == http.StatusOK {
			output := fmt.Sprintf("{\"node-%02d\": %d}", node, node)
			_, _ = w.Write([]byte(output))
		}
		return
	})
	server := &http.Server{Addr: addr, Handler: mux}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).WithField("node", node).Error("upstream node failed")
		}
	}()

	return server
}
