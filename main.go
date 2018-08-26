package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

type FileIntegrity struct {
	Digest   string `json:"digest"`
	FileName string `json:"file"`
	Tag      string `json:"tag"`
	Source   string `json:"source,omitempty"`
}

var (
	out    = flag.String("out", "sri.json", "File output")
	client *http.Client
)

func main() {
	flag.Parse()
	target := os.Args[1]
	client = &http.Client{Timeout: time.Second * 2}

	run(target, *out)
}

func run(target, out string) {
	var fileIntegrities []FileIntegrity
	var err error
	if _, err := url.ParseRequestURI(target); err == nil {
		fileIntegrities, err = handleDownload(target)
	} else if fi, err := os.Stat(target); err != nil || fi.Mode().IsRegular() {
		fileIntegrities, err = handleFile(target)
	} else {
		fileIntegrities, err = handleDir(target)
	}

	if err != nil {
		log.Fatal(err)
	}

	writeOutputToFile(out, fileIntegrities)
}

func handleDownload(target string) ([]FileIntegrity, error) {
	resp, err := client.Get(target)
	if err != nil {
		return nil, fmt.Errorf("Failure downloading script from %s. %s", target, err)
	}
	defer resp.Body.Close()

	return generateFileIntegrities(target, resp.Body)
}

func handleFile(target string) ([]FileIntegrity, error) {
	f, err := os.Open(target)
	if err != nil {
		return nil, err
	}

	return generateFileIntegrities(target, f)
}

func handleDir(target string) ([]FileIntegrity, error) {
	dir, err := ioutil.ReadDir(target)
	if err != nil {
		return nil, err
	}

	var outerErr error
	fisChan := make(chan []FileIntegrity, len(dir))
	for _, fi := range dir {
		go func(ifi os.FileInfo) {
			fi, err := handleFile(target + "/" + ifi.Name())
			if err != nil && outerErr == nil {
				outerErr = err
			}

			fisChan <- fi
		}(fi)
	}

	if outerErr != nil {
		return nil, err
	}

	combined := []FileIntegrity{}
	for i := 0; i < len(dir); i++ {
		combined = append(combined, <-fisChan...)
	}

	return combined, nil
}

func generateFileIntegrities(source string, r io.Reader) ([]FileIntegrity, error) {
	hashes := []io.Writer{
		sha256.New(),
		sha512.New384(),
		sha512.New(),
	}

	multiHasher := io.MultiWriter(hashes...)
	if _, err := io.Copy(multiHasher, r); err != nil {
		return nil, err
	}

	fiChan := make(chan string, len(hashes))
	for _, h := range hashes {
		go func(h hash.Hash) {
			fiChan <- fmt.Sprintf("sha%d-%s", h.Size()*8, base64.StdEncoding.EncodeToString(h.Sum(nil)))
		}(h.(hash.Hash))
	}

	fis := []FileIntegrity{}
	for i := 0; i < len(hashes); i++ {
		digest := <-fiChan
		fi := FileIntegrity{
			Digest:   digest,
			FileName: path.Base(source),
			Tag:      generateTag(source, digest),
		}

		if _, err := url.ParseRequestURI(source); err != nil {
			fi.Source = source
		}

		fis = append(fis, fi)
	}

	return fis, nil
}

func generateTag(source, digest string) string {
	if path.Ext(source) == ".css" {
		return fmt.Sprintf(`<link rel='stylesheet' href='%s' integrity='%s'>`, source, digest)
	}

	return fmt.Sprintf(`<script src='%s' integrity='%s'></script>`, source, digest)
}

func writeOutputToFile(path string, fis []FileIntegrity) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Unable to create file at location: %s. %s", path, err)
	}

	j := json.NewEncoder(f)
	j.SetEscapeHTML(false)
	j.SetIndent("", "\t")

	jsonFIs := constructJSONFileIntegrities(fis)
	j.Encode(jsonFIs)

	return nil
}

func constructJSONFileIntegrities(fis []FileIntegrity) map[string]interface{} {
	result := make(map[string]interface{})

	for _, fi := range fis {
		if result[fi.FileName] == nil {
			result[fi.FileName] = make(map[string]interface{})
		}

		fileNode := result[fi.FileName].(map[string]interface{})
		algo := strings.Split(fi.Digest, "-")[0]
		fileNode[algo] = make(map[string]interface{})
		algoNode := fileNode[algo].(map[string]interface{})

		algoNode["digest"] = fi.Digest
		algoNode["tag"] = fi.Tag
		if fi.Source != "" {
			algoNode["source"] = fi.Source
		}
	}

	return result
}
