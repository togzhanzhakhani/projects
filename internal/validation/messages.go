package validation

var validateMessages = map[string]string{
	"Name.required":      "Name is required",
	"Email.required":     "Email is required",
	"Email.email":        "Email must be a valid email address",
	"Email.unique":       "Email already exists",
	"Role.required":      "Role is required",
	"Role.oneof":         "Role must be one of: admin, manager, developer",
	
    "Description.required":"Description is required",
    "Description.max":     "Description must be at most 100 characters long",
	"StartDate.required":  "Start date is required",
	"EndDate.required":    "End date is required",
	"EndDate.gtfield":     "End date must be after start date",
	"ManagerID.required":  "Manager ID is required",
    "ManagerID.gt":            "Manager ID must be greater than 0",

    "Title.required":          "The title field is required.",
	"Priority.oneof":          "Invalid priority. Must be one of: low, medium, high.",
	"Status.oneof":            "Invalid status. Must be one of: todo, in_progress, done.",
	"AssigneeID.required":    "The assignee ID field is required.",
	"AssigneeID.gt":          "The assignee ID must be greater than 0.",
	"ProjectID.required":     "The project ID field is required.",
	"ProjectID.gt":           "The project ID must be greater than 0.",
    "CreatedAt.required":      "Start date is required",
	"CompletedAt.required":    "End date is required",
	"CompletedAt.gtfield":     "End date must be after start date",
}

func GetMessage(key string) string {
    return validateMessages[key]
}