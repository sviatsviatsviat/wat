package afterfileedit

import (
	"context"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

const cursorAfterFileEdit = `{
  "conversation_id": "c1",
  "generation_id": "g1",
  "hook_event_name": "afterFileEdit",
  "cursor_version": "1.7.2",
  "cwd": "/w",
  "file_path": "/w/main.go",
  "edits": [
    {"old_string": "foo", "new_string": "bar"}
  ]
}`

func TestDecode_AfterFileEdit(t *testing.T) {
	ev, err := testCodec.Decode([]byte(cursorAfterFileEdit))
	if err != nil {
		t.Fatal(err)
	}
	edit, ok := ev.(Event)
	if !ok {
		t.Fatalf("want AfterFileEdit, got %T", ev)
	}
	if edit.FilePath != "/w/main.go" || len(edit.Edits) != 1 || edit.Edits[0].OldString != "foo" || edit.Edits[0].NewString != "bar" {
		t.Fatalf("bad edit: %+v", edit)
	}
}

func TestRegisterHandler_ObserveOnly(t *testing.T) {
	c := hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)
	d := hookkit.NewDialect(c)

	var saw Event
	RegisterHandler(d, func(_ context.Context, hook Event) error {
		saw = hook
		return nil
	})

	handlers := d.HandlersFor(Event{}.EventName())
	if len(handlers) != 1 {
		t.Fatalf("handlers = %d, want 1", len(handlers))
	}

	decoded, err := c.Decode([]byte(cursorAfterFileEdit))
	if err != nil {
		t.Fatal(err)
	}
	out, err := handlers[0].Invoke(context.Background(), decoded)
	if err != nil {
		t.Fatal(err)
	}
	if out != nil {
		t.Fatalf("observe-only handler returned output %#v, want nil", out)
	}
	if saw.FilePath != "/w/main.go" || len(saw.Edits) != 1 || saw.Edits[0].NewString != "bar" {
		t.Fatalf("handler saw %#v", saw)
	}

	RegisterHandler(d, nil)
	if len(d.HandlersFor(Event{}.EventName())) != 1 {
		t.Fatalf("nil fn should not append handlers, got %d", len(d.HandlersFor(Event{}.EventName())))
	}
}

func init() {
	register(testCodec)
}
