// Package pluginSiteRedirect a SiteRedirect plugin.
package redirector

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Config the plugin configuration.
type Config struct {
	BaseURL         string            `json:"baseUrl,omitempty"`
	DefaultPath     string            `json:"defaultPath,omitempty"`
	RedirectCookies []string          `json:"redirectCookies,omitempty"`
	RedirectHeaders []string          `json:"redirectHeaders,omitempty"`
	RedirectionMap  map[string]string `json:"redirectionMap,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		BaseURL:         "",
		DefaultPath:     "",
		RedirectCookies: []string{},
		RedirectHeaders: []string{},
		RedirectionMap:  make(map[string]string),
	}
}

type SiteRedirect struct {
	config *Config
	next   http.Handler
	name   string
}

// New created a new SiteRedirect plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.RedirectCookies) == 0 && len(config.RedirectHeaders) == 0 {
		return nil, fmt.Errorf("you must specify at least one from RedirectCookies or RedirectHeaders")
	}

	if len(config.RedirectionMap) == 0 {
		return nil, fmt.Errorf("redirection map cannot be empty")
	}

	if config.BaseURL == "" {
		return nil, fmt.Errorf("baseUrl can't be empty")
	}
	config.BaseURL = strings.TrimSuffix(config.BaseURL, "/")

	if config.DefaultPath == "" {
		return nil, fmt.Errorf("default path can't be empty")
	}

	return &SiteRedirect{
		config: config,
		next:   next,
		name:   name,
	}, nil
}

func (s *SiteRedirect) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" || req.URL.Path == "" {

		for _, cookieName := range s.config.RedirectCookies {
			cookie, err := req.Cookie(cookieName)
			if err != nil {
				// this will only ever be http.ErrNoCookie
				continue
			}
			if s.config.RedirectionMap[cookie.Value] != "" {
				http.Redirect(
					rw,
					req,
					fmt.Sprintf("%s/%s/", s.config.BaseURL, s.config.RedirectionMap[cookie.Value]),
					301,
				)
				return
			}
		}

		for _, hdr := range s.config.RedirectHeaders {
			headerArr := req.Header[hdr]
			if len(headerArr) != 1 {
				continue
			}
			if s.config.RedirectionMap[headerArr[0]] != "" {
				http.Redirect(
					rw,
					req,
					fmt.Sprintf("%s/%s/", s.config.BaseURL, s.config.RedirectionMap[headerArr[0]]),
					301,
				)
				return
			}
		}
		http.Redirect(rw, req, fmt.Sprintf("%s/%s", s.config.BaseURL, s.config.DefaultPath), 301)
		return
	}
	s.next.ServeHTTP(rw, req)
}
