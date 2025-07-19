module github.com/gorkagg10/lovify/lovify-api

go 1.24.3

toolchain go1.24.5

require golang.org/x/crypto v0.37.0

require (
	github.com/gorilla/mux v1.8.1
	github.com/gorkagg10/lovify/lovify-authentication-service v0.0.0-20250717212628-fcd841f6bf4f
	github.com/gorkagg10/lovify/lovify-matching-service v0.0.0-20250718153318-4c1b78cd1d0f
	github.com/gorkagg10/lovify/lovify-user-service v0.0.0-20250717211827-da5b8800d754
	github.com/rs/cors v1.11.1
	google.golang.org/grpc v1.73.0
)

require github.com/gorkagg10/lovify/lovify-messaging-service v0.0.0-20250718162554-ec688207a148

require (
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6
)
