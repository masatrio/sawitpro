/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(60) NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    phone VARCHAR(13) UNIQUE,
    login_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE user_activity_logs (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    activity_type VARCHAR(13) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
