package repo

import (
	"os"
	"path/filepath"
	"testing"
)

func createTestPackage(t *testing.T, dir string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Join(dir, "templates"), 0755); nil != err {
		t.Fatal(err)
	}

	hullYAML := `apiVersion: v1
name: testpkg
version: 1.0.0
description: A test package
`
	if err := os.WriteFile(filepath.Join(dir, "hull.yaml"), []byte(hullYAML), 0644); nil != err {
		t.Fatal(err)
	}

	valuesYAML := `replicaCount: 1
image:
  repository: nginx
  tag: latest
`
	if err := os.WriteFile(filepath.Join(dir, "values.yaml"), []byte(valuesYAML), 0644); nil != err {
		t.Fatal(err)
	}

	template := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Values.name }}
`
	if err := os.WriteFile(filepath.Join(dir, "templates", "deployment.yaml"), []byte(template), 0644); nil != err {
		t.Fatal(err)
	}
}

func TestPackageArchive(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()

	createTestPackage(t, srcDir)

	archivePath, err := PackageArchive(srcDir, destDir, "")
	if nil != err {
		t.Fatalf("PackageArchive failed: %v", err)
	}

	expectedName := "testpkg-1.0.0.hull.tgz"
	if filepath.Base(archivePath) != expectedName {
		t.Errorf("expected archive name %s, got %s", expectedName, filepath.Base(archivePath))
	}

	info, err := os.Stat(archivePath)
	if nil != err {
		t.Fatalf("archive file not found: %v", err)
	}
	if info.Size() == 0 {
		t.Error("archive file is empty")
	}
}

func TestPackageArchiveVersionOverride(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()

	createTestPackage(t, srcDir)

	archivePath, err := PackageArchive(srcDir, destDir, "2.0.0")
	if nil != err {
		t.Fatalf("PackageArchive failed: %v", err)
	}

	expectedName := "testpkg-2.0.0.hull.tgz"
	if filepath.Base(archivePath) != expectedName {
		t.Errorf("expected archive name %s, got %s", expectedName, filepath.Base(archivePath))
	}
}

func TestExtractArchive(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()
	extractDir := t.TempDir()

	createTestPackage(t, srcDir)

	archivePath, err := PackageArchive(srcDir, destDir, "")
	if nil != err {
		t.Fatalf("PackageArchive failed: %v", err)
	}

	if err := ExtractArchive(archivePath, extractDir); nil != err {
		t.Fatalf("ExtractArchive failed: %v", err)
	}

	// Verify hull.yaml exists inside the extracted package dir
	hullYAMLPath := filepath.Join(extractDir, "testpkg", "hull.yaml")
	if _, err := os.Stat(hullYAMLPath); nil != err {
		t.Errorf("expected hull.yaml at %s: %v", hullYAMLPath, err)
	}

	// Verify templates directory
	templatePath := filepath.Join(extractDir, "testpkg", "templates", "deployment.yaml")
	if _, err := os.Stat(templatePath); nil != err {
		t.Errorf("expected deployment.yaml at %s: %v", templatePath, err)
	}
}

func TestHullignore(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()
	extractDir := t.TempDir()

	createTestPackage(t, srcDir)

	// Create .hullignore file
	ignoreContent := "# Comment line\nvalues.yaml\n"
	if err := os.WriteFile(filepath.Join(srcDir, ".hullignore"), []byte(ignoreContent), 0644); nil != err {
		t.Fatal(err)
	}

	archivePath, err := PackageArchive(srcDir, destDir, "")
	if nil != err {
		t.Fatalf("PackageArchive failed: %v", err)
	}

	if err := ExtractArchive(archivePath, extractDir); nil != err {
		t.Fatalf("ExtractArchive failed: %v", err)
	}

	// values.yaml should be ignored
	valuesPath := filepath.Join(extractDir, "testpkg", "values.yaml")
	if _, err := os.Stat(valuesPath); !os.IsNotExist(err) {
		t.Error("values.yaml should have been ignored but was included in archive")
	}

	// hull.yaml should still be present
	hullPath := filepath.Join(extractDir, "testpkg", "hull.yaml")
	if _, err := os.Stat(hullPath); nil != err {
		t.Errorf("hull.yaml should be present: %v", err)
	}

	// .hullignore should itself be ignored (default pattern)
	ignorePath := filepath.Join(extractDir, "testpkg", ".hullignore")
	if _, err := os.Stat(ignorePath); !os.IsNotExist(err) {
		t.Error(".hullignore should have been ignored but was included in archive")
	}
}

func TestHullignoreNegation(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()
	extractDir := t.TempDir()

	createTestPackage(t, srcDir)

	// Create a .txt file in templates that we would normally include
	if err := os.WriteFile(filepath.Join(srcDir, "templates", "notes.txt"), []byte("notes"), 0644); nil != err {
		t.Fatal(err)
	}

	// Ignore all .txt files but negate notes.txt
	ignoreContent := "*.txt\n!notes.txt\n"
	if err := os.WriteFile(filepath.Join(srcDir, ".hullignore"), []byte(ignoreContent), 0644); nil != err {
		t.Fatal(err)
	}

	archivePath, err := PackageArchive(srcDir, destDir, "")
	if nil != err {
		t.Fatalf("PackageArchive failed: %v", err)
	}

	if err := ExtractArchive(archivePath, extractDir); nil != err {
		t.Fatalf("ExtractArchive failed: %v", err)
	}

	// notes.txt should be present (negation)
	notesPath := filepath.Join(extractDir, "testpkg", "templates", "notes.txt")
	if _, err := os.Stat(notesPath); nil != err {
		t.Errorf("notes.txt should be present due to negation: %v", err)
	}
}

func TestShouldIgnore(t *testing.T) {
	patterns := []string{".git/", "*.hull.tgz", "secret.yaml"}

	tests := []struct {
		path     string
		isDir    bool
		expected bool
	}{
		{".git", true, true},
		{".git/objects", false, false},
		{"templates/deployment.yaml", false, false},
		{"mypackage-1.0.0.hull.tgz", false, true},
		{"secret.yaml", false, true},
		{"templates/secret.yaml", false, true},
	}

	for _, tc := range tests {
		result := shouldIgnore(tc.path, tc.isDir, patterns)
		if result != tc.expected {
			t.Errorf("shouldIgnore(%q, %v) = %v, want %v", tc.path, tc.isDir, result, tc.expected)
		}
	}
}

func TestDefaultIgnorePatterns(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()
	extractDir := t.TempDir()

	createTestPackage(t, srcDir)

	// Create a .git directory (should be ignored by default)
	gitDir := filepath.Join(srcDir, ".git")
	if err := os.MkdirAll(gitDir, 0755); nil != err {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(gitDir, "config"), []byte("gitconfig"), 0644); nil != err {
		t.Fatal(err)
	}

	archivePath, err := PackageArchive(srcDir, destDir, "")
	if nil != err {
		t.Fatalf("PackageArchive failed: %v", err)
	}

	if err := ExtractArchive(archivePath, extractDir); nil != err {
		t.Fatalf("ExtractArchive failed: %v", err)
	}

	// .git should not be in the archive
	gitPath := filepath.Join(extractDir, "testpkg", ".git")
	if _, err := os.Stat(gitPath); !os.IsNotExist(err) {
		t.Error(".git directory should have been ignored but was included in archive")
	}
}
