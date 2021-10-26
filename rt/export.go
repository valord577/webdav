package rt

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

// @author valor.

var c cfg

func ServAddr() string {
	return c.server.address
}

func UseTLS() bool {
	return c.server.useTLS
}

func CerFilepath() string {
	return c.server.cerFile
}

func KeyFilepath() string {
	return c.server.keyFile
}

func MinVersionTLS() uint16 {
	return c.server.minTLS
}

func CheckUserAndPass(user, pass string) bool {
	a, ok := c.accounts[user]
	if !ok {
		return false
	}
	return a.pass == pass
}

func ServeWebDAV(user string, w http.ResponseWriter, r *http.Request) {
	a, ok := c.accounts[user]
	if !ok {
		return
	}
	a.scope.ServeHTTP(w, r)
}

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
