package validation

import (
	"github.com/gin-gonic/gin"
	"net/http"
    "github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
    validate = validator.New()
}

func GetValidator() *validator.Validate {
    return validate
}

func ValidateStruct(c *gin.Context, obj interface{}) bool {
	validate := GetValidator()
	if err := validate.Struct(obj); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			translation := GetMessage(field + "." + err.Tag())
			validationErrors = append(validationErrors, translation)
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return false
	}
	return true
}


