# JavaScript/TypeScript Support Design Document

## Executive Summary

This document outlines the comprehensive design for extending Puku to support JavaScript and TypeScript projects with equivalent functionality to the existing Go implementation. The design leverages Puku's existing language-agnostic architecture to add JS/TS support through strategic extension points while maintaining full backward compatibility.

## Project Context

### Current State
Puku is a BUILD file maintenance tool for Go projects using the Please build system. It automatically generates and updates BUILD files, handles dependencies, and supports custom rule types through hierarchical configuration.

### Objective
Extend Puku to provide equivalent JavaScript/TypeScript support including:
- File analysis and import parsing
- Dependency resolution via npm registry
- BUILD rule generation for JS/TS projects
- Third-party package management
- Watch mode integration
- Configuration system extension

### Success Criteria
- Feature parity with existing Go functionality
- Performance comparable to Go implementation
- Zero impact on existing Go features
- Seamless integration with existing CLI and configuration

## Architectural Analysis

### Current Puku Architecture

**Core Components:**
- **CLI Layer** (`cmd/puku/puku.go`): Command routing and execution
- **Configuration** (`config/config.go`): Hierarchical puku.json system
- **File Analysis** (`generate/import.go`): Language-specific parsing
- **Dependency Resolution** (`generate/deps.go`): Third-party package handling
- **BUILD Management** (`edit/edit.go`): Rule creation and modification
- **Rule Types** (`kinds/kinds.go`): Build rule type system
- **Watch System** (`watch/watch.go`): File system monitoring

**Key Architectural Strengths:**
1. **Language-agnostic design**: Core logic abstracted from Go-specific details
2. **Extensible configuration**: Custom rule types via JSON configuration
3. **Pluggable dependency resolution**: Abstract `Proxy` interface
4. **Modular file processing**: Clear separation between parsing and rule generation

### Extension Strategy

**Approach: Extend Rather Than Replace**
- Leverage existing architectural patterns
- Add JS/TS support alongside Go functionality
- Reuse infrastructure (BUILD editing, watching, configuration)
- Maintain existing command-line interface

## Technical Design

### 1. File Analysis System

#### JavaScript/TypeScript Parser Integration

**Primary Parser: Tree-sitter**
```go
// New file: generate/js_import.go
type JSFile struct {
    PackageName string
    Imports     []Import
    FileType    FileType
    ModuleType  ModuleType // CommonJS, ES6, TypeScript
}

type Import struct {
    Path         string
    Type         ImportType // relative, bare, builtin
    ImportedName string
    IsTypeOnly   bool // TypeScript type imports
}

func ImportJSDir(dirPath string) ([]JSFile, error) {
    // Use tree-sitter for fast parsing
    // Detect .js, .ts, .jsx, .tsx files
    // Extract imports and classify file types
}
```

**Secondary Parser: TypeScript Compiler API**
- Used for complex module resolution
- Handles TypeScript-specific features
- Provides authoritative dependency resolution

#### File Type Classification

**File Types:**
- **Library files**: Regular source files → `js_library`, `ts_library`
- **Test files**: Files matching test patterns → `js_test`, `jest_test`
- **Binary files**: Entry point files → `js_binary`, `ts_binary`
- **Asset files**: CSS, JSON, etc. → appropriate asset rules

**Detection Patterns:**
```go
func (f *JSFile) IsTest() bool {
    // Check for test patterns: *.test.js, *.spec.ts, __tests__/*, etc.
}

func (f *JSFile) IsCmd() bool {
    // Check for binary patterns: bin/ directory, main field in package.json
}
```

### 2. Dependency Resolution

#### Local Dependency Resolution

**Relative Import Handling:**
```go
func resolveLocalImport(importPath string, currentDir string) (string, error) {
    // Resolve ./foo, ../bar imports to build targets
    // Handle index.js implicit resolution
    // Support TypeScript path mapping
}
```

**Module Discovery:**
- Scan project for package.json files
- Build dependency graph of local packages
- Generate targets for unbuildable packages

#### Third-Party Dependency Resolution

