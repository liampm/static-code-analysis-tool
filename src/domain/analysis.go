package domain

import uuid "github.com/satori/go.uuid"

type Analysis struct {
	Id       uuid.UUID   `json:"id"`
	JobId    uuid.UUID   `json:"jobId"`
	TaskId   uuid.UUID   `json:"task"`
	TargetId uuid.UUID   `json:"targetId"`
	Report   interface{} `json:"report"`
}

type AnalysisReference struct {
	Id       uuid.UUID   `json:"id"`
	JobId    uuid.UUID   `json:"jobId"`
	TaskId   uuid.UUID   `json:"task"`
	TargetId uuid.UUID   `json:"targetId"`
	Report   interface{} `json:"report"`
}

func fromTask(job *Job, target *Target, task task, taskId uuid.UUID) Analysis {
	return Analysis{
		Id:       uuid.NewV4(),
		TaskId:   taskId,
		TargetId: target.Id,
		JobId:    job.Id,
		Report:   task.analyse(target),
	}
}
