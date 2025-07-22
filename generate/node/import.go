package node

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/typescript/tsx"
	"github.com/smacker/go-tree-sitter/typescript/typescript"

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

// ImportDir parses all JavaScript/TypeScript files in a directory
func ImportDir(dir string) (map[string]*File, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]*File, len(files))
	for _, info := range files {
		if !info.Type().IsRegular() {
			continue
		}

		if !isJavaScriptFile(info.Name()) {
			continue
		}

		f, err := parseFile(dir, info.Name())
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", info.Name(), err)
		}
		ret[info.Name()] = f
	}

	return ret, nil
}

// isJavaScriptFile returns true if the file is a JavaScript or TypeScript file
func isJavaScriptFile(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".js" || ext == ".jsx" || ext == ".ts" || ext == ".tsx" || ext == ".mjs" || ext == ".cjs"
}

// parseFile parses a single JavaScript/TypeScript file and extracts imports
func parseFile(dir, filename string) (*File, error) {
	filePath := filepath.Join(dir, filename)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	parser := sitter.NewParser()
	
	// Choose appropriate grammar based on file extension
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".ts":
		parser.SetLanguage(typescript.GetLanguage())
	case ".tsx":
		parser.SetLanguage(tsx.GetLanguage())
	case ".js", ".jsx", ".mjs", ".cjs":
		parser.SetLanguage(javascript.GetLanguage())
	default:
		parser.SetLanguage(javascript.GetLanguage()) // fallback
	}

	tree, err := parser.ParseCtx(context.Background(), nil, content)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", ext, err)
	}
	defer tree.Close()

	file := &File{
		Name:       inferPackageName(dir),
		FileName:   filename,
		FileType:   classifyFileType(filename, content),
		Imports:    []Import{},
		HasDefault: false,
	}

	// Extract imports and exports from the syntax tree
	rootNode := tree.RootNode()
	extractImportsAndExports(rootNode, content, file)

	return file, nil
}

// inferPackageName tries to determine the package name from the directory
func inferPackageName(dir string) string {
	return filepath.Base(dir)
}

// classifyFileType determines if a file is a library, test, or binary
func classifyFileType(filename string, content []byte) FileType {
	base := strings.ToLower(filename)
	
	// Test files - check for common test patterns
	if strings.Contains(base, ".test.") || strings.Contains(base, ".spec.") ||
		strings.HasSuffix(base, "_test.js") || strings.HasSuffix(base, "_test.ts") ||
		strings.HasSuffix(base, ".test.js") || strings.HasSuffix(base, ".test.ts") ||
		strings.HasSuffix(base, ".spec.js") || strings.HasSuffix(base, ".spec.ts") {
		return Test
	}

	// Binary files - only if file starts with shebang
	if len(content) > 2 && string(content[0:2]) == "#!" {
		return Binary
	}

	// Everything else is a library (including index.js, main.js, etc.)
	return Library
}

// extractImportsAndExports traverses the syntax tree to find import and export statements
func extractImportsAndExports(node *sitter.Node, source []byte, file *File) {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		
		switch child.Type() {
		case "import_statement":
			extractImportStatement(child, source, file)
		case "export_statement":
			checkForDefaultExport(child, file)
		default:
			// Recursively process child nodes
			extractImportsAndExports(child, source, file)
		}
	}
}

// extractImportStatement processes import statements and adds them to the file
func extractImportStatement(node *sitter.Node, source []byte, file *File) {
	var importPath string
	var importedNames []string
	var hasDefault bool

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		
		switch child.Type() {
		case "string":
			// Extract the import path
			pathStr := child.Content(source)
			importPath = strings.Trim(pathStr, `"'`)
		case "import_clause":
			importedNames, hasDefault = extractImportClause(child, source)
		}
	}

	if importPath != "" {
		importType := classifyImportType(importPath)
		
		if len(importedNames) == 0 {
			// Side-effect import like `import "./styles.css"`
			file.Imports = append(file.Imports, Import{
				Path:         importPath,
				ImportedName: "",
				Type:         importType,
				IsDefault:    false,
			})
		} else {
			// Named or default imports
			for _, name := range importedNames {
				file.Imports = append(file.Imports, Import{
					Path:         importPath,
					ImportedName: name,
					Type:         importType,
					IsDefault:    hasDefault && name == importedNames[0], // First import is default if any
				})
			}
		}
	}
}

// extractImportClause processes the import clause to get imported names
func extractImportClause(node *sitter.Node, source []byte) ([]string, bool) {
	var names []string
	var hasDefault bool

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		
		switch child.Type() {
		case "identifier":
			// Default import
			names = append(names, child.Content(source))
			hasDefault = true
		case "namespace_import":
			// import * as name
			for j := 0; j < int(child.ChildCount()); j++ {
				if grandchild := child.Child(j); grandchild.Type() == "identifier" {
					names = append(names, grandchild.Content(source))
					break
				}
			}
		case "named_imports":
			// import { a, b }
			extractNamedImports(child, source, &names)
		}
	}

	return names, hasDefault
}

