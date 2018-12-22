package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"flag"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	sha256Algo = "sha256"
	sha384Algo = "sha384"
	sha512Algo = "sha512"
	allHashes  = "all"
)

var (
	compare = flag.Bool("compare", false, "Run in comparison mode")

	hashAlgo = flag.String("hash", "sha256", "Hashing algorithm")
	outPath  = flag.String("out", "", "Name of output file")

	hashes = map[string]func() hash.Hash{
		sha256Algo: func() hash.Hash { return sha256.New() },
		sha384Algo: func() hash.Hash { return sha512.New384() },
		sha512Algo: func() hash.Hash { return sha512.New() },
	}

	client = &http.Client{Timeout: time.Second * 2}
)

func main() {
	flag.Parse()

	if err := validateHash(*hashAlgo); err != nil {
		log.Fatalf("[sri] Invalid value for flag '-hash'. %q", err)
	}

	// If comparison flag specified, run a comparison between the two targets and exit(0) if match or exit(1) if
	// the digests differ. Digests are also printed to stdout in both cases.
	if *compare {
		if err := validateCompare(flag.Args()); err != nil {
			log.Fatalf("[sri] Unable to perform comparison. %q", err)
		}

		result, a, b, err := comparison(flag.Arg(0), flag.Arg(1))
		if err != nil {
			log.Fatalf("[sri] An error occured during comparison. %q", err)
		}

		fmt.Printf("%s - %s\n", flag.Arg(0), a)
		fmt.Printf("%s - %s\n", flag.Arg(1), b)

		if !result {
			fmt.Println("Digests did not match")
			os.Exit(1)
		} else {
			fmt.Println("Digests match")
			os.Exit(0)
		}
	}

	// If we aren't in comparison mode, we are in 'generate' mode, and will attempt to produce
	// SRIs for our given target and write them to either stdout or a file (if an outfile was provided)
	if err := validateGenerate(flag.Args()); err != nil {
		log.Fatalf("[sri] Unable to generate SRI output. %q", err)
	}

	fis, err := generate(flag.Arg(0), *hashAlgo)
	if err != nil {
		log.Fatalf("[sri] An error occured to generating SRI output. %q", err)
	}

	if *outPath != "" {
		writeOutputToFile(fis, *outPath)
	} else {
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")

		if err := enc.Encode(fis); err != nil {
			log.Fatalf("An error occured writing SRIs to stdout")
		}
	}

	os.Exit(0)
}

func generate(target, hashName string) ([]fileIntegrity, error) {
	var fileIntegrities []fileIntegrity
	var err error
	if _, err := url.ParseRequestURI(target); err == nil {
		fileIntegrities, err = handleDownload(target, hashName)
	} else if fi, err := os.Stat(target); err == nil && fi != nil && fi.Size() > 0 && fi.Mode().IsRegular() {
		fileIntegrities, err = handleFile(target, hashName)
	} else {
		fileIntegrities, err = handleDir(target, hashName)
	}

	if err != nil {
		return nil, err
	}

	if fileIntegrities == nil || len(fileIntegrities) == 0 {
		return nil, fmt.Errorf("No file integrities generated from target '%s'", target)
	}

	return fileIntegrities, nil
}

func handleDownload(target, hashName string) ([]fileIntegrity, error) {
	resp, err := client.Get(target)
	if err != nil {
		return nil, fmt.Errorf("Failure downloading script from %s. %s", target, err)
	}
	defer resp.Body.Close()

	return generateFileIntegrities(target, hashName, resp.Body)
}

func handleFile(target, hashName string) ([]fileIntegrity, error) {
	f, err := os.Open(target)
	if err != nil {
		return nil, err
	}

	return generateFileIntegrities(target, hashName, f)
}

func handleDir(target, hashName string) ([]fileIntegrity, error) {
	dir, err := ioutil.ReadDir(target)
	if err != nil {
		return nil, err
	}

	var outerErr error
	fisChan := make(chan []fileIntegrity, len(dir))
	for _, fi := range dir {
		go func(ifi os.FileInfo) {
			fi, err := handleFile(target+"/"+ifi.Name(), hashName)
			if err != nil && outerErr == nil {
				outerErr = err
			}

			fisChan <- fi
		}(fi)
	}

	if outerErr != nil {
		return nil, err
	}

	combined := []fileIntegrity{}
	for i := 0; i < len(dir); i++ {
		combined = append(combined, <-fisChan...)
	}

	return combined, nil
}

func validateHash(hashName string) error {
	v, ok := hashes[hashName]
	if !ok || v == nil {
		return fmt.Errorf("Invalid hashing algorithm '%s'. Expected one of 'sha256', 'sha384' or 'sha512'", hashName)
	}

	return nil
}

func validateGenerate(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected only a single target for SRI generation, received %d", len(args))
	}

	if args[0] == "" {
		return fmt.Errorf("Received an empty target for SRI generation")
	}

	return nil
}
