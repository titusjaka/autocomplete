package server

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"

	"github.com/titusjaka/autocomplete"
)

// Server listens TCP port and runs HTTP server
type Server struct {
	proxy  autocomplete.AVSProxy
	cache  autocomplete.Database
	logger *log.Logger
}

// New returns new Server structure
func New(proxy autocomplete.AVSProxy, cache autocomplete.Database) *Server {
	return &Server{
		proxy:  proxy,
		cache:  cache,
		logger: log.New(ioutil.Discard, "", 0),
	}
}

// SetLogger sets logger
func (h *Server) SetLogger(logger *log.Logger) error {
	if logger == nil {
		return errors.New("logger is nil")
	}
	h.logger = logger
	h.logger.Printf("[DEBUG] logger set successfully")
	return nil
}

// RunServer starts HTTP server on 'listen' address
func (h *Server) RunServer(ctx context.Context, listen string, debug bool) error {
	e := echo.New()
	e.GET("/search", h.searchPlaces)

	if debug {
		e.Use(middleware.Logger())
		e.Debug = true
	}

	errs := make(chan error, 1)
	go func() {
		errs <- e.Start(listen)
	}()

	select {
	case <-ctx.Done():
		return e.Shutdown(ctx)
	case err := <-errs:
		return err
	}
}

func (h *Server) searchPlaces(c echo.Context) error {
	rawQuery := c.Request().RequestURI
	cache, err := h.cache.Get(rawQuery)

	if err == nil {
		h.logger.Printf("[DEBUG] Got data for '%s' query from cache DB", rawQuery)
		return c.JSON(http.StatusOK, cache)
	}

	params := c.QueryParams()

	wd, err := h.proxy.SearchPlaces(params)
	if err != nil {
		h.logger.Printf("[ERROR] Failed to fetch data from proxy: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	h.logger.Printf("[DEBUG] Got data for '%s' query from proxy", rawQuery)
	if err = h.cache.Write(rawQuery, wd); err != nil {
		h.logger.Printf("[ERROR] Cannot save data to cache DB: %v", err)
	}

	return c.JSON(http.StatusOK, wd)
}
