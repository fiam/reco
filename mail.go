package reco

import (
	"context"
	"net/smtp"
	"strings"
)

const (
	htmlMIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

func Mail(server string, to string, from string) Renderer {
	var auth smtp.Auth
	var hostAndPort string
	if sp := strings.IndexByte(server, '@'); sp >= 0 {
		usernameAndPassword := server[:sp]
		hostAndPort = server[sp+1:]

		var username string
		var password string
		if sp := strings.IndexByte(usernameAndPassword, ':'); sp >= 0 {
			username = usernameAndPassword[:sp]
			password = usernameAndPassword[sp+1:]
		} else {
			username = usernameAndPassword
		}

		host := hostAndPort
		if sp := strings.IndexByte(host, ':'); sp >= 0 {
			host = hostAndPort[:sp]
		}
		auth = smtp.PlainAuth("", username, password, host)
	} else {
		hostAndPort = server
	}

	return func(ctx context.Context, rec *Recovery) error {
		htmlContents, err := HTML(ctx, rec)
		if err != nil {
			return err
		}
		body := "To: " + to + "\r\nSubject: " + rec.Title + "\r\n" + htmlMIME + "\r\n" + string(htmlContents)
		return smtp.SendMail(hostAndPort, auth, from, []string{to}, []byte(body))
	}
}
