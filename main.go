package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xswordsx/hddwatcher/lib"

	"github.com/dustin/go-humanize"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "c", "config.toml", "Config file (TOML format)")
	flag.Parse()
}

type mailInput struct {
	Hostname string
	Letter   string

	Free, Total, Used string
	UsedPercent       float32
	FreePercent       float32
}

func main() {
	// Read config
	if configFile == "" {
		fmt.Println("Must pass config file.")
		flag.Usage()
		os.Exit(1)
	}
	absPath, _ := filepath.Abs(configFile)
	log.Printf("Reading config file %q", absPath)
	cfg, err := configFromFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Validating configuration sanity")
	acceptableLanguages := make([]string, 0, len(subjects))
	for k := range subjects {
		acceptableLanguages = append(acceptableLanguages, k[:])
	}
	if err := cfg.validate(acceptableLanguages); err != nil {
		log.Fatal(err)
	}

	// Check disk space
	log.Printf("Checking disk space for %q", cfg.Drive.Path)
	_, total, free, err := lib.GetSpace(cfg.Drive.Path)
	if err != nil {
		// TODO
		log.Fatal(err)
	}
	if free > cfg.Drive.Limit {
		log.Println("Free space", humanize.Bytes(uint64(free)), "is greater than the limit", humanize.Bytes(uint64(cfg.Drive.Limit)))
		return
	}

	// MAIL

	funcsMap := template.FuncMap{
		"round": func(x float32) string {
			return fmt.Sprintf("%.2f", x)
		},
	}
	t := template.Must(template.New("mail").Funcs(funcsMap).Parse(mailTemplate))
	b := bytes.Buffer{}
	hn, _ := os.Hostname()
	i := mailInput{
		Hostname:    hn,
		Letter:      cfg.Drive.Letter,
		Free:        humanize.Bytes(uint64(free)),
		Total:       humanize.Bytes(uint64(total)),
		Used:        humanize.Bytes(uint64(total - free)),
		FreePercent: 100.0 * float32(free) / float32(total),
	}
	i.UsedPercent = 100 - i.FreePercent
	t.Execute(&b, i)

	auth := smtp.PlainAuth("", cfg.Mail.Username, cfg.Mail.Password, cfg.Mail.Server)
	base := strings.Join([]string{
		fmt.Sprintf("From: %s", cfg.Mail.Sender),
		fmt.Sprintf("To: %s", strings.Join(cfg.Mail.RecepientList, ",")),
		fmt.Sprintf("Subject: Дисковото пространство критично"),
		fmt.Sprintf("Content-Type: text/html; charset=utf-8"),
		"",
	}, "\r\n")

	body := []byte(base + b.String())

	addr := fmt.Sprintf("%s:%d", cfg.Mail.Server, cfg.Mail.Port)
	err = smtp.SendMail(addr, auth, cfg.Mail.Sender, cfg.Mail.RecepientList, body)
	if err != nil {
		log.Fatalf("error sending message: %v", err)
	}
}
