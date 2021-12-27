package account

type Account struct {
	ID       int64  `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}
