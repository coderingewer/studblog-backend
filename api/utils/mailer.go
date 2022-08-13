package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/k3a/html2text"
	"gopkg.in/mail.v2"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

func SendMail(to string, data EmailData, templateName string) error {

	template, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}
	var body bytes.Buffer

	template = template.Lookup("base.html")
	template.Execute(&body, &data)
	fmt.Println(template.Name())

	m := mail.NewMessage()
	m.SetHeader("From", os.Getenv("MAIL_USER"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)

	m.SetBody("text/html", data.URL)

	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))
	d := mail.NewDialer(os.Getenv("MAIL_SERVER"), 587, os.Getenv("MAIL_USERNAME"), os.Getenv("MAIL_PASSWORD"))
	err = d.DialAndSend(m)
	if err != nil {
		return err
	}
	return nil
}

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	fmt.Println("Am parsing templates...")

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}
