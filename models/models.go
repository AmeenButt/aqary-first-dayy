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
type Property struct {
	ID           int64     `json:"id"`
	UserId       int64     `json:"user_id"`
	SizeInSqFeet int64     `json:"sizeInSqFeet"`
	Location     string    `json:"location"`
	Demand       string    `json:"demand"`
	Status       string    `json:"status"`
	Images       []string  `json:"images"`
	User         UserModel `json:"user"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

type InputUserModel struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"  binding:"required"`
	Password       string `json:"password"  binding:"required"`
	ProfilePicture string `json:"profile_picture"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type InputUserWallet struct {
	ID        int32   `json:"id"`
	UserID    int32   `json:"user_id"  binding:"required"`
	Amount    float64 `json:"amount"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type InputUserTransaction struct {
	ID                int32   `json:"id"`
	UserWalletID      int32   `json:"user_wallet_id"  binding:"required"`
	TransactionAmount float64 `json:"amount"  binding:"required"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}
type InputProperty struct {
	ID           int64    `json:"id"`
	UserId       int64    `json:"user_id"  binding:"required"`
	SizeInSqFeet int64    `json:"sizeInSqFeet"  binding:"required"`
	Location     string   `json:"location"  binding:"required"`
	Demand       string   `json:"demand"  binding:"required"`
	Status       string   `json:"status"`
	Images       []string `json:"images"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}
