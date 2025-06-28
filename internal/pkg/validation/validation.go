package validation

import (
	"time"

	"github.com/alxand/nalo-workspace/internal/domain/models"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Init initializes the validator with custom validation rules
func Init() {
	// Register custom validation functions
	validate.RegisterValidation("date_format", validateDateFormat)
	validate.RegisterValidation("time_range", validateTimeRange)
	validate.RegisterValidation("status", validateStatus)
	validate.RegisterValidation("company_size", validateCompanySize)
}

// ValidateDailyTask validates a DailyTask model
func ValidateDailyTask(task *models.DailyTask) error {
	// Create a struct for validation that excludes the User relationship
	type DailyTaskValidation struct {
		Day               string    `json:"day" validate:"required"`
		Date              time.Time `json:"date" validate:"required,date_format"`
		StartTime         time.Time `json:"start_time" validate:"required,date_format"`
		EndTime           time.Time `json:"end_time" validate:"required,date_format"`
		Status            string    `json:"status" validate:"required,status"`
		Score             int       `json:"score" validate:"min=0,max=10"`
		ProductivityScore int       `json:"productivity_score" validate:"min=0,max=100"`
	}

	validationStruct := DailyTaskValidation{
		Day:               task.Day,
		Date:              task.Date,
		StartTime:         task.StartTime,
		EndTime:           task.EndTime,
		Status:            task.Status,
		Score:             task.Score,
		ProductivityScore: task.ProductivityScore,
	}

	return validate.Struct(validationStruct)
}

// ValidateContinent validates a Continent model
func ValidateContinent(continent *models.Continent) error {
	return validate.Struct(continent)
}

// ValidateCountry validates a Country model
func ValidateCountry(country *models.Country) error {
	return validate.Struct(country)
}

// ValidateCompany validates a Company model
func ValidateCompany(company *models.Company) error {
	return validate.Struct(company)
}

// validateDateFormat validates that the date is in the correct format
func validateDateFormat(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	// Check if date is not zero
	return !date.IsZero()
}

// validateTimeRange validates that start time is before end time
func validateTimeRange(fl validator.FieldLevel) bool {
	// This would need to be implemented based on your specific requirements
	// For now, return true as a placeholder
	return true
}

// validateStatus validates that the status is one of the allowed values
func validateStatus(fl validator.FieldLevel) bool {
	status, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	allowedStatuses := []string{"pending", "in_progress", "completed", "cancelled"}
	for _, allowed := range allowedStatuses {
		if status == allowed {
			return true
		}
	}
	return false
}

// validateCompanySize validates that the company size is one of the allowed values
func validateCompanySize(fl validator.FieldLevel) bool {
	size, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	allowedSizes := []string{"small", "medium", "large", "enterprise"}
	for _, allowed := range allowedSizes {
		if size == allowed {
			return true
		}
	}
	return false
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	return validate
}
