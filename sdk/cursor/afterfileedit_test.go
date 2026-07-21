package cursor

import "testing"

const cursorAfterFileEdit = `{
  "conversation_id": "c1",
  "hook_event_name": "afterFileEdit",
  "cursor_version": "1.7.2",
  "cwd": "/w",
  "file_path": "main.go",
  "edits": [
    {"old_string": "foo", "new_string": "bar"}
  ]
}`

func TestDecode_AfterFileEdit(t *testing.T) {
	ev, err := codec.Decode([]byte(cursorAfterFileEdit))
	if err != nil {
		t.Fatal(err)
	}
	edit, ok := ev.(AfterFileEdit)
	if !ok {
		t.Fatalf("want AfterFileEdit, got %T", ev)
	}
	if edit.FilePath != "main.go" || len(edit.Edits) != 1 || edit.Edits[0].OldString != "foo" {
		t.Fatalf("bad edit: %+v", edit)
	}
}
