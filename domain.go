package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/james-mcallister/may/database"
	"github.com/james-mcallister/may/entity"
	"github.com/james-mcallister/may/form"
)

type MayPage interface {
	InitTemplate() error
}

type Domain struct {
	db        *sql.DB
	Templates *template.Template
}

func NewDomain() Domain {
	d, err := database.NewDB()
	if err != nil {
		panic(err)
	}

	dbConn, err := d.Connect()
	if err != nil {
		panic(err)
	}

	return Domain{
		db: dbConn,
	}
}

func (d Domain) Init() {
	if err := database.InitDB(d.db); err != nil {
		panic(err)
	}
}

func home(arg string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>Home Page - " + arg + "</h1>"))
	})
}

func Logger(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Printf("%s %s", r.Method, r.URL)
			next.ServeHTTP(w, r)
		})
	}
}

func initRoutes(mux *http.ServeMux, logger *log.Logger, d Domain) error {
	middlewareLog := Logger(logger)

	mux.Handle("GET /home/", middlewareLog(home("Test")))

	// need to map the template to the route somehow
	mux.Handle("GET /employees/", middlewareLog(entity.Employees(d.Templates, d.db)))
	mux.Handle("POST /employees/", middlewareLog(form.NewEmployee(d.db)))
	mux.Handle("GET /employees/{id}/", middlewareLog(form.Employee(d.Templates, d.db)))
	mux.Handle("PUT /employees/{id}/", middlewareLog(form.UpdateEmployee(d.db)))
	mux.Handle("DELETE /employees/{id}/", middlewareLog(form.DeleteEmployee(d.db)))
	return nil
}
