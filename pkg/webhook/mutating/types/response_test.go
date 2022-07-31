package types

import "testing"

func Test_MutatingResponseFromFile_Allowed(t *testing.T) {
	r, err := MutatingResponseFromFile("testdata/response/good_allow.json")

	if err != nil {
		t.Fatalf("MutatingResponse should be loaded from file: %v", err)
	}

	if r == nil {
		t.Fatalf("MutatingResponse should not be nil")
	}

	if !r.Allowed {
		t.Fatalf("MutatingResponse should have allowed=true: %#v", r)
	}
}

func Test_MutatingResponseFromFile_AllowedWithWarnings(t *testing.T) {
	r, err := MutatingResponseFromFile("testdata/response/good_allow_warnings.json")

	if err != nil {
		t.Fatalf("MutatingResponse should be loaded from file: %v", err)
	}

	if r == nil {
		t.Fatalf("MutatingResponse should not be nil")
	}

	if !r.Allowed {
		t.Fatalf("MutatingResponse should have allowed=true: %#v", r)
	}

	if len(r.Warnings) != 2 {
		t.Fatalf("MutatingResponse should have warnings: %#v", r)
	}
}

func Test_MutatingResponseFromFile_NotAllowed_WithMessage(t *testing.T) {
	r, err := MutatingResponseFromFile("testdata/response/good_deny.json")

	if err != nil {
		t.Fatalf("MutatingResponse should be loaded from file: %v", err)
	}

	if r == nil {
		t.Fatalf("MutatingResponse should not be nil")
	}

	if r.Allowed {
		t.Fatalf("MutatingResponse should have allowed=false: %#v", r)
	}

	if r.Message == "" {
		t.Fatalf("MutatingResponse should have message: %#v", r)
	}
}

func Test_MutatingResponseFromFile_NotAllowed_WithoutMessage(t *testing.T) {
	r, err := MutatingResponseFromFile("testdata/response/good_deny_quiet.json")

	if err != nil {
		t.Fatalf("MutatingResponse should be loaded from file: %v", err)
	}

	if r == nil {
		t.Fatalf("MutatingResponse should not be nil")
	}

	if r.Allowed {
		t.Fatalf("MutatingResponse should have allowed=false: %#v", r)
	}

	if r.Message != "" {
		t.Fatalf("MutatingResponse should have no message: %#v", r)
	}
}
