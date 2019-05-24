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
	"github.com/go-chi/cors"
)

func main() {
	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "3001"
	}

	router := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:7878"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	router.Use(cors.Handler)

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
	targetController := controller.TargetController{
		targetReadRepo,
		targetRepo,
	}
	taskController := controller.TaskController{
		taskReadRepo,
		taskRepo,
	}
	jobController := controller.NewJobController(projectRepo, read.JobRepo(db), write.JobRepo(db))

	router.Get("/project", projectController.All())
	router.Get("/project/{id}", projectController.ById())
	router.Post("/project", projectController.Create())
	router.Get("/project/{projectId}/task", taskController.All())
	router.Get("/project/{projectId}/task/{id}", taskController.ById())
	router.Post("/project/{projectId}/task", taskController.Create())
	router.Get("/project/{projectId}/job", jobController.ByProjectId())
	router.Post("/project/{projectId}/job", jobController.Initiate())
	router.Get("/project/{projectId}/target", targetController.AllForProject())
	router.Get("/project/{projectId}/target/{id}", targetController.ById())
	router.Post("/project/{projectId}/target", targetController.Create())

	err = http.ListenAndServe(":"+PORT, router)

	if err != nil {
		panic(err)
	}
}
