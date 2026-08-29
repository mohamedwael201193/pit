package cli

import "github.com/mohamedwael201193/pit/internal/policy"

func ActivePolicy(dir string) policy.Policy {
	st, err := Load(dir)
	if err != nil {
		return policy.Default()
	}
	return policy.Load(dir, st.WorkspaceID)
}
