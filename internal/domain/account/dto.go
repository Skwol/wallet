package account

type AccountDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

func acountToDTO(account *Account) *AccountDTO {
	return &AccountDTO{
		ID:       account.ID,
		Username: account.Username,
	}
}

func accountsToDTO(accounts []*Account) []*AccountDTO {
	result := make([]*AccountDTO, len(accounts))
	for i, account := range accounts {
		result[i] = acountToDTO(account)
	}
	return result
}

type CreateAccountDTO struct {
	Username string `json:"username"`
}
