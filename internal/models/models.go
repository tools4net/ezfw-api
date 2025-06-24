package models

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
