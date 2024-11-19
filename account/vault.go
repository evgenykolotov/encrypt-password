package account

import (
	"demo/password/encrypter"
	"demo/password/output"
	"encoding/json"
	"strings"
	"time"

	"github.com/fatih/color"
)

type DataBase interface {
	Write([]byte)
	Read() ([]byte, error)
}

type Vault struct {
	Accounts  []Account `json:"accounts"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type VaultWithDataBase struct {
	Vault
	database  DataBase
	encrypter encrypter.Encrypter
}

func NewVault(database DataBase, encrypter encrypter.Encrypter) *VaultWithDataBase {
	file, err := database.Read()
	if err != nil {
		return &VaultWithDataBase{
			Vault: Vault{
				Accounts:  []Account{},
				UpdatedAt: time.Now(),
			},
			database:  database,
			encrypter: encrypter,
		}
	}
	data := encrypter.Decrypt(file)
	var vault Vault
	err = json.Unmarshal(data, &vault)
	color.Cyan("Найдено %d аккаунтов", len(vault.Accounts))
	if err != nil {
		output.PrintError("Не удалось разобрать файл data.vault")
		return &VaultWithDataBase{
			Vault: Vault{
				Accounts:  []Account{},
				UpdatedAt: time.Now(),
			},
			database:  database,
			encrypter: encrypter,
		}
	}
	return &VaultWithDataBase{
		Vault:     vault,
		database:  database,
		encrypter: encrypter,
	}
}

func (vault *VaultWithDataBase) AddAccount(acc *Account) {
	vault.Accounts = append(vault.Accounts, *acc)
	vault.save()
}

func (vault *Vault) ToBytes() ([]byte, error) {
	file, err := json.Marshal(vault)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (vault *VaultWithDataBase) FindAccounts(str string, checker func(Account, string) bool) []Account {
	var accounts []Account
	for _, account := range vault.Accounts {
		isMatched := checker(account, str)
		if isMatched {
			accounts = append(accounts, account)
		}
	}
	return accounts
}

func (vault *VaultWithDataBase) DeleteAccountByUrl(url string) bool {
	isDeleted := false
	var accounts []Account
	for _, account := range vault.Accounts {
		isMatched := strings.Contains(account.Url, url)
		if !isMatched {
			accounts = append(accounts, account)
			continue
		}
		isDeleted = true
	}
	vault.Accounts = accounts
	vault.save()
	return isDeleted
}

func (vault *VaultWithDataBase) save() {
	vault.UpdatedAt = time.Now()
	data, err := vault.Vault.ToBytes()
	encData := vault.encrypter.Encrypt(data)
	if err != nil {
		output.PrintError("Не удалось преобразовать файл data.json")
	}
	vault.database.Write(encData)
}
