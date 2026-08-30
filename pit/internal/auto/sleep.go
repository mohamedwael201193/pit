package auto

import "strings"

const (
	SleepDraft           = "DRAFT"
	SleepReady           = "READY"
	SleepConfirming      = "CONFIRMING"
	SleepArmed           = "ARMED"
	SleepWatching        = "WATCHING"
	SleepResearching     = "RESEARCHING"
	SleepEvaluating       = "EVALUATING"
	SleepReadyToExecute  = "READY_TO_EXECUTE"
	SleepExecuting       = "EXECUTING"
	SleepReconciling     = "RECONCILING"
	SleepProved          = "PROVED"
	SleepLearned         = "LEARNED"
	SleepStopped         = "STOPPED"
	SleepExpired         = "EXPIRED"
	SleepBlocked         = "BLOCKED"
	SleepFailed          = "FAILED"
	SleepRecovering      = "RECOVERING"
)

func SleepFromStage(stage, life, lastStop string) string {
	if life == "BLOCKED" {
		return SleepBlocked
	}
	s := strings.ToLower(strings.TrimSpace(stage))
	switch s {
	case "recovering":
		return SleepRecovering
	case "researching":
		return SleepResearching
	case "ranked", "evaluating":
		return SleepEvaluating
	case "eligible", "ready_eligible":
		return SleepReadyToExecute
	case "executing":
		return SleepExecuting
	case "executed", "resting", "reconciling":
		return SleepReconciling
	case "proved":
		return SleepProved
	case "learned":
		return SleepLearned
	case "exec_failed", "research_failed":
		return SleepFailed
	case "stopped":
		return SleepStopped
	case "starting", "scanning", "waiting", "empty", "searching", "execution-blocked", "exec blocked":
		return SleepWatching
	}
	switch lastStop {
	case "deadline", "autonomy_expired":
		return SleepExpired
	case "kill_switch":
		return SleepBlocked
	case "user_stop", "chat_stop":
		return SleepStopped
	}
	if life == "ACTIVE" {
		return SleepWatching
	}
	if life == "STOPPED" {
		return SleepStopped
	}
	if life == "READY" {
		return SleepReady
	}
	return SleepDraft
}

func SleepLoopAfterLearn() string {
	return SleepWatching
}
