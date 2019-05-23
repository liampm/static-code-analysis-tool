package read

import (
	"database/sql"
	"github.com/liampm/static-code-analysis-tool/domain"
	"github.com/satori/go.uuid"
)

type PostgresTargetReadRepo struct {
	db *sql.DB
}

func TargetRepo(db *sql.DB) *PostgresTargetReadRepo {
	return &PostgresTargetReadRepo{db: db}
}

func (repo *PostgresTargetReadRepo) Find(id uuid.UUID) (domain.Target, error) {
	target := domain.Target{}

	row := repo.db.QueryRow("SELECT * FROM target WHERE id = $1", id)

	// Populate the entity with the information from the row
	err := row.Scan(&target.Id, &target.ProjectId, &target.Name)

	if err == sql.ErrNoRows {
		return target, err // Return the error so that it can be dealt with
	} else if err != nil {
		panic(err) // Panic whilst we're in development
	}

	return target, nil
}

func (repo *PostgresTargetReadRepo) AllForProject(id uuid.UUID) []domain.Target {
	rows, err := repo.db.Query("SELECT * FROM target WHERE projectId = $1", id)

	if err != nil {
		panic(err) // Panic whilst we're in development
	}

	targets := []domain.Target{}

	for rows.Next() {
		target := domain.Target{}
		rows.Scan(&target.Id,  &target.ProjectId, &target.Name)
		targets = append(targets, target)
	}

	return targets
}
