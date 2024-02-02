package usecases

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-mail/mail"
	"stori/internal/core/domain"
)

func sendEmail(info *domain.Response, to, from, pass string) error {
	m := mail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Summary")
	m.SetBody("text/html",
		"Hello this is the summary <br>"+
			"Total balance is "+fmt.Sprintf("%v", info.TotalBalance)+"<br>"+
			translateMonth(info.NumberTransactions)+
			"Average debit amount "+fmt.Sprintf("%v", info.AverageDebit)+"<br>"+
			"Average credit amount "+fmt.Sprintf("%v", info.AverageCredit)+"<br><br>",
	)
	d := mail.NewDialer("smtp-mail.outlook.com", 587, from, pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("error al enviar email ", err)
		errors.New(fmt.Sprintf("%v %v", "el email no se pudo enviar", err))
	}
	fmt.Println("se ha enviado el email correctamente")
	return nil
}

func translateMonth(numberTransactions map[string]int) string {
	response := ""
	for i, value := range numberTransactions {
		switch i {
		case "1":
			response = response + "Number of transactions in enero " + fmt.Sprintf("%v", value) + " <br>"
		case "2":
			response = response + "Number of transactions in febrero " + fmt.Sprintf("%v", value) + " <br>"
		case "3":
			response = response + "Number of transactions in marzo " + fmt.Sprintf("%v", value) + " <br>"
		case "4":
			response = response + "Number of transactions in abril " + fmt.Sprintf("%v", value) + " <br>"
		case "5":
			response = response + "Number of transactions in mayo " + fmt.Sprintf("%v", value) + " <br>"
		case "6":
			response = response + "Number of transactions in junio " + fmt.Sprintf("%v", value) + " <br>"
		case "7":
			response = response + "Number of transactions in julio " + fmt.Sprintf("%v", value) + " <br>"
		case "8":
			response = response + "Number of transactions in agosto " + fmt.Sprintf("%v", value) + " <br>"
		case "9":
			response = response + "Number of transactions in septiembre " + fmt.Sprintf("%v", value) + " <br>"
		case "10":
			response = response + "Number of transactions in octubre " + fmt.Sprintf("%v", value) + " <br>"
		case "11":
			response = response + "Number of transactions in noviembre " + fmt.Sprintf("%v", value) + " <br>"
		case "12":
			response = response + "Number of transactions in diciembre " + fmt.Sprintf("%v", value) + " <br>"
		}
	}
	return response
}
