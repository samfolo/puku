# JavaScript/TypeScript Implementation Reference

## Architecture Extension Points

This document provides detailed technical guidance for implementing JavaScript/TypeScript support in Puku, mapping specific extension points in the existing codebase to the new functionality.

## File Analysis Implementation

### Current Go Implementation Patterns

**File: `generate/import.go`**
```go
// Current Go file parsing pattern (lines 20-35)
type GoFile struct {
    PackageName string
    Imports     []string
    HasMain     bool
}

func ImportDir(dir string) ([]GoFile, error) {
    // Uses go/parser for AST analysis
    // Extracts package name and imports
    // Classifies file types (test, binary, library)
}
```

### JS/TS Extension Pattern

**New File: `generate/js_import.go`**
```go
type JSFile struct {
    PackageName string      // From package.json or inferred
    Imports     []JSImport  // Enhanced import structure
    FileType    FileType    // test, binary, library
    ModuleType  ModuleType  // CommonJS, ES6, TypeScript
    HasDefault  bool        // Has default export
}

type JSImport struct {
    Path         string     // Import path
    ImportedName string     // Imported identifier
    Type         ImportType // relative, bare, builtin
    IsTypeOnly   bool       // TypeScript type-only imports
    IsDynamic    bool       // Dynamic import()
}

type ImportType int
const (
    RelativeImport ImportType = iota // ./foo, ../bar
    BareImport                      // lodash, @types/node
    BuiltinImport                   // fs, path (Node.js builtins)
    AssetImport                     // .css, .json files
)
```

**Tree-sitter Integration:**
```go
func parseJSFile(filePath string) (*JSFile, error) {
    // Use tree-sitter for fast, reliable parsing
    // Language detection: JavaScript vs TypeScript
    // Extract import/require statements
    // Classify module system (CommonJS vs ES6)
}

func extractImports(node *sitter.Node, source []byte) []JSImport {
    // Parse import statements: import, require, dynamic import
    // Handle TypeScript-specific syntax: type imports
    // Extract export information for dependency resolution
}
```

## Dependency Resolution Extension

### Current Go Resolution Pattern

**File: `generate/deps.go`**
```go
// Current dependency resolution (lines 22-50)
func resolveImport(pkg string, knownTargets map[string]string) (string, error) {
    // 1. Check known targets mapping
    // 2. Look for local packages
    // 3. Query module proxy for third-party
    // 4. Generate go_repo rule if needed
}
```

### JS/TS Resolution Extension

**New File: `generate/js_deps.go`**
```go
type JSResolver struct {
    config       *config.JSConfig
    npmProxy     *NPMProxy
    knownTargets map[string]string
    packageCache map[string]*PackageInfo
}

func (r *JSResolver) ResolveImport(imp JSImport, currentDir string) (string, error) {
    switch imp.Type {
    case RelativeImport:
        return r.resolveRelativeImport(imp.Path, currentDir)
    case BareImport:
        return r.resolveBareImport(imp.Path)
    case BuiltinImport:
        return "", nil // No build target needed for Node.js builtins
    case AssetImport:
        return r.resolveAssetImport(imp.Path, currentDir)
    }
}

func (r *JSResolver) resolveRelativeImport(path, currentDir string) (string, error) {
    // Resolve ./foo, ../bar to build targets
    // Handle index.js implicit resolution
    // Support TypeScript path mapping from tsconfig.json
    
    resolved := filepath.Join(currentDir, path)
    
    // Check for index file resolution
    if !hasFileExtension(path) {
        if exists(resolved + "/index.js") {
            resolved = resolved + "/index.js"
        } else if exists(resolved + "/index.ts") {
            resolved = resolved + "/index.ts"
        }
    }
    
    return buildTargetFromPath(resolved), nil
}

func (r *JSResolver) resolveBareImport(packageName string) (string, error) {
    // Check known targets first
    if target, exists := r.knownTargets[packageName]; exists {
        return target, nil
    }
    
    // Query npm registry
    pkg, err := r.npmProxy.GetPackage(packageName)
    if err != nil {
        return "", err
    }
    
    // Generate npm_package rule
    return r.generateNpmPackageRule(pkg)
}
```

