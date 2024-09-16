package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrNotEnoughParameters = errors.New("not enough parameters")
)

type app interface {
	Crop(width, height int, url string, headers map[string][]string) ([]byte, error)
	Resize(width, height int, url string, headers map[string][]string) ([]byte, error)
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
	mux.HandleFunc("/crop/{width}/{height}/{url...}", s.cropHandler)
	mux.HandleFunc("/resize/{width}/{height}/{url...}", s.resizeHandler)

	// Configure server.
	s.server = &http.Server{
		Addr:         s.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
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

// Crop handler.
func (s *Server) cropHandler(w http.ResponseWriter, r *http.Request) {
	// Get request parameters.
	width, height, url, err := getParameters(r.PathValue("width"), r.PathValue("height"), r.PathValue("url"))
	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process image.
	image, err := s.app.Crop(width, height, url, r.Header)
	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return image.
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(image)
	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Resize handler.
func (s *Server) resizeHandler(w http.ResponseWriter, r *http.Request) {
	// Get request parameters.
	width, height, url, err := getParameters(r.PathValue("width"), r.PathValue("height"), r.PathValue("url"))
	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process image.
	image, err := s.app.Resize(width, height, url, r.Header)
	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return image.
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(image)
	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Get request parameters.
func getParameters(width, height, imageUrl string) (int, int, string, error) {
	// Check if parameters are not empty.
	if width == "" || height == "" || imageUrl == "" {
		return 0, 0, "", ErrNotEnoughParameters
	}

	// Get width.
	w, err := strconv.Atoi(width)
	if err != nil {
		return 0, 0, "", err
	}

	// Get heigth.
	h, err := strconv.Atoi(height)
	if err != nil {
		return 0, 0, "", err
	}

	// Add scheme to url.
	u := "http://" + imageUrl

	return w, h, u, nil
}
