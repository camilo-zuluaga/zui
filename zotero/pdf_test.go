package zotero

import "testing"

func TestSystemPDFOpener(t *testing.T) {
	var called bool
	var gotCmd string
	var gotPath string

	mockRunner := func(name, path string) error {
		called = true
		gotCmd = name
		gotPath = path
		return nil
	}

	opener := &SystemPDFOpener{
		SystemPath:    "/tmp/Zotero/storage",
		CommandOpener: "open",
		RunCmd:        mockRunner,
	}

	err := opener.Open("mycollection", "file.pdf")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !called {
		t.Fatal("expected RunCmd to get called")
	}

	if gotCmd != "open" {
		t.Fatalf("expected command 'open', got %s", gotCmd)
	}

	if gotPath != "/tmp/Zotero/storage/mycollection/file.pdf" {
		t.Errorf("unexpected filepath: %s", gotPath)
	}
}
