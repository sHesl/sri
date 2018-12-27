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
	"sort"
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
		allHashes:  func() hash.Hash { return nil },
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

	fis, err := generate(flag.Args(), *hashAlgo)
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

func generate(targets []string, hashName string) ([]fileIntegrity, error) {
	var outerErr error
	fisChan := make(chan []fileIntegrity, len(targets))

	for _, target := range targets {
		go func(target string) {
			var fis []fileIntegrity
			if _, err := url.ParseRequestURI(target); err == nil {
				fis, outerErr = handleDownload(target, hashName)
			} else if fi, err := os.Stat(target); err == nil && fi != nil && fi.Size() > 0 && fi.Mode().IsRegular() {
				fis, outerErr = handleFile(target, hashName)
			} else {
				fis, outerErr = handleDir(target, hashName)
			}

			fisChan <- fis
		}(target)
	}

	if outerErr != nil {
		return nil, outerErr
	}

	combined := integrities{}
	for i := 0; i < len(targets); i++ {
		combined = append(combined, <-fisChan...)
	}

	if combined == nil || len(combined) == 0 {
		return nil, fmt.Errorf("No file integrities generated from targets '%q'", targets)
	}

	sort.Sort(combined)

	return combined, nil
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
		return nil, outerErr
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
		return fmt.Errorf("Invalid hashing algorithm '%s'. Expected one of 'sha256', 'sha384', 'sha512' or 'all'", hashName)
	}

	return nil
}

func validateGenerate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("No target specified for SRI generation")
	}

	if args[0] == "" {
		return fmt.Errorf("Received an empty target for SRI generation")
	}

	return nil
}
