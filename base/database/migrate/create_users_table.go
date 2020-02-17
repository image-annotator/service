package migrate

import (
	"fmt"

	"github.com/go-pg/migrations"
)

const users = `
CREATE TYPE user_role AS ENUM ('admin', 'editor', 'labeller');

CREATE TABLE users (
user_id serial NOT NULL,
created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
username text NOT NULL UNIQUE,
passcode text NOT NULL,
user_role user_role,
PRIMARY KEY (user_id)
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
