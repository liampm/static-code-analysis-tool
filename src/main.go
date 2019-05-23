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
	defer db.Close()

	if err != nil {
		panic(err)
	}

	taskRepo := write.TaskRepo(db)
	targetRepo := write.TargetRepo(db)
	taskReadRepo := read.TaskRepo(db)
	targetReadRepo := read.TargetRepo(db)
	projectRepo := write.ProjectRepo(db, taskReadRepo, targetReadRepo)


	projectController := controller.ProjectController{
		read.ProjectRepo(db),
		projectRepo,
	}
	taskController := controller.TaskController{
		taskReadRepo,
		taskRepo,
	}
	jobController := controller.NewJobController(projectRepo, write.JobRepo(db))

	router.Get("/project", projectController.All())
	router.Get("/project/{id}", projectController.ById())
	router.Post("/project", projectController.Create())
	router.Get("/project/{projectId}/task", taskController.All())
	router.Get("/project/{projectId}/task/{id}", taskController.ById())
	router.Post("/project/{projectId}/task", taskController.Create())
	router.Post("/project/{projectId}/job", jobController.Initiate())

	targetController := controller.TargetController{
		targetReadRepo,
		targetRepo,
	}

	router.Get("/project/{projectId}/target", targetController.AllForProject())
	router.Get("/project/{projectId}/target/{id}", targetController.ById())
	router.Post("/project/{projectId}/target", targetController.Create())

	err = http.ListenAndServe(":"+PORT, router)

	if err != nil {
		panic(err)
	}
}
