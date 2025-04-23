package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Resume struct to match the YAML structure
type Resume struct {
	Name    string `yaml:"name"`
	Contact struct {
		Website  string `yaml:"website"`
		Linkedin string `yaml:"linkedin"`
		Github   string `yaml:"github"`
		Email    string `yaml:"email"`
	} `yaml:"contact"`
	Experience []struct {
		Title            string   `yaml:"title"`
		Company          string   `yaml:"company"`
		Location         string   `yaml:"location"`
		Period           string   `yaml:"period"`
		Responsibilities []string `yaml:"responsibilities"`
	} `yaml:"experience"`
	Skills struct {
		ProgrammingLanguages []string `yaml:"programming_languages"`
		CloudInfrastructure  []string `yaml:"cloud_infrastructure"`
		Tools                []string `yaml:"tools"`
	} `yaml:"skills"`
	Education []struct {
		Institution string  `yaml:"institution"`
		Degree      string  `yaml:"degree"`
		CGPA        float64 `yaml:"cgpa"`
		Location    string  `yaml:"location"`
		Period      string  `yaml:"period"`
	} `yaml:"education"`
	Projects []struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		Link        string `yaml:"link"`
	} `yaml:"projects"`
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

}

func main() {

	dataPath := os.Getenv("DATA_PATH")
	if dataPath == "" {
		log.Fatal("DATA_PATH SHOULD BE SET")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address := fmt.Sprintf(":%s", port)

	server := &http.Server{Addr: address}

	http.HandleFunc("/", showHomePage)

	go func() {
		log.Infof("Server running on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// handle shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Errorf("Server shutdown failed: %v", err)
	} else {
		log.Info("Server gracefully stopped")
	}
}

// func LoadData(filePath string) (Resume, error) {

// 	yamlFile, err := os.ReadFile("test.yaml")
// 	if err != nil {
// 		// http.Error(w, "Error reading YAML file", http.StatusInternalServerError)
// 		return nil, err
// 	}

// 	var resume Resume
// 	err = yaml.Unmarshal(yamlFile, &resume)
// 	if err != nil {
// 		return nil, err
// 	}

// }

func showHomePage(w http.ResponseWriter, r *http.Request) {
	// Read and parse YAML file

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid Method", http.StatusBadGateway)
	}

	yamlFile, err := os.ReadFile(os.Getenv("DATA_PATH"))
	if err != nil {
		fmt.Println(yamlFile)
		log.Error("Error while reading the yaml file", err)
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}
	var resume Resume
	err = yaml.Unmarshal(yamlFile, &resume)
	if err != nil {
		log.Error("Error while parsing the yaml file", err)
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}

	// Parse the HTML template
	tmpl, err := template.ParseFiles("internal/views/resume.html")
	if err != nil {
		log.Error("Error while parsing the template file", err)
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}

	// Execute the template with resume data
	err = tmpl.Execute(w, resume)
	if err != nil {
		log.Error("Error while parsing the executing the template file", err)
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}
}
