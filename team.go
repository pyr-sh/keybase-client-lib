package keybase

import (
	"encoding/json"
	"errors"
	"fmt"
)

// teamAPIOut sends JSON requests to the team API and returns its response.
func teamAPIOut(k *Keybase, t TeamAPI) (TeamAPI, error) {
	jsonBytes, _ := json.Marshal(t)

	cmdOut, err := k.Exec("team", "api", "-m", string(jsonBytes))
	if err != nil {
		return TeamAPI{}, err
	}

	var r TeamAPI
	if err := json.Unmarshal(cmdOut, &r); err != nil {
		return TeamAPI{}, err
	}
	if r.Error != nil {
		return TeamAPI{}, errors.New(r.Error.Message)
	}

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

	r, err := teamAPIOut(t.keybase, m)
	if err == nil && r.Error == nil {
		r, err = t.MemberList()
	}
	return r, err
}

// RemoveUser removes a member from a team
func (t Team) RemoveUser(user string) (TeamAPI, error) {
	m := TeamAPI{
		Params: &tParams{},
	}
	m.Method = "remove-member"
	m.Params.Options.Team = t.Name
	m.Params.Options.Username = user

	r, err := teamAPIOut(t.keybase, m)
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

	r, err := teamAPIOut(t.keybase, m)
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

	r, err := teamAPIOut(t.keybase, m)
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

	r, err := teamAPIOut(t.keybase, m)
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

	r, err := teamAPIOut(t.keybase, m)
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

	r, err := teamAPIOut(t.keybase, m)
	return r, err
}

// CreateSubteam creates a subteam
func (t Team) CreateSubteam(name string) (TeamAPI, error) {
	m := TeamAPI{
		Params: &tParams{},
	}
	m.Method = "create-team"
	m.Params.Options.Team = fmt.Sprintf("%s.%s", t.Name, name)

	r, err := teamAPIOut(t.keybase, m)
	return r, err
}

// CreateTeam creates a new team
func (k *Keybase) CreateTeam(name string) (TeamAPI, error) {
	m := TeamAPI{
		Params: &tParams{},
	}
	m.Method = "create-team"
	m.Params.Options.Team = name

	r, err := teamAPIOut(k, m)
	return r, err
}
