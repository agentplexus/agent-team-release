// Package detect provides language detection for repositories.
package detect

import (
	"os"
	"path/filepath"
)

// Language represents a detected programming language.
type Language string

const (
	Go         Language = "go"
	TypeScript Language = "typescript"
	JavaScript Language = "javascript"
	Python     Language = "python"
	Rust       Language = "rust"
	Swift      Language = "swift"
)

// Detection holds information about a detected language.
type Detection struct {
	Language Language
	Path     string   // Directory where detected
	Files    []string // Indicator files found
}

// Detect scans a directory and returns all detected languages.
func Detect(dir string) ([]Detection, error) {
	var detections []Detection

	// Walk the directory looking for language indicators
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories and common non-source directories
		if d.IsDir() {
			name := d.Name()
			if name[0] == '.' || name == "node_modules" || name == "vendor" || name == "__pycache__" {
				return filepath.SkipDir
			}
			return nil
		}

		relDir := filepath.Dir(path)
		if relDir == "." {
			relDir = dir
		}

		// Check for language indicators
		switch d.Name() {
		case "go.mod":
			detections = appendIfNew(detections, Detection{
				Language: Go,
				Path:     relDir,
				Files:    []string{path},
			})
		case "package.json":
			// Check if it's TypeScript or JavaScript
			lang := JavaScript
			tsConfig := filepath.Join(relDir, "tsconfig.json")
			if _, err := os.Stat(tsConfig); err == nil {
				lang = TypeScript
			}
			detections = appendIfNew(detections, Detection{
				Language: lang,
				Path:     relDir,
				Files:    []string{path},
			})
		case "Cargo.toml":
			detections = appendIfNew(detections, Detection{
				Language: Rust,
				Path:     relDir,
				Files:    []string{path},
			})
		case "Package.swift":
			detections = appendIfNew(detections, Detection{
				Language: Swift,
				Path:     relDir,
				Files:    []string{path},
			})
		case "pyproject.toml", "setup.py", "requirements.txt":
			detections = appendIfNew(detections, Detection{
				Language: Python,
				Path:     relDir,
				Files:    []string{path},
			})
		}

		return nil
	})

	return detections, err
}

// appendIfNew adds a detection if the path isn't already detected for that language.
func appendIfNew(detections []Detection, d Detection) []Detection {
	for i, existing := range detections {
		if existing.Language == d.Language && existing.Path == d.Path {
			// Merge files
			detections[i].Files = append(existing.Files, d.Files...)
			return detections
		}
	}
	return append(detections, d)
}

// HasLanguage checks if a specific language was detected.
func HasLanguage(detections []Detection, lang Language) bool {
	for _, d := range detections {
		if d.Language == lang {
			return true
		}
	}
	return false
}

// GetByLanguage returns all detections for a specific language.
func GetByLanguage(detections []Detection, lang Language) []Detection {
	var result []Detection
	for _, d := range detections {
		if d.Language == lang {
			result = append(result, d)
		}
	}
	return result
}
