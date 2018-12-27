package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerate(t *testing.T) {
	// Initialise a stub JS file server to test downloading scripts
	serve := func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("console.log('hello world!');")) }
	stubJSFileHandler := http.HandlerFunc(serve)
	mockServer := httptest.NewServer(stubJSFileHandler)
	client = mockServer.Client()

	type testCase struct {
		targets []string
		exp     map[string]map[string]string
	}

	testCases := []testCase{
		{
			targets: []string{mockServer.URL}, // test file downloading
			exp: map[string]map[string]string{
				mockServer.URL: map[string]string{
					"sha256": "sha256-lClGOfcWqtQdAvO3zCRzZEg/4RmOMbr9/V54QO76j/A=",
					"sha384": "sha384-3Zn0DhQDSbiCfvVo1SIqZ0jy9ybVafdjeIRnqOOil7SXoC86q2Avs4w8xnN96fC2",
					"sha512": "sha512-gzbGfS1swNgrzjRJK75UMtYICNYdffO3ReSaRyFE6HiFlqn5Vvnw8OoNllTjFOdUZ622tZqukf5+p0OTRAL2Qg==",
				},
			},
		},
		{
			targets: []string{"test/test.js"}, // test single file
			exp: map[string]map[string]string{
				"test.js": map[string]string{
					"sha256": "sha256-jEUM4jWrIiMerWo9zYrx6XwQ5eI77uzuETBptBvPlRQ=",
					"sha384": "sha384-zBTHeP/UZLYRhjvTi7r3Dx7MTCNf/ddGENI26AacmrgqzH8YOkA+EJ14MXpwD4wL",
					"sha512": "sha512-RmToDhq62z0wvdoBN8yl5SfjLtbF84USifqZuJNyJ1K99b27Jo/BE16veNzzvTHVBY8NWvvD2M0Vc1NDJWf2Yw==",
				},
			},
		},
		{
			targets: []string{"test/test.js", "test/test.min.js"}, // test two files
			exp: map[string]map[string]string{
				"test.js": map[string]string{
					"sha256": "sha256-jEUM4jWrIiMerWo9zYrx6XwQ5eI77uzuETBptBvPlRQ=",
					"sha384": "sha384-zBTHeP/UZLYRhjvTi7r3Dx7MTCNf/ddGENI26AacmrgqzH8YOkA+EJ14MXpwD4wL",
					"sha512": "sha512-RmToDhq62z0wvdoBN8yl5SfjLtbF84USifqZuJNyJ1K99b27Jo/BE16veNzzvTHVBY8NWvvD2M0Vc1NDJWf2Yw==",
				},
				"test.min.js": map[string]string{
					"sha256": "sha256-ODBnPrz8p2bs/l/ffyD4jUqpRkTvzlmFu8WDCuYNYms=",
					"sha384": "sha384-GFSKzS/+oGDIT70dABnjqACvEFXH8kCG4tW9e3athjSbADyCkj3Mlfk1a2mmtAWa",
					"sha512": "sha512-KNFlyFMpPl1IsSOPottAYfbxAj9RZsp1gOw4usysmMjszS2OPMTff7RU/7AB6V6PgImA1SZg1RAm8GF9Q5tvFg==",
				},
			},
		},
		{
			targets: []string{"./test"}, // test directory
			exp: map[string]map[string]string{
				"test.js": map[string]string{
					"sha256": "sha256-jEUM4jWrIiMerWo9zYrx6XwQ5eI77uzuETBptBvPlRQ=",
					"sha384": "sha384-zBTHeP/UZLYRhjvTi7r3Dx7MTCNf/ddGENI26AacmrgqzH8YOkA+EJ14MXpwD4wL",
					"sha512": "sha512-RmToDhq62z0wvdoBN8yl5SfjLtbF84USifqZuJNyJ1K99b27Jo/BE16veNzzvTHVBY8NWvvD2M0Vc1NDJWf2Yw==",
				},
				"test.min.js": map[string]string{
					"sha256": "sha256-ODBnPrz8p2bs/l/ffyD4jUqpRkTvzlmFu8WDCuYNYms=",
					"sha384": "sha384-GFSKzS/+oGDIT70dABnjqACvEFXH8kCG4tW9e3athjSbADyCkj3Mlfk1a2mmtAWa",
					"sha512": "sha512-KNFlyFMpPl1IsSOPottAYfbxAj9RZsp1gOw4usysmMjszS2OPMTff7RU/7AB6V6PgImA1SZg1RAm8GF9Q5tvFg==",
				},
				"test.css": map[string]string{
					"sha256": "sha256-ckxnbs3D4win9ik/Eh1/55cPi1yJ4xBVTU5npga+uw8=",
					"sha384": "sha384-MTTquAJ9el7bwG1gCco7oa3yfC5RTZV8zsXxwrbLX3VM/dgKHYmvE8M1OZNgzDrU",
					"sha512": "sha512-w8UmgrBX4zom3NhAo/VY6LMLwf3ZW+0tbsR4HvpRHyLJYDux3eoBZW37Iqej5AS8+oCisLQ+PIKtlU4sVLRecA==",
				},
				"compare-diff-a.js": map[string]string{
					"sha256": "sha256-3IQq9U6JsK64FtWbunydrhJ4J7ZhUvqkLa7ZymJHWwE=",
					"sha384": "sha384-rSIo3wZrhBJ1wSxglziVHO8NkORbjEutoXUWShAdMjlNqAwGUlQTOoP1R05p4SyI",
					"sha512": "sha512-mjYhv+MbTRDos3zT2OUkMjXL7ttNco6zecBUWkqcqx1TsWSL68Fy8mMO2XrzsBYQUeG4HW8v4iGpMC10JHTslw==",
				},
				"compare-diff-b.js": map[string]string{
					"sha256": "sha256-BGvN/h+hPgaFcujuoAfEMVaeFX6JosMdgwnAqPSVgdQ=",
					"sha384": "sha384-mTwqjfCmkZz49qB9Peg0+cz9i+AjXApFWksuOL2CEMfu+wYJqPeK+GVOfZlyJ0K0",
					"sha512": "sha512-hRgOtwmZRi6A6KxOTLQeQSW9unQJfOio71gQjAaHh0u9ycbUOrN2e2pSE1TmgDFFS1yzlpkdBs3yU4foW4rVcw==",
				},
				"compare-same-a.js": map[string]string{
					"sha256": "sha256-hwj4HVJ7OFOzPES8HffZ4IySCiQq7P/+1RT9YQJMAXs=",
					"sha384": "sha384-6yKRlQqq9r5dU0GEivGoDai04RH+ufhfO1htclXkbjdJ+184pc1rRrsWhk3aDf3D",
					"sha512": "sha512-rk6z+muPQjmOe6aC/kKX9yYmaaLs42TIBxh02xwV/vkrIpD8RJzjt4cB4q9MXOeYkwD4wSyEMNvyvLNMJ2xQYQ==",
				},
				"compare-same-b.js": map[string]string{
					"sha256": "sha256-hwj4HVJ7OFOzPES8HffZ4IySCiQq7P/+1RT9YQJMAXs=",
					"sha384": "sha384-6yKRlQqq9r5dU0GEivGoDai04RH+ufhfO1htclXkbjdJ+184pc1rRrsWhk3aDf3D",
					"sha512": "sha512-rk6z+muPQjmOe6aC/kKX9yYmaaLs42TIBxh02xwV/vkrIpD8RJzjt4cB4q9MXOeYkwD4wSyEMNvyvLNMJ2xQYQ==",
				},
			},
		},
	}

	for _, tc := range testCases {
		// Test against each of our hash options
		for _, h := range []string{sha256Algo, sha384Algo, sha512Algo, allHashes} {
			fis, err := generate(tc.targets, h)
			if err != nil {
				t.Fatalf("Unexpected error from generate call (targets: %q). %q", tc.targets, err)
			}

			for _, fi := range fis {
				hash := fi.Digest[:6]

				if tc.targets[0] == mockServer.URL && tc.exp[fi.Source][hash] != fi.Digest {
					t.Fatalf("Expected integrity from %s to have %s digest of '%s'. Got '%s'",
						fi.Source, hash, fi.Digest, tc.exp[fi.Source][hash])
				}

				if tc.targets[0] != mockServer.URL && tc.exp[fi.FileName][hash] != fi.Digest {
					t.Fatalf("Expected integrity from %s to have %s digest of '%s'. Got '%s'",
						fi.FileName, hash, fi.Digest, tc.exp[fi.FileName][hash])
				}
			}
		}
	}
}
