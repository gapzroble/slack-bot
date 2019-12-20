package main

import (
	"testing"
)

// TestParse test
func TestParse(t *testing.T) {
	response := "[[[\"\\u003chttps://www.youtube.com/playlist?list\\u003dPLp4wchugWzHsEkxd-K_8zEcGIkQEIF0hM\\u003e To you others who did not have this day! \",\"\\u003chttps://www.youtube.com/playlist?list\\u003dPLp4wchugWzHsEkxd-K_8zEcGIkQEIF0hM\\u003e Till er andra som inte hann med denna dag!\",null,null,3,null,null,null,[[[\"d417779c06d67b45f16785426e87bf85\",\"GermanicB_afdafyisiwlbnosvyi_en_2019q2.md\"]\n]\n]\n]\n,[\"Ping \\u003c@ UA14R112N\\u003e\",\"Ping \\u003c@UA14R112N\\u003e\",null,null,3,null,null,null,[[[\"d417779c06d67b45f16785426e87bf85\",\"GermanicB_afdafyisiwlbnosvyi_en_2019q2.md\"]\n]\n]\n]\n]\n,null,\"sv\",null,null,null,null,[]\n]\n"

	expected := "<https://www.youtube.com/playlist?list=PLp4wchugWzHsEkxd-K_8zEcGIkQEIF0hM> To you others who did not have this day! Ping <@UA14R112N>"
	if parse(response) != expected {
		t.Error("Parse error")
	}
}

// TestCommand test
func TestCommand(t *testing.T) {
	input := "token=5eCkEuoxbbR2rY2FaaMjtWBN&team_id=TA18P4GSZ&team_domain=tiqqe&channel_id=DHS7SGTDX&channel_name=directmessage&user_id=UHS7SGL0Z&user_name=r.roble&command=%2Ftranslatebot&text=yes&response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FTA18P4GSZ%2F878512279365%2FerSNLBPKgPxV45gYHKNUPPkg&trigger_id=867034496323.341295152917.ef307a9be122ee80a0dac76e9065a489"

	if _, ok := command(input); !ok {
		t.Error("Expecting command")
	}
}

// TestTransform test
func TestTransform(t *testing.T) {
	tests := map[string]string{
		" <@ User>": " <@User>",
		// " : Laughing:":                       " :laughing:",
		// " <# CCR0E62H0 | tech discussions> ": " <#CCR0E62H0|tech-discussions> ",
		// " : heart:.":                         " :heart:.",
	}
	for test, expected := range tests {
		actual := transform(test)
		if actual != expected {
			t.Errorf("Failed to transform: %s", test)
		}
	}
}

// TestExtractReplace test
func TestExtractReplace(t *testing.T) {
	message := " Lol, <@Userxxx> is  :laughing: at <#CCR0E62H0|tech-discussions> :heart:.. järn: <https://www.cisco.com/c/en/us/support/docs/field-notices/704/fn70489.html> <https://support.hpe.com/hpsc/doc/public/display?docId=emr_na-a00092491en_us>"
	expected := message

	replacements := extract(&message)
	if len(replacements) != 6 {
		t.Errorf("Expecting 6 replacements, got %d", len(replacements))
	}

	replace(&message, replacements)
	if message != expected {
		t.Error("Not replaced")
	}
}
