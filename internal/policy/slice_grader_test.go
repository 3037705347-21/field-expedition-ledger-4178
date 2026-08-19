package policy

import "testing"

func TestCleanListDoesNotAliasCallerStorage(t *testing.T) {
	input := make([]string, 2, 4)
	input[0] = "fresh"
	input[1] = "flow"
	cleaned := CleanList(input)
	cleaned[0] = "changed"
	if input[0] != "fresh" {
		t.Fatalf("cleaned list aliases caller storage: input=%v cleaned=%v", input, cleaned)
	}
}
