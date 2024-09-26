package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	writeTimeout = time.Second * 15
	readTimeout  = time.Second * 15
	idleTimeout  = time.Second * 60
)

type app interface {
	Fill(width, height, url string, headers map[string][]string) ([]byte, error)
}

type logger interface {
	Info(string)
	Warn(string)
	Error(string)
	Debug(string)
}
type Server struct {
	addr   string
	app    app
	logger logger
	server *http.Server
}

func New(port string, app app, logg logger) *Server {
	return &Server{
		addr:   ":" + port,
		app:    app,
		logger: logg,
	}
}

func (s *Server) Start() error {
	// Create new router.
	mux := http.NewServeMux()

	// Configure router.
	mux.HandleFunc("/fill/{width}/{height}/{url...}", s.fillHandler)

	// Configure server.
	s.server = &http.Server{
		Addr:         s.addr,
		Handler:      mux,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
	}

	// Run server.
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	return nil
}

// Resize handler.
func (s *Server) fillHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("incoming request: " + r.URL.String())

	// Process image.
	resizedImage, err := s.app.Fill(r.PathValue("width"), r.PathValue("height"), r.PathValue("url"), r.Header)
	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// Return image.
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(resizedImage)
	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.logger.Debug("request " + r.URL.String() + " successfully processed")
}
