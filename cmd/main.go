package main

import (
	"log"

	"github.com/shuhrat-shokirov/wallet/pkg/wallet"
)

func main() {
	s := wallet.Service{}

	_, err := s.RegisterAccount("+992935626274")
	if err != nil {
		log.Println(err)
		return
	}

	err = s.ExportToFile("1.txt")
	if err != nil {
		log.Println(err)
		return
	}

	err = s.ImportFromFile("1.txt")
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(s.FindAccountByID(1))
}
