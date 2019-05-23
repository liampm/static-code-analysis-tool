package domain

import uuid "github.com/satori/go.uuid"

type ProjectRepo interface {
	Save(project Project)
	Find(id uuid.UUID) (Project, error)
}
type TaskRepo interface {
	Save(task TaskInstance)
	Find(id uuid.UUID) (TaskInstance, error)
}
type TargetRepo interface {
	Save(target Target)
	Find(id uuid.UUID) (Target, error)
}
type JobRepo interface {
	Save(job Job)
}
type AnalysisRepo interface {
	Save(analysis Analysis)
}

type ProjectReadRepo interface {
	Find(id uuid.UUID) (ProjectReference, error)
	All() []ProjectReference
}
type TaskReadRepo interface {
	Find(id uuid.UUID) (TaskInstance, error)
	AllForProject(projectId uuid.UUID) []TaskInstance
}
type TargetReadRepo interface {
	Find(id uuid.UUID) (Target, error)
	AllForProject(projectId uuid.UUID) ([]Target, error)
}
type JobReadRepo interface {
	Find(id uuid.UUID) (Job, error)
	AllForProject(project Project) []Job
}
type AnalysisReadRepo interface {
	Find(id uuid.UUID) (Analysis, error)
	AllForJob(job Job) []Analysis
}
