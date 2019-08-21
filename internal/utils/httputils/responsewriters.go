package httputils

import (
	"net/http"
)

type CatchAllResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *CatchAllResponseWriter) WriteHeader(status int) {
	w.status = status
	if status != http.StatusNotFound {
		w.ResponseWriter.WriteHeader(status)
	}
}

func (w *CatchAllResponseWriter) Write(p []byte) (int, error) {
	if w.status != http.StatusNotFound {
		return w.ResponseWriter.Write(p)
	}
	return len(p), nil
}

func CatchAllHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nfrw := &CatchAllResponseWriter{ResponseWriter: w}
		h.ServeHTTP(nfrw, r)
		if nfrw.status == http.StatusNotFound {
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}
