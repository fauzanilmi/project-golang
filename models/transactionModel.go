package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	CustomerId  uint   `gorm:"size:255;not null;" json:"customer_id"`
	Amount      uint   `gorm:"size:255;not null;" json:"amount"`
	Balance     uint   `gorm:"size:255;not null;" json:"balance"`
	Description string `gorm:"size:255;not null;" json:"denscription"`
}

func (tx *Transaction) UpdateBalance() (*Transaction, error) {

	var topupErr error = DB.Create(tx).Error
	if topupErr != nil {
		return &Transaction{}, topupErr
	}

	var updateBalanceErr error = DB.Model(&Customer{}).Where("id = ?", tx.CustomerId).Update("balance", tx.Balance).Error
	if updateBalanceErr != nil {
		return &Transaction{}, updateBalanceErr
	}

	return tx, nil
}
