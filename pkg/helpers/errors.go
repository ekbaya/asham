package helpers

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Common error types that can be used throughout the application
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicateKey   = errors.New("duplicate key violation")
	ErrForeignKey     = errors.New("foreign key violation")
	ErrConstraint     = errors.New("constraint violation")
	ErrInvalidData    = errors.New("invalid data")
	ErrDatabase       = errors.New("database error")
)

// IsNotFoundError checks if the error is a record not found error
func IsNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrRecordNotFound)
}

// IsDuplicateKeyError checks if the error is a duplicate key violation
func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrDuplicateKey) || strings.Contains(strings.ToLower(err.Error()), "duplicate") ||
		strings.Contains(strings.ToLower(err.Error()), "unique")
}

// IsForeignKeyError checks if the error is a foreign key violation
func IsForeignKeyError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrForeignKey) || strings.Contains(strings.ToLower(err.Error()), "foreign key")
}

// IsConstraintError checks if the error is a constraint violation
func IsConstraintError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrConstraint) ||
		strings.Contains(strings.ToLower(err.Error()), "constraint") ||
		IsForeignKeyError(err) ||
		IsDuplicateKeyError(err)
}

// HandleError wraps database errors with more context
func HandleError(err error, operation string, model interface{}) error {
	if err == nil {
		return nil
	}

	modelName := fmt.Sprintf("%T", model)

	// Handle specific error types
	switch {
	case IsNotFoundError(err):
		return fmt.Errorf("%w: %s not found during %s", ErrRecordNotFound, modelName, operation)
	case IsDuplicateKeyError(err):
		return fmt.Errorf("%w: duplicate entry for %s during %s", ErrDuplicateKey, modelName, operation)
	case IsForeignKeyError(err):
		return fmt.Errorf("%w: foreign key constraint failed for %s during %s", ErrForeignKey, modelName, operation)
	case IsConstraintError(err):
		return fmt.Errorf("%w: constraint violation for %s during %s", ErrConstraint, modelName, operation)
	default:
		return fmt.Errorf("%w: failed to %s %s: %v", ErrDatabase, operation, modelName, err)
	}
}

// CreateRecord attempts to create a record and handles errors
func CreateRecord(db *gorm.DB, model interface{}) error {
	result := db.Create(model)
	if result.Error != nil {
		return HandleError(result.Error, "create", model)
	}
	return nil
}

// FindRecord attempts to find a record by its ID and handles errors
func FindRecord(db *gorm.DB, model interface{}, id interface{}) error {
	result := db.First(model, id)
	if result.Error != nil {
		return HandleError(result.Error, "find", model)
	}
	return nil
}

// UpdateRecord attempts to update a record and handles errors
func UpdateRecord(db *gorm.DB, model interface{}) error {
	result := db.Save(model)
	if result.Error != nil {
		return HandleError(result.Error, "update", model)
	}
	return nil
}

// DeleteRecord attempts to delete a record and handles errors
func DeleteRecord(db *gorm.DB, model interface{}) error {
	result := db.Delete(model)
	if result.Error != nil {
		return HandleError(result.Error, "delete", model)
	}
	return nil
}

// Transaction executes the given function within a transaction
// It automatically handles commit/rollback based on the function's result
func Transaction(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := fn(tx); err != nil {
			return err
		}
		return nil
	})
}

// QueryBuilder represents a common query building pattern
type QueryBuilder struct {
	DB       *gorm.DB
	Model    interface{}
	PageSize int
	Page     int
}

// NewQueryBuilder creates a new query builder with default pagination
func NewQueryBuilder(db *gorm.DB, model interface{}) *QueryBuilder {
	return &QueryBuilder{
		DB:       db.Model(model),
		Model:    model,
		PageSize: 10,
		Page:     1,
	}
}

// SetPagination sets pagination parameters
func (qb *QueryBuilder) SetPagination(page, pageSize int) *QueryBuilder {
	if page > 0 {
		qb.Page = page
	}
	if pageSize > 0 {
		qb.PageSize = pageSize
	}
	return qb
}

// WithPreload adds preloading to the query
func (qb *QueryBuilder) WithPreload(preloads ...string) *QueryBuilder {
	for _, preload := range preloads {
		qb.DB = qb.DB.Preload(preload)
	}
	return qb
}

// WithCondition adds a where condition to the query
func (qb *QueryBuilder) WithCondition(condition string, args ...interface{}) *QueryBuilder {
	qb.DB = qb.DB.Where(condition, args...)
	return qb
}

// Find executes the query and populates the given destination
func (qb *QueryBuilder) Find(dest interface{}) error {
	offset := (qb.Page - 1) * qb.PageSize

	result := qb.DB.Offset(offset).Limit(qb.PageSize).Find(dest)
	if result.Error != nil {
		return HandleError(result.Error, "find", qb.Model)
	}
	return nil
}

// Count returns the total count of records matching the query
func (qb *QueryBuilder) Count() (int64, error) {
	var count int64
	result := qb.DB.Count(&count)
	if result.Error != nil {
		return 0, HandleError(result.Error, "count", qb.Model)
	}
	return count, nil
}
