package detect

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetect_Go(t *testing.T) {
	// Create temp directory with go.mod
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test"), 0600); err != nil {
		t.Fatal(err)
	}

	detections, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !HasLanguage(detections, Go) {
		t.Error("expected Go to be detected")
	}
}

func TestDetect_TypeScript(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "tsconfig.json"), []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}

	detections, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !HasLanguage(detections, TypeScript) {
		t.Error("expected TypeScript to be detected")
	}
}

func TestDetect_JavaScript(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}

	detections, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !HasLanguage(detections, JavaScript) {
		t.Error("expected JavaScript to be detected")
	}
}

func TestDetect_Rust(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "Cargo.toml"), []byte(""), 0600); err != nil {
		t.Fatal(err)
	}

	detections, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !HasLanguage(detections, Rust) {
		t.Error("expected Rust to be detected")
	}
}

func TestDetect_Swift(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "Package.swift"), []byte(""), 0600); err != nil {
		t.Fatal(err)
	}

	detections, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !HasLanguage(detections, Swift) {
		t.Error("expected Swift to be detected")
	}
}

func TestDetect_Python(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{
		{"pyproject.toml", "pyproject.toml"},
		{"setup.py", "setup.py"},
		{"requirements.txt", "requirements.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, tt.file), []byte(""), 0600); err != nil {
				t.Fatal(err)
			}

			detections, err := Detect(dir)
			if err != nil {
				t.Fatalf("Detect failed: %v", err)
			}

			if !HasLanguage(detections, Python) {
				t.Errorf("expected Python to be detected with %s", tt.file)
			}
		})
	}
}

func TestDetect_Empty(t *testing.T) {
	dir := t.TempDir()

	detections, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if len(detections) != 0 {
		t.Errorf("expected 0 detections, got %d", len(detections))
	}
}

func TestDetect_MultiLanguage(t *testing.T) {
	dir := t.TempDir()

	// Create Go project in backend/
	backend := filepath.Join(dir, "backend")
	if err := os.MkdirAll(backend, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(backend, "go.mod"), []byte("module test"), 0600); err != nil {
		t.Fatal(err)
	}

	// Create TypeScript project in frontend/
	frontend := filepath.Join(dir, "frontend")
	if err := os.MkdirAll(frontend, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(frontend, "package.json"), []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(frontend, "tsconfig.json"), []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}

	detections, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if !HasLanguage(detections, Go) {
		t.Error("expected Go to be detected")
	}
	if !HasLanguage(detections, TypeScript) {
		t.Error("expected TypeScript to be detected")
	}
}

func TestGetByLanguage(t *testing.T) {
	detections := []Detection{
		{Language: Go, Path: "backend"},
		{Language: TypeScript, Path: "frontend"},
		{Language: Go, Path: "tools"},
	}

	goDetections := GetByLanguage(detections, Go)
	if len(goDetections) != 2 {
		t.Errorf("expected 2 Go detections, got %d", len(goDetections))
	}

	tsDetections := GetByLanguage(detections, TypeScript)
	if len(tsDetections) != 1 {
		t.Errorf("expected 1 TypeScript detection, got %d", len(tsDetections))
	}
}

func TestHasLanguage(t *testing.T) {
	detections := []Detection{
		{Language: Go, Path: "."},
	}

	if !HasLanguage(detections, Go) {
		t.Error("expected HasLanguage to return true for Go")
	}
	if HasLanguage(detections, Python) {
		t.Error("expected HasLanguage to return false for Python")
	}
}
