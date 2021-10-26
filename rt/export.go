package rt

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

// @author valor.

var c cfg

// ServAddr returns the address of webdav server.
func ServAddr() string {
	return c.server.address
}

// UseTLS returns a boolean value
// that whether to use tls for webdav server.
func UseTLS() bool {
	return c.server.useTLS
}

// CerFilepath returns the path of cert file.
func CerFilepath() string {
	return c.server.cerFile
}

// KeyFilepath returns the path of key file.
func KeyFilepath() string {
	return c.server.keyFile
}

// MinVersionTLS returns the minimum supported version of webdav server.
func MinVersionTLS() uint16 {
	return c.server.minTLS
}

// CheckUserAndPass verifies the correctness of users and passwords.
func CheckUserAndPass(user, pass string) bool {
	a, ok := c.accounts[user]
	if !ok {
		return false
	}
	return a.pass == pass
}

// ServeWebDAV processes the client's request
// according to the webdav protocol.
func ServeWebDAV(user string, w http.ResponseWriter, r *http.Request) {
	a, ok := c.accounts[user]
	if !ok {
		return
	}
	a.scope.ServeHTTP(w, r)
}

// ReadInFile reads the configuration file of webdav server.
func ReadInFile(filepath string) error {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0444)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	bs, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	jsonBs, err := ignoreComments(bs)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBs, &c)
}
