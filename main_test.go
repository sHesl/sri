package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	testDownloadOutputPath = "testdownload.json"
	testFileOutputPath     = "testfile.json"
	testDirOutputPath      = "testdir.json"
)

func TestMain(m *testing.M) {
	m.Run()
	os.Remove(testDownloadOutputPath)
	os.Remove(testFileOutputPath)
	os.Remove(testDirOutputPath)
}

func TestDownload(t *testing.T) {
	stubJSFileHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("console.log('hello world!');"))
	})

	mockServer := httptest.NewServer(stubJSFileHandler)
	client = mockServer.Client()

	run(mockServer.URL, testDownloadOutputPath)

	assertFileContent(t, testDownloadOutputPath, mockServer.URL[7:], "sha256", "sha256-lClGOfcWqtQdAvO3zCRzZEg/4RmOMbr9/V54QO76j/A=")
	assertFileContent(t, testDownloadOutputPath, mockServer.URL[7:], "sha384", "sha384-3Zn0DhQDSbiCfvVo1SIqZ0jy9ybVafdjeIRnqOOil7SXoC86q2Avs4w8xnN96fC2")
	assertFileContent(t, testDownloadOutputPath, mockServer.URL[7:], "sha512", "sha512-gzbGfS1swNgrzjRJK75UMtYICNYdffO3ReSaRyFE6HiFlqn5Vvnw8OoNllTjFOdUZ622tZqukf5+p0OTRAL2Qg==")
}

func TestFile(t *testing.T) {
	run("test/test.js", testFileOutputPath)

	assertFileContent(t, testFileOutputPath, "test.js", "sha256", "sha256-jEUM4jWrIiMerWo9zYrx6XwQ5eI77uzuETBptBvPlRQ=")
	assertFileContent(t, testFileOutputPath, "test.js", "sha384", "sha384-zBTHeP/UZLYRhjvTi7r3Dx7MTCNf/ddGENI26AacmrgqzH8YOkA+EJ14MXpwD4wL")
	assertFileContent(t, testFileOutputPath, "test.js", "sha512", "sha512-RmToDhq62z0wvdoBN8yl5SfjLtbF84USifqZuJNyJ1K99b27Jo/BE16veNzzvTHVBY8NWvvD2M0Vc1NDJWf2Yw==")
}

func TestDir(t *testing.T) {
	run("test", testDirOutputPath)

	assertFileContent(t, testDirOutputPath, "test.js", "sha256", "sha256-jEUM4jWrIiMerWo9zYrx6XwQ5eI77uzuETBptBvPlRQ=")
	assertFileContent(t, testDirOutputPath, "test.js", "sha384", "sha384-zBTHeP/UZLYRhjvTi7r3Dx7MTCNf/ddGENI26AacmrgqzH8YOkA+EJ14MXpwD4wL")
	assertFileContent(t, testDirOutputPath, "test.js", "sha512", "sha512-RmToDhq62z0wvdoBN8yl5SfjLtbF84USifqZuJNyJ1K99b27Jo/BE16veNzzvTHVBY8NWvvD2M0Vc1NDJWf2Yw==")

	assertFileContent(t, testDirOutputPath, "test.min.js", "sha256", "sha256-ODBnPrz8p2bs/l/ffyD4jUqpRkTvzlmFu8WDCuYNYms=")
	assertFileContent(t, testDirOutputPath, "test.min.js", "sha384", "sha384-GFSKzS/+oGDIT70dABnjqACvEFXH8kCG4tW9e3athjSbADyCkj3Mlfk1a2mmtAWa")
	assertFileContent(t, testDirOutputPath, "test.min.js", "sha512", "sha512-KNFlyFMpPl1IsSOPottAYfbxAj9RZsp1gOw4usysmMjszS2OPMTff7RU/7AB6V6PgImA1SZg1RAm8GF9Q5tvFg==")

	assertFileContent(t, testDirOutputPath, "test.css", "sha256", "sha256-ckxnbs3D4win9ik/Eh1/55cPi1yJ4xBVTU5npga+uw8=")
	assertFileContent(t, testDirOutputPath, "test.css", "sha384", "sha384-MTTquAJ9el7bwG1gCco7oa3yfC5RTZV8zsXxwrbLX3VM/dgKHYmvE8M1OZNgzDrU")
	assertFileContent(t, testDirOutputPath, "test.css", "sha512", "sha512-w8UmgrBX4zom3NhAo/VY6LMLwf3ZW+0tbsR4HvpRHyLJYDux3eoBZW37Iqej5AS8+oCisLQ+PIKtlU4sVLRecA==")
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
