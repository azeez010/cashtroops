package fn

import (
	"github.com/adigunhammedolalekan/cashtroops/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendEmail(t *testing.T) {
	err := SendEmail(&types.MailRequest{
		User:  "Lekan Adigun",
		Email: "adigunadunfe@gmail.com",
		Title: "Test - CashTroops",
		Body:  "Hi, a new email from CashTroops",
	})
	t.Log(err)
	assert.Nil(t, err)
}
