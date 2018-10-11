package helperhttp

import "testing"

func TestGetString(t *testing.T) {
	oneValue := "one value"
	var values = map[string][]string{
		"none":     {},
		"single":   {oneValue},
		"multiple": {"one", "two", "three"},
	}

	var originalLength = len(values)

	// Test key does not exist
	if value, err := URLQueryParam(values).GetString("not found"); err != ParamNotFound {
		t.Error("Bad error returned")
	} else if value != "" {
		t.Error("Empty value attended")
	}

	// Test key with empty value
	if value, err := URLQueryParam(values).GetString("none"); err != nil {
		t.Error("Nil error attended")
	} else if value != "" {
		t.Error("Empty value attended")
	}

	// Test key with single value
	if value, err := URLQueryParam(values).GetString("single"); err != nil {
		t.Error("Nil error attended")
	} else if value != oneValue {
		t.Error("Bad value returned")
	}

	// Test key with multiple value
	if value, err := URLQueryParam(values).GetString("multiple"); err != ParamMultipleValues {
		t.Error("Bad error returned")
	} else if value != "" {
		t.Error("Empty value attended")
	}

	// Check length of map not modified
	if len(values) != originalLength {
		t.Error("Map values modified")
	}
}

func TestGetBool(t *testing.T) {
	var valuesTrue = map[string][]string{
		"1":    {"1"},
		"true": {"true"},
		"True": {"True"},
		"TRUE": {"TRUE"},
		"TRuE": {"TRuE"},
	}

	var valuesFalse = map[string][]string{
		"0":     {"0"},
		"false": {"false"},
		"False": {"False"},
		"FALSE": {"FALSE"},
		"FALsE": {"FALsE"},
	}

	var valuesBad = map[string][]string{
		"bad":   {"bad"},
		"empty": {""},
	}

	// Test true values
	for k, v := range valuesTrue {
		if value, err := URLQueryParam(valuesTrue).GetBool(k); err != nil {
			t.Errorf("Error returned for value %v", v)
		} else if value == false {
			t.Errorf("true attended for value %v", v)
		}
	}

	// Test false values
	for k, v := range valuesFalse {
		if value, err := URLQueryParam(valuesFalse).GetBool(k); err != nil {
			t.Errorf("Error returned for value %v", v)
		} else if value == true {
			t.Errorf("false attended for value %v", v)
		}
	}

	// Test bad values
	for k, v := range valuesBad {
		if _, err := URLQueryParam(valuesBad).GetBool(k); err == nil {
			t.Errorf("Nil error returned for value %v", v)
		}
	}
}

func TestGetInt(t *testing.T) {
	var values = map[string][]string{
		"15":     {"15"},
		"-19782": {"-19782"},
	}

	// Test positive value
	if value, err := URLQueryParam(values).GetInt("15"); err != nil {
		t.Error("Bad error returned")
	} else if value != 15 {
		t.Error("Value 15 attended")
	}

	// Test positive value
	if value, err := URLQueryParam(values).GetInt("-19782"); err != nil {
		t.Error("Bad error returned")
	} else if value != -19782 {
		t.Error("Value -19782 attended")
	}

	var valuesBad = map[string][]string{
		"empty": {""},
		"bad":   {"bad"},
		"true":  {"true"},
	}

	// Test bad values
	for k, v := range valuesBad {
		if _, err := URLQueryParam(valuesBad).GetInt(k); err == nil {
			t.Errorf("Nil error returned for value %v", v)
		}
	}
}
