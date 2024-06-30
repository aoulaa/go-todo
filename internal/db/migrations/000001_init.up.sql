CREATE TABLE users (
                       id UUID PRIMARY KEY,
                       username VARCHAR(255) NOT NULL UNIQUE,
                       first_name VARCHAR(255),
                       last_name VARCHAR(255),
                       email VARCHAR(255) NOT NULL UNIQUE,
                       password VARCHAR(255) NOT NULL,
                       created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP(3)
);

CREATE TABLE todo (
                      id UUID PRIMARY KEY,
                      user_id UUID NOT NULL,
                      title VARCHAR(255),
                      description TEXT,
                      completed BOOLEAN NOT NULL DEFAULT FALSE,
                      created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
                      updated_at TIMESTAMP(3) NOT NULL,
                      FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
