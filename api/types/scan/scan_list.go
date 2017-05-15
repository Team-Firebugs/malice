package scan

import "github.com/maliceio/malice/api/types"

// ScansListOKBody scans list ok body
type ScansListOKBody struct {

	// List of scans
	// Required: true
	Scans []*types.Scan `json:"Scans"`

	// Warnings that occurred when fetching the list of scans
	// Required: true
	Warnings []string `json:"Warnings"`
}
