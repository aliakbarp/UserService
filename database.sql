/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */


CREATE TABLE users (
	id serial PRIMARY KEY,
  phone_number VARCHAR ( 17 ) UNIQUE NOT NULL,
	full_name VARCHAR ( 50 ) NOT NULL,
	hashed_pass VARCHAR ( 74 ) NOT NULL,
  count int NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP
);


