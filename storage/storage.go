package storage

import (
	"database/sql"
	"fmt"

	"github.com/djumpen/go-rest-admin/apperrors"
	"github.com/jinzhu/gorm"
)

//sqlErrorHandle Wrapping db errors into apperrors
func sqlErrHandle(err error) error {
	switch err {
	case nil:
		return nil
	case sql.ErrNoRows:
		return apperrors.NewNoRows(err)
	}
	// TODO: add duplicate entry err
	// if mysqldb.CheckError(err, mysqldb.ER_DUP_ENTRY) {
	// 	return apierrors.NewDuplicateEntry(err)
	// }
	if err == gorm.ErrRecordNotFound {
		return apperrors.NewNoRows(err)
	}
	return err
}

// wrapErr Wrapping db errors into apperrors
func wrapErr(err error) error {

	switch err {
	case nil:
		return nil
	case sql.ErrNoRows:
		return apperrors.NewNoRows(err)
	}
	// TODO: add duplicate entry err
	// if mysqldb.CheckError(err, mysqldb.ER_DUP_ENTRY) {
	// 	return apierrors.NewDuplicateEntry(err)
	// }
	if err == gorm.ErrRecordNotFound {
		return apperrors.NewNoRows(err)
	}
	return err
}

// For query debugging
func scopeStoreSqlAndVars(scope *gorm.Scope) {
	scope.DB().InstantSet("sql", scope.SQL)
	scope.DB().InstantSet("sqlVars", scope.SQLVars)
}

// Middleware for Preload() func to enum fields in SELECT
func gSelect(fields ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(fields)
	}
}

// Middleware with table prefix
func gSelectWithPrefix(table string, fields ...string) func(db *gorm.DB) *gorm.DB {
	for i, v := range fields {
		fields[i] = fmt.Sprintf("%s.%s", table, v)
	}
	return gSelect(fields...)
}
