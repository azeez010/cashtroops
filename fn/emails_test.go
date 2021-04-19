package fn

import (
	"github.com/adigunhammedolalekan/cashtroops/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendEmail(t *testing.T) {
	// os.Setenv("SG_KEY", "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")
	err := SendEmail(&types.MailRequest{
		User:  "Lekan Adigun",
		Email: "adigunadunfe@gmail.com",
		Title: "Test - CashTroops",
		Body:  "Hi, a new email from CashTroops",
	})
	t.Log(err)
	assert.Nil(t, err)
}
