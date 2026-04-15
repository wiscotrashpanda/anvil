package manifest

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type Metadata struct {
	Name string `yaml:"name"`
}

type Envelope struct {
	APIVersion string    `yaml:"apiVersion"`
	Kind       string    `yaml:"kind"`
	Metadata   Metadata  `yaml:"metadata"`
	Spec       yaml.Node `yaml:"spec"`
}

type GitHubRepositorySpec struct {
	Owner       string `yaml:"owner"`
	Name        string `yaml:"name"`
	Visibility  string `yaml:"visibility"`
	Description string `yaml:"description"`
	AutoInit    bool   `yaml:"autoInit"`
}

type GitHubRepositoryManifest struct {
	Path       string
	APIVersion string
	Metadata   Metadata
	Spec       GitHubRepositorySpec
}

type Result struct {
	GitHubRepositories []GitHubRepositoryManifest
}

func LoadDir(path string) (Result, error) {
	info, err := os.Stat(path)
	if err != nil {
		return Result{}, fmt.Errorf("stat manifests path %q: %w", path, err)
	}

	if !info.IsDir() {
		return Result{}, fmt.Errorf("manifests path must be a directory: %s", path)
	}

	var manifestPaths []string

	err = filepath.WalkDir(path, func(currentPath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if entry.IsDir() {
			return nil
		}

		if !isYAMLFile(entry.Name()) {
			return nil
		}

		manifestPaths = append(manifestPaths, currentPath)
		return nil
	})
	if err != nil {
		return Result{}, fmt.Errorf("walk manifests directory %q: %w", path, err)
	}

	sort.Strings(manifestPaths)

	var result Result

	for _, manifestPath := range manifestPaths {
		repository, err := loadGitHubRepositoryManifest(manifestPath)
		if err != nil {
			return Result{}, err
		}

		result.GitHubRepositories = append(result.GitHubRepositories, repository)
	}

	return result, nil
}

func loadGitHubRepositoryManifest(path string) (GitHubRepositoryManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return GitHubRepositoryManifest{}, fmt.Errorf("read manifest %q: %w", path, err)
	}

	var envelope Envelope
	if err := yaml.Unmarshal(data, &envelope); err != nil {
		return GitHubRepositoryManifest{}, fmt.Errorf("parse manifest %q: %w", path, err)
	}

	if envelope.APIVersion == "" {
		return GitHubRepositoryManifest{}, fmt.Errorf("manifest %q missing apiVersion", path)
	}

	if envelope.Kind == "" {
		return GitHubRepositoryManifest{}, fmt.Errorf("manifest %q missing kind", path)
	}

	if envelope.Metadata.Name == "" {
		return GitHubRepositoryManifest{}, fmt.Errorf("manifest %q missing metadata.name", path)
	}

	switch envelope.Kind {
	case "GitHubRepository":
		var spec GitHubRepositorySpec
		if err := envelope.Spec.Decode(&spec); err != nil {
			return GitHubRepositoryManifest{}, fmt.Errorf("decode GitHubRepository spec in %q: %w", path, err)
		}

		if spec.Owner == "" {
			return GitHubRepositoryManifest{}, fmt.Errorf("manifest %q missing spec.owner", path)
		}

		if spec.Name == "" {
			return GitHubRepositoryManifest{}, fmt.Errorf("manifest %q missing spec.name", path)
		}

		return GitHubRepositoryManifest{
			Path:       path,
			APIVersion: envelope.APIVersion,
			Metadata:   envelope.Metadata,
			Spec:       spec,
		}, nil
	default:
		return GitHubRepositoryManifest{}, fmt.Errorf("manifest %q has unsupported kind %q", path, envelope.Kind)
	}
}

func isYAMLFile(name string) bool {
	lowerName := strings.ToLower(name)
	return strings.HasSuffix(lowerName, ".yaml") || strings.HasSuffix(lowerName, ".yml")
}
