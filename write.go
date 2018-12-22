package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func writeOutputToFile(fis []fileIntegrity, outPath string) error {
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("Unable to create file at location: %s. %s", outPath, err)
	}

	j := json.NewEncoder(f)
	j.SetEscapeHTML(false)
	j.SetIndent("", "\t")

	jsonFIs := constructJSONFileIntegrities(fis)
	j.Encode(jsonFIs)

	return nil
}

func constructJSONFileIntegrities(fis []fileIntegrity) map[string]interface{} {
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
