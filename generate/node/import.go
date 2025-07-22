package node

import (
	"github.com/please-build/puku/kinds"
)

// File represents a single JavaScript/TypeScript file in a package
type File struct {
	// Name is the package name from package.json or inferred from directory
	Name, FileName string
	// Imports are the imports of this file
	Imports []Import
	// FileType indicates if this is a library, test, or binary file
	FileType FileType
	// HasDefault indicates if this file has a default export
	HasDefault bool
}

// Import represents an import statement in JavaScript/TypeScript
type Import struct {
	// Path is the import path (e.g., "./utils", "lodash", "fs")
	Path string
	// ImportedName is the imported identifier (e.g., "React", "useState")
	ImportedName string
	// Type classifies the import as relative, bare, builtin, or asset
	Type ImportType
	// IsDefault indicates if this is a default import
	IsDefault bool
}

// FileType represents the type of JavaScript/TypeScript file
type FileType int

const (
	Library FileType = iota
	Test
	Binary
)

// ImportType represents the type of import
type ImportType int

const (
	RelativeImport ImportType = iota // ./foo, ../bar
	BareImport                       // lodash, @types/node
	BuiltinImport                    // fs, path (Node.js builtins)
	AssetImport                      // .css, .json files
)

// IsTest returns whether the Node.js file is a test
func (f *File) IsTest() bool {
	return f.FileType == Test
}

// IsCmd returns whether the Node.js file is a binary/command
func (f *File) IsCmd() bool {
	return f.FileType == Binary
}

// KindType returns the kinds.Type for this Node.js file
func (f *File) KindType() kinds.Type {
	if f.IsTest() {
		return kinds.Test
	}
	if f.IsCmd() {
		return kinds.Bin
	}
	return kinds.Lib
}