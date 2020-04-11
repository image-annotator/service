package migrate

import (
	"fmt"

	"github.com/go-pg/migrations"
)

const accountTable = `
CREATE TABLE accounts (
id serial NOT NULL,
created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
updated_at timestamp with time zone DEFAULT current_timestamp,
last_login timestamp with time zone NOT NULL DEFAULT current_timestamp,
email text NOT NULL UNIQUE,
name text NOT NULL,
active boolean NOT NULL DEFAULT TRUE,
roles text[] NOT NULL DEFAULT '{"user"}',
PRIMARY KEY (id)
)`

const tokenTable = `
CREATE TABLE tokens (
id serial NOT NULL,
created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
account_id int NOT NULL REFERENCES accounts(id),
token text NOT NULL UNIQUE,
expiry timestamp with time zone NOT NULL,
mobile boolean NOT NULL DEFAULT FALSE,
identifier text,
PRIMARY KEY (id)
)`

const imageTable = `
CREATE TABLE images (
image_id serial NOT NULL,
created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
labeled boolean NOT NULL DEFAULT FALSE,
image_path VARCHAR (100) NOT NULL UNIQUE,
dataset VARCHAR (100) NOT NULL,
file_name VARCHAR (100) NOT NULL,
PRIMARY KEY (image_id)
)`

const userTable = `
CREATE TABLE users (
user_id serial NOT NULL,
created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
username VARCHAR (50) NOT NULL UNIQUE,
cookie VARCHAR (50) NOT NULL UNIQUE,
passcode VARCHAR (50) NOT NULL,
user_role VARCHAR (50) NOT NULL,
PRIMARY KEY (user_id)
)`

const accessControlTable = `
CREATE TABLE access_controls(
    image_id integer NOT NULL,
    user_id integer NOT NULL,
    timeout timestamp with time zone,
    CONSTRAINT access_controls_pkey PRIMARY KEY (image_id),
    CONSTRAINT fkimage_id FOREIGN KEY (image_id)
        REFERENCES images (image_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT fkuser_id FOREIGN KEY (user_id)
        REFERENCES users (user_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
)`

const contentTable = `
CREATE TABLE contents
(
    label_contents_id serial NOT NULL,
    content_name VARCHAR (50) NOT NULL,
	created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
	updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
    PRIMARY KEY (label_contents_id),
    UNIQUE (content_name)
)`

const labelTable = `
CREATE TABLE labels
(
    label_id serial NOT NULL,
    image_id integer NOT NULL,
    label_x_center double precision,
    label_y_center double precision,
    label_width double precision,
    label_height double precision,
    label_content_id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    PRIMARY KEY (label_id),
    CONSTRAINT fkimage_id FOREIGN KEY (image_id)
        REFERENCES images (image_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT fklabel_content_id FOREIGN KEY (label_content_id)
        REFERENCES contents (label_contents_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
)`

const insertAdmin = `
INSERT INTO users(
	username, cookie, passcode, user_role)
	VALUES ('adminone', 'admincookie', 'password', 'admin') ON CONFLICT DO NOTHING;
`

func init() {
	up := []string{
		accountTable,
		tokenTable,
		userTable,
		imageTable,
		accessControlTable,
		contentTable,
		labelTable,
		insertAdmin,
	}

	down := []string{
		`DROP TABLE IF EXISTS tokens`,
		`DROP TABLE IF EXISTS accounts`,
		`DROP TABLE IF EXISTS access_controls`,
		`DROP TABLE IF EXISTS labels`,
		`DROP TABLE IF EXISTS images`,
		`DROP TABLE IF EXISTS users`,
		`DROP TABLE IF EXISTS contents`,
	}

	migrations.Register(func(db migrations.DB) error {
		fmt.Println("creating initial tables")
		for _, q := range up {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	}, func(db migrations.DB) error {
		fmt.Println("dropping initial tables")
		for _, q := range down {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
