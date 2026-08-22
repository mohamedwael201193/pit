package chain8004

import "fmt"

func PortableAcrossNetworks(from, to string) error {
	if from == "" || to == "" {
		return fmt.Errorf("missing_network")
	}
	if from != to {
		return fmt.Errorf("ids_not_portable")
	}
	return nil
}
