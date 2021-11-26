package rt

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	log "github.com/valord577/webdav/logger"
	"golang.org/x/net/webdav"
)

// @author valor.

// cfg is the sum of all configurations.
// For more configuration information, please check 'app.jsonc'.
type cfg struct {
	server   serv             // json:"server"
	accounts map[string]*acct // json:"accounts"
}

// serv defines the behavior of the HTTP server.
type serv struct {
	address string // json:"address"

	useTLS  bool   // json:"useTLS"
	cerFile string // json:"cerFile"
	keyFile string // json:"keyFile"
	minTLS  uint16 // json:"minTLS"
}

// acct contains user information and a valid handler for WebDAV.
type acct struct {
	user  string          // json:"username"
	pass  string          // json:"passcode"
	scope *webdav.Handler // json:"scope"
}

// UnmarshalJSON can decode 'app.jsonc' into cfg
func (c *cfg) UnmarshalJSON(bs []byte) error {
	length := len(bs)
	if length == 0 {
		return errors.New("empty configuration file")
	}

	ddl := appJsonc{}
	err := json.Unmarshal(bs, &ddl)
	if err != nil {
		return err
	}

	// logger
	switch ddl.Logger.Output {
	case "logfile":
		log.InitLogfile(
			ddl.Logger.TimeFmt, ddl.Logger.LogfileName, ddl.Logger.MaxLineNum, ddl.Logger.Level)

	default:
		log.InitConsole(ddl.Logger.TimeFmt, ddl.Logger.Highlight, ddl.Logger.Level)
	}

	// accounts
	acctLength := len(ddl.Accounts)
	if acctLength == 0 {
		return errors.New("at least one account is required")
	}

	as := make(map[string]*acct, acctLength)
	for _, ddlAcct := range ddl.Accounts {
		fs := &webdav.Handler{
			FileSystem: webdav.Dir(ddlAcct.Scope),
			LockSystem: webdav.NewMemLS(),

			Logger: getWebdavLogHandler(ddl.Logger.Verbose),
		}

		a := acct{
			user:  ddlAcct.User,
			pass:  ddlAcct.Pass,
			scope: fs,
		}
		as[a.user] = &a
	}

	// server
	s := serv{
		address: ddl.Server.Address,
		useTLS:  ddl.Server.UseTLS,
		cerFile: ddl.Server.CerFile,
		keyFile: ddl.Server.KeyFile,
	}
	switch ddl.Server.MinTLS {
	case "v1.0":
		s.minTLS = tls.VersionTLS10
	case "v1.1":
		s.minTLS = tls.VersionTLS11
	case "v1.2":
		s.minTLS = tls.VersionTLS12
	case "v1.3":
		s.minTLS = tls.VersionTLS13
	default:
		s.minTLS = tls.VersionTLS12
	}

	c.accounts = as
	c.server = s
	return nil
}

func getWebdavLogHandler(verbose bool) func(*http.Request, error) {
	return func(r *http.Request, e error) {
		if e != nil {
			log.Errorf("webdav server error: %s\n%s", e.Error(), webdavHttpMessage(r))
		} else {
			if verbose {
				log.Infof("webdav server http request messages: \n%s", webdavHttpMessage(r))
			}
		}
	}
}

func webdavHttpMessage(r *http.Request) string {

	var (
		bs  []byte = nil
		err error  = nil
	)

	if r.ContentLength > 0 && r.Body != nil {
		bs, err = io.ReadAll(r.Body)
		if err != nil {
			log.Errorf("read webdav request body err: %s", err.Error())
		}
	}

	b := &strings.Builder{}
	b.WriteString(">>>>>>>>>>>>>>>>>>\n")

	b.WriteString(r.Method)
	b.WriteString(" ")
	b.WriteString(r.Proto)
	b.WriteString("\n")

	b.WriteString("HOST: ")
	b.WriteString(r.Host)
	b.WriteString("\n")

	b.WriteString("URI: ")
	b.WriteString(r.RequestURI)
	b.WriteString("\n")

	b.WriteString("FROM: ")
	b.WriteString(r.RemoteAddr)
	b.WriteString("\n\n")

	headers := r.Header
	for key, values := range headers {
		for _, v := range values {
			b.WriteString("Header[")
			b.WriteString(key)
			b.WriteString("]: ")
			b.WriteString(v)
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")

	if bs != nil {
		b.WriteString("BODY: ")
		b.Write(bs)
		b.WriteString("\n\n")
	}

	b.WriteString("<<<<<<<<<<<<<<<<<<")
	return b.String()
}

func ignoreComments(bs []byte) ([]byte, error) {
	length := len(bs)
	if length == 0 {
		return []byte{}, nil
	}

	b := &bytes.Buffer{}

	// state      example        description
	//   0           -          Initial state
	//
	//                          Encountered '/' in state 0,
	//   1      int a = b; /    which means you may encounter a comment,
	//                          then enter state 1
	//
	//                          Encountered '/' in state 1,
	//   4      int a = b; //   which means entering the comment part,
	//                          then enter state 4
	//
	//                          Encountered '"' in state 0,
	//   7      char s[] = "    which means entering the string constant,
	//                          then enter state 7
	state := 0

	for i := 0; i < length; i++ {
		switch state {
		case 0:
			switch bs[i] {
			case '/':
				state = 1
			case '"':
				state = 7
				b.WriteByte(bs[i])
			case ' ':
			case '\r':
			case '\n':
			default:
				b.WriteByte(bs[i])
			}

		case 1:
			switch bs[i] {
			case '/':
				state = 4
			default:
				state = 0

				b.WriteByte(bs[i-1])
				b.WriteByte(bs[i])
			}

		case 4:
			if bs[i] == '\n' {
				state = 0
			}
			if bs[i] == '\r' {
				state = 0

				if bs[i+1] == '\n' {
					i++
				}
			}

		case 7:
			switch bs[i] {
			case '"':
				state = 0
				b.WriteByte(bs[i])
			default:
				b.WriteByte(bs[i])
			}

		default:
			return nil, errors.New("invalid state")
		}
	}

	return b.Bytes(), nil
}
