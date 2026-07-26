package workspaceopen

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

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

func TestDecode_WorkspaceOpen(t *testing.T) {
	e := mustDecode[Event](t, `{"hook_event_name":"workspaceOpen","cursor_version":"1.7.2","workspace_roots":["/w"],"user_email":null}`)
	if e.CursorVersion != "1.7.2" || len(e.WorkspaceRoots) != 1 || e.WorkspaceRoots[0] != "/w" {
		t.Fatalf("event=%+v", e)
	}
}

func TestEncode_WorkspaceOpen_pluginPaths(t *testing.T) {
	out, code, err := results{}.PluginPaths([]string{"/plugins/a", "/plugins/b"}).Encode()
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	var got struct {
		PluginPaths []string `json:"pluginPaths"`
	}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	want := []string{"/plugins/a", "/plugins/b"}
	if !reflect.DeepEqual(got.PluginPaths, want) {
		t.Fatalf("pluginPaths = %v, want %v", got.PluginPaths, want)
	}
}

func TestEncode_WorkspaceOpen_noopIsSilent(t *testing.T) {
	out, code, err := results{}.Noop().Encode()
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	if out != nil {
		t.Fatalf("noop stdout = %s, want nil", out)
	}
}

func TestMerge_WorkspaceOpen_pluginPathsTakeLast(t *testing.T) {
	first := []string{"/plugins/a"}
	second := []string{"/plugins/b", "/plugins/c"}
	a := results{}.PluginPaths(first)
	b := results{}.PluginPaths(second)

	merged, warnings, err := a.Merge(b.(hookkit.Output))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 1 || warnings[0] == "" {
		t.Fatalf("warnings = %v", warnings)
	}
	out := merged.(output)
	if !reflect.DeepEqual(out.pluginPaths, second) {
		t.Fatalf("pluginPaths = %v, want %v", out.pluginPaths, second)
	}
	if merged.Stop() {
		t.Fatal("pluginPaths should not stop")
	}

	// Merge must clone; mutating inputs or the merged slice must not cross-contaminate.
	second[0] = "/mutated"
	out.pluginPaths[1] = "/also-mutated"
	if !reflect.DeepEqual(first, []string{"/plugins/a"}) {
		t.Fatalf("first input mutated: %v", first)
	}
	if !reflect.DeepEqual(second, []string{"/mutated", "/plugins/c"}) {
		t.Fatalf("second input unexpectedly reshaped: %v", second)
	}
	if reflect.DeepEqual(out.pluginPaths, second) {
		t.Fatal("merged slice shares backing with caller input")
	}
}

func TestMerge_WorkspaceOpen_emptyKeepsPrior(t *testing.T) {
	a := results{}.PluginPaths([]string{"/plugins/a"})
	b := results{}.Noop()
	merged, warnings, err := a.Merge(b.(hookkit.Output))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 0 {
		t.Fatalf("warnings = %v", warnings)
	}
	out := merged.(output)
	if !reflect.DeepEqual(out.pluginPaths, []string{"/plugins/a"}) {
		t.Fatalf("pluginPaths = %v", out.pluginPaths)
	}
}

func TestRegisterHandler_RegistersDecoder(t *testing.T) {
	c := hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)
	d := hookkit.NewDialect(c)
	raw := []byte(`{"hook_event_name":"workspaceOpen","cursor_version":"1.7.2","workspace_roots":["/w"]}`)

	_, err := c.Decode(raw)
	if err == nil {
		t.Fatal("expected unknown event before RegisterHandler")
	}

	RegisterHandler(d, func(context.Context, Event, Results) (Output, error) {
		return results{}.Noop(), nil
	})

	ev, err := c.Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if ev.EventName() != "workspaceOpen" {
		t.Fatalf("EventName = %q", ev.EventName())
	}

	RegisterHandler(d, nil)
	if len(d.HandlersFor("workspaceOpen")) != 1 {
		t.Fatalf("nil fn should not append handlers, got %d", len(d.HandlersFor("workspaceOpen")))
	}
}

func init() {
	register(testCodec)
}
