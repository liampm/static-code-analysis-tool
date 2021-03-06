package write

import (
	"database/sql"
	"encoding/json"
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
	_, err := repo.Find(target.Id)

	if err == sql.ErrNoRows {
		repo.insert(target)
		return
	} else if err == nil {
		repo.update(target)
		return
	}

	panic(err) // Panic whilst we're in development
}

func (repo *PostgresTargetRepo) Find(id uuid.UUID) (target domain.Target, err error) {
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
	jsonConfig, _ := json.Marshal(target.Config)
	_, err := repo.db.Exec("INSERT INTO target (id, project_id, name, config) VALUES ($1, $2, $3, $4)", target.Id, target.ProjectId, target.Name, jsonConfig)

	if err != nil {
		panic(err) // Panic whilst we're in development
	}
}

func (repo *PostgresTargetRepo) update(target domain.Target) {
	jsonConfig, _ := json.Marshal(target.Config)
	_, err := repo.db.Exec("UPDATE target SET project_id = $2, name = $3, config = $4 WHERE id = $1", target.Id, target.ProjectId, target.Name, jsonConfig)

	if err != nil {
		panic(err) // Panic whilst we're in development
	}
}
