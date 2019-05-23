CREATE DATABASE scat;

CREATE TABLE project (
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
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
    report JSON
);