package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/liampm/static-code-analysis-tool/adapter/controller"
	"github.com/liampm/static-code-analysis-tool/adapter/pg/read"
	"github.com/liampm/static-code-analysis-tool/adapter/pg/write"
	_ "github.com/lib/pq"
	"net/http"
	"os"
)

func main() {
	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "3001"
	}

	router := chi.NewRouter()

	db, err := sql.Open("postgres", fmt.Sprintf("host=db user=%s dbname=%s password=%s sslmode=disable", "dev", "scat", "dev"))

	if err != nil {
		panic(err)
	}

	projectController := controller.ProjectController{
		read.ProjectRepo(db),
		write.ProjectRepo(db),
	}

	router.Get("/project", projectController.All())
	router.Get("/project/{id}", projectController.ById())
	router.Post("/project", projectController.Create())

	targetController := controller.TargetController{
		read.TargetRepo(db),
		write.TargetRepo(db),
	}

	router.Get("/target/project/{projectId}", targetController.AllForProject())
	router.Get("/target/{id}", targetController.ById())
	router.Post("/target", targetController.Create())

	err = http.ListenAndServe(":"+PORT, router)

	if err != nil {
		panic(err)
	}
}
