package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/LidoHon/LetsGO-snippetBox.git/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type application struct {
	infoLog  		*log.Logger
	errorLog 		*log.Logger
	snippets 		*models.SnippetModel
	templateCache 	map[string]*template.Template
	formDecoder 	*form.Decoder
	sessionManager 	*scs.SessionManager
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


	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime =12 * time.Hour

	// Initialize application struct
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB:db},
		templateCache: templateCache,
		formDecoder: formDecoder,
		sessionManager: sessionManager,
	}

	// Set up HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000" 
	}


	// only elliptic curves with
// assembly implementations are used.since the others are very cpu intensive
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		// can configure the minimum and maximum TLS versions as well if we know all computers in our user base support tls 1.2 or whatever the case is

		// MinVersion: tls.VersionTLS12,
		// MaxVersion: tls.VersionTLS12,


			/* cipher suites are a set of algorithms that are used to encrypt and decrypt data in a secure way */

		// we can also restrict cipher suites  for example 
		// CipherSuites: []uint16{
		// 	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		// 	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		// 	tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		// 	tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		// 	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		// 	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		// 	},

	
	}

	srv := &http.Server{
		Addr:     ":" + port,
		ErrorLog: errorLog,
		Handler:  app.routes(),
		TLSConfig: tlsConfig,
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		// we can also limit maximum header length..but Go always adds an additional 4096 bytes
		// MaxHeaderBytes: 524288,
	}

	infoLog.Print("Server is running on port ", port)
	err = srv.ListenAndServeTLS("../../tls/cert.pem", "../../tls/key.pem")
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