// extractNamedImports processes named import specifiers
func extractNamedImports(node *sitter.Node, source []byte, names *[]string) {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "import_specifier" {
			// Get the imported name (could be aliased)
			for j := 0; j < int(child.ChildCount()); j++ {
				grandchild := child.Child(j)
				if grandchild.Type() == "identifier" {
					*names = append(*names, grandchild.Content(source))
					break
				}
			}
		}
	}
}

// checkForDefaultExport checks if the export statement is a default export
func checkForDefaultExport(node *sitter.Node, file *File) {
	nodeText := node.Content(nil)
	if strings.Contains(nodeText, "default") {
		file.HasDefault = true
	}
}

// classifyImportType determines the type of import based on the import path
func classifyImportType(importPath string) ImportType {
	// Node.js builtins
	if IsNodeBuiltin(importPath) {
		return BuiltinImport
	}

	// Relative imports - includes ./foo, ../bar, ., ..
	if strings.HasPrefix(importPath, "./") || strings.HasPrefix(importPath, "../") ||
		importPath == "." || importPath == ".." {
		return RelativeImport
	}


	// Everything else is a bare import (third-party packages, aliases, etc.)
	// Resolution will happen later during dependency resolution phase
	return BareImport
}

// ResolveRelativeImport resolves a relative import path to an absolute path
// within the project structure, handling ./foo, ../bar, and bare filenames
func ResolveRelativeImport(currentDir, importPath string) (string, error) {
	if importPath == "" {
		return "", fmt.Errorf("empty import path")
	}

	// Handle relative paths
	if strings.HasPrefix(importPath, "./") || strings.HasPrefix(importPath, "../") ||
		importPath == "." || importPath == ".." {
		
		absPath := filepath.Join(currentDir, importPath)
		cleanPath := filepath.Clean(absPath)
		
		// Check if the resolved path exists as a directory or file
		if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
			// Try with common JS/TS extensions
			for _, ext := range []string{".js", ".ts", ".jsx", ".tsx", ".mjs", ".cjs"} {
				if _, err := os.Stat(cleanPath + ext); err == nil {
					return cleanPath + ext, nil
				}
			}
			
			// Check for index files in directory
			indexPath := filepath.Join(cleanPath, "index")
			for _, ext := range []string{".js", ".ts", ".jsx", ".tsx", ".mjs", ".cjs"} {
				if _, err := os.Stat(indexPath + ext); err == nil {
					return indexPath + ext, nil
				}
			}
			
			return "", fmt.Errorf("could not resolve import %q from %q", importPath, currentDir)
		}
		
		return cleanPath, nil
	}
	
	// For bare imports (no ./ or ../), they should be handled by package resolution
	return "", fmt.Errorf("not a relative import: %q", importPath)
}

// FindPackageTarget looks for an existing js_library target in the given directory's BUILD file
// This mirrors the localDep functionality from the Go implementation
func FindPackageTarget(dir string) (string, error) {
	buildFile := filepath.Join(dir, "BUILD")
	
	// Check if BUILD file exists
	if _, err := os.Stat(buildFile); os.IsNotExist(err) {
		// No BUILD file means no existing target, but could be generated later
		return "", nil
	}
	
	// For now, we'll return a predictable target name
	// In the future, this should parse the BUILD file to find actual js_library targets
	packageName := filepath.Base(dir)
	return fmt.Sprintf("//%s", packageName), nil
}

// ResolveDependency resolves a single import to a build target
// This will be used during build file generation
func (f *File) ResolveDependency(imp Import, currentDir string) (string, error) {
	switch imp.Type {
	case BuiltinImport:
		// Node.js builtins don't need build targets
		return "", nil
		
	case RelativeImport:
		// Resolve relative path and find target
		resolvedPath, err := ResolveRelativeImport(currentDir, imp.Path)
		if err != nil {
			return "", fmt.Errorf("resolving relative import %q: %w", imp.Path, err)
		}
		
		// Get the directory containing the resolved file/package
		targetDir := resolvedPath
		if !strings.HasSuffix(resolvedPath, "/") {
			// If it's a file, get its directory
			targetDir = filepath.Dir(resolvedPath)
		}
		
		target, err := FindPackageTarget(targetDir)
		if err != nil {
			return "", fmt.Errorf("finding package target in %q: %w", targetDir, err)
		}
		
		return target, nil
		
	case BareImport:
		// Bare imports need third-party dependency resolution
		// This will be handled by the main dependency resolution system
		return "", fmt.Errorf("bare import resolution not implemented: %q", imp.Path)
		
	default:
		return "", fmt.Errorf("unknown import type for %q", imp.Path)
	}
}