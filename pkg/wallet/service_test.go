package wallet

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/shuhrat-shokirov/wallet/pkg/types"
)

func TestService_Register(t *testing.T) {
	svc := Service{}
	_, err := svc.RegisterAccount("+992000000000")
	if err != nil {
		t.Error(err)
	}

	_, err = svc.RegisterAccount("+992000000000")
	if err != ErrPhoneNumberRegistred {
		t.Error(err)
	}
}

func TestService_Deposit(t *testing.T) {
	svc := Service{}
	err := svc.Deposit(1, 0)
	if err != ErrAmountMustBePositive {
		t.Error(err)
	}

	err = svc.Deposit(1, 1)
	if err != ErrAccountNotFound {
		t.Error(err)
	}

	account, err := svc.RegisterAccount("+992000000010")
	if err != nil {
		t.Error(err)
	}

	err = svc.Deposit(account.ID, 1)
	if err != nil {
		t.Error(err)
	}
}

func TestService_Pay(t *testing.T) {
	svc := Service{}
	_, err := svc.Pay(1, 0, "auto")
	if err != ErrAmountMustBePositive {
		t.Error(err)
	}

	_, err = svc.Pay(1, 1, "auto")
	if err != ErrAccountNotFound {
		t.Error(err)
	}

	account, err := svc.RegisterAccount("+992000000000")
	if err != nil {
		t.Error(err)
	}

	_, err = svc.Pay(account.ID, 1, "auto")
	if err != ErrNotEnoughBalance {
		t.Error(err)
	}

	account.Balance = 100

	_, err = svc.Pay(account.ID, 100, "auto")
	if err != nil {
		t.Error(err)
	}
}

func TestService_FindbyAccountById_success(t *testing.T) {
	svc := Service{}
	account, err := svc.RegisterAccount("+992000000000")
	if err != nil {
		t.Error(err)
	}

	_, err = svc.FindAccountByID(account.ID)
	if err != nil {
		t.Error(err)
	}
}

func TestService_FindByAccountByID_notFound(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+992000000000")
	_, err := svc.FindAccountByID(2)
	// тут даст false, так как err (уже имеет что то внутри)
	if err != ErrAccountNotFound {
		t.Error(err)
	}
}

func TestFindPaymentByID_success(t *testing.T) {
	svc := &Service{}

	phone := types.Phone("+992000000000")

	account, err := svc.RegisterAccount(phone)
	if err != nil {
		t.Error(err)
		return
	}

	err = svc.Deposit(account.ID, 1000)
	if err != nil {
		t.Error(err)
		return
	}

	pay, err := svc.Pay(account.ID, 500, "auto")
	if err != nil {
		t.Error(err)
		return
	}

	got, err := svc.FindPaymentByID(pay.ID)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(got, pay) {
		t.Error(err)
		return
	}
}

