package main

import "fmt"

// comparison runs a sha256 comparison against the two provided targets, returning the equality of their hashes,
// as well as their individual digests and any resulting errors.
func comparison(a, b string) (bool, string, string, error) {
	afis, err := generate(a, sha256Algo)
	if err != nil {
		return false, "", "", err
	}

	if len(afis) == 0 {
		return false, "", "", fmt.Errorf("Unable to produce integrities for %s", a)
	}

	bfis, err := generate(b, sha256Algo)
	if err != nil {
		return false, "", "", err
	}

	if len(bfis) == 0 {
		return false, "", "", fmt.Errorf("Unable to produce integrities for %s", b)
	}

	return afis[0].Digest == bfis[0].Digest, afis[0].Digest, bfis[0].Digest, nil
}

func validateCompare(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("Expected two targets to be specified for comparison")
	}

	if args[0] == "" || args[1] == "" {
		return fmt.Errorf("Received an empty target for comparison")
	}

	if args[0] == args[1] {
		return fmt.Errorf("Received two indentical inputs for comparison")
	}

	return nil
}
