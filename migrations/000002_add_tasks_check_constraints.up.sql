ALTER TABLE tasks ADD CONSTRAINT tasks_due_date_check CHECK (due_date > created_at);
ALTER TABLE tasks ADD CONSTRAINT tasks_priority_check CHECK (priority IN ('high', 'medium', 'low'));




