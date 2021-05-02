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

var (
	configFile   string
	printVersion bool
)

func init() {
	flag.StringVar(&configFile, "c", "config.toml", "Config file (TOML format)")
	flag.BoolVar(&printVersion, "v", false, "Print version information and exit")
	flag.Parse()
}

type mailTemplateInput struct {
	Hostname string
	Letter   string
	Limit    string

	Free, Total, Used string

	UsedPercent float32
	FreePercent float32

	Version string
}

func main() {
	if printVersion {
		fmt.Printf("hddwatcher version %s built at %s commit %s\n", version, builtAt, commit)
		return
	}

	logger := log.New(os.Stdout, "[hddwatcher] ", log.LstdFlags|log.Lmsgprefix)
	// Read config
	if configFile == "" {
		fmt.Println("Must pass config file.")
		flag.Usage()
		os.Exit(1)
	}
	absPath, _ := filepath.Abs(configFile)
	logger.Printf("Reading config file %q", absPath)
	cfg, err := configFromFile(configFile)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Println("Validating configuration sanity")
	acceptableLanguages := make([]string, 0, len(subjects))
	for k := range subjects {
		acceptableLanguages = append(acceptableLanguages, k[:])
	}
	if err := cfg.validate(acceptableLanguages); err != nil {
		logger.Fatal(err)
	}
	logger.Println(" > Config OK")

	// Check disk space
	logger.Printf("Checking disk space for %q", cfg.Drive.Path)
	_, total, free, err := lib.GetSpace(cfg.Drive.Path)
	if err != nil {
		// TODO
		logger.Fatal(err)
	}
	logger.Printf(" > Free space:  %s", humanize.IBytes(uint64(free)))
	logger.Printf(" > Total space: %s", humanize.IBytes(uint64(total)))
	if cfg.Drive.LimitBytes < free {
		logger.Printf(
			"Limit (< %s on %q) not reached - skipping email notification",
			humanize.IBytes(uint64(cfg.Drive.LimitBytes)),
			cfg.Drive.Path,
		)
		logger.Println("Done")
		return
	}

	// MAIL

	t := templates.Lookup(cfg.Mail.Language + ".html")
	if t == nil {
		log.Fatalf("no template for language %q", cfg.Mail.Language)
	}
	b := bytes.Buffer{}
	hn, _ := os.Hostname()
	i := mailTemplateInput{
		Hostname:    hn,
		Letter:      cfg.Drive.Path,
		Limit:       humanize.IBytes(cfg.Drive.LimitBytes),
		Free:        humanize.IBytes(free),
		Total:       humanize.IBytes(total),
		Used:        humanize.IBytes(total - free),
		FreePercent: 100.0 * float32(free) / float32(total),
		Version:     version,
	}
	i.UsedPercent = 100 - i.FreePercent

	if err := t.Execute(&b, i); err != nil {
		logger.Fatalf("Cannot execute template %q: %v", t.Name(), err)
	}

	start := time.Now()
	logger.Printf("Sending notification to %d recepients", len(cfg.Mail.RecepientList))
	auth := smtp.PlainAuth("", cfg.Mail.Username, cfg.Mail.Password, cfg.Mail.Server)
	// The "To:" header is omited so that all recepients
	// receive a Bcc.
	base := strings.Join([]string{
		fmt.Sprintf("From: %s", cfg.Mail.Sender),
		fmt.Sprintf("Subject: " + subjects[cfg.Mail.Language]),
		"Content-Type: text/html; charset=utf-8",
		"",
	}, "\r\n")

	body := []byte(base + b.String())

	addr := fmt.Sprintf("%s:%d", cfg.Mail.Server, cfg.Mail.Port)
	err = smtp.SendMail(addr, auth, cfg.Mail.Sender, cfg.Mail.RecepientList, body)
	if err != nil {
		log.Fatalf("error sending message: %v", err)
	}
	logger.Printf("Sending emails took %v", time.Since(start))
	logger.Printf("Done")
}
