# Claude Working Standards

This document establishes autonomous working guidelines for contributing to the Puku codebase.

## Git Workflow Standards

### Branch Naming Conventions
- `feature/description-of-work` - New features or enhancements
- `fix/issue-description` - Bug fixes
- `refactor/component-name` - Code refactoring
- `docs/topic` - Documentation updates
- `test/component-name` - Test improvements
- `setup/infrastructure-work` - Development infrastructure

### Commit Standards
- **Frequency**: Commit logical units of work - typically every 15-30 minutes of focused development
- **Message Format**: 
  ```
  type: brief description (50 chars max)
  
  Optional detailed explanation of what and why,
  not how. Reference issues/decisions as needed.
  ```
- **Types**: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`
- **Content**: Focus on WHY changes were made, not just what changed
- **History**: Each commit should represent a working state

### Development Branch Rules
- Never work directly on `master` branch
- Create feature branches from latest `master`
- Keep branches focused - one concern per branch
- Rebase onto master before creating PRs
- Delete merged branches to keep repository clean

## Code Quality Standards

### Before Writing Code
1. Study existing patterns in the relevant package
2. Check test files to understand expected behavior
3. Review similar implementations in the codebase
4. Understand the Please build system conventions
5. Identify integration points and dependencies

### Code Implementation Standards
- **Consistency**: Follow existing code style and patterns
- **Go Conventions**: Adhere to Go best practices and idioms
- **Please Patterns**: Follow Please build system conventions
- **Error Handling**: Implement comprehensive error handling
- **Edge Cases**: Consider and handle boundary conditions
- **Documentation**: Comment public interfaces and complex logic

### Testing Requirements
- Write tests for all new functionality
- Follow existing test patterns and conventions
- Use table-driven tests where appropriate
- Mock external dependencies appropriately
- Ensure tests are deterministic and isolated
- Run full test suite before committing: `plz test //...`

### Code Review Readiness
- Every change must be self-reviewed before commit
- Ensure code tells a clear story through commits
- Validate that each commit compiles and tests pass
- Check that changes don't break existing functionality
- Verify consistent error handling and logging

## Comment and Documentation Standards

### When to Comment
- Public interfaces and exported functions
- Complex algorithms or business logic
- Non-obvious design decisions
- Integration points with Please build system
- Configuration file formats and options

### Comment Style
```go
// PublicFunction performs X operation by doing Y.
// It returns Z when conditions A and B are met.
// 
// Example usage:
//   result, err := PublicFunction(input)
//   if err != nil { ... }
func PublicFunction(input string) (Result, error) {
```

### Documentation Updates
- Update CODEBASE-KNOWLEDGE.md with architectural insights
- Document new patterns or conventions discovered
- Update PROGRESS-TRACKER.md with work status
- Keep CLAUDE.md current with build/test commands

## Code Quality Principles

### Leave Code Better Than Found
- Fix minor issues encountered while working
- Improve variable names and comments when touching code
- Remove unused imports and dead code
- Ensure consistent formatting and style
- Update tests when modifying behavior

### Standards for New Code
- Prefer explicit over implicit
- Handle errors appropriately - don't ignore them
- Use meaningful variable and function names
- Keep functions focused and single-purpose
- Minimize dependencies and coupling
- Follow the existing architecture patterns

### Refactoring Guidelines
- Make incremental improvements
- Maintain backward compatibility unless explicitly changing APIs
- Ensure comprehensive test coverage before refactoring
- Document architectural changes in CODEBASE-KNOWLEDGE.md
- Consider impact on other packages and dependencies

## Error Handling Standards

### Error Handling Approach
- Return errors explicitly - don't panic unnecessarily  
- Wrap errors with context: `fmt.Errorf("operation failed: %w", err)`
- Log errors at appropriate levels
- Handle errors close to where they occur
- Provide meaningful error messages for users

### Edge Cases and Validation
- Validate input parameters
- Handle nil pointers and empty collections
- Consider concurrent access where applicable  
- Test boundary conditions
- Plan for partial failures and recovery

## Integration Standards

### Please Build System Integration
- Understand BUILD file structure and conventions
- Use appropriate build rule types (go_library, go_test, go_binary)
- Maintain proper dependency declarations
- Follow visibility rules and package organization
- Test changes with `plz build` and `plz test`

### Puku-Specific Considerations
- Understand configuration hierarchy (puku.json files)
- Respect existing kind mappings and patterns
- Consider impact on dependency resolution
- Test with various Go module configurations
- Validate BUILD file generation correctness

## Quality Assurance Checklist

Before every commit:
- [ ] Code compiles: `plz build //...`
- [ ] Tests pass: `plz test //...`
- [ ] Code follows existing patterns
- [ ] Error handling is appropriate
- [ ] Comments updated for public interfaces
- [ ] No obvious edge cases missed
- [ ] Progress documented in PROGRESS-TRACKER.md
- [ ] Architecture insights captured in CODEBASE-KNOWLEDGE.md

## Working Session Standards

### Session Start
1. Review PROGRESS-TRACKER.md for current status
2. Check out correct branch and sync with master
3. Run tests to ensure clean starting state
4. Update todo list with current session goals

### During Work
- Commit frequently with clear messages
- Update progress tracker as work proceeds
- Document discoveries in CODEBASE-KNOWLEDGE.md
- Test changes incrementally

### Session End  
- Commit all work with clear status
- Update PROGRESS-TRACKER.md with current state
- Document any blocking issues or decisions needed
- Push branch to preserve work across sessions