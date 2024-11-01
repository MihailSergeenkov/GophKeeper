BEGIN TRANSACTION;

DROP INDEX user_data_user_id_index;
DROP TABLE user_data;
DROP TYPE user_data_type;

DROP INDEX users_login_index;
DROP TABLE users;

COMMIT;
