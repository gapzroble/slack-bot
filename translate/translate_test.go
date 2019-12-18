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
