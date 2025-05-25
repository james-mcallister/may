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

	mux.Handle("GET /employees/", middlewareLog(entity.Employees(d.Templates, d.db)))
	mux.Handle("POST /employees/", middlewareLog(form.NewEmployee(d.db)))
	mux.Handle("GET /employees/{id}/", middlewareLog(form.Employee(d.Templates, d.db)))
	mux.Handle("PUT /employees/{id}/", middlewareLog(form.UpdateEmployee(d.db)))
	mux.Handle("DELETE /employees/{id}/", middlewareLog(form.DeleteEmployee(d.db)))

	mux.Handle("GET /compensation/", middlewareLog(entity.Compensation(d.Templates, d.db)))
	mux.Handle("POST /compensation/", middlewareLog(form.NewCompensation(d.db)))
	mux.Handle("GET /compensation/{id}/", middlewareLog(form.Compensation(d.Templates, d.db)))
	mux.Handle("PUT /compensation/{id}/", middlewareLog(form.UpdateCompensation(d.db)))
	mux.Handle("DELETE /compensation/{id}/", middlewareLog(form.DeleteCompensation(d.db)))

	mux.Handle("GET /ipts/", middlewareLog(entity.Ipts(d.Templates, d.db)))
	mux.Handle("POST /ipts/", middlewareLog(form.NewIpt(d.db)))
	mux.Handle("GET /ipts/{id}/", middlewareLog(form.Ipt(d.Templates, d.db)))
	mux.Handle("PUT /ipts/{id}/", middlewareLog(form.UpdateIpt(d.db)))
	mux.Handle("DELETE /ipts/{id}/", middlewareLog(form.DeleteIpt(d.db)))

	mux.Handle("GET /material/", middlewareLog(entity.Material(d.Templates, d.db)))
	mux.Handle("POST /material/", middlewareLog(form.NewMaterial(d.db)))
	mux.Handle("GET /material/{id}/", middlewareLog(form.Material(d.Templates, d.db)))
	mux.Handle("PUT /material/{id}/", middlewareLog(form.UpdateMaterial(d.db)))
	mux.Handle("DELETE /material/{id}/", middlewareLog(form.DeleteMaterial(d.db)))

	mux.Handle("GET /networks/", middlewareLog(entity.Networks(d.Templates, d.db)))
	mux.Handle("POST /networks/", middlewareLog(form.NewNetwork(d.db)))
	mux.Handle("GET /networks/{id}/", middlewareLog(form.Network(d.Templates, d.db)))
	mux.Handle("PUT /networks/{id}/", middlewareLog(form.UpdateNetwork(d.db)))
	mux.Handle("DELETE /networks/{id}/", middlewareLog(form.DeleteNetwork(d.db)))

	mux.Handle("GET /projects/", middlewareLog(entity.Projects(d.Templates, d.db)))
	mux.Handle("POST /projects/", middlewareLog(form.NewProject(d.db)))
	mux.Handle("GET /projects/{id}/", middlewareLog(form.Project(d.Templates, d.db)))
	mux.Handle("PUT /projects/{id}/", middlewareLog(form.UpdateProject(d.db)))
	mux.Handle("DELETE /projects/{id}/", middlewareLog(form.DeleteProject(d.db)))
	return nil
}
