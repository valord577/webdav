package serve

import (
	"net/http"

	log "github.com/valord577/webdav/logger"
	"github.com/valord577/webdav/rt"
)

// @author valor.

// WebDAVServ returns HTTP server for webdav.
func WebDAVServ() *http.Server {
	// Use net/http DefaultServeMux
	http.HandleFunc("/", globalHandler)

	addr := rt.ServAddr()
	serv := &http.Server{
		Addr:    addr,
		Handler: http.DefaultServeMux,
	}
	log.Infof("webdav server is starting at [%s]", addr)
	return serv
}

func globalHandler(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok {
		basicAuthRequired(w)
		log.Warnf("basic auth required from [%s]", r.RemoteAddr)
		return
	}

	ok = rt.CheckUserAndPass(user, pass)
	if !ok {
		basicAuthRequired(w)
		log.Warnf("basic auth failed: [%s]@[%s] from [%s]", user, pass, r.RemoteAddr)
		return
	}

	log.Infof("webdav works for [%s] from [%s]", user, r.RemoteAddr)
	rt.ServeWebDAV(user, w, r)
}

func basicAuthRequired(w http.ResponseWriter) {
	status := http.StatusUnauthorized

	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	w.WriteHeader(status)
	w.Write([]byte(http.StatusText(status)))
}
