package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/LidoHon/LetsGO-snippetBox.git/internal/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type application struct {
	infoLog  		*log.Logger
	errorLog 		*log.Logger
	snippets 		*models.SnippetModel
	templateCache 	map[string]*template.Template
}

func main() {

		// addr :=flag.String("addr", ":"+port, "HTTP network address")

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set in the environment")
	}

	// Define the DSN with the loaded DB_URL
	dsn := flag.String("dsn", dbURL, "MySQL data source name")
	flag.Parse()

	// Set up loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Open database connection
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err !=nil{
		errorLog.Fatal(err)
	}

	// Initialize application struct
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB:db},
		templateCache: templateCache,
	}

	// Set up HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000" 
	}

	srv := &http.Server{
		Addr:     ":" + port,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Print("Server is running on port ", port)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// openDB opens a connection to the database and checks if it's alive
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
