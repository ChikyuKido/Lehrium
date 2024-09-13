package main

import(
    "github.com/gin-gonic/gin"
    "os"
    "fmt"
    "strconv"
) 

func getEnvOrDefault(key string, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func main() {
	portVal := getEnvOrDefault("GIN_PORT", "8080")
	_, portErr := strconv.Atoi(portVal)
	if portErr != nil {
		portVal = "8081"
	}
	r := gin.Default()
	fmt.Println("Starting on port " + portVal)
	r.Run(fmt.Sprintf(":%s", portVal))
}