**NPM Registry Integration:**
```go
// New file: generate/npm_proxy.go
type NPMProxy struct {
    registryURL string
    cache       map[string]*PackageInfo
}

func (p *NPMProxy) ResolvePackage(name, version string) (*PackageInfo, error) {
    // Query npm registry for package information
    // Handle version range resolution
    // Cache results for performance
}
```

**Package.json Integration:**
```go
func parsePackageJSON(path string) (*PackageConfig, error) {
    // Extract dependencies, devDependencies
    // Handle workspace configurations
    // Support package.json exports field
}
```

### 3. Build Rule Generation

#### Rule Type Mapping

**Standard Rules:**
```json
{
  "libKinds": {
    "js_library": {
      "srcsArg": "srcs",
      "providedDeps": ["//third_party/npm:node_modules"]
    },
    "ts_library": {
      "srcsArg": "srcs",
      "providedDeps": ["//third_party/npm:typescript"]
    }
  },
  "testKinds": {
    "js_test": {},
    "jest_test": {
      "providedDeps": ["//third_party/npm:jest"]
    }
  },
  "binKinds": {
    "js_binary": {},
    "ts_binary": {}
  }
}
```

**Rule Generation Logic:**
```go
// Extend edit/edit.go
func NewJSLibraryRule(name string, srcs []string) *build.Rule {
    // Create js_library rule with appropriate attributes
    // Handle TypeScript-specific configuration
    // Set visibility and dependencies
}

func NewNpmPackageRule(pkg *PackageInfo) *build.Rule {
    // Create npm_package rule for third-party dependencies
    // Handle version constraints and peer dependencies
}
```

### 4. Configuration Extension

#### Hierarchical Configuration

**Extended puku.json Schema:**
```json
{
  "base": "../puku.json",
  "jsConfig": {
    "packageManager": "npm|yarn|pnpm",
    "testPatterns": ["*.test.js", "*.spec.ts", "__tests__/*"],
    "assetPatterns": ["*.css", "*.json"],
    "thirdPartyDir": "third_party/npm"
  },
  "tsConfig": {
    "configFile": "tsconfig.json",
    "pathMapping": true,
    "typeOnlyImports": true
  }
}
```

**Configuration Loading:**
```go
// Extend config/config.go
type Config struct {
    // ... existing fields
    JSConfig *JSConfig `json:"jsConfig,omitempty"`
    TSConfig *TSConfig `json:"tsConfig,omitempty"`
}

type JSConfig struct {
    PackageManager  string   `json:"packageManager"`
    TestPatterns    []string `json:"testPatterns"`
    AssetPatterns   []string `json:"assetPatterns"`
    ThirdPartyDir   string   `json:"thirdPartyDir"`
}
```

### 5. Third-Party Package Management

#### NPM Package Rules

**Rule Structure:**
```go
func NewNpmPackageRule(name, version string) *build.Rule {
    return &build.Rule{
        Kind: "npm_package",
        Name: strings.ReplaceAll(name, "/", "_"),
        AttrStrings: map[string]string{
            "package": name,
            "version": version,
        },
    }
}
```

**Synchronization Command:**
```go
// New command: puku sync-npm
func syncNPM(cfg *config.Config, pleaseConfig *please.Config, wd string) error {
    // Parse package.json dependencies
    // Generate/update npm_package rules
    // Handle lock file integration
    // Update BUILD file third-party dependencies
}
```

### 6. Watch Mode Integration

#### File System Monitoring

**Extended File Patterns:**
```go
// Extend watch/watch.go
var jsExtensions = []string{".js", ".ts", ".jsx", ".tsx"}
var configFiles = []string{"package.json", "tsconfig.json"}

func (w *Watcher) shouldWatch(path string) bool {
    // Existing Go file logic
    // Add JS/TS file extensions
    // Include configuration files
}
```

**Change Detection:**
```go
func (w *Watcher) handleJSChange(path string) {
    // Use tree-sitter for incremental parsing
    // Detect import changes
    // Trigger appropriate regeneration
}
```

## Implementation Phases

### Phase 1: Core Infrastructure (Weeks 1-2)

**Goals:**
- Basic JS/TS file analysis
- Local dependency resolution
- Simple rule generation

