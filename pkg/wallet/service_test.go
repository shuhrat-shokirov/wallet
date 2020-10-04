package wallet

import (
	"testing"
)

func TestService_FindbyAccountById_success(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+9929351007")
	account, err := svc.FindAccountByID(1)
	if err != nil {
		t.Errorf("не удалось найти аккаунт, получили: %v", account)
	}
	
}

func TestService_FindByAccountByID_notFound(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+992938151007")
	account, err := svc.FindAccountByID(2)
	// тут даст false, так как err (уже имеет что то внутри)
	if err != ErrAccountNotFound {
		t.Errorf("аккаунт не найден, аккаунт: %v", account)
	}
}
