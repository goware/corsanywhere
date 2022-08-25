package corsanywhere

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
)

func CORSAnywhereHandler() http.Handler {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(corsAnywhereUsage))
	})
	r.Handle("/*", corsProxy())
	return r
}

func corsProxy() http.Handler {
	director := func(req *http.Request) {
		corsURL := chi.URLParam(req, "*")

		u, err := url.Parse(corsURL)
		if err != nil {
			return
		}

		req.URL.Scheme = u.Scheme
		req.URL.Host = u.Host
		req.URL.Path = u.Path
		req.Host = u.Host

		// NOTE: the req.Query will already be set properly for us

		req.Header.Del("set-cookie")
		req.Header.Del("set-cookie2")
	}

	modifyResponse := func(resp *http.Response) error {
		resp.Header.Set("access-control-allow-origin", "*")
		resp.Header.Set("access-control-max-age", "3000000")
		return nil
	}

	proxy := &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifyResponse,
	}

	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		// TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corsURL := chi.URLParam(r, "*")

		// Verify cors proxy url is valid
		_, err := url.Parse(corsURL)
		if err != nil {
			respondError(w, r, "invalid cors proxy url")
			return
		}

		// Handle pre-flight
		if r.Method == "OPTIONS" {
			handlePreflight(w, r)
			return
		}

		// Require origin header
		if r.Header.Get("origin") == "" {
			respondError(w, r, "origin header is required on the request")
			return
		}

		proxy.ServeHTTP(w, r)
	})
}

func handlePreflight(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w, r)
	w.WriteHeader(200)
}

func respondError(w http.ResponseWriter, r *http.Request, body string) {
	setCORSHeaders(w, r)
	w.WriteHeader(422)
	w.Write([]byte(body))
}

func setCORSHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("access-control-allow-origin", "*")
	w.Header().Set("access-control-max-age", "3000000")

	if r.Header.Get("access-control-request-method") != "" {
		w.Header().Set("access-control-allow-methods", r.Header.Get("access-control-request-method"))
	}

	if r.Header.Get("access-control-request-headers") != "" {
		w.Header().Set("access-control-request-headers", r.Header.Get("access-control-request-headers"))
	}
}

var corsAnywhereUsage = `cors-anywhere usage:

http://localhost:<port>/http(s)://your-domain.com/endpoint

Inspired by https://github.com/Redocly/cors-anywhere
`
