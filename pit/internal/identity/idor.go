package identity

import "fmt"

func SameWorkspace(got, want string) error {
	a, err := ParseWorkspaceID(got)
	if err != nil {
		return err
	}
	b, err := ParseWorkspaceID(want)
	if err != nil {
		return err
	}
	if a != b {
		return fmt.Errorf("wrong_workspace")
	}
	return nil
}
