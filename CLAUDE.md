# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Puku is a tool for maintaining Go build rules, similar to Gazelle but for the Please build system. It automatically generates and updates BUILD files for Go packages, handles third-party dependencies, and supports custom build rule types through configuration.

## Build System Commands

This project uses the Please build system:

- `plz build //...` - Build all targets
- `plz test //...` - Run all tests  
- `plz build //cmd/puku` - Build the main puku binary
- `plz test //generate:all` - Run tests for the generate package
- `./pleasew` - Use the wrapper script if plz is not installed

## Testing

- Individual package tests: `plz test //package_name:all`
- All tests: `plz test //...`
- E2E tests are located in `e2e/tests/`
- Test files follow standard Go `*_test.go` naming convention

## Architecture

### Core Components

- **cmd/puku** - Main CLI entry point and command definitions
- **generate/** - Core logic for generating BUILD rules, parsing Go imports, and dependency resolution
- **edit/** - BUILD file parsing and modification logic using Please's buildtools
- **config/** - Configuration file parsing and management (puku.json files)
- **graph/** - Dependency graph construction and analysis
- **work/** - File system traversal and package discovery
- **sync/** - go.mod synchronization with BUILD files
- **migrate/** - Migration from go_module to go_repo rules
- **watch/** - File system watching for automatic updates

### Key Files

- **generate/generate.go** - Main generation logic that orchestrates the entire process
- **generate/deps.go** - Dependency resolution and third-party package handling
- **generate/import.go** - Go import parsing and resolution
- **edit/edit.go** - BUILD file modification and rule management
- **config/config.go** - Configuration system with hierarchical puku.json support

### Configuration System

Puku uses hierarchical `puku.json` files that can be placed at any directory level:
- Supports custom rule kinds (libKinds, testKinds, binKinds)
- Third-party dependency directory configuration
- Known target mappings for special cases
- Non-Go source handling for proto/generated files

### Dependency Resolution Flow

1. Parse Go source files for imports
2. Resolve imports against known targets and installed packages
3. For go_repo: check module proxy or go.mod file
4. Generate new go_repo rules as needed
5. Update BUILD files with resolved dependencies

## Development Workflow

The codebase follows Please build conventions with BUILD files in each package defining build targets. When adding new functionality:

1. Add source files to appropriate package
2. Update BUILD file if needed (or let puku handle it)
3. Add tests following Go conventions
4. Run tests with `plz test //package:all`

## File Structure Patterns

- Each Go package has its own directory with a BUILD file
- Test files are co-located with source files
- E2E tests have their own test repositories under `e2e/tests/`
- Third-party dependencies are managed in `third_party/go/`

## Claude Working Environment

**CRITICAL**: ALWAYS consult `docs/claude-work/WORKING-STANDARDS.md` before making ANY changes to understand:
- Conventional commit format (feat:, fix:, chore:, docs:, test:, refactor:)
- Branch naming conventions
- Code quality standards
- Git workflow requirements

For Claude Code instances working on this codebase, comprehensive working documentation is maintained in `docs/claude-work/`:

- **WORKING-STANDARDS.md** - Git workflow, code quality standards, and development guidelines (CONSULT FIRST)
- **DEVELOPMENT-PROCESS.md** - Systematic development cycle and quality assurance processes
- **PROGRESS-TRACKER.md** - Session continuity, current work status, and context for resuming work
- **CODEBASE-KNOWLEDGE.md** - Growing architectural insights, patterns, and development conventions

These files provide full context for autonomous development and should be consulted at the start of each working session and before any commits.

## JavaScript/TypeScript Extension Project

**Overarching Task**: Extend Puku to support JavaScript and TypeScript projects with equivalent functionality to the existing Go implementation.

**Comprehensive Design Documentation**: Located in `docs/js-ts-extension/`:

- **DESIGN.md** - Complete architectural analysis, technical design, and implementation strategy for adding JS/TS support
- **IMPLEMENTATION-REFERENCE.md** - Detailed technical patterns, extension points, and code examples for implementation

This represents a major architectural extension that will add multi-language support to Puku while maintaining full backward compatibility with existing Go functionality. The design leverages Puku's language-agnostic architecture and extensible configuration system to provide seamless JS/TS integration.