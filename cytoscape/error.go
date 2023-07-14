package cytoscape

import "github.com/layer5io/meshkit/errors"

const (
	ErrPatternFromCytoscapeCode = "11093"
)
func ErrPatternFromCytoscape(err error) error {
	return errors.New(ErrPatternFromCytoscapeCode, errors.Alert, []string{"Could not create pattern file from given cytoscape"}, []string{err.Error()}, []string{"Invalid cytoscape body", "Service name is empty for one or more services", "_data does not have correct data"}, []string{"Make sure cytoscape is valid", "Check if valid service name was passed in the request", "Make sure _data field has \"settings\" field"})
}
