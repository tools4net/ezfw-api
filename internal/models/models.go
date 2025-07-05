package models

import "fmt"

// This file can be used for models that are shared across different proxy types
// or for general API request/response models not specific to a configuration.

// For now, it's kept minimal to avoid conflicts with specific proxy models
// like xray.go and singbox.go.

// Example of a shared model if needed in the future:
// type APIErrorResponse struct {
// Code    string `json:"code"`
// Message string `json:"message"`
// }

// SingBoxConfig related structs are in singbox.go
// XrayConfig related structs are in xray.go

// ErrorResponse represents a generic error response for API calls.
type ErrorResponse struct {
	Error string `json:"error" example:"Detailed error message"`
}

// Common errors for V2 models
var (
	ErrInvalidServiceType = fmt.Errorf("invalid service type")
	ErrNodeNotFound       = fmt.Errorf("node not found")
	ErrServiceNotFound    = fmt.Errorf("service instance not found")
)
