package env

import "os"

func GetDifyBaseUrl() string {
	return os.Getenv("DIFY_BASE_URL")
}

func GetDifyApiKey() string {
	return os.Getenv("DIFY_API_KEY")
}
