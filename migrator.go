package oracle

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/migrator"
)

type Migrator struct {
	migrator.Migrator
}

func (m Migrator) CurrentDatabase() (name string) {
	_ = m.DB.Raw("SELECT ORA_DATABASE_NAME FROM DUAL").Row().Scan(&name)

	return
}

func (m Migrator) DropTable(values ...interface{}) error {
	values = m.ReorderModels(values, false)

	return m.DB.Connection(func(tx *gorm.DB) error {
		for i := len(values) - 1; i >= 0; i-- {
			value := values[i]

			if !m.HasTable(value) {
				continue
			}

			if err := m.RunWithValue(values[i], func(stmt *gorm.Statement) error {
				return tx.Exec("DROP TABLE ? CASCADE CONSTRAINTS", clause.Table{Name: stmt.Table}).Error
			}); err != nil {
				return err
			}
		}

		return nil
	})
}

func (m Migrator) HasTable(value interface{}) bool {
	var count int

	_ = m.RunWithValue(value, func(stmt *gorm.Statement) error {
		return m.DB.Raw("SELECT COUNT(*) FROM USER_TABLES WHERE TABLE_NAME = ?", stmt.Table).Row().Scan(&count)
	})

	return count > 0
}

func (m Migrator) RenameTable(oldName, newName interface{}) error {
	var oldTable, newTable interface{}

	if v, ok := oldName.(string); ok {
		oldTable = clause.Table{Name: v}
	} else {
		stmt := &gorm.Statement{DB: m.DB}
		if err := stmt.Parse(oldName); err == nil {
			oldTable = m.CurrentTable(stmt)
		} else {
			return err
		}
	}

	if v, ok := newName.(string); ok {
		newTable = clause.Table{Name: v}
	} else {
		stmt := &gorm.Statement{DB: m.DB}
		if err := stmt.Parse(newName); err == nil {
			newTable = m.CurrentTable(stmt)
		} else {
			return err
		}
	}

	return m.DB.Exec("RENAME ? TO ?", oldTable, newTable).Error
}

func (m Migrator) GetTables() (tableList []string, err error) {
	err = m.DB.Raw("SELECT TABLE_NAME FROM USER_TABLES").Scan(&tableList).Error

	return
}

func (m Migrator) AlterColumn(value interface{}, field string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if field := stmt.Schema.LookUpField(field); field != nil {
			return m.DB.Exec(
				"ALTER TABLE ? MODIFY ? ?",
				m.CurrentTable(stmt), clause.Column{Name: field.DBName}, m.FullDataTypeOf(field),
			).Error
		}

		return fmt.Errorf("failed to look up field with name: %s", field)
	})
}

func (m Migrator) HasColumn(value interface{}, field string) bool {
	var count int

	_ = m.RunWithValue(value, func(stmt *gorm.Statement) error {
		return m.DB.Raw(
			"SELECT COUNT(*) FROM USER_TAB_COLS WHERE TABLE_NAME = ? AND COLUMN_NAME",
			stmt.Table, field,
		).Row().Scan(&count)
	})

	return count > 0
}

func (m Migrator) HasConstraint(value interface{}, name string) bool {
	var count int

	_ = m.RunWithValue(value, func(stmt *gorm.Statement) error {
		return m.DB.Raw(
			"SELECT COUNT(*) FROM USER_CONSTRAINTS WHERE TABLE_NAME = ? AND CONSTRAINT_NAME",
			stmt.Table, name,
		).Row().Scan(&count)
	})

	return count > 0
}

func (m Migrator) DropIndex(value interface{}, name string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if idx := stmt.Schema.LookIndex(name); idx != nil {
			name = idx.Name
		}

		return m.DB.Exec("DROP INDEX ?", clause.Column{Name: name}).Error
	})
}

func (m Migrator) HasIndex(value interface{}, name string) bool {
	var count int

	_ = m.RunWithValue(value, func(stmt *gorm.Statement) error {
		return m.DB.Raw(
			"SELECT COUNT(*) FROM USER_INDEXES WHERE TABLE_NAME = ? AND INDEX_NAME = ?",
			stmt.Table, name,
		).Row().Scan(&count)
	})

	return count > 0
}

func (m Migrator) RenameIndex(value interface{}, oldName, newName string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		return m.DB.Exec("ALTER INDEX ? RENAME TO ?", clause.Column{Name: oldName}, clause.Column{Name: newName}).Error
	})
}
