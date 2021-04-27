package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/baza04/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	addr      string
	staticDir string
	dsn       string
}

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	snippets *mysql.SnippetModel
}

func main() {
	conf := config{}
	// Handle flags value to variable
	flag.StringVar(&conf.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&conf.staticDir, "static", "../../ui/static/", "Path to static assets")
	flag.StringVar(&conf.dsn, "dsn", login+":"+pass+"@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(conf.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := application{infoLog, errorLog, &mysql.SnippetModel{DB: db}}
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

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
