package cli

import (
	"encoding/json"
	"fmt"
	"io"
)

func WantJSON(args []string) (bool, []string) {
	out := make([]string, 0, len(args))
	want := false
	for _, a := range args {
		if a == "--json" {
			want = true
			continue
		}
		out = append(out, a)
	}
	return want, out
}

func Emit(w io.Writer, asJSON bool, human string, obj any) {
	if asJSON {
		_ = json.NewEncoder(w).Encode(obj)
		return
	}
	fmt.Fprint(w, human)
}
