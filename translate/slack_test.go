package main

import "testing"

// TestQuote test
func TestQuote(t *testing.T) {
	tests := []struct {
		message string
		sender  string
		quoted  string
	}{
		{
			message: "test",
			sender:  "test",
			quoted:  "<@test>:\n> test",
		},
		{
			message: "test\nhello",
			sender:  "test",
			quoted:  "<@test>:\n> test\n>hello",
		},
		{
			message: "> test\nhello",
			sender:  "test",
			quoted:  "<@test>:\n> test\n>hello",
		},
	}
	for _, test := range tests {
		actual := quote(test.message, test.sender)
		if actual != test.quoted {
			t.Errorf("Expecting '%s', got '%s'", test.quoted, actual)
		}
	}
}
