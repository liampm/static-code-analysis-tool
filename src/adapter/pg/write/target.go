package write

import (
	"database/sql"
	"github.com/liampm/static-code-analysis-tool/domain"
	uuid "github.com/satori/go.uuid"
)

type PostgresTargetRepo struct {
	db *sql.DB
}

func TargetRepo(db *sql.DB) *PostgresTargetRepo {
	return &PostgresTargetRepo{db: db}
}

func (repo *PostgresTargetRepo) Save(target domain.Target) {
	_, err := repo.find(target.Id)

	if err == sql.ErrNoRows {
		repo.insert(target)
		return
	} else if err == nil {
		repo.update(target)
		return
	}

	panic(err) // Panic whilst we're in development
}

func (repo *PostgresTargetRepo) find(id uuid.UUID) (target domain.Target, err error) {
	target = domain.Target{}

	row := repo.db.QueryRow("SELECT * FROM target WHERE id = $1", id)

	// Populate the entity with the information from the row
	err = row.Scan(&target.Id, &target.ProjectId, &target.Name)

	if err == sql.ErrNoRows {
		return target, err // Return the error so that it can be dealt with
	} else if err != nil {
		panic(err) // Panic whilst we're in development
	}

	return target, nil
}

func (repo *PostgresTargetRepo) insert(target domain.Target) {
	_, err := repo.db.Exec("INSERT INTO target (id, name) VALUES ($1, $2, $3)", target.Id, target.ProjectId, target.Name)

	if err != nil {
		panic(err) // Panic whilst we're in development
	}
}

func (repo *PostgresTargetRepo) update(target domain.Target) {
	_, err := repo.db.Exec("UPDATE target SET id = $1, projectId = $2, name = $3", target.Id, target.ProjectId, target.Name)

	if err != nil {
		panic(err) // Panic whilst we're in development
	}
}
