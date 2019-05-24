package read

import (
	"database/sql"
	"github.com/liampm/static-code-analysis-tool/domain"
	"github.com/satori/go.uuid"
)

type PostgresJobReadRepo struct {
	db *sql.DB
}

func JobRepo(db *sql.DB) *PostgresJobReadRepo {
	return &PostgresJobReadRepo{db: db}
}

func (repo *PostgresJobReadRepo) Find(id uuid.UUID) (domain.JobReference, error) {
	job := domain.JobReference{}

	row := repo.db.QueryRow("SELECT * FROM job WHERE id = $1", id)

	// Populate the entity with the information from the row
	err := row.Scan(&job.Id, &job.ProjectId, &job.Analyses)

	if err == sql.ErrNoRows {
		return job, err // Return the error so that it can be dealt with
	} else if err != nil {
		panic(err) // Panic whilst we're in development
	}

	return job, nil
}

func (repo *PostgresJobReadRepo) AllForProject(project domain.Project) []domain.JobReference {
	rows, err := repo.db.Query("SELECT * FROM job WHERE project_id = $1", project.Id)

	if err != nil {
		panic(err) // Panic whilst we're in development
	}

	jobs := []domain.JobReference{}

	for rows.Next() {
		job := domain.JobReference{}
		rows.Scan(&job.Id, &job.ProjectId, &job.Analyses)
		jobs = append(jobs, job)
	}

	return jobs
}
