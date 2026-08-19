package config

import "testing"

func TestEnvironmentGetBaseURL(t *testing.T) {
	// api-dev.eka.care does not resolve; the development host is api.dev.eka.care.
	for _, tc := range []struct {
		env  Environment
		want string
	}{
		{EnvironmentProduction, "https://api.eka.care"},
		{EnvironmentDevelopment, "https://api.dev.eka.care"},
		{Environment("bogus"), "https://api.eka.care"},
	} {
		if got := tc.env.GetBaseURL(); got != tc.want {
			t.Errorf("Environment(%q).GetBaseURL() = %q, want %q", tc.env, got, tc.want)
		}
	}
}
