package migrate

import (
	"fmt"

	"github.com/go-pg/migrations"
)

const users = `
CREATE TABLE users (
user_id serial NOT NULL,
created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
username string NOT NULL UNIQUE,
passcode string NOT NULL,
user_role ENUM (admin, editor, labeller),
PRIMARY KEY (image_id)
)`

func init() {
	up := []string{
		users,
	}

	down := []string{
		`DROP TABLE users`,
	}

	migrations.Register(func(db migrations.DB) error {
		fmt.Println("create user table")
		for _, q := range up {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	}, func(db migrations.DB) error {
		fmt.Println("dropping user tables")
		for _, q := range down {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