**Deliverables:**
- `generate/js_import.go` with tree-sitter integration
- Extended configuration schema
- Basic `js_library` rule generation
- Unit tests for core functionality

**Success Criteria:**
- Parse JS/TS files and extract imports
- Generate basic BUILD rules for simple projects
- Pass integration tests with sample projects

### Phase 2: Third-Party Integration (Weeks 3-4)

**Goals:**
- NPM registry integration
- Package.json synchronization
- Third-party rule generation

**Deliverables:**
- `generate/npm_proxy.go` with registry integration
- `sync-npm` command implementation
- `npm_package` rule generation
- E2E tests with real npm packages

**Success Criteria:**
- Resolve and generate rules for npm dependencies
- Synchronize package.json to BUILD files
- Handle common package manager configurations

### Phase 3: Advanced Features (Weeks 5-6)

**Goals:**
- TypeScript-specific features
- Watch mode integration
- Performance optimization

**Deliverables:**
- TypeScript compiler API integration
- Watch mode for JS/TS files
- Advanced caching strategies
- Comprehensive test coverage

**Success Criteria:**
- Handle complex TypeScript projects
- Watch mode performance comparable to Go
- Production-ready error handling

### Phase 4: Production Readiness (Weeks 7-8)

**Goals:**
- Performance optimization
- Edge case handling
- Documentation and examples

**Deliverables:**
- Performance benchmarks
- Comprehensive error handling
- User documentation
- Migration guides

**Success Criteria:**
- Performance targets met
- Robust error recovery
- Complete feature parity with Go implementation

## Performance Specifications

### Response Time Targets

**Initial Analysis:**
- Small projects (< 100 files): < 1 second
- Medium projects (< 1000 files): < 5 seconds
- Large projects (< 10000 files): < 30 seconds

**Incremental Updates:**
- Single file changes: < 200ms
- Configuration changes: < 1 second
- Dependency updates: < 5 seconds

**Memory Usage:**
- Base overhead: < 50MB
- Per-file overhead: < 1KB
- Cache size: Configurable, default 100MB

### Optimization Strategies

**Caching:**
- Parse result caching with file modification time
- Dependency resolution caching
- NPM registry response caching

**Incremental Processing:**
- Tree-sitter incremental parsing
- Change detection for minimal reprocessing
- Dependency graph differential updates

## Risk Assessment

### Technical Risks

**High Risk:**
- **Node.js module resolution complexity**: Mitigation through TypeScript compiler API
- **NPM registry reliability**: Mitigation through caching and fallback strategies
- **Performance with large projects**: Mitigation through incremental processing

**Medium Risk:**
- **TypeScript configuration complexity**: Mitigation through comprehensive testing
- **Package manager ecosystem differences**: Mitigation through abstraction layers

**Low Risk:**
- **BUILD rule generation**: Leverages existing proven patterns
- **Configuration system extension**: Well-defined extension points

### Compatibility Risks

**Backward Compatibility:**
- Zero risk to existing Go functionality
- Configuration files remain compatible
- CLI interface unchanged

**Forward Compatibility:**
- Design allows for additional language support
- Configuration schema versioning
- Extensible rule type system

## Success Metrics

### Functional Requirements
- Successfully analyze 95% of popular JS/TS project structures
- Generate correct BUILD rules for standard project layouts
- Resolve dependencies accurately for npm packages
- Integrate seamlessly with existing Puku commands

### Performance Requirements
- Initial analysis within target response times
- Incremental updates faster than full regeneration
- Memory usage scaling linearly with project size
- Cache hit rates > 80% for typical workflows

### Quality Requirements
- Test coverage > 90% for new code
- Zero regressions in existing functionality
- Clear error messages for common failure modes
- Comprehensive documentation and examples

## Conclusion

This design provides a comprehensive roadmap for extending Puku to support JavaScript and TypeScript while maintaining its architectural integrity and performance characteristics. The phased implementation approach ensures incremental delivery of value while managing technical complexity and risk.

The design leverages Puku's existing strengths in configuration management, BUILD file editing, and dependency resolution while adding language-specific capabilities through well-defined extension points. This approach ensures that JS/TS support feels like a natural extension of Puku's capabilities rather than a bolted-on addition.