package main

import (
	"bytes"
	"fmt"
	"github.com/andybalholm/brotli"
	"html/template"
	"io/fs"
	"mime"
	"net/http"
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

func unescapeJavaScript(content []byte) []byte {
	strContent := string(content)
	strContent = strings.ReplaceAll(strContent, "&lt;", "<")
	strContent = strings.ReplaceAll(strContent, "&gt;", ">")
	strContent = strings.ReplaceAll(strContent, "&amp;", "&")
	return []byte(strContent)
}

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
		if strings.HasSuffix(filepath, ".js") {
			content = unescapeJavaScript(content)
		}
	} else {
		content, _ = os.ReadFile(filepath)
	}
	if !strings.Contains(filepath, "imgs") {
		var compressedContent bytes.Buffer
		//writer, _ := gzip.NewWriterLevel(&compressedContent, gzip.BestCompression)
		writer := brotli.NewWriterLevel(&compressedContent, brotli.BestCompression)
		_, err := writer.Write(content)
		if err != nil {
			return nil
		}
		writer.Close()

		compressedData := compressedContent.Bytes()
		cache[path] = compressedData
	} else {
		cache[path] = content
	}
	return cache[path]
}

func serveDirectory(rootPath string, baseDir string, data any) {
	filepath.Walk(baseDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath, _ := filepath.Rel(baseDir, path)
			urlPath := rootPath + relativePath
			servePage(urlPath, path, data)
		}
		return nil
	})
}
func servePage(path string, diskPath string, data any) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		content := getCachedContent(path, diskPath, data)
		contentType := mime.TypeByExtension(filepath.Ext(diskPath))

		// Compress everything except images
		if !strings.Contains(diskPath, "imgs") {
			w.Header().Set("Content-Encoding", "br")
		}
		//if !debug && filepath.Ext(diskPath) == ".css" {
		//	w.Header().Set("Cache-Control", "public, max-age=3600")
		//}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", contentType)
		w.Write(content)
	})
}

func getEnvOrDefault(key string, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

type SliderData struct {
	Min        int
	Max        int
	InitialMin int
	InitialMax int
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

	serveDirectory("/css/", "./css", nil)
	serveDirectory("/js/", "./js", settings)
	serveDirectory("/", "./sites", nil)
	serveDirectory("/imgs/", "./imgs", nil)
	servePage("/", "./sites/index.html", nil)
	servePage("/auth/login", "./sites/auth/login.html", nil)
	servePage("/auth/register", "./sites/auth/register.html", nil)
	servePage("/auth/verifyEmail", "./sites/auth/verifyEmail.html", nil)

	var teachersPageData struct {
		SliderData SliderData
	}
	slider := SliderData{
		Min:        1,
		Max:        10,
		InitialMin: 1,
		InitialMax: 10,
	}
	teachersPageData.SliderData = slider

	servePage("/teachers", "./sites/teacher/teachers.html", teachersPageData)
	servePage("/teacher/:id", "./sites/teacher/teacher.html", nil)
	//	servePage("/admin/dashboard", "./sites/admin/dashboard.html", sitesGroup) //TODO: optional

	fmt.Println("Starting on port " + portVal)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", portVal), nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
