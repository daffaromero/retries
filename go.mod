module github.com/daffaromero/retries

go 1.22.5

replace github.com/daffaromero/retries/services/common => ./services/common

require (
	github.com/daffaromero/retries/services/common v0.0.0
	github.com/go-playground/validator/v10 v10.22.0
	github.com/gofiber/fiber/v3 v3.0.0-beta.3
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.6.0
	google.golang.org/grpc v1.65.0
)

require (
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/gabriel-vasile/mimetype v1.4.5 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/gofiber/utils/v2 v2.0.0-beta.6 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.55.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240808171019-573a1156607a // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
