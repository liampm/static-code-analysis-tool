package write

import (
	"database/sql"
	"github.com/liampm/static-code-analysis-tool/domain"
	uuid "github.com/satori/go.uuid"
)

type PostgresTaskRepo struct {
	db *sql.DB
}

func TaskRepo(db *sql.DB) *PostgresTaskRepo {
	return &PostgresTaskRepo{db: db}
}

func (repo *PostgresTaskRepo) Save(task domain.TaskInstance) {
	_, err := repo.Find(task.Id)

	if err == sql.ErrNoRows {
		repo.insert(task)
		return
	} else if err == nil {
		repo.update(task)
		return
	}

	panic(err) // Panic whilst we're in development
}

func (repo *PostgresTaskRepo) Find(id uuid.UUID) (task domain.TaskInstance, err error) {
	task = domain.TaskInstance{}

	row := repo.db.QueryRow("SELECT * FROM task WHERE id = $1", id)

	// Populate the entity with the information from the row
	err = row.Scan(&task.Id, &task.Name)

	if err == sql.ErrNoRows {
		return task, err // Return the error so that it can be dealt with
	} else if err != nil {
		panic(err) // Panic whilst we're in development
	}

	return task, nil
}

func (repo *PostgresTaskRepo) insert(task domain.TaskInstance) {
	_, err := repo.db.Exec("INSERT INTO task (id, project_id, task) VALUES ($1, $2, $3)", task.Id, task.ProjectId, task.Name)

	if err != nil {
		panic(err) // Panic whilst we're in development
	}
}

func (repo *PostgresTaskRepo) update(task domain.TaskInstance) {
	_, err := repo.db.Exec("UPDATE task SET project_id = $2, task = $3 WHERE id = $1", task.Id, task.ProjectId, task.Name)

	if err != nil {
		panic(err) // Panic whilst we're in development
	}
}
