package node

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/please-build/puku/kinds"
)

func TestFile_KindType(t *testing.T) {
	tests := []struct {
		name     string
		fileType FileType
		expected kinds.Type
	}{
		{
			name:     "library file maps to Lib kind",
			fileType: Library,
			expected: kinds.Lib,
		},
		{
			name:     "test file maps to Test kind",
			fileType: Test,
			expected: kinds.Test,
		},
		{
			name:     "binary file maps to Bin kind",
			fileType: Binary,
			expected: kinds.Bin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{FileType: tt.fileType}
			assert.Equal(t, tt.expected, f.KindType())
		})
	}
}