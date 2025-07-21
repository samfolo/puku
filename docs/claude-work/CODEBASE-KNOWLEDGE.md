# Puku Codebase Knowledge

This document captures growing understanding of the Puku codebase, architectural insights, and development patterns.

## Project Architecture Overview

### Core Purpose
Puku is a BUILD file maintenance tool for Go projects using the Please build system. It automatically generates and updates BUILD files, manages dependencies, and supports custom build rule types.

### Key Components

#### generate/ - Core Generation Logic
- **generate.go**: Main orchestration logic for BUILD file generation
- **deps.go**: Dependency resolution and third-party package handling  
- **import.go**: Go import parsing and resolution to BUILD targets
- **Key Patterns**: Processes packages in dependency order, resolves imports through multiple fallback strategies

#### edit/ - BUILD File Manipulation
- **edit.go**: BUILD file parsing and modification using Please's buildtools
- **build_targets.go**: Target creation and management logic
- **rule.go**: Individual rule manipulation and formatting
- **Key Patterns**: Preserves existing formatting and comments, uses AST manipulation

#### config/ - Configuration Management
- **config.go**: Hierarchical puku.json configuration loading
- **Key Patterns**: Configuration inheritance from parent directories, supports custom rule kinds

#### work/ - Package Discovery
- **work.go**: File system traversal and Go package discovery
- **Key Patterns**: Respects .gitignore and puku.json stop directives

#### sync/ - Module Synchronization  
- **sync.go**: go.mod to BUILD file synchronization
- **Key Patterns**: Bridges standard Go tooling with Please build system

#### watch/ - File System Monitoring
- **watch.go**: File system watching for automatic BUILD updates
- **Key Patterns**: Uses fsnotify for efficient file monitoring

#### migrate/ - Rule Migration
- **migrate.go**: Migration from go_module to go_repo rules
- **Key Patterns**: Preserves existing behavior during migration

## Code Style and Conventions

### Go Style Patterns
- Standard Go conventions followed throughout
- Consistent error handling with wrapped errors
- Table-driven tests extensively used
- Context-aware operations where appropriate

### Please Build Integration
- All packages have BUILD files with explicit dependencies
- Uses go_library, go_test, go_binary rule types appropriately
- Proper visibility declarations (PUBLIC for exported packages)
- Third-party dependencies managed in third_party/go/

### Configuration Patterns
- Hierarchical configuration through puku.json files
- Support for custom rule kinds (libKinds, testKinds, binKinds)
- Sensible defaults with override capabilities
- Stop directive support for excluding directories

### Testing Patterns
- Comprehensive unit tests with table-driven approach
- E2E tests with dedicated test repositories
- Test fixtures in dedicated test_data/ directories
- Integration tests for go.mod synchronization

## Architectural Insights

### Dependency Resolution Strategy
1. Known imports (explicit configuration)
2. Installed packages (go_module/go_repo rules)
3. Module naming convention inference
4. Go module proxy/go.mod consultation
5. Automatic go_repo rule generation

### BUILD File Generation Process
1. **Discovery**: Find Go packages in directory tree
2. **Analysis**: Parse Go source files for imports and package types
3. **Resolution**: Resolve imports to BUILD targets
4. **Generation**: Create or update BUILD rules
5. **Formatting**: Preserve formatting and add necessary rules

### Configuration Inheritance
- puku.json files loaded hierarchically from root to leaf
- Later configurations override earlier ones
- Supports package-specific customizations
- Stop directive prevents processing of subdirectories

## Integration Points

### Please Build System
- Relies on Please's buildtools for AST manipulation
- Uses Please's query system for existing target discovery
- Integrates with Please's plugin system for Go rules
- Respects Please's visibility and naming conventions

### Go Module System
- Optional go.mod integration for robust dependency resolution
- Supports standard go get workflow for adding dependencies
- Syncs go.mod changes to BUILD files automatically
- Handles module versioning and replacement directives

### File System Integration
- Respects .gitignore patterns
- Monitors file system changes for automatic updates
- Handles symbolic links appropriately
- Supports glob patterns for selective processing

## Key Design Decisions

### Configuration Over Convention
- Extensive configurability through puku.json files
- Support for custom build rule types
- Override mechanisms for special cases
- Balanced with sensible defaults

### Incremental Processing
- Only processes changed packages when possible
- Preserves existing BUILD file formatting and comments
- Minimal disruption to existing configurations
- Support for partial updates (specific packages/modules)

### Error Handling Philosophy
- Comprehensive error messages with context
- Graceful degradation when possible
- Clear indication of what failed and why
- Recovery mechanisms for common issues

## Extension Patterns

### Adding New Rule Types
- Define kind configuration in puku.json
- Specify source allocation strategy (library/test/binary)
- Configure provided dependencies
- Set default visibility rules

### Custom Import Resolution
- Add entries to knownTargets configuration
- Implement custom resolution in import.go
- Support for non-standard package layouts
- Handle generated code and proto files

### New Command Implementation
- Follow existing CLI structure in cmd/puku/
- Use please build system for dependency management
- Implement proper error handling and logging
- Add appropriate tests and documentation

## Anti-Patterns and Gotchas

### What to Avoid
- Modifying BUILD files directly without using edit/ package
- Ignoring configuration hierarchy
- Hard-coding import paths without considering customization
- Breaking existing BUILD file formatting unnecessarily

### Common Pitfalls
- Not handling go.mod absent scenarios
- Assuming standard Go package layout
- Missing edge cases in import resolution
- Not preserving user customizations in BUILD files

### Testing Considerations
- E2E tests require careful setup of test repositories
- Mock external dependencies (module proxy, file system)
- Test configuration inheritance scenarios
- Validate BUILD file generation doesn't break builds

## Performance Considerations

### Optimization Strategies
- Caches parsed GO source files
- Batches file system operations
- Uses efficient data structures (tries for import matching)
- Minimizes BUILD file rewrites

### Scalability Patterns
- Processes packages in dependency order
- Supports incremental updates
- Handles large repositories efficiently
- Memory-conscious parsing and generation

## Development Workflow Integration

### Testing Strategy
- Unit tests: `plz test //package:all`
- Full test suite: `plz test //...`
- E2E tests: `plz test //e2e/...`
- Integration tests for specific scenarios

### Build Commands
- Build all: `plz build //...`
- Build main binary: `plz build //cmd/puku`
- Development builds for testing changes

### Development Tools
- Please wrapper: `./pleasew` for consistent environment
- Configuration validation through puku.json schema
- Debug logging available through standard Go logging

---

**Last Updated**: 2025-07-21 - Initial architecture analysis during setup
**Next Review**: After first feature implementation to validate understanding