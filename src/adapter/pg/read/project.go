package read

import (
	"database/sql"
	"github.com/liampm/static-code-analysis-tool/domain"
	"github.com/satori/go.uuid"
)

type PostgresProjectReadRepo struct {
	db *sql.DB
}

func ProjectRepo(db *sql.DB) *PostgresProjectReadRepo {
	return &PostgresProjectReadRepo{db: db}
}

func (repo *PostgresProjectReadRepo) Find(id uuid.UUID) (domain.ProjectReference, error) {
	project := domain.ProjectReference{}

	row := repo.db.QueryRow("SELECT * FROM project WHERE id = $1", id)

	// Populate the entity with the information from the row
	err := row.Scan(&project.Id, &project.Name)

	if err == sql.ErrNoRows {
		return project, err // Return the error so that it can be dealt with
	} else if err != nil {
		panic(err) // Panic whilst we're in development
	}

	return project, nil
}

func (repo *PostgresProjectReadRepo) All() []domain.ProjectReference {
	rows, err := repo.db.Query("SELECT * FROM project")

	if err != nil {
		panic(err) // Panic whilst we're in development
	}

	projects := []domain.ProjectReference{}

	for rows.Next() {
		project := domain.ProjectReference{}
		rows.Scan(&project.Id, &project.Name)
		projects = append(projects, project)
	}

	return projects
}
