# Claude Development Process

This document defines the established development cycle and systematic approach for working on the Puku codebase.

## Development Cycle Overview

### 1. Work Initiation
**Before starting any work:**
- Review `PROGRESS-TRACKER.md` for current status and context
- Check out correct branch and sync with master if needed
- Run `plz test //...` to ensure clean starting state
- Update todo list with current session goals
- Review relevant sections in `CODEBASE-KNOWLEDGE.md`

### 2. Requirements Analysis
**For new features or changes:**
- Understand the problem space and user requirements
- Study existing similar functionality in the codebase
- Identify integration points and affected components
- Plan implementation approach using existing patterns
- Document initial approach in `PROGRESS-TRACKER.md`

### 3. Code Study Phase
**Before implementing anything:**
- Read existing code in the relevant packages
- Understand current patterns and conventions
- Study test files to understand expected behavior
- Identify extension points and integration patterns
- Document insights in `CODEBASE-KNOWLEDGE.md`

### 4. Implementation Planning
**Create detailed implementation plan:**
- Break work into logical, testable increments
- Identify test cases needed
- Plan commit structure for reviewable history
- Consider backward compatibility and migration needs
- Update todo list with specific implementation tasks

### 5. Incremental Development
**Development loop (repeat for each logical increment):**
- Implement smallest possible working change
- Write/update tests for the change
- Run affected tests: `plz test //package:all`
- Validate with full test suite if significant: `plz test //...`
- Commit with clear, descriptive message following conventional commit format (feat:, fix:, chore:, docs:, test:, refactor:)
- Update progress tracker with status

### 6. Integration Validation
**Before considering work complete:**
- Run full test suite: `plz test //...`
- Build all targets: `plz build //...`
- Test with representative use cases
- Validate no existing functionality broken
- Check integration with Please build system

### 7. Documentation and Knowledge Capture
**Finalize work session:**
- Update `CODEBASE-KNOWLEDGE.md` with new insights
- Document any architectural decisions made
- Update `PROGRESS-TRACKER.md` with current status
- Ensure commit messages tell a clear story
- Push branch to preserve work across sessions

## Code Study Methodology

### Understanding Existing Patterns
1. **Package Structure Analysis**
   - Read package-level documentation and comments
   - Understand public vs private API boundaries
   - Identify main types and their relationships
   - Map dependencies and integration points

2. **Implementation Pattern Discovery**
   - Study similar functionality in other packages
   - Understand error handling patterns
   - Identify configuration mechanisms used
   - Note testing approaches and conventions

3. **Please Build System Integration**
   - Review BUILD files for dependency patterns
   - Understand rule types and visibility declarations
   - Check build and test target organization
   - Validate proper Please conventions followed

### Test-Driven Understanding
- Read test files before implementation files
- Understand expected behavior from test cases
- Identify edge cases and error conditions tested
- Use tests to understand API usage patterns
- Follow table-driven test patterns where established

## Change Validation Process

### Pre-Commit Validation
**Every commit must pass:**
- [ ] Code compiles: `plz build //affected/packages:all`
- [ ] Tests pass: `plz test //affected/packages:all`
- [ ] Full build works: `plz build //...` (for significant changes)
- [ ] Full test suite passes: `plz test //...` (for significant changes)
- [ ] Code follows existing patterns and conventions
- [ ] Error handling is comprehensive and consistent

### Post-Implementation Review
**After completing a feature:**
- Self-review all changes for quality and consistency
- Verify commit history tells a clear story
- Ensure documentation is updated appropriately
- Validate no unintended side effects introduced
- Check that Please build system integration is correct

## Testing Strategy

### Test Categories
1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test package interactions
3. **E2E Tests**: Test complete workflows with test repositories
4. **Build System Tests**: Validate Please build integration

### Test Development Approach
- Write tests alongside implementation, not after
- Follow existing test patterns and conventions
- Use table-driven tests for comprehensive coverage
- Mock external dependencies appropriately
- Ensure tests are deterministic and isolated

### Test Execution Strategy
- Run affected package tests during development
- Run full test suite before major commits
- Use E2E tests to validate end-to-end workflows
- Test both success and failure scenarios

## Architectural Decision Process

### Making Design Decisions
1. **Research Phase**
   - Study existing similar implementations
   - Understand current architectural patterns
   - Consider integration with Please build system
   - Research Go best practices for the domain

2. **Decision Documentation**
   - Document decision rationale in `CODEBASE-KNOWLEDGE.md`
   - Include alternatives considered and why rejected
   - Note any trade-offs or limitations introduced
   - Reference relevant commit or PR for context

3. **Validation Process**
   - Implement incrementally to validate approach
   - Test with realistic scenarios and edge cases
   - Get feedback through code review if possible
   - Be prepared to iterate based on learning

## Consistency Maintenance

### Code Style Consistency
- Follow existing formatting and naming conventions
- Use similar error handling patterns
- Maintain consistent logging and debugging approaches
- Follow established configuration patterns

### Architectural Consistency  
- Use existing data structures and abstractions
- Follow established package organization patterns
- Maintain consistent API design approaches
- Preserve existing integration patterns

### Build System Consistency
- Follow Please build rule naming conventions
- Maintain proper dependency declarations
- Use appropriate visibility settings
- Follow package organization patterns

## Progress Tracking Standards

### Session Management
- Start each session by reviewing progress tracker
- Update progress tracker as work proceeds
- End sessions with clear status documentation
- Note any blockers or decisions needed

### Cross-Session Continuity
- Document context needed to resume work
- Reference specific commits and branches
- Note any environmental setup requirements
- Capture any insights or discoveries made

### Knowledge Building
- Continuously update codebase knowledge document
- Document patterns and conventions discovered
- Note any anti-patterns or gotchas encountered
- Build searchable institutional knowledge

## Quality Assurance Integration

### Continuous Quality Checks
- Code review all changes before committing
- Run tests frequently during development
- Validate integration points regularly
- Check for consistent error handling

### Pre-Completion Validation
- Full test suite execution
- Build system validation
- Integration testing with realistic scenarios
- Documentation and knowledge capture review

---

**Process Version**: 1.0  
**Last Updated**: 2025-07-21 - Initial development process establishment  
**Next Review**: After first feature implementation to refine based on experience