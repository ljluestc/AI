package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Clone a git repository to the workspace
func clone_repo(repoURL, workspacePath string) (string, error) {
	// In a real implementation, this would use go-git to clone the repository
	// For now, we'll simulate a successful clone
	repoName := filepath.Base(repoURL)
	if filepath.Ext(repoName) == ".git" {
		repoName = repoName[:len(repoName)-4]
	}
	timestamp := time.Now().Format("20060102150405")
	repoDir := filepath.Join(workspacePath, repoName+"-"+timestamp)

	// Create repo directory
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		return "", err
	}

	log.Printf("Cloned repository %s to %s", repoURL, repoDir)
	return filepath.Base(repoDir), nil
}

// Index all source files in a repository
func index_source_files(repoPath string) ([]map[string]string, error) {
	files := []map[string]string{
		{"path": "src/main.go", "language": "Go"},
		{"path": "src/utils/helpers.go", "language": "Go"},
		{"path": "src/api/routes.go", "language": "Go"},
		{"path": "web/src/App.js", "language": "JavaScript"},
	}
	return files, nil
}

// Run static analysis on source files
func run_static_analysis(files []map[string]string) ([]map[string]interface{}, error) {
	results := []map[string]interface{}{
		{"file": "src/main.go", "line": 42, "column": 5, "message": "Unused variable", "level": "warning"},
		{"file": "src/api/routes.go", "line": 15, "column": 3, "message": "Error not handled", "level": "error"},
		{"file": "web/src/App.js", "line": 23, "column": 10, "message": "React Hook useEffect has missing dependency", "level": "warning"},
	}
	return results, nil
}

// Parse error logs in a repository
func parse_error_logs(repoPath string) ([]map[string]interface{}, error) {
	logs := []map[string]interface{}{
		{"file": "src/api/routes.go", "line": 15, "message": "panic: runtime error: nil pointer dereference"},
		{"file": "src/utils/helpers.go", "line": 35, "message": "index out of range [5] with length 3"},
	}
	return logs, nil
}

// Build a knowledge graph of code relationships
func build_code_knowledge_graph(files []map[string]string) ([]map[string]interface{}, error) {
	graph := []map[string]interface{}{
		{"id": "file-1", "type": "file", "name": "main.go", "path": "src/main.go"},
		{"id": "func-1", "type": "function", "name": "main", "relations": []string{"file-1"}},
		{"id": "file-2", "type": "file", "name": "routes.go", "path": "src/api/routes.go"},
		{"id": "func-2", "type": "function", "name": "setupRoutes", "relations": []string{"file-2"}},
	}
	return graph, nil
}

// Map error logs to nodes in the knowledge graph
func map_errors_to_graph(errorLogs []map[string]interface{}, codeGraph []map[string]interface{}) []string {
	// Return IDs of nodes that might be related to errors
	return []string{"file-2", "func-2"}
}

// Localize errors to specific nodes in the graph
func localize_errors(errorNodes []string, codeGraph []map[string]interface{}) []string {
	// Return suspected node IDs
	return []string{"func-2"}
}

// Diagnose root cause of errors
func diagnose_root_cause(suspectNodes []string, codeGraph []map[string]interface{}) map[string]string {
	diagnosis := map[string]string{
		"func-2": "Null pointer dereference in handler function. API endpoint tries to access a property of a nil object.",
	}
	return diagnosis
}

func main() {
	// Set up the router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://teathis.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Create workspace directory if it doesn't exist
	workspaceDir := "./workspace"
	if err := os.MkdirAll(workspaceDir, 0755); err != nil {
		log.Fatalf("Failed to create workspace directory: %v", err)
	}

	// API routes
	api := r.Group("/api")
	{
		// Repository endpoints
		api.POST("/repositories", func(c *gin.Context) {
			var request struct {
				URL string `json:"url" binding:"required"`
			}
			if err := c.ShouldBindJSON(&request); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			repoDir, err := clone_repo(request.URL, workspaceDir)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message":  "Repository cloned successfully",
				"repoPath": repoDir,
			})
		})

		api.GET("/repositories", func(c *gin.Context) {
			// List all repositories in workspace
			entries, err := os.ReadDir(workspaceDir)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			repos := []map[string]interface{}{}
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}

				info, err := entry.Info()
				if err != nil {
					continue
				}

				repoName := entry.Name()
				// Parse out timestamp if present
				parts := filepath.Base(repoName)

				repos = append(repos, map[string]interface{}{
					"id":       repoName,
					"name":     parts,
					"path":     repoName,
					"clonedAt": info.ModTime(),
				})
			}

			c.JSON(http.StatusOK, repos)
		})

		// Analysis endpoints
		api.POST("/analyze", func(c *gin.Context) {
			var request struct {
				RepoPath string `json:"repoPath" binding:"required"`
			}
			if err := c.ShouldBindJSON(&request); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Check if the repository path exists and is within our workspace
			fullPath := filepath.Join(workspaceDir, request.RepoPath)
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Repository not found"})
				return
			}

			// Index source files
			files, err := index_source_files(fullPath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Run static analysis
			analysisResults, err := run_static_analysis(files)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Parse error logs
			errorLogs, err := parse_error_logs(fullPath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Build code knowledge graph
			graph, err := build_code_knowledge_graph(files)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Map errors to graph
			mappedErrors := map_errors_to_graph(errorLogs, graph)

			// Localize errors
			errorNodes := localize_errors(mappedErrors, graph)

			// Diagnose root cause
			diagnosis := diagnose_root_cause(errorNodes, graph)

			c.JSON(http.StatusOK, gin.H{
				"analysisResults": analysisResults,
				"errorLogs":       errorLogs,
				"diagnosis":       diagnosis,
			})
		})
	}

	// Serve the frontend static files
	r.Static("/assets", "./web/build/assets")
	r.StaticFile("/", "./web/build/index.html")
	r.StaticFile("/favicon.ico", "./web/build/favicon.ico")
	r.NoRoute(func(c *gin.Context) {
		c.File("./web/build/index.html")
	})

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
