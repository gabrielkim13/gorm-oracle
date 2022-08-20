package callbacks

import (
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
)

func create(db *gorm.DB) {
	if db.Error != nil {
		return
	}

	outVars := make(OutVars)

	if db.Statement.Schema != nil {
		if !db.Statement.Unscoped {
			for _, c := range db.Statement.Schema.CreateClauses {
				db.Statement.AddClause(c)
			}
		}

		if len(db.Statement.Schema.FieldsWithDefaultDBValue) > 0 {
			if _, ok := db.Statement.Clauses["RETURNING"]; !ok {
				fromColumns := make([]clause.Column, 0, len(db.Statement.Schema.FieldsWithDefaultDBValue))
				for _, field := range db.Statement.Schema.FieldsWithDefaultDBValue {
					fromColumns = append(fromColumns, clause.Column{Name: field.DBName})
					outVars.AddField(field)
				}
				db.Statement.AddClause(clause.Returning{Columns: fromColumns})
			}
		}
	}

	if db.Statement.SQL.Len() == 0 {
		db.Statement.SQL.Grow(180)
		db.Statement.AddClauseIfNotExists(clause.Insert{})
		db.Statement.AddClause(callbacks.ConvertToCreateValues(db.Statement))

		db.Statement.Build(db.Statement.BuildClauses...)
		outVars.AppendVars(db.Statement)
	}

	isDryRun := !db.DryRun && db.Error == nil
	if !isDryRun {
		return
	}

	result, err := db.Statement.ConnPool.ExecContext(
		db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...,
	)

	if err != nil {
		_ = db.AddError(err)
		return
	}

	db.RowsAffected, _ = result.RowsAffected()
	if db.RowsAffected != 0 && db.Statement.Schema != nil &&
		db.Statement.Schema.PrioritizedPrimaryField != nil &&
		db.Statement.Schema.PrioritizedPrimaryField.HasDefaultValue {
		insertID := outVars.ValueOf(db.Statement.Schema.PrioritizedPrimaryField.Name)

		switch db.Statement.ReflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
				rv := db.Statement.ReflectValue.Index(i)
				if reflect.Indirect(rv).Kind() != reflect.Struct {
					break
				}

				if _, isZero := db.Statement.Schema.PrioritizedPrimaryField.ValueOf(db.Statement.Context, rv); isZero {
					_ = db.AddError(db.Statement.Schema.PrioritizedPrimaryField.Set(db.Statement.Context, rv, insertID))
				}
			}
		case reflect.Struct:
			_, isZero := db.Statement.Schema.PrioritizedPrimaryField.ValueOf(db.Statement.Context, db.Statement.ReflectValue)
			if isZero {
				_ = db.AddError(db.Statement.Schema.PrioritizedPrimaryField.Set(db.Statement.Context, db.Statement.ReflectValue, insertID))
			}
		}
	}

	return
}