**NPM Registry Integration:**
```go
type NPMProxy struct {
    registryURL string
    client      *http.Client
    cache       *PackageCache
}

type PackageInfo struct {
    Name         string            `json:"name"`
    Version      string            `json:"version"`
    Dependencies map[string]string `json:"dependencies"`
    PeerDeps     map[string]string `json:"peerDependencies"`
    Exports      interface{}       `json:"exports"`
}

func (p *NPMProxy) GetPackage(name string) (*PackageInfo, error) {
    // Query npm registry API
    // Handle scoped packages (@types/node)
    // Cache responses for performance
    // Support version range resolution
}
```

## BUILD Rule Generation

### Current Rule Generation Pattern

**File: `edit/edit.go`**
```go
// Current Go rule generation (lines 94-140)
func NewGoRepoRule(name, importPath, sum string) *build.Rule {
    return &build.Rule{
        Kind: "go_repo",
        Name: name,
        AttrStrings: map[string]string{
            "importpath": importPath,
            "sum":        sum,
        },
    }
}
```

### JS/TS Rule Generation Extension

**Extend: `edit/edit.go`**
```go
func NewJSLibraryRule(name string, srcs []string, deps []string) *build.Rule {
    rule := &build.Rule{
        Kind: "js_library",
        Name: name,
    }
    
    // Set sources
    if len(srcs) > 0 {
        rule.AttrStrings = map[string]string{"srcs": formatStringList(srcs)}
    }
    
    // Set dependencies
    if len(deps) > 0 {
        rule.AttrStrings["deps"] = formatStringList(deps)
    }
    
    return rule
}

func NewTSLibraryRule(name string, srcs []string, deps []string, tsconfig string) *build.Rule {
    rule := NewJSLibraryRule(name, srcs, deps)
    rule.Kind = "ts_library"
    
    if tsconfig != "" {
        rule.AttrStrings["tsconfig"] = tsconfig
    }
    
    return rule
}

func NewNpmPackageRule(pkg *PackageInfo) *build.Rule {
    // Sanitize package name for Please rule names
    ruleName := strings.ReplaceAll(pkg.Name, "/", "_")
    ruleName = strings.ReplaceAll(ruleName, "@", "")
    
    rule := &build.Rule{
        Kind: "npm_package",
        Name: ruleName,
        AttrStrings: map[string]string{
            "package": pkg.Name,
            "version": pkg.Version,
        },
    }
    
    // Handle peer dependencies
    if len(pkg.PeerDeps) > 0 {
        var peerDeps []string
        for dep := range pkg.PeerDeps {
            peerDeps = append(peerDeps, fmt.Sprintf("//third_party/npm:%s", sanitizeName(dep)))
        }
        rule.AttrStrings["peer_deps"] = formatStringList(peerDeps)
    }
    
    return rule
}
```

## Configuration System Extension

### Current Configuration Pattern

**File: `config/config.go`**
```go
// Current configuration structure (lines 12-45)
type Config struct {
    Base        string              `json:"base,omitempty"`
    LibKinds    map[string]*Kind    `json:"libKinds,omitempty"`
    TestKinds   map[string]*Kind    `json:"testKinds,omitempty"`
    BinKinds    map[string]*Kind    `json:"binKinds,omitempty"`
    KnownTargets map[string]string  `json:"knownTargets,omitempty"`
}
```

### JS/TS Configuration Extension

**Extend: `config/config.go`**
```go
type Config struct {
    // ... existing fields
    JSConfig *JSConfig `json:"jsConfig,omitempty"`
    TSConfig *TSConfig `json:"tsConfig,omitempty"`
}

type JSConfig struct {
    PackageManager    string   `json:"packageManager"`     // npm, yarn, pnpm
    TestPatterns      []string `json:"testPatterns"`       // ["*.test.js", "*.spec.ts"]
    AssetPatterns     []string `json:"assetPatterns"`      // ["*.css", "*.json"]
    ThirdPartyDir     string   `json:"thirdPartyDir"`      // "third_party/npm"
    NodeModulesDir    string   `json:"nodeModulesDir"`     // "node_modules"
    IgnorePatterns    []string `json:"ignorePatterns"`     // ["node_modules/**"]
}

type TSConfig struct {
    ConfigFile        string `json:"configFile"`          // "tsconfig.json"
    PathMapping       bool   `json:"pathMapping"`         // Support path mapping
    TypeOnlyImports   bool   `json:"typeOnlyImports"`     // Handle type-only imports
    DeclarationFiles  bool   `json:"declarationFiles"`    // Include .d.ts files
}

// Default configurations
func DefaultJSConfig() *JSConfig {
    return &JSConfig{
        PackageManager: "npm",
        TestPatterns:   []string{"*.test.js", "*.spec.js", "__tests__/**"},
        AssetPatterns:  []string{"*.css", "*.json", "*.md"},
        ThirdPartyDir:  "third_party/npm",
        NodeModulesDir: "node_modules",
        IgnorePatterns: []string{"node_modules/**", "dist/**", "build/**"},
    }
}
```

