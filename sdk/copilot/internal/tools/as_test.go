package tools

import (
	"encoding/json"
	"testing"
)

func TestAsAccessors_RuntimeAndClaudeNames(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name   string
		tool   string
		raw    string
		check  func(Input, *testing.T)
		reject string
	}{
		{
			name: "bash_runtime",
			tool: "bash",
			raw:  `{"command":"ls"}`,
			check: func(in Input, t *testing.T) {
				got, ok := in.AsBash()
				if !ok || got.Command != "ls" {
					t.Fatalf("AsBash = %+v, %v", got, ok)
				}
			},
			reject: "edit",
		},
		{
			name: "bash_claude",
			tool: "Bash",
			raw:  `{"command":"ls"}`,
			check: func(in Input, t *testing.T) {
				got, ok := in.AsBash()
				if !ok || got.Command != "ls" {
					t.Fatalf("AsBash = %+v, %v", got, ok)
				}
			},
		},
		{
			name: "view_read",
			tool: "Read",
			raw:  `{"file_path":"/a","offset":2,"limit":5}`,
			check: func(in Input, t *testing.T) {
				got, ok := in.AsView()
				if !ok || got.Path != "/a" || got.Offset != 2 || got.Limit != 5 {
					t.Fatalf("AsView = %+v, %v", got, ok)
				}
			},
			reject: "bash",
		},
		{
			name: "create_write",
			tool: "Write",
			raw:  `{"file_path":"/a","content":"x"}`,
			check: func(in Input, t *testing.T) {
				got, ok := in.AsCreate()
				if !ok || got.Path != "/a" || got.Content != "x" {
					t.Fatalf("AsCreate = %+v, %v", got, ok)
				}
			},
		},
		{
			name: "edit_alias",
			tool: "str_replace_editor",
			raw:  `{"path":"/a","old_str":"x","new_str":"y"}`,
			check: func(in Input, t *testing.T) {
				got, ok := in.AsEdit()
				if !ok || got.Path != "/a" || got.OldString != "x" || got.NewString != "y" {
					t.Fatalf("AsEdit = %+v, %v", got, ok)
				}
			},
		},
		{
			name: "edit_claude",
			tool: "Edit",
			raw:  `{"file_path":"/a","old_string":"x","new_string":"y"}`,
			check: func(in Input, t *testing.T) {
				got, ok := in.AsEdit()
				if !ok || got.Path != "/a" || got.OldString != "x" || got.NewString != "y" {
					t.Fatalf("AsEdit = %+v, %v", got, ok)
				}
			},
		},
		{
			name: "edit_apply_patch",
			tool: "apply_patch",
			raw:  `{"patch":"*** Begin Patch\n*** End Patch"}`,
			check: func(in Input, t *testing.T) {
				got, ok := in.AsEdit()
				if !ok || got.Patch == "" {
					t.Fatalf("AsEdit = %+v, %v", got, ok)
				}
			},
		},
		{
			name: "edit_file_text",
			tool: "str_replace_editor",
			raw:  `{"path":"/a","file_text":"hello"}`,
			check: func(in Input, t *testing.T) {
				got, ok := in.AsEdit()
				if !ok || got.Path != "/a" || got.Content != "hello" {
					t.Fatalf("AsEdit = %+v, %v", got, ok)
				}
			},
		},
		{
			name: "view_runtime",
			tool: "view",
			raw:  `{"path":"/b"}`,
			check: func(in Input, t *testing.T) {
				got, ok := in.AsView()
				if !ok || got.Path != "/b" {
					t.Fatalf("AsView = %+v, %v", got, ok)
				}
			},
		},
		{
			name: "create_runtime",
			tool: "create",
			raw:  `{"path":"/b","content":"z"}`,
			check: func(in Input, t *testing.T) {
				got, ok := in.AsCreate()
				if !ok || got.Path != "/b" || got.Content != "z" {
					t.Fatalf("AsCreate = %+v, %v", got, ok)
				}
			},
		},
		{
			name: "edit_runtime",
			tool: "edit",
			raw:  `{"path":"/b","content":"z"}`,
			check: func(in Input, t *testing.T) {
				got, ok := in.AsEdit()
				if !ok || got.Path != "/b" || got.Content != "z" {
					t.Fatalf("AsEdit = %+v, %v", got, ok)
				}
			},
		},
		{
			name: "glob",
			tool: "Glob",
			raw:  `{"pattern":"*.go","path":"/src"}`,
			check: func(in Input, t *testing.T) {
				got, ok := in.AsGlob()
				if !ok || got.Pattern != "*.go" || got.Path != "/src" {
					t.Fatalf("AsGlob = %+v, %v", got, ok)
				}
			},
		},
		{
			name: "grep_rg",
			tool: "rg",
			raw:  `{"pattern":"foo"}`,
			check: func(in Input, t *testing.T) {
				got, ok := in.AsGrep()
				if !ok || got.Pattern != "foo" {
					t.Fatalf("AsGrep = %+v, %v", got, ok)
				}
			},
		},
		{
			name: "web_fetch",
			tool: "WebFetch",
			raw:  `{"url":"https://example.com"}`,
			check: func(in Input, t *testing.T) {
				got, ok := in.AsWebFetch()
				if !ok || got.URL != "https://example.com" {
					t.Fatalf("AsWebFetch = %+v, %v", got, ok)
				}
			},
		},
		{
			name: "web_search",
			tool: "web_search",
			raw:  `{"query":"wat hooks"}`,
			check: func(in Input, t *testing.T) {
				got, ok := in.AsWebSearch()
				if !ok || got.Query != "wat hooks" {
					t.Fatalf("AsWebSearch = %+v, %v", got, ok)
				}
			},
		},
		{
			name: "task_agent",
			tool: "Agent",
			raw:  `{"prompt":"explore","subagent_type":"explore"}`,
			check: func(in Input, t *testing.T) {
				got, ok := in.AsTask()
				if !ok || got.Prompt != "explore" || got.SubagentType != "explore" {
					t.Fatalf("AsTask = %+v, %v", got, ok)
				}
			},
		},
		{
			name: "ask_user",
			tool: "AskUserQuestion",
			raw:  `{"questions":[{"question":"Pick?","multiSelect":true}]}`,
			check: func(in Input, t *testing.T) {
				got, ok := in.AsAskUser()
				if !ok || len(got.Questions) != 1 || !got.Questions[0].MultiSelect {
					t.Fatalf("AsAskUser = %+v, %v", got, ok)
				}
			},
		},
		{
			name: "update_todo",
			tool: "TodoWrite",
			raw:  `{"todos":[{"content":"ship","status":"pending","activeForm":"shipping"}]}`,
			check: func(in Input, t *testing.T) {
				got, ok := in.AsUpdateTodo()
				if !ok || len(got.Todos) != 1 || got.Todos[0].Content != "ship" {
					t.Fatalf("AsUpdateTodo = %+v, %v", got, ok)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			in := NewInput(tc.tool, json.RawMessage(tc.raw))
			tc.check(in, t)
			if tc.reject != "" {
				bad := NewInput(tc.reject, json.RawMessage(tc.raw))
				switch tc.name {
				case "bash_runtime":
					if _, ok := bad.AsBash(); ok {
						t.Fatal("AsBash should reject edit")
					}
				case "view_read":
					if _, ok := bad.AsView(); ok {
						t.Fatal("AsView should reject bash")
					}
				}
			}
		})
	}
}
