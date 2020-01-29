// Auto-generated to Go types using avdl-compiler v1.4.6 (https://github.com/keybase/node-avdl-compiler)
//   Input file: ../client/protocol/avdl/keybase1/teamsearch.avdl

package keybase1

type TeamSearchItem struct {
	Id           TeamID   `codec:"id" json:"id"`
	Name         string   `codec:"name" json:"name"`
	Description  *string  `codec:"description,omitempty" json:"description,omitempty"`
	MemberCount  int      `codec:"memberCount" json:"memberCount"`
	LastActive   Time     `codec:"lastActive" json:"lastActive"`
	InTeam       bool     `codec:"inTeam" json:"inTeam"`
	PublicAdmins []string `codec:"publicAdmins" json:"publicAdmins"`
}

func (o TeamSearchItem) DeepCopy() TeamSearchItem {
	return TeamSearchItem{
		Id:   o.Id.DeepCopy(),
		Name: o.Name,
		Description: (func(x *string) *string {
			if x == nil {
				return nil
			}
			tmp := (*x)
			return &tmp
		})(o.Description),
		MemberCount: o.MemberCount,
		LastActive:  o.LastActive.DeepCopy(),
		InTeam:      o.InTeam,
		PublicAdmins: (func(x []string) []string {
			if x == nil {
				return nil
			}
			ret := make([]string, len(x))
			for i, v := range x {
				vCopy := v
				ret[i] = vCopy
			}
			return ret
		})(o.PublicAdmins),
	}
}

type TeamSearchRes struct {
	Results []TeamSearchItem `codec:"results" json:"results"`
}

func (o TeamSearchRes) DeepCopy() TeamSearchRes {
	return TeamSearchRes{
		Results: (func(x []TeamSearchItem) []TeamSearchItem {
			if x == nil {
				return nil
			}
			ret := make([]TeamSearchItem, len(x))
			for i, v := range x {
				vCopy := v.DeepCopy()
				ret[i] = vCopy
			}
			return ret
		})(o.Results),
	}
}
