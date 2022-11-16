package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/mail"
	"net/smtp"
	"time"
)

func main() {
	wichtel := []string{
		"some",
		"mail",
		"addresses"}

	wichtel = shuffle(wichtel)

	for i := 0; i < len(wichtel); i++ {
		toAddress := wichtel[i]
		var partner string
		if i == len(wichtel)-1 {
			partner = wichtel[0]
		} else {
			partner = wichtel[i+1]
		}

		sendMail(toAddress, partner)
	}

}

func shuffle(src []string) []string {
	final := make([]string, len(src))
	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(len(src))

	for i, v := range perm {
		final[v] = src[i]
	}
	return final
}

func sendMail(toAddress string, partner string) {
	from := mail.Address{"", "from@mail.address"}
	to := mail.Address{"", toAddress}
	subj := "Auslosung TribeClub Wichtel"
	body := "Hi,\n\nDu hast " + partner + " als Partner!\n\nViel Spaß\n\nDer Wichtelhäuptling"

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := "some.smtp.address:465"

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", "someAuthUsername", "somePassword", host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()
}
