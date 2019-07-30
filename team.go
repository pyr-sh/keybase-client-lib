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

// AddUser adds a member to a team by username
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
	if err == nil && r.Error == nil {
		r, err = t.MemberList()
	}
	return r, err
}

// AddReaders adds members to a team by username, and sets their roles to Reader
func (t Team) AddReaders(users ...string) (TeamAPI, error) {
	m := TeamAPI{
		Params: &tParams{},
	}
	m.Method = "add-members"
	m.Params.Options.Team = t.Name
	addUsers := []usernames{}
	for _, u := range users {
		addUsers = append(addUsers, usernames{Username: u, Role: "reader"})
	}
	m.Params.Options.Usernames = addUsers

	r, err := teamAPIOut(t.keybase.Path, m)
	if err == nil && r.Error == nil {
		r, err = t.MemberList()
	}
	return r, err
}

// AddWriters adds members to a team by username, and sets their roles to Writer
func (t Team) AddWriters(users ...string) (TeamAPI, error) {
	m := TeamAPI{
		Params: &tParams{},
	}
	m.Method = "add-members"
	m.Params.Options.Team = t.Name
	addUsers := []usernames{}
	for _, u := range users {
		addUsers = append(addUsers, usernames{Username: u, Role: "writer"})
	}
	m.Params.Options.Usernames = addUsers

	r, err := teamAPIOut(t.keybase.Path, m)
	if err == nil && r.Error == nil {
		r, err = t.MemberList()
	}
	return r, err
}

// AddAdmins adds members to a team by username, and sets their roles to Writer
func (t Team) AddAdmins(users ...string) (TeamAPI, error) {
	m := TeamAPI{
		Params: &tParams{},
	}
	m.Method = "add-members"
	m.Params.Options.Team = t.Name
	addUsers := []usernames{}
	for _, u := range users {
		addUsers = append(addUsers, usernames{Username: u, Role: "admin"})
	}
	m.Params.Options.Usernames = addUsers

	r, err := teamAPIOut(t.keybase.Path, m)
	if err == nil && r.Error == nil {
		r, err = t.MemberList()
	}
	return r, err
}

// AddOwners adds members to a team by username, and sets their roles to Writer
func (t Team) AddOwners(users ...string) (TeamAPI, error) {
	m := TeamAPI{
		Params: &tParams{},
	}
	m.Method = "add-members"
	m.Params.Options.Team = t.Name
	addUsers := []usernames{}
	for _, u := range users {
		addUsers = append(addUsers, usernames{Username: u, Role: "owner"})
	}
	m.Params.Options.Usernames = addUsers

	r, err := teamAPIOut(t.keybase.Path, m)
	if err == nil && r.Error == nil {
		r, err = t.MemberList()
	}
	return r, err
}

// MemberList returns a list of a team's members
func (t Team) MemberList() (TeamAPI, error) {
	m := TeamAPI{
		Params: &tParams{},
	}
	m.Method = "list-team-memberships"
	m.Params.Options.Team = t.Name

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
