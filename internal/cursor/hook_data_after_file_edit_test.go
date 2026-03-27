package cursor

import "testing"

func TestNewHookDataAfterFileEdit_fullPayload(t *testing.T) {
	input := `{
		"hook_event_name": "afterFileEdit",
		"conversation_id": "conv-1",
		"file_path": "C:/repo/file.txt",
		"edits": [
			{"old_string":"old-1","new_string":"new-1"},
			{"old_string":"old-2","new_string":"new-2"}
		]
	}`

	commonData, err := NewHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	hookData, err := NewHookDataWithCommon[AfterFileEditFields]([]byte(input), commonData)
	if err != nil {
		t.Fatalf("NewHookDataWithCommon: %v", err)
	}

	assertEqual(t, "afterFileEdit", hookData.HookEventName)
	assertEqual(t, "conv-1", hookData.ConversationID)
	assertEqual(t, "C:/repo/file.txt", hookData.Fields.FilePath)
	if len(hookData.Fields.Edits) != 2 {
		t.Fatalf("Edits length: want 2, got %d", len(hookData.Fields.Edits))
	}
	assertEqual(t, "old-1", hookData.Fields.Edits[0].OldString)
	assertEqual(t, "new-1", hookData.Fields.Edits[0].NewString)
	assertEqual(t, "old-2", hookData.Fields.Edits[1].OldString)
	assertEqual(t, "new-2", hookData.Fields.Edits[1].NewString)
}

func TestNewHookDataAfterFileEdit_zeroValueFields(t *testing.T) {
	input := `{
		"hook_event_name": "afterFileEdit",
		"file_path": "",
		"edits": []
	}`

	commonData, err := NewHookDataCommon([]byte(input))
	if err != nil {
		t.Fatalf("NewHookDataCommon: %v", err)
	}
	hookData, err := NewHookDataWithCommon[AfterFileEditFields]([]byte(input), commonData)
	if err != nil {
		t.Fatalf("NewHookDataWithCommon: %v", err)
	}
	assertEqual(t, "", hookData.Fields.FilePath)
	assertIntEqual(t, 0, len(hookData.Fields.Edits))
}

func TestNewHookDataAfterFileEdit_invalidJSON(t *testing.T) {
	_, err := NewHookDataWithCommon[AfterFileEditFields]([]byte(`not json`), HookDataCommon{})
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
