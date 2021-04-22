dev:
	export DATABASE_URL=postgres://postgres:manman@localhost:5901/cashtroops?sslmode=disable ADDR=9001 SG_KEY=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 SESSION_CACHE=${HOME}/session BC_TOKEN=4dbf837667594004a79a40b30e34e9fc PS_KEY=sk_test_f4f968b46db4223943c7d277621bba8f4106a66c && go fmt ./... && go run server.go
