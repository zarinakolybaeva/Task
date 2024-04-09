"Do Task" 
The project aims to provide a platform for managing tasks efficiently. 
The core of the project is an API server that exposes endpoints for managing tasks.
It handles HTTP requests from clients and performs corresponding actions on tasks stored in the database.
The project integrates with a PostgreSQL database to store task-related data.
tasks (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    description text NOT NULL,
    due_date timestamp(0) with time zone NOT NULL,
    priority text NOT NULL,
    status text NOT NULL,
    category text NOT NULL,
    user_id bigserial,
    version integer NOT NULL DEFAULT 1
--     FOREIGN KEY (user_id) REFERENCES users(id)
);
