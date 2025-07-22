package node

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/please-build/puku/kinds"
)

func TestFile_KindType(t *testing.T) {
	tests := []struct {
		name     string
		fileType FileType
		expected kinds.Type
	}{
		{
			name:     "library file maps to Lib kind",
			fileType: Library,
			expected: kinds.Lib,
		},
		{
			name:     "test file maps to Test kind",
			fileType: Test,
			expected: kinds.Test,
		},
		{
			name:     "binary file maps to Bin kind",
			fileType: Binary,
			expected: kinds.Bin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{FileType: tt.fileType}
			assert.Equal(t, tt.expected, f.KindType())
		})
	}
}

func TestResolveRelativeImport(t *testing.T) {
	tests := []struct {
		name        string
		currentDir  string
		importPath  string
		expectError bool
		expectedMsg string
	}{
		{
			name:       "current directory",
			currentDir: "/project/src",
			importPath: "./utils",
			// Will fail since file doesn't exist, but path should be correct
			expectError: true,
			expectedMsg: `could not resolve import "./utils" from "/project/src"`,
		},
		{
			name:       "parent directory",
			currentDir: "/project/src/components",
			importPath: "../utils",
			expectError: true,
			expectedMsg: `could not resolve import "../utils" from "/project/src/components"`,
		},
		{
			name:        "empty import path",
			currentDir:  "/project/src",
			importPath:  "",
			expectError: true,
			expectedMsg: "empty import path",
		},
		{
			name:        "bare import should fail",
			currentDir:  "/project/src",
			importPath:  "lodash",
			expectError: true,
			expectedMsg: `not a relative import: "lodash"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ResolveRelativeImport(tt.currentDir, tt.importPath)
			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClassifyImportType(t *testing.T) {
	tests := []struct {
		name       string
		importPath string
		expected   ImportType
	}{
		{
			name:       "relative import with ./",
			importPath: "./utils",
			expected:   RelativeImport,
		},
		{
			name:       "relative import with ../",
			importPath: "../components",
			expected:   RelativeImport,
		},
		{
			name:       "current directory",
			importPath: ".",
			expected:   RelativeImport,
		},
		{
			name:       "parent directory",
			importPath: "..",
			expected:   RelativeImport,
		},
		{
			name:       "node builtin fs",
			importPath: "fs",
			expected:   BuiltinImport,
		},
		{
			name:       "node builtin path",
			importPath: "path",
			expected:   BuiltinImport,
		},
		{
			name:       "css relative import",
			importPath: "./styles.css",
			expected:   RelativeImport,
		},
		{
			name:       "json relative import",
			importPath: "./config.json",
			expected:   RelativeImport,
		},
		{
			name:       "third party package",
			importPath: "lodash",
			expected:   BareImport,
		},
		{
			name:       "scoped package",
			importPath: "@types/node",
			expected:   BareImport,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifyImportType(tt.importPath)
			assert.Equal(t, tt.expected, result)
		})
	}
}