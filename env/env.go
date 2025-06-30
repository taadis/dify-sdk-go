package env

import "os"

func GetDifyApiHost() string {
	return os.Getenv("DIFY_API_HOST")
}

func GetDifyApiKey() string {
	return os.Getenv("DIFY_API_KEY")
}
