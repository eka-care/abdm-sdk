// Package abdm provides ABDM (Ayushman Bharat Digital Mission) API services.
//
// This package contains all ABDM-related services organized under a single client.
// The ABDM client is designed to work with the main Eka SDK client which handles
// authentication, configuration, and HTTP client management.
package abdm

import (
	"github.com/eka-care/eka-sdk-go/internal/interfaces"
	"github.com/eka-care/eka-sdk-go/internal/utils"
	"github.com/eka-care/eka-sdk-go/services/abdm/abha/login"
	"github.com/eka-care/eka-sdk-go/services/abdm/abha/profile"
	"github.com/eka-care/eka-sdk-go/services/abdm/abha/registration"
)

// Client represents the ABDM services client
// It organizes all ABDM-related services under a single interface
type Client struct {
	loginService        *login.Service
	registrationService *registration.Service
	profileService      *profile.Service
	utilsService        *utils.Service
}

// NewClient creates a new ABDM client with the given configuration
// The configuration is managed by the main SDK client
func NewClient(config interfaces.Config) *Client {
	return &Client{
		loginService:        login.NewService(config),
		registrationService: registration.NewService(config),
		profileService:      profile.NewService(config),
		utilsService:        utils.NewService(config),
	}
}

// Login returns the ABDM login service
func (c *Client) Login() *login.Service {
	return c.loginService
}

// Registration returns the ABDM registration service
func (c *Client) Registration() *registration.Service {
	return c.registrationService
}

// Profile returns the ABDM profile service
func (c *Client) Profile() *profile.Service {
	return c.profileService
}

// Utils returns the ABDM utils service
func (c *Client) Utils() *utils.Service {
	return c.utilsService
}
