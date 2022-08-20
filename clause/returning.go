package clause

import (
	"fmt"

	"gorm.io/gorm/clause"
)

func rewriteReturningClause(c clause.Clause, builder clause.Builder) {
	var returningClause clause.Returning

	if expression, ok := c.Expression.(clause.Returning); ok {
		returningClause = expression
	} else {
		return
	}

	columns := returningClause.Columns

	if len(columns) == 0 || (len(columns) == 1 && columns[0].Name == "*") {
		return
	}

	_, _ = builder.WriteString("RETURNING ")

	for idx, column := range columns {
		if idx > 0 {
			_ = builder.WriteByte(',')
		}

		builder.WriteQuoted(column)
	}

	_, _ = builder.WriteString(" INTO ")

	for idx, column := range columns {
		if idx > 0 {
			_ = builder.WriteByte(',')
		}

		_, _ = builder.WriteString(fmt.Sprintf(":%s", column.Name))
	}
}
