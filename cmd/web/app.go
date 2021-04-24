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

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir(conf.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// f, err := os.OpenFile("../../info.log", os.O_RDWR|os.O_CREATE, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer f.Close()

	server := &http.Server{
		Addr:         conf.addr,
		Handler:      mux,
		ErrorLog:     errorLog,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	infoLog.Printf("Listening port%s\n", server.Addr)
	errorLog.Fatal(server.ListenAndServe())
}
