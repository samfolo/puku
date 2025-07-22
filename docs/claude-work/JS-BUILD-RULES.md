# JavaScript Build Rules Reference

**CRITICAL**: These are the actual JavaScript build rules used in this project. This information was provided once and must not be lost.

## Core Library Rule

### js_library
```python
def js_library(name, src = None, srcs = None, deps = None, visibility = None, test_only = False, labels:list = []):
```
- **Purpose**: Creates a JavaScript library filegroup
- **Key behavior**: Uses filegroup with `output_is_complete = False`
- **Sources**: Either `src` (single) or `srcs` (multiple), mutually exclusive
- **Labels**: Always includes "js_library" label
- **Requirements**: Requires ["js", "light_js"]

## Testing Rules

### vitest (MODERN DEFAULT - preferred for new tests)
```python
def vitest(name:str, entry:str, workspace:str, config:str="//build/js/vitest:config", deps:list=None, labels:list=[], cover:bool=True):
```
- **Purpose**: Modern testing with vitest
- **Key features**:
  - Creates internal `js_library` for test entry: `_{name}#entry`
  - Uses `yarn_cmd` to run vitest command
  - Uses `gentest` to parse JUnit XML results
  - Coverage support via `--coverage` flag
  - **Limitation**: gentest doesn't actually run test, just collects results

### jest (LEGACY - older option)
```python
def jest(name:str, entry:str, workspace:str, script_name:str="build-test", build_environment_vars:dict={}, run_environment_vars:dict={}, args:str="", deps:list=None, data:list=[], target:str="", labels:list=[], cover:bool=True, bundler:str="webpack", setup_files:list=[]):
```
- **Purpose**: Jest testing with bundling
- **Key features**:
  - Creates internal `js_library`: `_{name}#entry`
  - Bundles test with webpack or vite: `_{name}#bundle`
  - Supports build/run environment variables
  - Coverage via instrumentation
  - Setup files support

### mocha (OLDEST - rare)
```python
def mocha(name:str, entry:str, workspace:str, script_name:str="build-test", build_environment_vars:dict={}, run_environment_vars:dict={}, args:str="", deps:list=None, data:list=[], target:str="", labels:list=[], cover:bool=True):
```
- **Purpose**: Mocha testing (legacy)
- **Key features**:
  - Creates `_{name}#entry` and `_{name}#bundle`
  - Uses webpack bundling
  - NYC for coverage
  - JUnit reporter for Please integration

## Third-party Dependencies

### yarn_module
```python
def yarn_module(name:str, url:str, package_name:str, version:str, out:str=None, hashes:list=None, test_only:bool=False, patches:list=None, visibility:list=None, deps:list=[], licences:list=[]):
```
- **Purpose**: Install third-party packages from Yarn registry
- **Key features**:
  - Downloads from URL with hash verification
  - Creates filegroup with dependency metadata
  - Supports patches and licensing
  - Used with yarn_deps tool for BUILD file generation

## Rule Pattern Analysis

### Common Patterns:
1. **Internal rules**: Most test rules create `_{name}#entry` js_library rules
2. **Bundling**: Test rules often bundle sources: `_{name}#bundle`
3. **Workspace dependency**: All rules require a `workspace` parameter
4. **Coverage**: All test rules support optional coverage (`cover:bool=True`)
5. **gentest wrapper**: Test rules use gentest for Please integration

### Key Attributes:
- **entry**: Single test file entry point
- **workspace**: Required yarn_workspace target
- **config**: Configuration file dependencies
- **deps**: Additional dependencies
- **labels**: Rule categorization
- **test_only**: Marks rules as test-only

### Build System Integration:
- Uses yarn commands and workspaces
- JUnit XML for test results
- Coverage JSON for coverage reports
- File bungling for test execution
- Environment variable passing

## Rule Hierarchy for Puku:

1. **js_library** - Core library rule (maps to kinds.Lib)
2. **vitest/jest/mocha** - Test rules (maps to kinds.Test)  
3. **yarn_module** - Third-party deps (maps to kinds.ThirdParty)
4. **No explicit binary rule** - Would need to check for executable patterns

This is the definitive reference for JavaScript build rules in this project.