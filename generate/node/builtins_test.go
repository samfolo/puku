package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNodeBuiltin(t *testing.T) {
	tests := []struct {
		name       string
		modulePath string
		expected   bool
	}{
		// Standard Node.js builtins
		{name: "fs is builtin", modulePath: "fs", expected: true},
		{name: "path is builtin", modulePath: "path", expected: true},
		{name: "crypto is builtin", modulePath: "crypto", expected: true},
		
		// Node.js prefixed builtins
		{name: "node:fs is builtin", modulePath: "node:fs", expected: true},
		{name: "node:path is builtin", modulePath: "node:path", expected: true},
		{name: "node:crypto is builtin", modulePath: "node:crypto", expected: true},
		
		// Submodules
		{name: "fs/promises is builtin", modulePath: "fs/promises", expected: true},
		{name: "node:fs/promises is builtin", modulePath: "node:fs/promises", expected: true},
		{name: "path/posix is builtin", modulePath: "path/posix", expected: true},
		{name: "util/types is builtin", modulePath: "util/types", expected: true},
		
		// Third-party packages (not builtins)
		{name: "lodash is not builtin", modulePath: "lodash", expected: false},
		{name: "react is not builtin", modulePath: "react", expected: false},
		{name: "express is not builtin", modulePath: "express", expected: false},
		{name: "@types/node is not builtin", modulePath: "@types/node", expected: false},
		{name: "@babel/core is not builtin", modulePath: "@babel/core", expected: false},
		
		// Edge cases
		{name: "empty string is not builtin", modulePath: "", expected: false},
		{name: "invalid node: prefix", modulePath: "node:invalid", expected: false},
		{name: "partial match fails", modulePath: "f", expected: false},
		{name: "case sensitive", modulePath: "FS", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNodeBuiltin(tt.modulePath)
			assert.Equal(t, tt.expected, result)
		})
	}
}