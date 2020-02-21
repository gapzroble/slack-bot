package google

import (
	"testing"
)

// TestExtractReplace test
func TestExtractReplace(t *testing.T) {
	message := ` Lol, <@Userxxx> is  :laughing: at <#CCR0E62H0|tech-discussions> :heart:.. järn: <https://www.cisco.com/c/en/us/support/docs/field-notices/704/fn70489.html> <https://support.hpe.com/hpsc/doc/public/display?docId=emr_na-a00092491en_us>
	Hej alla :blush:
	a sätt. :sweat_smile:
	dig? :hugging_face:`
	expected := message

	replacements := extract(&message)
	if len(replacements) != 9 {
		t.Errorf("Expecting 9 replacements, got %d", len(replacements))
	}

	replace(&message, replacements)
	if message != expected {
		t.Error("Not replaced")
	}
}
