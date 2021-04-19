dev:
	export DATABASE_URL=postgres://postgres:manman@localhost:5901/cashtroops?sslmode=disable ADDR=9001 SG_KEY=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 SESSION_CACHE=${HOME}/sess BC_TOKEN=4dbf837667594004a79a40b30e34e9fc && go run server.go
