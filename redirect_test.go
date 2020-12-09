package redirector

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	redirectMaps = map[string]string{
		"A": "a",
		"B": "b",
		"b": "b",
		"C": "c",
	}
)

func TestDemo(t *testing.T) {
	tests := []struct {
		name       string
		Config     *Config
		cookie     *http.Cookie
		headers    http.Header
		next       func(t *testing.T) http.Handler
		assertFunc func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "no headers",
			Config: &Config{
				BaseURL: "https://www.mysite.com",
				RedirectHeaders: []string{
					"header1",
				},
				RedirectionMap:  redirectMaps,
				RedirectCookies: []string{},
				DefaultPath:     "default/",
			},
			next: func(t *testing.T) http.Handler {
				t.Helper()
				return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					t.Errorf("got called")
				})
			},
			assertFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				location := resp.HeaderMap["Location"][0]
				want := "https://www.mysite.com/default/"
				if location != want {
					t.Errorf("wanted url %s, got %s", want, location)
				}
			},
		},
		{
			name: "cf headers",
			headers: http.Header{
				"header2": []string{"A"},
			},
			Config: &Config{
				BaseURL: "https://www.mysite.com",
				RedirectHeaders: []string{
					"header2",
				},
				RedirectCookies: []string{},
				RedirectionMap:  redirectMaps,
				DefaultPath:     "default/",
			},
			next: func(t *testing.T) http.Handler {
				t.Helper()
				return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					t.Errorf("got called")
				})
			},
			assertFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				location := resp.HeaderMap["Location"][0]
				want := "https://www.mysite.com/a/"
				if location != want {
					t.Errorf("wanted url %s, got %s", want, location)
				}
			},
		},
		{
			headers: http.Header{
				"header1": []string{"b"},
			},
			name: "header",
			Config: &Config{
				BaseURL: "https://www.mysite.com",
				RedirectHeaders: []string{
					"header1",
				},
				RedirectionMap:  redirectMaps,
				DefaultPath:     "default/",
				RedirectCookies: []string{},
			},
			next: func(t *testing.T) http.Handler {
				t.Helper()
				return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					t.Errorf("got called")
				})
			},
			assertFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				location := resp.HeaderMap["Location"][0]
				want := "https://www.mysite.com/b/"
				if location != want {
					t.Errorf("wanted url %s, got %s", want, location)
				}
			},
		},
		{
			headers: http.Header{
				"header1": []string{"A"},
			},
			name: "baseurl with trailing slash",
			Config: &Config{
				BaseURL: "https://www.mysite.com/",
				RedirectHeaders: []string{
					"header1",
				},
				RedirectionMap:  redirectMaps,
				RedirectCookies: []string{},
				DefaultPath:     "default/",
			},
			next: func(t *testing.T) http.Handler {
				t.Helper()
				return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					t.Errorf("got called")
				})
			},
			assertFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				location := resp.HeaderMap["Location"][0]
				want := "https://www.mysite.com/a/"
				if location != want {
					t.Errorf("wanted url %s, got %s", want, location)
				}
			},
		},
		{
			headers: http.Header{
				"header1": []string{"D"},
			},
			name: "non-existent map",
			Config: &Config{
				BaseURL: "https://www.mysite.com/",
				RedirectHeaders: []string{
					"header1",
				},
				RedirectionMap:  redirectMaps,
				RedirectCookies: []string{},
				DefaultPath:     "default/",
			},
			next: func(t *testing.T) http.Handler {
				t.Helper()
				return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					t.Errorf("got called")
				})
			},
			assertFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				location := resp.HeaderMap["Location"][0]
				want := "https://www.mysite.com/default/"
				if location != want {
					t.Errorf("wanted url %s, got %s", want, location)
				}
			},
		},
		{
			headers: http.Header{
				"header1": []string{"B"},
				"header2": []string{"C"},
			},
			name: "multiple headers",
			Config: &Config{
				BaseURL: "https://www.mysite.com/",
				RedirectHeaders: []string{
					"header1",
					"header2",
				},
				RedirectionMap:  redirectMaps,
				RedirectCookies: []string{},
				DefaultPath:     "default/",
			},
			next: func(t *testing.T) http.Handler {
				t.Helper()
				return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					t.Errorf("got called")
				})
			},
			assertFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				location := resp.HeaderMap["Location"][0]
				want := "https://www.mysite.com/b/"
				if location != want {
					t.Errorf("wanted url %s, got %s", want, location)
				}
			},
		},
		{
			cookie: &http.Cookie{
				Name:  "header2",
				Value: "C",
			},
			headers: http.Header{
				"header1": []string{"b"},
			},
			name: "cookie headers",
			Config: &Config{
				BaseURL: "https://www.mysite.com/",
				RedirectHeaders: []string{
					"header1",
				},
				RedirectCookies: []string{
					"header2",
				},
				RedirectionMap: redirectMaps,
				DefaultPath:    "default/",
			},
			next: func(t *testing.T) http.Handler {
				t.Helper()
				return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					t.Errorf("got called")
				})
			},
			assertFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				location := resp.HeaderMap["Location"][0]
				want := "https://www.mysite.com/c/"
				if location != want {
					t.Errorf("wanted url %s, got %s", want, location)
				}
			},
		},
		{
			name: "no redirect",
			Config: &Config{
				BaseURL:         "https://www.mysite.com/noexist",
				RedirectHeaders: []string{"dummy"},
				RedirectionMap:  redirectMaps,
				RedirectCookies: []string{},
				DefaultPath:     "default/",
			},
			next: func(t *testing.T) http.Handler {
				t.Helper()
				return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					rw.Write([]byte("hello"))
				})
			},
			assertFunc: func(t *testing.T, resp *httptest.ResponseRecorder) {
				b := resp.Body
				str, _ := ioutil.ReadAll(b)
				if string(str) != "hello" {
					t.Errorf("wanted hello, got %s", string(str))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			handler, err := New(ctx, tt.next(t), tt.Config, "redirect-test")
			if err != nil {
				t.Fatalf("error with new redirect: %+v", err)
			}
			recorder := httptest.NewRecorder()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, tt.Config.BaseURL, nil)
			if err != nil {
				t.Errorf("error with new request: %+v", err)
			}
			req.Header = tt.headers
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			handler.ServeHTTP(recorder, req)

			tt.assertFunc(t, recorder)
		})
	}
}
