package clause

import "gorm.io/gorm"

func RewriteClauseBuilders(db *gorm.DB) {
	db.ClauseBuilders["LIMIT"] = rewriteLimitClause
	db.ClauseBuilders["FOR"] = rewriteLockingClause
	db.ClauseBuilders["RETURNING"] = rewriteReturningClause
}
