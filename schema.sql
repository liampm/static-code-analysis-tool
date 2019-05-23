CREATE DATABASE scat;

CREATE TABLE project (
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TYPE tasks AS ENUM ('fileCount');

CREATE TABLE task (
    id UUID NOT NULL PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES project(id),
    task tasks NOT NULL,
    UNIQUE(project_id, task)
);

CREATE TABLE target (
    id UUID NOT NULL PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES project(id),
    name VARCHAR(255) NOT NULL,
    config JSON
);

CREATE TABLE job (
    id UUID NOT NULL PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES project(id),
    date_initiated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
);

CREATE TABLE analysis (
    id UUID NOT NULL PRIMARY KEY,
    job_id UUID NOT NULL REFERENCES job(id),
    target_id UUID NOT NULL REFERENCES target(id),
    task_id UUID NOT NULL REFERENCES task(id),
    report JSON
);
