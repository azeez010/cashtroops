package main

import (
	"encoding/json"
	"fmt"
	"github.com/adigunhammedolalekan/cashtroops/config"
	"github.com/adigunhammedolalekan/cashtroops/database"
	"github.com/adigunhammedolalekan/cashtroops/http"
	"github.com/adigunhammedolalekan/cashtroops/libs/bc"
	"github.com/adigunhammedolalekan/cashtroops/libs/paystackclient"
	"github.com/adigunhammedolalekan/cashtroops/ops"
	"github.com/adigunhammedolalekan/cashtroops/session"
	"github.com/adigunhammedolalekan/cashtroops/types"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	nethttp "net/http"
	"os"
	"path/filepath"
)

func main() {
	cfg := config.New()
	router := chi.NewRouter()
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	db, err := database.Open(cfg.DatabaseUrl)
	if err != nil {
		logger.WithField("database_url", cfg.DatabaseUrl).
			WithError(err).Fatal("failed to open database")
	}
	sess, err := session.New(cfg.SessionCacheDir)
	if err != nil {
		logger.WithError(err).Fatal("failed to init session store")
	}
	bcClient, err := bc.New(cfg.BlockCypherToken, "bcy", "test", logger)
	if err != nil {
		logger.WithError(err).Fatal("failed to init BC client")
	}
	banks, err := loadBanks()
	if err != nil {
		logger.WithError(err).Fatal("failed to load bank data")
	}

	ps := paystackclient.New(cfg.PayStackKey, logger)
	userOps := ops.NewUserOps(db, sess, logger)
	accountOps := ops.NewAccountOps(db, ps, logger)
	accountOps.SetBanks(banks)

	paymentOpts := ops.NewPaymentOps(db, bcClient, userOps, accountOps, ps, logger)
	userHandler := http.NewUserHandler(userOps, logger)
	accountHandler := http.NewAccountHandler(accountOps, userOps, logger)
	paymentHandler := http.NewPaymentHandler(paymentOpts, userOps, logger)

	if err := paymentOpts.InitRate("USD-NGN", 490); err != nil {
		logger.WithError(err).Fatal("failed to init rate")
	}

	router.Route("/api", func(r chi.Router) {
		r.Post("/user/new", userHandler.CreateUser)
		r.Post("/user/authenticate", userHandler.AuthenticateUser)
		r.Post("/user/activate", userHandler.ActivateAccount)
		r.Get("/me", userHandler.Me)
		r.Get("/user/{email}/resetpassword", userHandler.RequestPasswordReset)
		r.Post("/user/verifypasswordreset", userHandler.VerifyPasswordResetRequest)
		r.Post("/user/changepassword", userHandler.ResetPassword)
		r.Put("/me/changepassword", userHandler.ChangePassword)
		r.Post("/me/beneficiary/new", accountHandler.AddBeneficiary)
		r.Delete("/me/beneficiary/{id}/remove", accountHandler.RemoveBeneficiary)
		r.Get("/me/beneficiaries", accountHandler.ListBeneficiaries)
		r.Post("/payment/init", paymentHandler.InitializePayment)
		r.Post("/txn/events", paymentHandler.TxnEventHandler)
		r.Get("/me/payments", paymentHandler.ListPayments)
		r.Post("/transfer/events", paymentHandler.TransferEventHandler)
		r.Get("/banks", accountHandler.Banks)
	})

	addr := fmt.Sprintf(":%s", cfg.Addr)
	logger.Infof("API running at %s", addr)
	if err := nethttp.ListenAndServe(addr, router); err != nil {
		logger.WithError(err).Fatal("failed to start API server")
	}
}

func loadBanks() ([]types.Bank, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(wd, "banks.json")
	data, err := ioutil.ReadFile(dir)
	if err != nil {
		return nil, err
	}
	values := make([]types.Bank, 0)
	if err := json.Unmarshal(data, &values); err != nil {
		return nil, err
	}
	return values, nil
}
