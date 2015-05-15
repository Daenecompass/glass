package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"

	"github.com/advanderveer/timer/model"
)

var Version = "0.0.0"
var Build = "gobuild"

var mbu = flag.Duration("mbu", time.Minute*6, "The minimal billable unit")
var bind = flag.String("bind", ":0", "Address to bind the Daemon to")

func main() {
	flag.Parse()

	timer := NewTimer(*mbu)
	svr, err := NewServer(*bind, timer)
	if err != nil {
		log.Fatal(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(errwrap.Wrapf("Failed to fetch current working dir: {{err}}", err))
	}

	m := model.New(dir)
	info := model.NewDeamon(dir, svr.Addr())
	err = m.UpsertDaemonInfo(info)
	if err != nil {
		log.Fatal(errwrap.Wrapf("Failed write Daemon info: {{err}}", err))
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		svr.Stop()
	}()

	log.Printf("Listening on '%s'", svr.Addr())
	err = svr.Start()
	if err != nil && !strings.Contains(err.Error(), "closed network connection") {
		log.Fatal(err)
	}

	log.Printf("Writing information to database...")
	info.Addr = ""
	err = m.UpsertDaemonInfo(info)
	if err != nil {
		log.Fatal(errwrap.Wrapf("Failed write Daemon info: {{err}}", err))
	}

	log.Printf("Done")
}
