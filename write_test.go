package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"
)

func TestWriteOutputToFile(t *testing.T) {
	fis, err := generate("./test", allHashes)
	if err != nil {
		t.Fatalf("Unexpected error from generate call (target: './test'). %q", err)
	}

	if err := writeOutputToFile(fis, testWriteFileOutputPath); err != nil {
		t.Fatalf("Unexpected error from writeOutputToFile call. %q", err)
	}

	assertFileContent(t, testWriteFileOutputPath, "test.js", "sha256", "sha256-jEUM4jWrIiMerWo9zYrx6XwQ5eI77uzuETBptBvPlRQ=")
	assertFileContent(t, testWriteFileOutputPath, "test.js", "sha384", "sha384-zBTHeP/UZLYRhjvTi7r3Dx7MTCNf/ddGENI26AacmrgqzH8YOkA+EJ14MXpwD4wL")
	assertFileContent(t, testWriteFileOutputPath, "test.js", "sha512", "sha512-RmToDhq62z0wvdoBN8yl5SfjLtbF84USifqZuJNyJ1K99b27Jo/BE16veNzzvTHVBY8NWvvD2M0Vc1NDJWf2Yw==")

	assertFileContent(t, testWriteFileOutputPath, "test.min.js", "sha256", "sha256-ODBnPrz8p2bs/l/ffyD4jUqpRkTvzlmFu8WDCuYNYms=")
	assertFileContent(t, testWriteFileOutputPath, "test.min.js", "sha384", "sha384-GFSKzS/+oGDIT70dABnjqACvEFXH8kCG4tW9e3athjSbADyCkj3Mlfk1a2mmtAWa")
	assertFileContent(t, testWriteFileOutputPath, "test.min.js", "sha512", "sha512-KNFlyFMpPl1IsSOPottAYfbxAj9RZsp1gOw4usysmMjszS2OPMTff7RU/7AB6V6PgImA1SZg1RAm8GF9Q5tvFg==")

	assertFileContent(t, testWriteFileOutputPath, "test.css", "sha256", "sha256-ckxnbs3D4win9ik/Eh1/55cPi1yJ4xBVTU5npga+uw8=")
	assertFileContent(t, testWriteFileOutputPath, "test.css", "sha384", "sha384-MTTquAJ9el7bwG1gCco7oa3yfC5RTZV8zsXxwrbLX3VM/dgKHYmvE8M1OZNgzDrU")
	assertFileContent(t, testWriteFileOutputPath, "test.css", "sha512", "sha512-w8UmgrBX4zom3NhAo/VY6LMLwf3ZW+0tbsR4HvpRHyLJYDux3eoBZW37Iqej5AS8+oCisLQ+PIKtlU4sVLRecA==")
}

func assertFileContent(t *testing.T, jsonFile, fileName, algo, digest string) {
	b, _ := ioutil.ReadFile(jsonFile)

	var output map[string]interface{}
	json.Unmarshal(b, &output)

	fileNode := output[fileName].(map[string]interface{})
	algoNode := fileNode[algo].(map[string]interface{})

	if result := algoNode["digest"]; result != digest {
		t.Fatalf("Got digest %s. Expected: %s", result, digest)
	}

	if tag := algoNode["tag"].(string); !strings.Contains(tag, digest) {
		t.Fatalf("Tag did not contain expected digest. Got %s", tag)
	}
}
