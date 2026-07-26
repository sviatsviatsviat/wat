package aftertabfileedit

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func mustDecode[E any](t *testing.T, raw string) E {
	t.Helper()
	ev, err := testCodec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if ev.EventName() == "" {
		t.Fatal("EventName empty")
	}
	typed, ok := ev.(E)
	if !ok {
		t.Fatalf("want %T, got %T", *new(E), ev)
	}
	return typed
}

// cursorAfterTabFileEdit mirrors Cursor Hooks docs for afterTabFileEdit,
// including Tab-specific range / old_line / new_line edit detail.
const cursorAfterTabFileEdit = `{
  "conversation_id": "c1",
  "generation_id": "g1",
  "hook_event_name": "afterTabFileEdit",
  "cursor_version": "1.7.2",
  "workspace_roots": ["/w"],
  "user_email": "dev@example.com",
  "transcript_path": "/tmp/t.jsonl",
  "cwd": "/w",
  "file_path": "/w/main.go",
  "edits": [
    {
      "old_string": "foo",
      "new_string": "bar",
      "range": {
        "start_line_number": 10,
        "start_column": 5,
        "end_line_number": 10,
        "end_column": 20
      },
      "old_line": "\tfoo := 1",
      "new_line": "\tbar := 1"
    }
  ]
}`

func TestDecode_AfterTabFileEdit(t *testing.T) {
	ev := mustDecode[Event](t, cursorAfterTabFileEdit)

	if ev.FilePath != "/w/main.go" {
		t.Errorf("FilePath = %q, want /w/main.go", ev.FilePath)
	}
	if len(ev.Edits) != 1 {
		t.Fatalf("len(Edits) = %d, want 1", len(ev.Edits))
	}

	edit := ev.Edits[0]
	if edit.OldString != "foo" {
		t.Errorf("OldString = %q, want foo", edit.OldString)
	}
	if edit.NewString != "bar" {
		t.Errorf("NewString = %q, want bar", edit.NewString)
	}
	if edit.OldLine != "\tfoo := 1" {
		t.Errorf("OldLine = %q, want %q", edit.OldLine, "\tfoo := 1")
	}
	if edit.NewLine != "\tbar := 1" {
		t.Errorf("NewLine = %q, want %q", edit.NewLine, "\tbar := 1")
	}

	wantRange := event.EditRange{
		StartLineNumber: 10,
		StartColumn:     5,
		EndLineNumber:   10,
		EndColumn:       20,
	}
	if edit.Range != wantRange {
		t.Errorf("Range = %+v, want %+v", edit.Range, wantRange)
	}
}

func TestDecode_AfterTabFileEdit_emptyEdits(t *testing.T) {
	ev := mustDecode[Event](t, `{"hook_event_name":"afterTabFileEdit","conversation_id":"c1","file_path":"x.go","edits":[]}`)
	if ev.FilePath != "x.go" {
		t.Errorf("FilePath = %q, want x.go", ev.FilePath)
	}
	if len(ev.Edits) != 0 {
		t.Fatalf("len(Edits) = %d, want 0", len(ev.Edits))
	}
}

func init() {
	register(testCodec)
}
