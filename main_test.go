package main

import (
	"os"
	"testing"
)

var (
	testWriteFileOutputPath = "test.json"
)

func TestMain(m *testing.M) {
	m.Run()
	os.Remove(testWriteFileOutputPath)
}

func TestValidateHash(t *testing.T) {
	for _, h := range []string{sha256Algo, sha384Algo, sha512Algo, allHashes} {
		if err := validateHash(h); err != nil {
			t.Fatalf("Expected %s to be a valid hash value", h)
		}
	}

	if err := validateHash("not a real hash"); err == nil {
		t.Fatalf("Expected invalid hash value to produce an error")
	}
}

func TestValidateGenerate(t *testing.T) {
	type testCase struct {
		inputs []string
		errMsg string
	}

	testCases := []testCase{
		{
			inputs: []string{},
			errMsg: "No target specified for SRI generation",
		},
		{
			inputs: []string{""},
			errMsg: "Received an empty target for SRI generation",
		},
		{
			inputs: []string{"one input"},
			errMsg: "",
		},
	}

	for _, tc := range testCases {
		err := validateGenerate(tc.inputs)

		if err == nil && tc.errMsg != "" {
			t.Fatalf("Expected an error from validateGenerate call")
		}

		if tc.errMsg == "" && err != nil {
			t.Fatalf("Unexpected an error from validateCompare call. %q", err)
		}

		if err != nil && err.Error() != tc.errMsg {
			t.Fatalf("Expected error message of %s. Got %q", tc.errMsg, err)
		}
	}
}
