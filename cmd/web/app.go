package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	addr      string
	staticDir string
}

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	conf := &config{}
	// Handle flags value to variable
	flag.StringVar(&conf.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&conf.staticDir, "static", "../../ui/static/", "Path to static assets")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app := application{infoLog, errorLog}

	server := &http.Server{
		Addr:         conf.addr,
		Handler:      app.routes(conf.staticDir),
		ErrorLog:     errorLog,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	infoLog.Printf("Listening port%s\n", server.Addr)
	errorLog.Fatal(server.ListenAndServe())
}
