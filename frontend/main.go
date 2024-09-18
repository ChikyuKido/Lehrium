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

func getCachedContent(path string, filepath string) []byte {
	if debug {
		content, err := os.ReadFile(filepath)
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
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil
	}
	var compressedContent bytes.Buffer
	writer, _ := gzip.NewWriterLevel(&compressedContent, gzip.BestCompression)
	_, err = writer.Write(content)
	if err != nil {
		return nil
	}
	writer.Close()

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
	debugVal := getEnvOrDefault("LEHIUM_FRONTEND_DEBUG", "true")
	debug = debugVal == "true"
	portVal := getEnvOrDefault("GIN_PORT", "8081")
	_, portErr := strconv.Atoi(portVal)
	if portErr != nil {
		portVal = "8081"
	}
	fmt.Println("Start with debug: " + debugVal)
	r := gin.Default()
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
	r.Run(fmt.Sprintf(":%s", portVal))
}


