package fixtures

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

const (
	sqlFileExtension = ".sql"
)

func Load(db *sql.DB, path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		fileName := file.Name()
		ctx := context.Background()

		if filepath.Ext(fileName) != sqlFileExtension {
			continue
		}

		script, err := ioutil.ReadFile(filepath.Join(path, fileName))
		if err != nil {
			return err
		}

		err = func() error {
			tx, err := db.Begin() // nolint:govet
			if err != nil {
				return err
			}

			defer func() {
				_ = tx.Rollback()
			}()

			_, err = tx.ExecContext(ctx, string(script))
			if err != nil {
				return err
			}

			return tx.Commit()
		}()

		if err != nil {
			return err
		}
		fmt.Printf("creating fixture: %s\n", fileName)
	}

	return nil
}
