package settings

import "os"

var _settings map[string]string

func setSetting(key string, defaultValue string) {
	if value, present := os.LookupEnv(key); present {
		_settings[key] = value
	} else {
		_settings[key] = defaultValue
	}
}

func init() {
	_settings = make(map[string]string)
	setSetting("PORT", "8000")
	setSetting("SURVEY_RUNNER_URL", "http://localhost:5000")
	setSetting("SURVEY_RUNNER_SCHEMA_URL", Get("SURVEY_RUNNER_URL"))
	setSetting("SCHEMA_VALIDATOR_URL", "")
	setSetting("SURVEY_REGISTER_URL", "")
	setSetting("JWT_ENCRYPTION_KEY_PATH", "jwt-test-keys/sdc-user-authentication-encryption-sr-public-key.pem")
	setSetting("JWT_SIGNING_KEY_PATH", "jwt-test-keys/sdc-user-authentication-signing-launcher-private-key.pem")
	setSetting("CASE_API_URL", "http://localhost:8161")
}

// Get returns the value for the specified named setting
func Get(name string) string {
	return _settings[name]
}
