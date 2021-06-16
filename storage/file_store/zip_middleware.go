package file_store

import (
	"compress/gzip"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

type GzipHandler struct {
}

func (g *GzipHandler) GzipMiddleware(next httprouter.Handle) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			wrw := NewWrappedResponseWriter(w)
			wrw.Header().Set("Content-Encoding", "gzip")

			next(wrw, r, ps)
			defer wrw.Flush()

			return
		}

		next(w, r, ps)
	}
}

type WrappedResponseWriter struct {
	rw http.ResponseWriter
	gw *gzip.Writer
}

func NewWrappedResponseWriter(rw http.ResponseWriter) *WrappedResponseWriter {
	gw := gzip.NewWriter(rw)

	return &WrappedResponseWriter{rw: rw, gw: gw}
}

func (wr *WrappedResponseWriter) Header() http.Header {
	return wr.rw.Header()
}

func (wr *WrappedResponseWriter) Write(d []byte) (int, error) {
	return wr.gw.Write(d)
}

func (wr *WrappedResponseWriter) WriteHeader(statusCode int) {
	wr.rw.WriteHeader(statusCode)
}

func (wr *WrappedResponseWriter) Flush() {
	wr.gw.Flush()
	wr.gw.Close()
}
