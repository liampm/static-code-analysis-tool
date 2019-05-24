package read

import (
	"database/sql"
	"encoding/json"
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
	config := ""

	row := repo.db.QueryRow("SELECT * FROM target WHERE id = $1", id)

	// Populate the entity with the information from the row
	err := row.Scan(&target.Id, &target.ProjectId, &target.Name, &config)

	if err == sql.ErrNoRows {
		return target, err // Return the error so that it can be dealt with
	} else if err != nil {
		panic(err) // Panic whilst we're in development
	}

	target.Config = loadConfig(config)

	return target, nil
}

func (repo *PostgresTargetReadRepo) AllForProject(id uuid.UUID) ([]domain.Target, error) {
	rows, err := repo.db.Query("SELECT id, project_id, name, config FROM target WHERE project_id = $1", id)

	if err != nil {
		return nil, err // Panic whilst we're in development
	}

	if err = rows.Err(); err != nil {
		panic(err) // Error related to the iteration of rows
	}

	targets := []domain.Target{}

	for rows.Next() {
		target := domain.Target{}
		dbConfig := ""
		rows.Scan(&target.Id, &target.ProjectId, &target.Name, &dbConfig)
		target.Config = loadConfig(dbConfig)

		targets = append(targets, target)
	}

	return targets, nil
}

type repoConfig struct {
	Type domain.TargetType `json:"type"`
	Details domain.RepoDetails `json:"details"`
}


func loadConfig (configString string) domain.TargetConfiguration {
	config := repoConfig{}
	_ = json.Unmarshal([]byte(configString), &config)

	return domain.TargetConfiguration{
		Type:    config.Type,
		Details: config.Details,
	}
}