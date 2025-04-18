package helpers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Filter represents a database filter
type Filter struct {
	Field       string
	Operator    string
	Value       interface{}
	ValueType   string
	IsDateRange bool
}

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page   int
	Limit  int
	Offset int
}


// GetPaginationParams extracts pagination parameters from context
func GetPaginationParams(c *gin.Context) PaginationParams {
	page, _ := c.GetQuery("page")
	limit, _ := c.GetQuery("limit")
	
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}
	
	offset := (pageInt - 1) * limitInt
	
	return PaginationParams{
		Page:   pageInt,
		Limit:  limitInt,
		Offset: offset,
	}
}


func ApplyFilters(query *gorm.DB, filters []Filter) *gorm.DB {
	for _, filter := range filters {
		if filter.IsDateRange {
			dateRange := strings.Split(filter.Value.(string), ",")
			if len(dateRange) == 2 {
				query = query.Where("DATE("+filter.Field+") >= ? AND DATE("+filter.Field+") <= ?", dateRange[0], dateRange[1])
			} else {
				query = query.Where("DATE("+filter.Field+") >= ?", dateRange[0])
			}
		} else if filter.Operator == "in" {
			query = query.Where(filter.Field+" IN (?)", filter.Value)
		} else if filter.Operator == "like" {
			query = query.Where(filter.Field+" LIKE ?", "%"+filter.Value.(string)+"%")
		} else {
			query = query.Where(filter.Field+" "+filter.Operator+" ?", filter.Value)
		}
	}
	return query
}


// BuildMultiSearchFilter builds a filter for searching across multiple fields
func BuildMultiSearchFilter(query *gorm.DB, searchValue string, fields []string) *gorm.DB {
	var conditions []string
	var values []interface{}

	// Normalize the search value for phone numbers
	normalizedSearchValue := strings.ReplaceAll(searchValue, " ", "")
	if strings.HasPrefix(normalizedSearchValue, "0") {
		normalizedSearchValue = "254" + normalizedSearchValue[1:]
	}

	for _, field := range fields {
		conditions = append(conditions, "LOWER("+field+") = ?")
		values = append(values, strings.ToLower(searchValue))

		// Add condition for normalized phone number
		conditions = append(conditions, "LOWER("+field+") = ?")
		values = append(values, strings.ToLower(normalizedSearchValue))
	}

	whereClause := strings.Join(conditions, " OR ")
	args := make([]interface{}, len(values))
	for i, v := range values {
		args[i] = v
	}

	return query.Where(whereClause, args...)
}


// ExecutePaginatedQuery executes a query with pagination and returns total count
func ExecutePaginatedQuery(query *gorm.DB, pagination PaginationParams, result interface{}) (int64, error) {
	var total int64
	countQuery := query
	err := countQuery.Count(&total).Error
	if err != nil {
		return 0, err
	}
	
	err = query.Offset(pagination.Offset).Limit(pagination.Limit).Order("created_at DESC").Find(result).Error
	if err != nil {
		return 0, err
	}
	
	return total, nil
}


// ValidateAndBindJSON validates and binds JSON request to a struct
/*func ValidateAndBindJSON(c *gin.Context, payload interface{}) bool {
	if err := c.ShouldBindJSON(payload); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			utilities.ShowError(c, http.StatusBadRequest, utilities.FormatValidationError(validationErrors))
		} else {
			utilities.ShowError(c, http.StatusBadRequest, err.Error())
		}
		return false
	}
	return true
}*/

/*func ValidateAndBindQuery(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		//ShowErrorResponse(c, http.StatusBadRequest, FormatValidationError(err))

		utilities.ShowError(c, http.StatusBadRequest, FormatValidationError(err))
		return false
	}
	return true
}*/

	


// FormatValidationError formats validation errors (placeholder - implement as needed)
func FormatValidationError(err error) string {
	return err.Error()
}


// ShowPaginatedResponse shows a paginated response
func ShowPaginatedResponse(c *gin.Context, data interface{}, total int64, page int, limit int) {
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data": gin.H{
			"records":      data,
			"total":        total,
			"current_page": page,
			"per_page":     limit,
			"last_page":    (int(total) + limit - 1) / limit,
		},
	})
}

// ShowErrorResponse shows an error response
func ShowErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"status":  false,
		"message": message,
	})
}