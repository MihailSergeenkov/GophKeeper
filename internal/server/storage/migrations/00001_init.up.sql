BEGIN TRANSACTION;

CREATE TABLE users(
	id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	login VARCHAR(200) NOT NULL,
	password BYTEA NOT NULL
);
CREATE UNIQUE INDEX users_login_index ON users(login);

CREATE TYPE user_data_type AS ENUM ('password', 'card', 'text', 'file');

CREATE TABLE user_data(
	id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	user_id INT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
  data BYTEA NOT NULL,
	type user_data_type NOT NULL,
	mark VARCHAR(100) NOT NULL,
	description VARCHAR(3000) NOT NULL
);
CREATE INDEX user_data_user_id_index ON user_data(user_id);

COMMIT;
