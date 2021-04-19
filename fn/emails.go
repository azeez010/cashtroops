package fn

import (
	"errors"
	"fmt"
	"github.com/adigunhammedolalekan/cashtroops/types"
	"github.com/matcornic/hermes/v2"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"os"
)

func GenerateResetPasswordEmail(code string) (string, error) {
	h := hermes.Hermes{
		Product: hermes.Product{
			Name:        "CashTroops",
			Link:        "https://cashtroops.africa",
			Logo:        "",
			Copyright:   "cashtroops.africa",
			TroubleText: "Contact: hello@cashtroops.africa",
		},
	}
	e := hermes.Email{
		Body: hermes.Body{
			Intros: []string{
				"You have received this email because a password reset request for your CashTroops account was received.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Please use the code below to reset your password:",
					Button: hermes.Button{
						Text: code,
					},
				},
			},
			Outros: []string{
				"If you did not request a password reset, no further action is required on your part.",
			},
			Signature: "Thanks",
		},
	}
	return h.GenerateHTML(e)
}

func GenerateWelcomeEmail(accountName, code string) (string, error) {
	h := hermes.Hermes{
		Product: hermes.Product{
			Name:        "CashTroops",
			Link:        "https://cashtroops.africa",
			Logo:        "",
			Copyright:   "cashtroops.africa",
			TroubleText: "Contact: hello@cashtroops.africa",
		},
	}
	e := hermes.Email{
		Body: hermes.Body{
			Name: accountName,
			Intros: []string{
				"Welcome to CashTroops",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Use the code below to activate your account",
					Button: hermes.Button{
						Text: code,
						Link: fmt.Sprintf(""),
					},
				},
			},
			Outros: []string{
				"We're looking forward to giving you the best experience",
			},
			Signature: "Thanks",
		},
	}
	return h.GenerateHTML(e)
}

func GenerateCoinReceivedEmail() (string, error) {
	panic("")
}

func GenerateNairaReceivedEmail() (string, error) {
	panic("")
}

func GenerateDealCompletedEmail() (string, error) {
	panic("")
}

func SendEmail(req *types.MailRequest) error {
	log.Println("sending mail to ", req.Email)
	from := mail.NewEmail("CashTroops", "adigunhammed.lekan@gmail.com")
	subject := req.Title
	to := mail.NewEmail(req.User, req.Email)
	message := mail.NewSingleEmail(from, subject, to, req.Body, req.Body)
	client := sendgrid.NewSendClient(os.Getenv("SG_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return err
	}

	log.Printf("email sent to %s; Response => %v\n", req.Email, response)
	if response.StatusCode > 299 {
		return errors.New(response.Body)
	}
	return nil
}
