package clause

import (
	"fmt"

	"gorm.io/gorm/clause"
)

func rewriteLimitClause(c clause.Clause, builder clause.Builder) {
	var limitClause clause.Limit

	if expression, ok := c.Expression.(clause.Limit); ok {
		limitClause = expression
	} else {
		return
	}

	if offset := limitClause.Offset; offset > 0 {
		_, _ = builder.WriteString(fmt.Sprintf(" OFFSET %d ROWS", offset))
	}

	if limit := limitClause.Limit; limit > 0 {
		_, _ = builder.WriteString(fmt.Sprintf(" FETCH NEXT %d ROWS ONLY", limit))
	}
}
