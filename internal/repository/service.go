package repository

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// RepoInfo contains information about a repository
type RepoInfo struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	URL      string    `json:"url"`
	ClonedAt time.Time `json:"clonedAt"`
}

// FileInfo contains information about a file in a repository
type FileInfo struct {
	Path     string `json:"path"`
	Language string `json:"language"`
	Size     int64  `json:"size"`
}

// Service provides repository operations
type Service struct {
	WorkspacePath string
}

// NewService creates a new repository service
func NewService(workspacePath string) (*Service, error) {
	// Create workspace directory if it doesn't exist
	if err := os.MkdirAll(workspacePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create workspace directory: %w", err)
	}

	return &Service{
		WorkspacePath: workspacePath,
	}, nil
}

// CloneRepository clones a Git repository
func (s *Service) CloneRepository(url string) (*RepoInfo, error) {
	// Extract repository name from URL
	repoName := filepath.Base(url)
	if filepath.Ext(repoName) == ".git" {
		repoName = repoName[:len(repoName)-4]
	}

	// Create a unique directory name
	timestamp := time.Now().Format("20060102150405")
	dirName := fmt.Sprintf("%s-%s", repoName, timestamp)
	repoPath := filepath.Join(s.WorkspacePath, dirName)

	// Clone the repository
	repo, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to clone repository: %w", err)
	}

	// Get repository info
	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get repository head: %w", err)
	}

	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get head commit: %w", err)
	}

	return &RepoInfo{
		ID:       dirName,
		Name:     repoName,
		Path:     dirName,
		URL:      url,
		ClonedAt: commit.Committer.When,
	}, nil
}

// ListRepositories lists all repositories in the workspace
func (s *Service) ListRepositories() ([]RepoInfo, error) {
	entries, err := os.ReadDir(s.WorkspacePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read workspace directory: %w", err)
	}

	repos := []RepoInfo{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		repoPath := filepath.Join(s.WorkspacePath, entry.Name())
		repo, err := git.PlainOpen(repoPath)
		if err != nil {
			// Not a Git repository, skip
			continue
		}

		config, err := repo.Config()
		if err != nil {
			continue
		}

		url := ""
		if config.Remotes["origin"] != nil && len(config.Remotes["origin"].URLs) > 0 {
			url = config.Remotes["origin"].URLs[0]
		}

		repoName := entry.Name()
		// Remove timestamp if present
		if len(repoName) > 15 && repoName[len(repoName)-15] == '-' {
			repoName = repoName[:len(repoName)-15]
		}

		repos = append(repos, RepoInfo{
			ID:       entry.Name(),
			Name:     repoName,
			Path:     entry.Name(),
			URL:      url,
			ClonedAt: info.ModTime(),
		})
	}

	return repos, nil
}

// GetRepository gets a repository by ID
func (s *Service) GetRepository(id string) (*RepoInfo, error) {
	repoPath := filepath.Join(s.WorkspacePath, id)

	// Check if the repository exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return nil, errors.New("repository not found")
	}

	// Open the repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Get repository info
	config, err := repo.Config()
	if err != nil {
		return nil, fmt.Errorf("failed to get repository config: %w", err)
	}

	url := ""
	if config.Remotes["origin"] != nil && len(config.Remotes["origin"].URLs) > 0 {
		url = config.Remotes["origin"].URLs[0]
	}

	// Get repository name
	repoName := id
	if len(repoName) > 15 && repoName[len(repoName)-15] == '-' {
		repoName = repoName[:len(repoName)-15]
	}

	// Get repository creation time
	fileInfo, err := os.Stat(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository info: %w", err)
	}

	return &RepoInfo{
		ID:       id,
		Name:     repoName,
		Path:     id,
		URL:      url,
		ClonedAt: fileInfo.ModTime(),
	}, nil
}

// ListFiles lists all files in a repository
func (s *Service) ListFiles(repoID string) ([]FileInfo, error) {
	repoPath := filepath.Join(s.WorkspacePath, repoID)

	// Check if the repository exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return nil, errors.New("repository not found")
	}

	// Open the repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Get the HEAD reference
	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get repository head: %w", err)
	}

	// Get the commit object for the HEAD reference
	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get head commit: %w", err)
	}

	// Get the tree for the commit
	tree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit tree: %w", err)
	}

	// Create a list to store file info
	files := []FileInfo{}

	// Iterate through the tree
	err = tree.Files().ForEach(func(f *object.File) error {
		// Skip directories
		if f.Mode&0040000 == 0040000 {
			return nil
		}

		// Get file info
		language := detectLanguage(f.Name)
		size := f.Size

		files = append(files, FileInfo{
			Path:     f.Name,
			Language: language,
			Size:     size,
		})

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate through files: %w", err)
	}

	return files, nil
}

// detectLanguage detects the language of a file based on its extension
func detectLanguage(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".go":
		return "Go"
	case ".js", ".jsx":
		return "JavaScript"
	case ".ts", ".tsx":
		return "TypeScript"
	case ".py":
		return "Python"
	case ".java":
		return "Java"
	case ".c", ".cpp", ".cc", ".h", ".hpp":
		return "C/C++"
	case ".html", ".htm":
		return "HTML"
	case ".css":
		return "CSS"
	case ".md":
		return "Markdown"
	case ".json":
		return "JSON"
	case ".yaml", ".yml":
		return "YAML"
	default:
		return "Unknown"
	}
}
