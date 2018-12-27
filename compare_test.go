package main

import "testing"

func TestCompare(t *testing.T) {
	same, a, b, err := comparison("test/compare-same-a.js", "test/compare-same-b.js")

	if err != nil {
		t.Fatalf("Unexpected error from comparison call. %q", err)
	}

	if !same {
		t.Fatalf("Expected comparison of compare-same-a.js and compare-same-b.js to be true")
	}

	exp := "sha256-hwj4HVJ7OFOzPES8HffZ4IySCiQq7P/+1RT9YQJMAXs="
	if a != exp || b != exp {
		t.Fatalf("Expected 'compare-same-a.js' and 'compare-same-b.js' to have sha256 digest %s. Got %s", exp, a)
	}
}

func TestCompareFail(t *testing.T) {
	same, a, b, err := comparison("test/compare-diff-a.js", "test/compare-diff-b.js")

	if err != nil {
		t.Fatalf("Unexpected error from comparison call. %q", err)
	}

	if same {
		t.Fatalf("Expected comparison of compare-diff-a.js and compare-diff-b.js to be false")
	}

	expectedA := "sha256-3IQq9U6JsK64FtWbunydrhJ4J7ZhUvqkLa7ZymJHWwE="
	expectedB := "sha256-BGvN/h+hPgaFcujuoAfEMVaeFX6JosMdgwnAqPSVgdQ="

	if a != expectedA {
		t.Fatalf("Expected both 'compare-diff-a.js' to have digest %s. Got %s", expectedA, a)
	}

	if b != expectedB {
		t.Fatalf("Expected both 'compare-diff-b.js' to have digest %s. Got %s", expectedB, b)
	}
}

func TestCompareHandlesErrors(t *testing.T) {
	type testCase struct {
		aInput string
		bInput string
		errMsg string
	}

	testCases := []testCase{
		{"test/not-real.js", "test/compare-same-b.js", "Unable to produce both integrities for [\"test/not-real.js\" \"test/compare-same-b.js\"]"},
		{"test/compare-same-a.js", "test/not-real.js", "Unable to produce both integrities for [\"test/compare-same-a.js\" \"test/not-real.js\"]"},
	}

	for _, tc := range testCases {
		same, a, b, err := comparison(tc.aInput, tc.bInput)

		if err == nil {
			t.Fatalf("Expected an error from comparison call")
		}

		if err.Error() != tc.errMsg {
			t.Fatalf("Expected error message of %s. Got %q", tc.errMsg, err)
		}

		if same {
			t.Fatalf("Expected comparison call to return false on error")
		}

		if a != "" || b != "" {
			t.Fatalf("Expected failed comparison to not return digests")
		}
	}
}

func TestValidateCompare(t *testing.T) {
	type testCase struct {
		inputs []string
		errMsg string
	}

	testCases := []testCase{
		{
			inputs: []string{"only one input"},
			errMsg: "Expected two targets to be specified for comparison",
		},
		{
			inputs: []string{"two inputs", "two different inputs"},
			errMsg: "",
		},
		{
			inputs: []string{"two inputs", ""},
			errMsg: "Received an empty target for comparison",
		},
		{
			inputs: []string{"two identical inputs", "two identical inputs"},
			errMsg: "Received two indentical inputs for comparison",
		},
		{
			inputs: []string{"", "two inputs"},
			errMsg: "Received an empty target for comparison",
		},
		{
			inputs: []string{"three inputs", "three inputs", "three inputs"},
			errMsg: "Expected two targets to be specified for comparison",
		},
	}

	for _, tc := range testCases {
		err := validateCompare(tc.inputs)

		if err == nil && tc.errMsg != "" {
			t.Fatalf("Expected an error from validateCompare call")
		}

		if tc.errMsg == "" && err != nil {
			t.Fatalf("Unexpected an error from validateCompare call. %q", err)
		}

		if err != nil && err.Error() != tc.errMsg {
			t.Fatalf("Expected error message of %s. Got %q", tc.errMsg, err)
		}
	}
}
