package schema

// import "github.com/layer5io/meshkit/errors"

// const (
// 	ErrGetSchemaCode    = "11091"
// 		ErrUpdateSchemaCode = "11092"

// )


// func ErrGetSchema(err error) error {
// 	return errors.New(ErrGetSchemaCode, errors.Alert, []string{"Could not get schema for the given CRD"}, []string{err.Error()}, []string{"Unable to marshal from cue value to JSON", "Unable to unmarshal from JSON to Go type"}, []string{"Verify CRD has valid schema.", "Malformed JSON provided", "CUE path to propery doesn't exist"})
// }

// func ErrUpdateSchema(err error, obj string) error {
// 	return errors.New(ErrUpdateSchemaCode, errors.Alert, []string{"Failed to update schema properties for ", obj}, []string{err.Error()}, []string{"Incorrect type assertion", "Selector.Unquoted might have been invoked on non-string label", "error during conversion from cue.Selector to string"}, []string{"Ensure correct type assertion", "Perform appropriate conversion from cue.Selector to string", "Verify CRD has valid schema"})
// }
