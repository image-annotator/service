package migrate

import (
	"fmt"

	"github.com/go-pg/migrations"
)

const images = `
CREATE TABLE images (
image_id serial NOT NULL,
created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
image_path string NOT NULL UNIQUE,
PRIMARY KEY (image_id)
)`

func init() {
	up := []string{
		images,
	}

	down := []string{
		`DROP TABLE images`,
	}

	migrations.Register(func(db migrations.DB) error {
		fmt.Println("create image table")
		for _, q := range up {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	}, func(db migrations.DB) error {
		fmt.Println("dropping image tables")
		for _, q := range down {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
