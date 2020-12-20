package parser_test

import (
	"testing"

	"github.com/reesilva/quake-log/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name      string
		want      string
		parameter string
	}{
		{
			name:      "Este teste tem que falhar",
			want:      "Este cara é o Felipe Neto",
			parameter: "Glauber Rocha",
		},
		{
			name:      "Este teste tem que passar",
			want:      "Este cara é o Felipe Neto",
			parameter: "Felipe Neto",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.ParseLine(tt.parameter)
			assert.Equal(t, tt.want, got)
		})
	}
}