**Configuration Loading:**
```go
func loadConfig(dir string) (*Config, error) {
    // Existing config loading logic
    cfg, err := loadExistingConfig(dir)
    if err != nil {
        return nil, err
    }
    
    // Apply JS/TS defaults if not specified
    if cfg.JSConfig == nil {
        cfg.JSConfig = DefaultJSConfig()
    }
    if cfg.TSConfig == nil {
        cfg.TSConfig = DefaultTSConfig()
    }
    
    return cfg, nil
}
```

## Command-Line Integration

### Current CLI Pattern

**File: `cmd/puku/puku.go`**
```go
// Current command mapping (lines 79-149)
var funcs = map[string]func(*config.Config, *please.Config, string) error{
    "fmt":     generate.Update,
    "sync":    sync.Sync,
    "watch":   watch.Watch,
    "migrate": migrate.Migrate,
    "lint":    lint.Lint,
}
```

### JS/TS CLI Extension

**Extend: `cmd/puku/puku.go`**
```go
var funcs = map[string]func(*config.Config, *please.Config, string) error{
    // ... existing commands
    "sync-npm":    sync.SyncNPM,     // Sync package.json to BUILD files
    "fmt-js":      generate.UpdateJS, // Format JS/TS BUILD files only
    "migrate-js":  migrate.MigrateJS, // Convert webpack configs to BUILD
}

// Or extend existing commands to handle JS/TS
func enhancedUpdate(cfg *config.Config, pleaseConfig *please.Config, wd string) error {
    // Run existing Go update
    if err := generate.Update(cfg, pleaseConfig, wd); err != nil {
        return err
    }
    
    // Run JS/TS update if configured
    if cfg.JSConfig != nil || cfg.TSConfig != nil {
        return generate.UpdateJS(cfg, pleaseConfig, wd)
    }
    
    return nil
}
```

**New Commands Implementation:**
```go
// New file: sync/sync_npm.go
func SyncNPM(cfg *config.Config, pleaseConfig *please.Config, wd string) error {
    // Find package.json files
    packageFiles, err := findPackageJSONFiles(wd)
    if err != nil {
        return err
    }
    
    for _, pkgFile := range packageFiles {
        if err := syncPackageJSON(pkgFile, cfg); err != nil {
            return fmt.Errorf("failed to sync %s: %w", pkgFile, err)
        }
    }
    
    return nil
}

func syncPackageJSON(pkgPath string, cfg *config.Config) error {
    // Parse package.json
    pkg, err := parsePackageJSON(pkgPath)
    if err != nil {
        return err
    }
    
    // Generate npm_package rules for dependencies
    buildFile := filepath.Join(filepath.Dir(pkgPath), "BUILD")
    
    for depName, version := range pkg.Dependencies {
        rule := NewNpmPackageRule(&PackageInfo{
            Name:    depName,
            Version: version,
        })
        
        if err := addRuleToBuildFile(buildFile, rule); err != nil {
            return err
        }
    }
    
    return nil
}
```

## Watch Mode Integration

### Current Watch Pattern

**File: `watch/watch.go`**
```go
// Current file watching (lines 70-100)
func (w *Watcher) shouldWatch(path string) bool {
    return strings.HasSuffix(path, ".go") && 
           !strings.Contains(path, "/.git/")
}
```

### JS/TS Watch Extension

**Extend: `watch/watch.go`**
```go
var jsExtensions = []string{".js", ".ts", ".jsx", ".tsx"}
var jsConfigFiles = []string{"package.json", "tsconfig.json", "webpack.config.js"}

func (w *Watcher) shouldWatch(path string) bool {
    // Existing Go file logic
    if strings.HasSuffix(path, ".go") && !strings.Contains(path, "/.git/") {
        return true
    }
    
    // JS/TS file extensions
    for _, ext := range jsExtensions {
        if strings.HasSuffix(path, ext) {
            return !w.shouldIgnore(path)
        }
    }
    
    // JS/TS configuration files
    basename := filepath.Base(path)
    for _, configFile := range jsConfigFiles {
        if basename == configFile {
            return true
        }
    }
    
    return false
}

func (w *Watcher) shouldIgnore(path string) bool {
    // Check against ignore patterns from JS config
    if w.config.JSConfig != nil {
        for _, pattern := range w.config.JSConfig.IgnorePatterns {
            if matched, _ := filepath.Match(pattern, path); matched {
                return true
            }
        }
    }
    return false
}

func (w *Watcher) handleJSChange(path string) error {
    // Use tree-sitter for fast incremental parsing
    // Detect import changes
    // Trigger appropriate BUILD file regeneration
    
    if isConfigFile(path) {
        // Full regeneration for config changes
        return w.regenerateAll()
    }
    
    // Incremental update for source file changes
    return w.updateAffectedRules(path)
}
```

