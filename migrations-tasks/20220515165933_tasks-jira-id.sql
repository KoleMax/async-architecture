-- +goose Up
-- SQL in this section is executed when the migration is applied.

alter table tasks add jira_id text;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

alter table tasks drop column jira_id;

