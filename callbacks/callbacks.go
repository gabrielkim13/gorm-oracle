package callbacks

import "gorm.io/gorm"

func RewriteCallbacks(db *gorm.DB) (err error) {
	if err = db.Callback().Create().Replace("gorm:create", create); err != nil {
		return
	}

	return
}