## Testing Patterns

### Unit Test Structure

**Test File: `generate/js_import_test.go`**
```go
func TestParseJSFile(t *testing.T) {
    tests := []struct {
        name     string
        content  string
        expected *JSFile
    }{
        {
            name: "ES6 imports",
            content: `
                import React from 'react';
                import { useState } from 'react';
                import * as utils from './utils';
                
                export default function App() {}
            `,
            expected: &JSFile{
                Imports: []JSImport{
                    {Path: "react", ImportedName: "React", Type: BareImport},
                    {Path: "react", ImportedName: "useState", Type: BareImport},
                    {Path: "./utils", ImportedName: "*", Type: RelativeImport},
                },
                HasDefault: true,
                ModuleType: ES6Module,
            },
        },
        {
            name: "TypeScript type imports",
            content: `
                import type { User } from './types';
                import { api } from './api';
            `,
            expected: &JSFile{
                Imports: []JSImport{
                    {Path: "./types", ImportedName: "User", Type: RelativeImport, IsTypeOnly: true},
                    {Path: "./api", ImportedName: "api", Type: RelativeImport},
                },
                ModuleType: TypeScriptModule,
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := parseJSContent(tt.content)
            require.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Integration Test Pattern

**Test File: `e2e/tests/js_project_test.go`**
```go
func TestJSProjectGeneration(t *testing.T) {
    // Create temporary project structure
    projectDir := createTempJSProject(t)
    defer os.RemoveAll(projectDir)
    
    // Run puku fmt
    err := generate.UpdateJS(defaultConfig(), nil, projectDir)
    require.NoError(t, err)
    
    // Verify generated BUILD files
    buildContent := readFile(t, filepath.Join(projectDir, "BUILD"))
    assert.Contains(t, buildContent, "js_library(")
    assert.Contains(t, buildContent, `name = "lib"`)
    assert.Contains(t, buildContent, `srcs = ["index.js"]`)
}

func createTempJSProject(t *testing.T) string {
    dir := t.TempDir()
    
    // Create package.json
    writeFile(t, filepath.Join(dir, "package.json"), `{
        "name": "test-project",
        "dependencies": {
            "lodash": "^4.17.21"
        }
    }`)
    
    // Create source files
    writeFile(t, filepath.Join(dir, "index.js"), `
        import _ from 'lodash';
        export const utils = { map: _.map };
    `)
    
    return dir
}
```

## Performance Optimization Patterns

### Caching Strategy

```go
type ParseCache struct {
    mutex sync.RWMutex
    files map[string]*CacheEntry
}

type CacheEntry struct {
    ModTime time.Time
    Result  *JSFile
    Hash    string
}

func (c *ParseCache) Get(path string) (*JSFile, bool) {
    c.mutex.RLock()
    defer c.mutex.RUnlock()
    
    entry, exists := c.files[path]
    if !exists {
        return nil, false
    }
    
    // Check if file has been modified
    stat, err := os.Stat(path)
    if err != nil || !stat.ModTime().Equal(entry.ModTime) {
        delete(c.files, path)
        return nil, false
    }
    
    return entry.Result, true
}
```

### Incremental Processing

```go
func (r *JSResolver) UpdateIncrementally(changedFiles []string) error {
    // Only reprocess files that changed or depend on changed files
    affectedFiles := r.findAffectedFiles(changedFiles)
    
    for _, file := range affectedFiles {
        if err := r.processFile(file); err != nil {
            return err
        }
    }
    
    return nil
}

func (r *JSResolver) findAffectedFiles(changedFiles []string) []string {
    // Build reverse dependency graph
    // Find files that import from changed files
    // Return minimal set of files to reprocess
}
```

This implementation reference provides concrete patterns for extending Puku's architecture to support JavaScript and TypeScript while maintaining consistency with existing code patterns and performance characteristics.