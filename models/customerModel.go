package models

import (
	"fmt"
	"project/token"

	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null;" json:"password"`
	Name     string `gorm:"size:255;not null;" json:"name"`
	AuthUUID string `gorm:"size:255;not null;" json:"auth_uuid"`
	Balance  uint   `gorm:"size:255;not null;" json:"balance"`
}

func (c *Customer) SaveUser() (*Customer, error) {

	var err error = DB.Create(c).Error
	if err != nil {
		return &Customer{}, err
	}

	return c, nil
}

func (c *Customer) BeforeSave(tx *gorm.DB) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	c.Password = string(hashedPassword)
	return nil

}

func LoginCheck(username string, password string) (string, Customer, error) {

	var err error

	c := Customer{}

	err = DB.Model(Customer{}).Where("username = ?", username).Take(&c).Error

	if err != nil {
		message := fmt.Sprintf("Customer dengan %s tidak ditemukan", username)
		return message, c, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password))

	if err != nil {
		return "Password yang di inputkan salah", c, err
	}

	authUUID, err := CreateAuth(c.Username)

	if err != nil {
		return "Gagal create Auth", c, err
	}

	token, err := token.GenerateToken(c.ID, authUUID)

	if err != nil {
		return "", c, err
	}

	return token, c, nil
}

func CreateAuth(username string) (string, error) {
	c := Customer{}

	c.AuthUUID = uuid.NewV4().String()

	err := DB.Model(&Customer{}).Where("username = ?", username).Update("auth_uuid", c.AuthUUID).Error

	if err != nil {
		return "", err
	}
	return string(c.AuthUUID), nil
}

func CheckAuth(authDetails *token.AuthDetails) error {
	c := Customer{}

	err := DB.First(&Customer{}, "auth_uuid = ?", authDetails.AuthUuid).Take(&c).Error

	if err != nil {
		return err
	}
	return nil
}

func DeleteAuth(authDetails *token.AuthDetails) (string, error) {
	var c Customer

	err := DB.First(&c, authDetails.UserId).Update("auth_uuid", "").Error

	if err != nil {
		return "", err
	}
	return string(c.Name), nil
}

func GetBalancebyId(custId uint) (Customer, error) {
	var err error
	c := Customer{}

	err = DB.Model(Customer{}).Where("id = ?", custId).Take(&c).Error

	if err != nil {

		return c, err
	}

	return c, nil
}

func GetBalancebyUsername(username string) (Customer, error) {
	var err error
	c := Customer{}

	err = DB.Model(Customer{}).Where("username = ?", username).Take(&c).Error

	if err != nil {

		return c, err
	}

	return c, nil
}
