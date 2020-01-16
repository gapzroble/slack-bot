package main

import (
	"testing"
)

// TestCommand test
func TestCommand(t *testing.T) {
	input := "token=5eCkEuoxbbR2rY2FaaMjtWBN&team_id=TA18P4GSZ&team_domain=tiqqe&channel_id=DHS7SGTDX&channel_name=directmessage&user_id=UHS7SGL0Z&user_name=r.roble&command=%2Ftranslatebot&text=yes&response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FTA18P4GSZ%2F878512279365%2FerSNLBPKgPxV45gYHKNUPPkg&trigger_id=867034496323.341295152917.ef307a9be122ee80a0dac76e9065a489"

	if _, ok := command(input); !ok {
		t.Error("Expecting command")
	}
}

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
