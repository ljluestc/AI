package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SourceFile represents a source code file
type SourceFile struct {
	Path     string `json:"path"`
	Language string `json:"language"`
	Content  string `json:"content,omitempty"`
}

// AnalysisResult represents results from static analysis
type AnalysisResult struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Message string `json:"message"`
	Level   string `json:"level"` // error, warning, info
}

// ErrorLog represents an error from a log file
type ErrorLog struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Message string `json:"message"`
}

// GraphNode represents a node in the code knowledge graph
type GraphNode struct {
	ID        string   `json:"id"`
	Type      string   `json:"type"` // file, function, class, etc.
	Name      string   `json:"name"`
	Path      string   `json:"path,omitempty"`
	Relations []string `json:"relations,omitempty"`
}

// Service provides code analysis operations
type Service struct{}

// NewService creates a new analyzer service
func NewService() *Service {
	return &Service{}
}

// IndexSourceFiles finds and indexes all source files in a repository
func (s *Service) IndexSourceFiles(repoPath string) ([]SourceFile, error) {
	var files []SourceFile

	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Skip .git and other hidden directories
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		// Determine language based on file extension
		ext := strings.ToLower(filepath.Ext(path))
		var language string

		switch ext {
		case ".go":
			language = "Go"
		case ".js", ".jsx":
			language = "JavaScript"
		case ".ts", ".tsx":
			language = "TypeScript"
		case ".py":
			language = "Python"
		case ".java":
			language = "Java"
		case ".c", ".cpp", ".cc", ".h", ".hpp":
			language = "C/C++"
		default:
			// Skip non-source files
			return nil
		}

		// Get relative path from the repo root
		relPath, err := filepath.Rel(repoPath, path)
		if err != nil {
			return err
		}

		files = append(files, SourceFile{
			Path:     relPath,
			Language: language,
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to index source files: %w", err)
	}

	return files, nil
}

// RunStaticAnalysis performs static code analysis on the source files
func (s *Service) RunStaticAnalysis(files []SourceFile) ([]AnalysisResult, error) {
	// This would typically involve running actual static analysis tools
	// Here we're providing a simplified mock implementation
	var results []AnalysisResult

	for _, file := range files {
		// Mock analysis based on file extension
		switch strings.ToLower(filepath.Ext(file.Path)) {
		case ".go":
			results = append(results, AnalysisResult{
				File:    file.Path,
				Line:    10,
				Column:  5,
				Message: "Function exceeds maximum cyclomatic complexity",
				Level:   "warning",
			})
		case ".js", ".jsx", ".ts", ".tsx":
			results = append(results, AnalysisResult{
				File:    file.Path,
				Line:    20,
				Column:  15,
				Message: "Unused variable",
				Level:   "warning",
			})
		case ".py":
			results = append(results, AnalysisResult{
				File:    file.Path,
				Line:    30,
				Column:  1,
				Message: "Missing docstring",
				Level:   "info",
			})
		}
	}

	return results, nil
}

// ParseErrorLogs parses build or runtime error logs
func (s *Service) ParseErrorLogs(repoPath string) ([]ErrorLog, error) {
	// This would typically involve parsing actual error logs
	// Here we're providing a simplified mock implementation
	logs := []ErrorLog{
		{
			File:    "main.go",
			Line:    42,
			Message: "undefined: someVariable",
		},
		{
			File:    "handlers/user.go",
			Line:    15,
			Message: "cannot use x (type int) as type string",
		},
	}

	return logs, nil
}

// BuildCodeKnowledgeGraph builds a graph representation of the code
func (s *Service) BuildCodeKnowledgeGraph(files []SourceFile) ([]GraphNode, error) {
	// This would typically involve building an actual code graph
	// Here we're providing a simplified mock implementation
	var graph []GraphNode

	for i, file := range files {
		// Create a node for each file
		fileID := fmt.Sprintf("file-%d", i)
		graph = append(graph, GraphNode{
			ID:   fileID,
			Type: "file",
			Name: filepath.Base(file.Path),
			Path: file.Path,
		})

		// Add some mock function nodes for each file
		for j := 0; j < 3; j++ {
			funcID := fmt.Sprintf("func-%d-%d", i, j)
			graph = append(graph, GraphNode{
				ID:        funcID,
				Type:      "function",
				Name:      fmt.Sprintf("function%d", j),
				Relations: []string{fileID},
			})
		}
	}

	return graph, nil
}

// MapErrorsToGraph maps error logs to nodes in the code graph
func (s *Service) MapErrorsToGraph(logs []ErrorLog, graph []GraphNode) []string {
	// This would typically involve complex mapping logic
	// Here we're providing a simplified mock implementation
	var errorNodeIDs []string

	for _, log := range logs {
		for _, node := range graph {
			if node.Type == "file" && strings.HasSuffix(node.Path, log.File) {
				errorNodeIDs = append(errorNodeIDs, node.ID)
				break
			}
		}
	}

	return errorNodeIDs
}

// LocalizeErrors finds specific nodes in the graph responsible for errors
func (s *Service) LocalizeErrors(errorNodeIDs []string, graph []GraphNode) []string {
	// This would typically involve graph traversal algorithms
	// Here we're providing a simplified mock implementation
	var suspectNodeIDs []string

	// For simplicity, just collect related function nodes
	for _, errorID := range errorNodeIDs {
		for _, node := range graph {
			if node.Type == "function" {
				for _, relID := range node.Relations {
					if relID == errorID {
						suspectNodeIDs = append(suspectNodeIDs, node.ID)
					}
				}
			}
		}
	}

	return suspectNodeIDs
}

// DiagnoseRootCause identifies the root cause of the errors
func (s *Service) DiagnoseRootCause(suspectNodeIDs []string, graph []GraphNode) map[string]string {
	// This would typically involve sophisticated analysis
	// Here we're providing a simplified mock implementation
	diagnosis := make(map[string]string)

	for _, nodeID := range suspectNodeIDs {
		for _, node := range graph {
			if node.ID == nodeID {
				diagnosis[nodeID] = fmt.Sprintf("Potential issue in %s: incorrect API usage", node.Name)
			}
		}
	}

	return diagnosis
}
