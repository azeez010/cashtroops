package main

import (
	"fmt"
	"github.com/adigunhammedolalekan/cashtroops/config"
	"github.com/adigunhammedolalekan/cashtroops/database"
	"github.com/adigunhammedolalekan/cashtroops/http"
	"github.com/adigunhammedolalekan/cashtroops/libs/bc"
	"github.com/adigunhammedolalekan/cashtroops/ops"
	"github.com/adigunhammedolalekan/cashtroops/session"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	nethttp "net/http"
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

	userOps := ops.NewUserOps(db, sess, logger)
	accountOps := ops.NewAccountOps(db, logger)
	paymentOpts := ops.NewPaymentOps(db, bcClient, userOps, accountOps, logger)
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
	})

	addr := fmt.Sprintf(":%s", cfg.Addr)
	logger.Infof("API running at %s", addr)
	if err := nethttp.ListenAndServe(addr, router); err != nil {
		logger.WithError(err).Fatal("failed to start API server")
	}
}
