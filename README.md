go-postfixadmin
=================================

A simple ajax app+api in golang to manage a [postfixadmin](https://github.com/postfixadmin) flavoured postfix server. (WIP and not serious yet)

- It's being used in production running on a dedicated server, with a simple token authentication and plain text auth (legacy reasons!).
- The `/mailbox/*` like others is wip and experiment in returning imap mailboxes via json
- TODO: proper authentiction, encryption, true REST, use postfixadmin encryption login etc.
- Help+advice wanted from golang devs and postfixadmin users

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/daffodil/go-postfixadmin)

# Install

Copy the config file [etc/config-skel.yaml](etc/config-skel.yaml) to eg `config.yaml`

```bash
go get github.com/daffodil/go-postfixadmin
go build

./go-postfixadmin -c /path/to/config.yaml
```
Open [http://pfa.localhost:23456](http://pfa.localhost:23456) (default port)

Visit [main.go](https://github.com/daffodil/go-postfixadmin/blob/master/main.go#L64) for
the url endpoints (TODO: document the urls)

Contributions and feedback welcome

