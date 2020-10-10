package wallet

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/google/uuid"

	"github.com/shuhrat-shokirov/wallet/pkg/types"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

var (
	ErrPhoneNumberRegistred = errors.New("phone already registred")
	ErrAmountMustBePositive = errors.New("amount must be greater that zero")
	ErrAccountNotFound      = errors.New("account not found")
	ErrNotEnoughBalance     = errors.New("not enough balance")
	ErrPaymentNotFound      = errors.New("payment not found")
	ErrFavoriteNotFound     = errors.New("favorite not found")
)

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorite      []*types.Favorite
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneNumberRegistred
		}
	}

	s.nextAccountID++

	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}

	s.accounts = append(s.accounts, account)

	return account, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	account.Balance += amount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}

	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance

	}
	account.Balance -= amount

	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}

	s.payments = append(s.payments, payment)
	return payment, nil

}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}

	return nil, ErrAccountNotFound
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}

	return nil, ErrPaymentNotFound
}

func (s *Service) Reject(paymentID string) error {
	targetPayment, targetAccount, err := s.findPaymentAndAccountByPaymentID(paymentID)
	if err != nil {
		return err
	}

	targetPayment.Status = types.PaymentStatusFail
	targetAccount.Balance += targetPayment.Amount

	return nil
}

func (s *Service) findPaymentAndAccountByPaymentID(paymentID string) (*types.Payment, *types.Account, error) {
	var (
		targetPayment *types.Payment
		targetAccount *types.Account
	)

	for _, payment := range s.payments {
		if payment.ID == paymentID {
			targetPayment = payment
			break
		}
	}

	if targetPayment == nil {
		return nil, nil, ErrPaymentNotFound
	}

	for _, account := range s.accounts {
		if account.ID == targetPayment.AccountID {
			targetAccount = account
			break
		}
	}

	if targetAccount == nil {
		return nil, nil, ErrAccountNotFound
	}

	return targetPayment, targetAccount, nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	targetPayment, targetAccount, err := s.findPaymentAndAccountByPaymentID(paymentID)
	if err != nil {
		return nil, err
	}

	return s.Pay(targetAccount.ID, targetPayment.Amount, targetPayment.Category)
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	targetPayment, targetAccount, err := s.findPaymentAndAccountByPaymentID(paymentID)
	if err != nil {
		return nil, err
	}

	favorite := &types.Favorite{
		ID:        uuid.New().String(),
		AccountID: targetAccount.ID,
		Name:      name,
		Amount:    targetPayment.Amount,
		Category:  targetPayment.Category,
	}

	s.favorite = append(s.favorite, favorite)

	return favorite, nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	var targetFavorite *types.Favorite

	for _, favorite := range s.favorite {
		if favorite.ID == favoriteID {
			targetFavorite = favorite
			break
		}
	}

	if targetFavorite == nil {
		return nil, ErrFavoriteNotFound
	}

	payment, err := s.Pay(targetFavorite.AccountID, targetFavorite.Amount, targetFavorite.Category)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return err
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	result := ""
	for _, account := range s.accounts {
		result += strconv.Itoa(int(account.ID)) + "|"
		result += string(account.Phone) + "|"
		result += strconv.Itoa(int(account.Balance)) + ";"
	}

	_, err = file.Write([]byte(result))
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
