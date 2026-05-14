package validators

import "testing"

func TestLead_Valid(t *testing.T) {
	input := LeadInput{
		Name:    "  Иван Иванов  ",
		Phone:   "+7 (999) 123-45-67",
		Email:   "ivan@example.com",
		Comment: "Хочу заказать услугу",
	}
	out, errs := Lead(input)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if out.Name != "Иван Иванов" {
		t.Errorf("Name not trimmed: %q", out.Name)
	}
}

func TestLead_RequiredFields(t *testing.T) {
	_, errs := Lead(LeadInput{})
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors (name+phone), got %d: %v", len(errs), errs)
	}

	fields := map[string]bool{}
	for _, e := range errs {
		fields[e.Field] = true
	}

	if !fields["name"] {
		t.Error("expected error for field 'name'")
	}

	if !fields["phone"] {
		t.Error("expected error for field 'phone'")
	}
}

func TestLead_WhitespaceOnly(t *testing.T) {
	if _, errs := Lead(LeadInput{Name: "   ", Phone: "\t"}); len(errs) != 2 {
		t.Errorf("expected 2 errors for whitespace-only fields, got %d", len(errs))
	}
}

func TestLead_EmailOptional(t *testing.T) {
	// valid with no email
	if _, errs := Lead(LeadInput{Name: "Test", Phone: "123"}); len(errs) != 0 {
		t.Errorf("expected no errors when email is empty, got %v", errs)
	}
}

func TestLead_InvalidEmail(t *testing.T) {
	if _, errs := Lead(LeadInput{Name: "Test", Phone: "123", Email: "notanemail"}); len(errs) != 1 || errs[0].Field != "email" {
		t.Errorf("expected email error, got %v", errs)
	}
}

func TestLead_ValidEmail(t *testing.T) {
	if _, errs := Lead(LeadInput{Name: "Test", Phone: "123", Email: "test@domain.co.uk"}); len(errs) != 0 {
		t.Errorf("expected no errors for valid email, got %v", errs)
	}
}
