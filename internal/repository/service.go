package repository

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// RepoInfo represents information about a cloned repository
type RepoInfo struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Path      string    `json:"path"`
	ClonedAt  time.Time `json:"clonedAt"`
	LastScan  time.Time `json:"lastScan,omitempty"`
	ErrorsNum int       `json:"errorsNum"`
}

// Service provides repository operations
type Service struct {
	workspaceDir string
}

// NewService creates a new repository service
func NewService(workspaceDir string) *Service {
	return &Service{
		workspaceDir: workspaceDir,
	}
}

// CloneRepo clones a git repository to the workspace
func (s *Service) CloneRepo(repoURL string) (string, error) {
	// Extract repo name from URL
	urlParts := strings.Split(repoURL, "/")
	repoName := strings.TrimSuffix(urlParts[len(urlParts)-1], ".git")

	// Create a unique directory name
	timestamp := time.Now().Format("20060102150405")
	repoDir := fmt.Sprintf("%s-%s", repoName, timestamp)
	fullPath := filepath.Join(s.workspaceDir, repoDir)

	// Clone the repository
	_, err := git.PlainClone(fullPath, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})
	if err != nil {
		return "", fmt.Errorf("failed to clone repository: %w", err)
	}

	return repoDir, nil
}

// ListRepos returns information about all cloned repositories
func (s *Service) ListRepos() ([]RepoInfo, error) {
	var repos []RepoInfo

	// Read the workspace directory
	entries, err := os.ReadDir(s.workspaceDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read workspace directory: %w", err)
	}

	// Iterate through each entry and check if it's a git repository
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		fullPath := filepath.Join(s.workspaceDir, entry.Name())

		// Check if it's a git repository
		_, err := git.PlainOpen(fullPath)
		if err != nil {
			continue // Not a git repository, skip
		}

		// Get directory info
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Get repository name (without timestamp)
		nameParts := strings.Split(entry.Name(), "-")
		name := strings.Join(nameParts[:len(nameParts)-1], "-")

		repos = append(repos, RepoInfo{
			ID:       entry.Name(),
			Name:     name,
			Path:     entry.Name(),
			ClonedAt: info.ModTime(),
		})
	}

	return repos, nil
}

// GetRepoByID returns information about a specific repository
func (s *Service) GetRepoByID(repoID string) (*RepoInfo, error) {
	fullPath := filepath.Join(s.workspaceDir, repoID)

	// Check if the directory exists
	info, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("repository not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get repository info: %w", err)
	}

	// Check if it's a git repository
	_, err = git.PlainOpen(fullPath)
	if err != nil {
		return nil, fmt.Errorf("not a valid git repository: %w", err)
	}

	// Get repository name (without timestamp)
	nameParts := strings.Split(repoID, "-")
	name := strings.Join(nameParts[:len(nameParts)-1], "-")

	return &RepoInfo{
		ID:       repoID,
		Name:     name,
		Path:     repoID,
		ClonedAt: info.ModTime(),
	}, nil
}
