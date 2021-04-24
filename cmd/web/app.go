package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	addr, staticDir string
}

func main() {
	conf := &config{}
	flag.StringVar(&conf.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&conf.staticDir, "static", "../../ui/static/", "Path to static assets")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fileServer := http.FileServer(http.Dir(conf.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	server := http.Server{
		Addr:         conf.addr,
		Handler:      mux,
		ErrorLog:     errorLog,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	infoLog.Printf("Listening port%s\n", server.Addr)
	errorLog.Fatal(server.ListenAndServe())
}
