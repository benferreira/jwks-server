package test_helper

import "os"

func UnsetTestEnvironment() {
	os.Unsetenv("DEBUG")
	os.Unsetenv("PORT")
	os.Unsetenv("PRETTY_LOGGING")
	os.Unsetenv("TEST_MODE")
	os.Unsetenv("TLS")
	os.Unsetenv("TLS_PRIVATE_KEY_PATH")
	os.Unsetenv("TLS_CERT_PATH")
	os.Unsetenv("RSA_PUB_KEY")
	os.Unsetenv("RSA_KEYS_FILE")
}
