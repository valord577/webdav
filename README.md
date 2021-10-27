WebDAV
======

![Go Report Card](https://goreportcard.com/badge/github.com/valord577/webdav)
![Github Release](https://img.shields.io/github/v/release/valord577/webdav.svg)
![GitHub License](https://img.shields.io/github/license/valord577/webdav)

Web Distributed Authoring and Versioning (WebDAV) is an extension to Hypertext Transfer Protocol (HTTP) that defines how basic file functions such as copy, move, delete, and create are performed by using HTTP.

Usage
------

The binary file of WebDAV is a CLI application.

<details>
<summary>
- Print help information.
</summary>

```text
[user@host webdav]% ./out/bin/webdav -h
A Lightweight WebDAV Server.

Usage:
  webdav <command> [arguments]

The available commands are:
  info    Print information.
  serv    Startup webdav server.

Use "webdav <command> -h" for more information about a command.
```
</details>

<details>
<summary>
- Print build and version information.
</summary>

```text
[user@host webdav]% ./out/bin/webdav info
webdav v1.0 2021-10-27 go1.16.8 linux/amd64
```
</details>

<details>
<summary>
- Startup a WebDAV server.
</summary>

```text
[user@host webdav]% ./out/bin/webdav serv -h
Startup webdav server.

Usage:
  webdav serv [flags...]

The available flags are:
  -c                Declare the configuration file path

[user@host webdav]% ./out/bin/webdav serv -c out/cfg/app.jsonc
2021-10-27 19:28:22        INFO    cmd/serv.go:46  activated cfg file: out/cfg/app.jsonc
2021-10-27 19:28:22        INFO    serve/router.go:22      webdav server is starting at [:60080]

```
</details>

Configuration
------

Please refer to the ['app.jsonc'](rt/app.jsonc) for more configurations. There, you can customize the configurations for WebDAV server.

Systemd
------

Please refer to the ['webdav.service.example'](systemd/webdav.service.example) to manage WebDAV server as a service.

Development
------

<details>
<summary>
- Build project.
</summary>

```text
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o out/bin/webdav github.com/valord577/webdav
```
</details>

Changes
------

See the [CHANGES](CHANGE.md) for changes.

License
------

See the [LICENSE](LICENSE) for Rights and Limitations (MIT).
