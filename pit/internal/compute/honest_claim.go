package compute

import "fmt"

func HonestIndependence(label Independence) error {
	if label == Providers || label == Models {
		return nil
	}
	if label == EnvelopeOnly {
		return nil
	}
	return fmt.Errorf("unknown_independence")
}

func RefuseIndependentClaim(label Independence, claim string) error {
	if claim == "independent_providers" && label != Providers {
		return fmt.Errorf("honesty")
	}
	if claim == "independent_models" && label != Models {
		return fmt.Errorf("honesty")
	}
	return nil
}
