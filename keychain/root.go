package keychain

import (
	"errors"

	"github.com/keybase/go-keychain"
)

const serviceName = "bpass.cli"
const accountName = "master-password"

func StoreMasterPassword(masterPassword []byte) error {
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(serviceName)
	item.SetAccount(accountName)
	item.SetLabel("bpass cli master password")
	item.SetAccessGroup(serviceName)
	item.SetData(masterPassword)
	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleAfterFirstUnlock)

	err :=keychain.AddItem(item)
	return err
}

func RetrieveMasterPassword() ([]byte, error) {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(serviceName)
	query.SetAccount(accountName)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)

	item, err :=keychain.QueryItem(query)
	if err != nil {
		return nil, err
	}
	if len(item) == 0 {
		return nil, errors.New("master password not found in keychain")
	}
	return item[0].Data, nil
}

func DeleteMasterPassword() error {
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(serviceName)
	item.SetAccount(accountName)
	item.SetMatchLimit(keychain.MatchLimitOne)
	err := keychain.DeleteItem(item)
	return err
}
