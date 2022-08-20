package callbacks

import (
	"database/sql"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type OutVars map[string]sql.Out

func (o *OutVars) AddField(field *schema.Field) {
	(*o)[field.Name] = sql.Out{Dest: reflect.New(field.FieldType).Interface()}
}

func (o *OutVars) AppendVars(stmt *gorm.Statement) {
	for _, v := range *o {
		stmt.Vars = append(stmt.Vars, v)
	}
}

func (o *OutVars) ValueOf(name string) (value interface{}) {
	if v, ok := (*o)[name]; ok {
		return v.Dest
	}

	return
}
