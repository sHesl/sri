package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"net/url"
	"path"
)

type fileIntegrity struct {
	Digest   string `json:"digest"`
	FileName string `json:"file"`
	Tag      string `json:"tag"`
	Source   string `json:"source,omitempty"`
}

func generateFileIntegrities(source, hashAlgo string, r io.Reader) ([]fileIntegrity, error) {
	var hs []io.Writer
	if hashAlgo == allHashes {
		hs = []io.Writer{
			sha256.New(),
			sha512.New384(),
			sha512.New(),
		}
	} else {
		hs = []io.Writer{hashes[hashAlgo]()}
	}

	multiHasher := io.MultiWriter(hs...)
	if _, err := io.Copy(multiHasher, r); err != nil {
		return nil, err
	}

	fiChan := make(chan string, len(hs))
	for _, h := range hs {
		go func(h hash.Hash) {
			fiChan <- fmt.Sprintf("sha%d-%s", h.Size()*8, base64.StdEncoding.EncodeToString(h.Sum(nil)))
		}(h.(hash.Hash))
	}

	fis := []fileIntegrity{}
	for i := 0; i < len(hs); i++ {
		digest := <-fiChan
		fi := fileIntegrity{
			Digest:   digest,
			FileName: path.Base(source),
			Tag:      generateTag(source, digest),
		}

		if _, err := url.ParseRequestURI(source); err == nil {
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
