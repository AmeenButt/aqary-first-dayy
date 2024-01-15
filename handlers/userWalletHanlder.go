package handlers

import "github.com/jackc/pgx/v5"

type WalletHanlder struct {
	conn *pgx.Conn
}

func GetWalletHanlder(conn *pgx.Conn) *WalletHanlder {
	return &WalletHanlder{conn: conn}
}

// func (wallet)
