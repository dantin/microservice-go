package model

// Account is the main persist model.
type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ToString return a string that represent `Account`.
func (a *Account) ToString() string {
	return a.ID + " " + a.Name
}
