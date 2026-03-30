package shellparse

import (
	"maps"
	"slices"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/core"
)

// Obfuscated paths only (no real workspace locations). Asserts fast-path [core.CommandNode] pipeline shape
// for cmdlet-style lines similar to Get-Item | Select-Object, Set-Location; go build, and Test-Path; Get-ChildItem.
func TestPowerShellParser_parseFast_pipelineCommandAST_obfuscatedPaths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want []core.CommandNode
	}{
		{
			name: "get_item_pipe_select_object",
			raw:  `Get-Item 'X:\Contoso\Labs\release.exe' | Select-Object FullName, Length, LastWriteTime`,
			want: []core.CommandNode{
				{
					Name: "Get-Item", Args: []string{`X:\Contoso\Labs\release.exe`},
					Flags: map[string]string{}, PipeIndex: 0, PipeLength: 2,
				},
				{
					Name: "Select-Object", Args: []string{"FullName,", "Length,", "LastWriteTime"},
					Flags: map[string]string{}, PipeIndex: 1, PipeLength: 2,
				},
			},
		},
		{
			name: "set_location_semicolon_go_build",
			raw:  `Set-Location 'X:\Contoso\BuildTree'; go build ./cmd/wat`,
			want: []core.CommandNode{
				{
					Name: "Set-Location", Args: []string{`X:\Contoso\BuildTree`},
					Flags: map[string]string{}, PipeIndex: 0, PipeLength: 2,
				},
				{
					Name: "go", Args: []string{"build", "./cmd/wat"},
					Flags: map[string]string{}, PipeIndex: 1, PipeLength: 2,
				},
			},
		},
		{
			name: "test_path_semicolon_get_child_item_error_action",
			raw:  `Test-Path 'X:\Contoso\Labs\release.exe'; Get-ChildItem 'X:\Contoso\Labs\release.exe' -ErrorAction SilentlyContinue`,
			want: []core.CommandNode{
				{
					Name: "Test-Path", Args: []string{`X:\Contoso\Labs\release.exe`},
					Flags: map[string]string{}, PipeIndex: 0, PipeLength: 2,
				},
				{
					Name: "Get-ChildItem", Args: []string{`X:\Contoso\Labs\release.exe`},
					Flags:     map[string]string{"-ErrorAction": "SilentlyContinue"},
					PipeIndex: 1, PipeLength: 2,
				},
			},
		},
	}

	p := &PowerShellParser{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			res, err := p.Parse(tt.raw)
			if err != nil {
				t.Fatal(err)
			}
			if res.Dialect != core.DialectPowerShell {
				t.Fatalf("dialect: got %q", res.Dialect)
			}
			if res.Raw != tt.raw {
				t.Fatalf("Raw: got %q", res.Raw)
			}
			got := res.Pipeline
			if len(got) != len(tt.want) {
				t.Fatalf("pipeline len: got %d (%+v), want %d", len(got), got, len(tt.want))
			}
			for i := range got {
				g, w := got[i], tt.want[i]
				if g.Name != w.Name {
					t.Errorf("node[%d].Name: got %q, want %q", i, g.Name, w.Name)
				}
				if !slices.Equal(g.Args, w.Args) {
					t.Errorf("node[%d].Args: got %#v, want %#v", i, g.Args, w.Args)
				}
				if !maps.Equal(g.Flags, w.Flags) {
					t.Errorf("node[%d].Flags: got %#v, want %#v", i, g.Flags, w.Flags)
				}
				if !slices.Equal(g.Switches, w.Switches) {
					t.Errorf("node[%d].Switches: got %#v, want %#v", i, g.Switches, w.Switches)
				}
				if g.PipeIndex != w.PipeIndex || g.PipeLength != w.PipeLength {
					t.Errorf("node[%d].PipeIndex/PipeLength: got %d/%d, want %d/%d",
						i, g.PipeIndex, g.PipeLength, w.PipeIndex, w.PipeLength)
				}
			}
		})
	}
}