func TestService_Reject_success(t *testing.T) {
	svc := &Service{}

	phone := types.Phone("+992000000000")

	account, err := svc.RegisterAccount(phone)
	if err != nil {
		t.Error(err)
		return
	}

	err = svc.Deposit(account.ID, 1000)
	if err != nil {
		t.Error(err)
		return
	}

	pay, err := svc.Pay(account.ID, 500, "auto")
	if err != nil {
		t.Error(err)
		return
	}

	err = svc.Reject(pay.ID)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_Reject_fail(t *testing.T) {
	svc := Service{}

	svc.RegisterAccount("+992000000000")

	account, err := svc.FindAccountByID(1)
	if err != nil {
		t.Error(err)
	}

	err = svc.Deposit(account.ID, 1000_00)
	if err != nil {
		t.Error(err)
	}

	payment, err := svc.Pay(account.ID, 100_00, "auto")
	if err != nil {
		t.Error(err)
	}

	pay, err := svc.FindPaymentByID(payment.ID)
	if err != nil {
		t.Error(pay)
	}

	editPayID := "4"

	err = svc.Reject(editPayID)
	if err != ErrPaymentNotFound {
		t.Error(err)
	}
}

func TestService_Repeat_success(t *testing.T) {
	svc := &Service{}

	phone := types.Phone("+992000000000")

	account, err := svc.RegisterAccount(phone)
	if err != nil {
		t.Error(err)
		return
	}

	err = svc.Deposit(account.ID, 1000)
	if err != nil {
		t.Error(err)
		return
	}

	pay, err := svc.Pay(account.ID, 500, "auto")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = svc.Repeat(pay.ID)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_FavoritePayment_success(t *testing.T) {
	svc := &Service{}

	phone := types.Phone("+992000000000")

	account, err := svc.RegisterAccount(phone)
	if err != nil {
		t.Error(err)
		return
	}

	err = svc.Deposit(account.ID, 1000)
	if err != nil {
		t.Error(err)
		return
	}

	pay, err := svc.Pay(account.ID, 500, "auto")
	if err != nil {
		t.Error(err)
		return
	}

	favorite, err := svc.FavoritePayment(pay.ID, "pay")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(favorite)
}

func TestService_PayFromFavorite_success(t *testing.T) {
	svc := &Service{}

	phone := types.Phone("+992000000000")

	account, err := svc.RegisterAccount(phone)
	if err != nil {
		t.Error(err)
		return
	}

	err = svc.Deposit(account.ID, 1000)
	if err != nil {
		t.Error(err)
		return
	}

	pay, err := svc.Pay(account.ID, 500, "auto")
	if err != nil {
		t.Error(err)
		return
	}

	favorite, err := svc.FavoritePayment(pay.ID, "pay")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = svc.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_ExportToFile_EmptyData(t *testing.T) {
	svc := &Service{}

	err := svc.ExportToFile("1.txt")
	if err != nil {
		t.Error(err)
	}
	file, err := os.Open("1.txt")
	if err != nil {
		t.Error(err)
	}

	stats, err := file.Stat()
	if err != nil {
		t.Error(err)
	}

	if stats.Size() != 0 {
		t.Error("file must be zero")
	}
}

func TestService_ExportToFile(t *testing.T) {
	svc := &Service{}

	_, err := svc.RegisterAccount("+992000000000")
	if err != nil {
		t.Error(err)
	}

	err = svc.ExportToFile("1.txt")
	if err != nil {
		t.Error(err)
	}
	file, err := os.Open("1.txt")
	if err != nil {
		t.Error(err)
	}

	stats, err := file.Stat()
	if err != nil {
		t.Error(err)
	}

	if stats.Size() == 0 {
		t.Error("file must be zero")
	}
}

func TestService_ImportToFile(t *testing.T) {
	svc := &Service{}

	err := svc.ImportFromFile("1.txt")
	if err != nil {
		t.Error(err)
	}

	k := 0
	for _, account := range svc.accounts {
		if account.Phone == "+992000000000" {
			k++
		}
	}

	if k <= 0 {
		t.Error("incorrect func")
	}
}

func TestSetice_Export(t *testing.T) {
	svc := &Service{}

	account, err := svc.RegisterAccount("+992000000000")
	if err != nil {
		t.Error(err)
	}

	account.Balance = 100

	payment, err := svc.Pay(account.ID, 100, "auto")
	if err != nil {
		t.Error(err)
	}

	_, err = svc.FavoritePayment(payment.ID, "isbraniy")
	if err != nil {
		t.Error(err)
	}

	err = svc.Export(".")
	if err != nil {
		t.Error(err)
	}

	_, err = ioutil.ReadFile("accounts.dump")
	if err != nil {
		t.Error(err)
	}

	_, err = ioutil.ReadFile("payments.dump")
	if err != nil {
		t.Error(err)
	}

	_, err = ioutil.ReadFile("favorites.dump")
	if err != nil {
		t.Error(err)
	}
}

func TestService_Import(t *testing.T) {
	svc := &Service{}

	err := svc.Import(".")
	if err != nil {
		t.Error(err)
	}

	if svc.accounts[0].Phone != "+992000000000" {
		t.Error("incorrect func")
	}
}
func TestService_Import_IfHaveData(t *testing.T) {
	svc := &Service{}

	account, err := svc.RegisterAccount("+992000000000")
	if err != nil {
		t.Error(err)
	}

	account.Balance = 100

	payment, err := svc.Pay(account.ID, 100, "auto")
	if err != nil {
		t.Error(err)
	}

	_, err = svc.FavoritePayment(payment.ID, "isbraniy")
	if err != nil {
		t.Error(err)
	}

	err = svc.Import(".")
	if err != nil {
		t.Error(err)
	}

	if account.Phone == "+992" {
		t.Error("incorrect func")
	}
}
