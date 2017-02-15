//package GoWorld
package main

import (
    "fmt"
    "log"
    "net"
    "net/mail"
    "net/smtp"
    "crypto/tls"
    "flag"
)

var (
    server = flag.String("mail.server", "", "SMTP Server to send email.")
    from = flag.String("mail.from", "", "sender's email.")
    password = flag.String("mail.password", "", "sender's password")
    to = flag.String("mail.to", "", "reciepent's email address.")
)

// Send SSL/TLS Email
func SendMail(serverName string, sender string, password string, receiver string, subject string, bodyMsg string) {
    //convert from, to into mail address
    from := mail.Address{"", sender}
    to   := mail.Address{"", receiver}

    //headers
    headers := make(map[string]string)
    headers["From"] = from.String()
    headers["To"] = to.String()
    headers["Subject"] = subject

    //message
    bodyMessage := ""
    for k,v := range headers {
        bodyMessage += fmt.Sprintf("%s: %s\r\n", k, v)
    }
    bodyMessage += "\r\n" + bodyMsg

    host, _, _ := net.SplitHostPort(serverName)

    //config tls
    tlsconfig := &tls.Config {
        InsecureSkipVerify: true,
        ServerName: host,
    }

    //call tls.Dial for smtp servers running on 465 that require an ssl connection
    connection, err := tls.Dial("tcp", serverName, tlsconfig)
    if err != nil {
        log.Panic(err)
    }
    //get client
    c, err := smtp.NewClient(connection, host)
    if err != nil {
        log.Panic(err)
    }

    auth := smtp.PlainAuth("",sender, password, host)

    //Auth
    if err = c.Auth(auth); err != nil {
        log.Panic(err)
    }

    //From
    if err = c.Mail(from.Address); err != nil {
        log.Panic(err)
    }
    //To
    if err = c.Rcpt(to.Address); err != nil {
        log.Panic(err)
    }

    // Data
    w, err := c.Data()
    if err != nil {
        log.Panic(err)
    }
    //Write body text
    _, err = w.Write([]byte(bodyMessage))
    if err != nil {
        log.Panic(err)
    }
    //close writer
    err = w.Close()
    if err != nil {
        log.Panic(err)
    }
    //exit
    c.Quit()
}

func main() {
    flag.Parse()
    SendMail(*server, *from, *password, *to, "test", "asdfasdfasdfasdf")
}
