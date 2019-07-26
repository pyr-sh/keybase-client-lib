package keybase

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// teamAPIOut sends JSON requests to the team API and returns its response.
func teamAPIOut(keybasePath string, w TeamAPI) (TeamAPI, error) {
	jsonBytes, _ := json.Marshal(w)

	cmd := exec.Command(keybasePath, "team", "api", "-m", string(jsonBytes))
	cmdOut, err := cmd.Output()
	if err != nil {
		return TeamAPI{}, err
	}

	var r TeamAPI
	json.Unmarshal(cmdOut, &r)

	return r, nil
}

// CreateSubteam creates a subteam
func (t Team) CreateSubteam(name string) (TeamAPI, error) {
	m := TeamAPI{
		Params: &tParams{},
	}
	m.Method = "create-team"
	m.Params.Options.Team = fmt.Sprintf("%s.%s", t.Name, name)

	r, err := teamAPIOut(t.keybase.Path, m)
	return r, err
}
