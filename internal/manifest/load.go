package manifest

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	manifestv1alpha1 "github.com/emkaytec/alloy/manifest/v1alpha1"
	"gopkg.in/yaml.v3"
)

type LoadedGitHubRepositoryManifest struct {
	Path     string
	Manifest manifestv1alpha1.GitHubRepositoryManifest
}

type LoadedHCPTerraformWorkspaceManifest struct {
	Path     string
	Manifest manifestv1alpha1.HCPTerraformWorkspaceManifest
}

type Result struct {
	GitHubRepositories     []LoadedGitHubRepositoryManifest
	HCPTerraformWorkspaces []LoadedHCPTerraformWorkspaceManifest
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
		envelope, err := loadManifestEnvelope(manifestPath)
		if err != nil {
			return Result{}, err
		}

		switch envelope.Kind {
		case manifestv1alpha1.KindGitHubRepository:
			repository, err := decodeGitHubRepositoryManifest(manifestPath, envelope)
			if err != nil {
				return Result{}, err
			}

			result.GitHubRepositories = append(result.GitHubRepositories, repository)
		case manifestv1alpha1.KindHCPTerraformWorkspace:
			workspace, err := decodeHCPTerraformWorkspaceManifest(manifestPath, envelope)
			if err != nil {
				return Result{}, err
			}

			result.HCPTerraformWorkspaces = append(result.HCPTerraformWorkspaces, workspace)
		default:
			return Result{}, fmt.Errorf("manifest %q has unsupported kind %q", manifestPath, envelope.Kind)
		}
	}

	return result, nil
}

func loadManifestEnvelope(path string) (manifestv1alpha1.Envelope, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return manifestv1alpha1.Envelope{}, fmt.Errorf("read manifest %q: %w", path, err)
	}

	var envelope manifestv1alpha1.Envelope
	if err := yaml.Unmarshal(data, &envelope); err != nil {
		return manifestv1alpha1.Envelope{}, fmt.Errorf("parse manifest %q: %w", path, err)
	}

	if envelope.APIVersion == "" {
		return manifestv1alpha1.Envelope{}, fmt.Errorf("manifest %q missing apiVersion", path)
	}

	if envelope.Kind == "" {
		return manifestv1alpha1.Envelope{}, fmt.Errorf("manifest %q missing kind", path)
	}

	if envelope.Metadata.Name == "" {
		return manifestv1alpha1.Envelope{}, fmt.Errorf("manifest %q missing metadata.name", path)
	}

	return envelope, nil
}

func decodeGitHubRepositoryManifest(path string, envelope manifestv1alpha1.Envelope) (LoadedGitHubRepositoryManifest, error) {
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
}

func decodeHCPTerraformWorkspaceManifest(path string, envelope manifestv1alpha1.Envelope) (LoadedHCPTerraformWorkspaceManifest, error) {
	var spec manifestv1alpha1.HCPTerraformWorkspaceSpec
	if err := envelope.Spec.Decode(&spec); err != nil {
		return LoadedHCPTerraformWorkspaceManifest{}, fmt.Errorf("decode HCPTerraformWorkspace spec in %q: %w", path, err)
	}

	manifest := manifestv1alpha1.NewHCPTerraformWorkspaceManifest(envelope.Metadata, spec)
	manifest.APIVersion = envelope.APIVersion

	if err := manifest.Validate(); err != nil {
		return LoadedHCPTerraformWorkspaceManifest{}, fmt.Errorf("manifest %q %w", path, err)
	}

	return LoadedHCPTerraformWorkspaceManifest{
		Path:     path,
		Manifest: manifest,
	}, nil
}

func isYAMLFile(name string) bool {
	lowerName := strings.ToLower(name)
	return strings.HasSuffix(lowerName, ".yaml") || strings.HasSuffix(lowerName, ".yml")
}
