package migrate

import (
	"rcsp/internal/database"
	"rcsp/internal/model"

	"github.com/go-pg/pg/orm"
)

func CreateSchema(db *database.Database) error {
	models := []interface{}{
		(*model.Order)(nil),
	}

	for _, model := range models {
		op := orm.CreateTableOptions{}
		err := db.DB.Model(model).CreateTable(&op)
		if err != nil {
			return err
		}
	}
	return nil
}
