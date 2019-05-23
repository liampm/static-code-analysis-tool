package read

import (
	"database/sql"
	"github.com/liampm/static-code-analysis-tool/domain"
	"github.com/satori/go.uuid"
)

type PostgresTaskReadRepo struct {
	db *sql.DB
}

func TaskRepo(db *sql.DB) *PostgresTaskReadRepo {
	return &PostgresTaskReadRepo{db: db}
}

func (repo *PostgresTaskReadRepo) Find(id uuid.UUID) (domain.TaskInstance, error) {

	task := domain.TaskInstance{}

	row := repo.db.QueryRow("SELECT * FROM task WHERE id = $1", id)

	// Populate the entity with the information from the row
	err := row.Scan(&task.Id, &task.ProjectId, &task.Name)

	if err == sql.ErrNoRows {
		return task, err // Return the error so that it can be dealt with
	} else if err != nil {
		panic(err) // Panic whilst we're in development
	}

	return task, nil
}

func (repo *PostgresTaskReadRepo) AllForProject(projectId uuid.UUID) []domain.TaskInstance {

	rows, err := repo.db.Query("SELECT * FROM task WHERE project_id = $1", projectId)

	if err != nil {
		panic(err) // Panic whilst we're in development
	}

	tasks := []domain.TaskInstance{}

	for rows.Next() {
		task := domain.TaskInstance{}
		rows.Scan(&task.Id, &task.ProjectId, &task.Name)
		tasks = append(tasks, task)
	}

	return tasks
}


