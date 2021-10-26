package rt

// @author valor.

// appJsonc is a structure that matches 'app.jsonc'.
type appJsonc struct {
	Server   appServ   `json:"server"`
	Logger   logServ   `json:"logger"`
	Accounts []appAcct `json:"accounts"`
}

type appServ struct {
	Address string `json:"address"`

	UseTLS  bool   `json:"useTLS"`
	CerFile string `json:"cerFile"`
	KeyFile string `json:"keyFile"`
	MinTLS  string `json:"minTLS"`
}

type logServ struct {
	Output  string `json:"output"`
	Level   string `json:"level"`
	TimeFmt string `json:"tmfmt"`
	Verbose bool   `json:"verbose"`

	Highlight bool `json:"highlight"`

	MaxLineNum  int    `json:"maxLineNum"`
	LogfileName string `json:"logfile"`
}

type appAcct struct {
	User  string `json:"username"`
	Pass  string `json:"passcode"`
	Scope string `json:"scope"`
}
