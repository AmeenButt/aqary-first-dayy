package models

type UserModel struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	ProfilePicture string `json:"profile_picture"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type UserWallet struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
	Amount    float64   `json:"amount"`
	User      UserModel `json:"user"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type UserTransaction struct {
	ID                int32      `json:"id"`
	UserWalletID      int32      `json:"user_wallet_id"`
	TransactionAmount float64    `json:"amount"`
	UserWalletData    UserWallet `json:"user_wallet_data"`
	CreatedAt         string     `json:"created_at"`
	UpdatedAt         string     `json:"updated_at"`
}
