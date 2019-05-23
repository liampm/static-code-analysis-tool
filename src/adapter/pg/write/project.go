package write

import (
	"database/sql"
	"github.com/liampm/static-code-analysis-tool/domain"
	uuid "github.com/satori/go.uuid"
)

type PostgresProjectRepo struct {
	db *sql.DB
	taskRepo domain.TaskReadRepo
	targetRepo domain.TargetReadRepo
}

func ProjectRepo(db *sql.DB, taskRepo domain.TaskReadRepo, targetRepo domain.TargetReadRepo) *PostgresProjectRepo {
	return &PostgresProjectRepo{
		db: db,
		taskRepo: taskRepo,
		targetRepo: targetRepo,
	}
}

func (repo *PostgresProjectRepo) Save(project domain.Project) {
	_, err := repo.Find(project.Id)

	if err == sql.ErrNoRows {
		repo.insert(project)
		return
	} else if err == nil {
		repo.update(project)
		return
	}

	panic(err) // Panic whilst we're in development
}

func (repo *PostgresProjectRepo) Find(id uuid.UUID) (project domain.Project, err error) {
	project, err := repo.find(id)

	if err != nil {
		return domain.Project{}, err
	}

	project.Targets, _ = repo.targetRepo.AllForProject(project.Id)
	project.Tasks = repo.taskRepo.AllForProject(project.Id)

	return project, nil
}

func (repo *PostgresProjectRepo) find(id uuid.UUID) (project domain.Project, err error) {
	project = domain.Project{}

	row := repo.db.QueryRow("SELECT * FROM project WHERE id = $1", id)

	// Populate the entity with the information from the row
	err = row.Scan(&project.Id, &project.Name)

	if err == sql.ErrNoRows {
		return project, err // Return the error so that it can be dealt with
	} else if err != nil {
		panic(err) // Panic whilst we're in development
	}

	return project, nil
}

func (repo *PostgresProjectRepo) insert(project domain.Project) {
	_, err := repo.db.Exec("INSERT INTO project (id, name) VALUES ($1, $2)", project.Id, project.Name)

	if err != nil {
		panic(err) // Panic whilst we're in development
	}
}

func (repo *PostgresProjectRepo) update(project domain.Project) {
	_, err := repo.db.Exec("UPDATE project SET name = $2 WHERE id = $1", project.Id, project.Name)

	if err != nil {
		panic(err) // Panic whilst we're in development
	}
}
