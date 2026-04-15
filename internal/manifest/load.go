package manifest

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	manifestv1alpha1 "github.com/wiscotrashpanda/alloy/manifest/v1alpha1"
	"gopkg.in/yaml.v3"
)

type LoadedGitHubRepositoryManifest struct {
	Path     string
	Manifest manifestv1alpha1.GitHubRepositoryManifest
}

type Result struct {
	GitHubRepositories []LoadedGitHubRepositoryManifest
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

func loadGitHubRepositoryManifest(path string) (LoadedGitHubRepositoryManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return LoadedGitHubRepositoryManifest{}, fmt.Errorf("read manifest %q: %w", path, err)
	}

	var envelope manifestv1alpha1.Envelope
	if err := yaml.Unmarshal(data, &envelope); err != nil {
		return LoadedGitHubRepositoryManifest{}, fmt.Errorf("parse manifest %q: %w", path, err)
	}

	if envelope.APIVersion == "" {
		return LoadedGitHubRepositoryManifest{}, fmt.Errorf("manifest %q missing apiVersion", path)
	}

	if envelope.Kind == "" {
		return LoadedGitHubRepositoryManifest{}, fmt.Errorf("manifest %q missing kind", path)
	}

	if envelope.Metadata.Name == "" {
		return LoadedGitHubRepositoryManifest{}, fmt.Errorf("manifest %q missing metadata.name", path)
	}

	switch envelope.Kind {
	case manifestv1alpha1.KindGitHubRepository:
		var spec manifestv1alpha1.GitHubRepositorySpec
		if err := envelope.Spec.Decode(&spec); err != nil {
			return LoadedGitHubRepositoryManifest{}, fmt.Errorf("decode GitHubRepository spec in %q: %w", path, err)
		}

		manifest := manifestv1alpha1.NewGitHubRepositoryManifest(envelope.Metadata, spec)
		manifest.APIVersion = envelope.APIVersion

		if err := manifest.Validate(); err != nil {
			return LoadedGitHubRepositoryManifest{}, fmt.Errorf("manifest %q %w", path, err)
		}

		return LoadedGitHubRepositoryManifest{
			Path:     path,
			Manifest: manifest,
		}, nil
	default:
		return LoadedGitHubRepositoryManifest{}, fmt.Errorf("manifest %q has unsupported kind %q", path, envelope.Kind)
	}
}

func isYAMLFile(name string) bool {
	lowerName := strings.ToLower(name)
	return strings.HasSuffix(lowerName, ".yaml") || strings.HasSuffix(lowerName, ".yml")
}
