package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var (
	cache      = make(map[string][]byte)
	cacheMutex = &sync.Mutex{}
	debug      = true
)

func getTemplateFiles(directory string) []string {
	var files []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".html" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return files
}

func getCachedContent(path string, filepath string, data any) []byte {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	if content, found := cache[path]; found && !debug {
		return content
	}
	var content []byte
	if strings.HasSuffix(filepath, ".html") || strings.HasSuffix(filepath, ".js") {
		byteBuffer := bytes.NewBuffer(make([]byte, 0))
		templates := getTemplateFiles("./sites/partials")
		t, err := template.ParseFiles(append([]string{filepath}, templates...)...)
		if err != nil {
			fmt.Printf("Failed to parse template %v", err)
			return nil
		}
		err = t.Execute(byteBuffer, data)
		if err != nil {
			fmt.Printf("Failed to execute template %v", err)
			return nil
		}
		content = byteBuffer.Bytes()
	} else {
		content, _ = os.ReadFile(filepath)
	}
	var compressedContent bytes.Buffer
	writer, _ := gzip.NewWriterLevel(&compressedContent, gzip.BestCompression)
	_, err := writer.Write(content)
	if err != nil {
		return nil
	}
	writer.Close()

	compressedData := compressedContent.Bytes()
	cache[path] = compressedData
	return compressedData
}

func serveDirectory(rootPath string, baseDir string, r *gin.RouterGroup, data any) {
	filepath.Walk(baseDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath, _ := filepath.Rel(baseDir, path)
			urlPath := rootPath + relativePath
			servePage(urlPath, path, r, data)
		}
		return nil
	})
}
func servePage(path string, diskPath string, r *gin.RouterGroup, data any) {
	r.GET(path, func(c *gin.Context) {
		content := getCachedContent(path, diskPath, data)
		contentType := mime.TypeByExtension(filepath.Ext(diskPath))
		c.Header("Content-Encoding", "gzip")
		if !debug && filepath.Ext(diskPath) == ".css" {
			c.Header("Cache-Control", "public, max-age=3600")
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
	debugVal := getEnvOrDefault("LEHRIUM_FRONTEND_DEBUG", "true")
	debug = debugVal == "true"
	portVal := getEnvOrDefault("GIN_PORT", "8081")
	_, portErr := strconv.Atoi(portVal)
	if portErr != nil {
		portVal = "8081"
	}
	fmt.Println("Start with debug: " + debugVal)

	var settings struct {
		ApiBaseUrl string
	}
	settings.ApiBaseUrl = "http://localhost:8080"

	r := gin.Default()
	sitesGroup := r.Group("/")
	serveDirectory("/css/", "./css", sitesGroup, nil)
	serveDirectory("/js/", "./js", sitesGroup, settings)
	serveDirectory("/", "./sites", sitesGroup, nil)
	serveDirectory("/imgs/", "./imgs", sitesGroup, nil)
	servePage("/", "./sites/index.html", sitesGroup, nil)
	servePage("/auth/login", "./sites/auth/login.html", sitesGroup, nil)
	servePage("/auth/register", "./sites/auth/register.html", sitesGroup, nil)
	servePage("/auth/verifyEmail", "./sites/auth/verifyEmail.html", sitesGroup, nil)

	servePage("/teachers", "./sites/teacher/teachers.html", sitesGroup, nil)
	servePage("/teacher/:id", "./sites/auth/teacher.html", sitesGroup, nil)
	//	servePage("/admin/dashboard", "./sites/admin/dashboard.html", sitesGroup) //TODO: optional

	fmt.Println("Starting on port " + portVal)
	r.Run(fmt.Sprintf(":%s", portVal))
}
