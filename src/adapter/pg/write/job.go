package write

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/liampm/static-code-analysis-tool/domain"
	uuid "github.com/satori/go.uuid"
	"log"
)

type PostgresJobRepo struct {
	db *sql.DB
}

func JobRepo(db *sql.DB) *PostgresJobRepo {
	return &PostgresJobRepo{db: db}
}

func (repo *PostgresJobRepo) Save(job domain.Job) {
	_, err := repo.find(job.Id)

	if err == sql.ErrNoRows {
		repo.insert(job)
		return
	} else if err == nil {
		fmt.Println("Updating")
		fmt.Println(job)
		repo.update(job)
		return
	}

	panic(err) // Panic whilst we're in development
}

func (repo *PostgresJobRepo) find(id uuid.UUID) (job domain.Job, err error) {
	job = domain.Job{}

	row := repo.db.QueryRow("SELECT id FROM job WHERE id = $1", id)

	// Populate the entity with the information from the row
	err = row.Scan(&job.Id)

	if err == sql.ErrNoRows {
		return job, err // Return the error so that it can be dealt with
	} else if err != nil {
		panic(err) // Panic whilst we're in development
	}

	return job, nil
}

func (repo *PostgresJobRepo) insert(job domain.Job) {
	_, err := repo.db.Exec("INSERT INTO job (id, project_id) VALUES ($1, $2)", job.Id, job.ProjectId)

	if err != nil {
		panic(err) // Panic whilst we're in development
	}

	for _, analysis := range job.Analyses {
		repo.touchAnalysis(analysis)
	}
}

func (repo *PostgresJobRepo) update(job domain.Job) {
	_, err := repo.db.Exec("UPDATE job SET project_id = $2 WHERE id = $1", job.Id, job.ProjectId)

	if err != nil {
		panic(err) // Panic whilst we're in development
	}

	for _, analysis := range job.Analyses {
		repo.touchAnalysis(analysis)
	}
}

func (repo *PostgresJobRepo) touchAnalysis(analysis domain.Analysis) {
	log.Println("touch")
	log.Println(analysis.JobId)
	row := repo.db.QueryRow("SELECT id FROM analysis WHERE id = $1", analysis.Id)

	// Populate the entity with the information from the row
	err := row.Scan(&analysis.Id)

	jsonReport, _ := json.Marshal(analysis.Report)
	log.Println(analysis.Report)

	if err == sql.ErrNoRows {
		_, err = repo.db.Exec(
			"INSERT INTO analysis (id, job_id, target_id, task_id, report) VALUES ($1, $2, $3, $4, $5)",
			analysis.Id,
			analysis.JobId,
			analysis.TargetId,
			analysis.TaskId,
			jsonReport,
		)

		if err != nil {
			panic(err) // Panic whilst we're in development
		}
		return
	} else if err != nil {
		panic(err) // Panic whilst we're in development
	}

	_, err = repo.db.Exec(
		"UPDATE analysis SET job_id = $2, target_id = $3, task_id = $4, report = $5) WHERE id = $1",
		analysis.Id,
		analysis.JobId,
		analysis.TargetId,
		analysis.TaskId,
		jsonReport,
	)

	if err != nil {
		panic(err) // Panic whilst we're in development
	}

	return
}
