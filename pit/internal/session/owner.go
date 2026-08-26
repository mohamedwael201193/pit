package session

import (
	"fmt"
	"strings"
)

func BindWorkspace(sessionWS, callerWS string) error {
	if strings.TrimSpace(sessionWS) == "" || strings.TrimSpace(callerWS) == "" {
		return fmt.Errorf("wrong_workspace")
	}
	if sessionWS != callerWS {
		return fmt.Errorf("wrong_workspace")
	}
	return nil
}
