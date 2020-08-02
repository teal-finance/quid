package models

// User : base model
type User struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	PasswordHash string  `json:"-"`
	Groups       []Group `json:"groups,omitempty"`
	Namespace    string  `json:"namespace,omitempty"`
}

// GroupNames : get the user group names
func (user User) GroupNames() []string {
	u := []string{}
	for _, g := range user.Groups {
		u = append(u, g.Name)
	}
	return u
}
