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

// AddUser adds members to a team by username
func (t Team) AddUser(user, role string) (TeamAPI, error) {
	m := TeamAPI{
		Params: &tParams{},
	}
	m.Method = "add-members"
	m.Params.Options.Team = t.Name
	m.Params.Options.Usernames = []usernames{
		{
			Username: user,
			Role:     role,
		},
	}

	r, err := teamAPIOut(t.keybase.Path, m)
	return r, err
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

// CreateTeam creates a new team
func (k *Keybase) CreateTeam(name string) (TeamAPI, error) {
	m := TeamAPI{
		Params: &tParams{},
	}
	m.Method = "create-team"
	m.Params.Options.Team = name

	r, err := teamAPIOut(k.Path, m)
	return r, err
}
