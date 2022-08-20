package clause

import (
	"gorm.io/gorm/clause"
)

func rewriteLockingClause(c clause.Clause, builder clause.Builder) {
	if _, ok := c.Expression.(clause.Locking); ok {
		_, _ = builder.WriteString(" FOR UPDATE")
	}
}
