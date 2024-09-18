package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/fs"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var (
	cache      = make(map[string][]byte)
	cacheMutex = &sync.Mutex{}
	debug      = true
)

func getCachedContent(path string, diskPath string) []byte {
	if debug {
		content, err := os.ReadFile(diskPath)
		if err != nil {
			return nil
		}
		return content
	}

	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	if content, found := cache[path]; found {
		return content
	}

	content, err := os.ReadFile(diskPath)
	if err != nil {
		return nil
	}

	var compressedContent bytes.Buffer
	writer := gzip.NewWriter(&compressedContent)
	defer writer.Close()
	_, err = writer.Write(content)
	if err != nil {
		return nil
	}

	if err := writer.Close(); err != nil {
		return nil
	}

	compressedData := compressedContent.Bytes()
	cache[path] = compressedData
	return compressedData
}

func serveDirectory(rootPath string, baseDir string, r *gin.RouterGroup) {
	filepath.Walk(baseDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath, _ := filepath.Rel(baseDir, path)
			urlPath := rootPath + relativePath
			servePage(urlPath, path, r)
		}
		return nil
	})
}

func servePage(path string, diskPath string, r *gin.RouterGroup) {
	r.GET(path, func(c *gin.Context) {
		content := getCachedContent(path, diskPath)
		if content == nil {
			c.Status(404)
			return
		}
		contentType := mime.TypeByExtension(filepath.Ext(diskPath))
		if !debug {
			c.Header("Content-Encoding", "gzip")
		}
		c.Data(200, contentType, content)
	})
}

func getEnvOrDefault(key string, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func main() {
	r := gin.Default()

	debugVal := getEnvOrDefault("LEHIUM_FRONTEND_DEBUG", "true")
	debug = debugVal == "true"
	portVal := getEnvOrDefault("GIN_PORT", "8081")

	if _, portErr := strconv.Atoi(portVal); portErr != nil {
		portVal = "8081"
	}
	fmt.Println("Start with debug: " + debugVal)

	sitesGroup := r.Group("/")
	serveDirectory("/css/", "./css", sitesGroup)
	serveDirectory("/js/", "./js", sitesGroup)
	serveDirectory("/", "./sites", sitesGroup)
	serveDirectory("/imgs/", "./imgs", sitesGroup)
	servePage("/", "./sites/index.html", sitesGroup)
	servePage("/login", "./sites/login.html", sitesGroup)
	servePage("/register", "./sites/register.html", sitesGroup)
	servePage("/succesfullLogin", "./sites/succesfullLogin.html", sitesGroup)
	servePage("/succesfullRegister", "./sites/succesfullRegister.html", sitesGroup)

	fmt.Println("Starting on port " + portVal)
	if err := r.Run(fmt.Sprintf(":%s", portVal)); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
