package compute

import "fmt"

// RunCommittee seals researcher, then challenger, then risk. Any role fail stops the operation.
func RunCommittee(bin string, jobs []DirectJob) error {
	return RunCommitteeStages(bin, jobs, nil, nil)
}

func RunCommitteeStages(bin string, jobs []DirectJob, stage StageFn, stop func() bool) error {
	if len(jobs) != 3 {
		return fmt.Errorf("committee_incomplete")
	}
	if err := MustNativeSealer(bin); err != nil {
		return err
	}
	seen := map[Role]struct{}{}
	for i := range jobs {
		jobs[i].Bin = bin
		switch jobs[i].Role {
		case Researcher, Challenger, Risk:
		default:
			return fmt.Errorf("bad_role")
		}
		if _, ok := seen[jobs[i].Role]; ok {
			return fmt.Errorf("duplicate_role")
		}
		seen[jobs[i].Role] = struct{}{}
	}
	if _, ok := seen[Researcher]; !ok {
		return fmt.Errorf("committee_incomplete")
	}
	if _, ok := seen[Challenger]; !ok {
		return fmt.Errorf("committee_incomplete")
	}
	if _, ok := seen[Risk]; !ok {
		return fmt.Errorf("committee_incomplete")
	}
	labels := map[Role]string{
		Researcher: "RESEARCHER",
		Challenger: "CHALLENGER",
		Risk:       "RISK",
	}
	for i := range jobs {
		if stopped(stop) {
			return fmt.Errorf("research_cancelled")
		}
		if err := RunSealedAskCtl(jobs[i], stage, stop); err != nil {
			return err
		}
		notify(stage, labels[jobs[i].Role])
	}
	return nil
}

func CommitteeRoles() []Role {
	return []Role{Researcher, Challenger, Risk}
}
