module github.com/gorkagg10/lovify-api

go 1.24.1

require golang.org/x/crypto v0.37.0

require (
	github.com/gorilla/mux v1.8.1
	github.com/gorkagg10/lovify/lovify-user-service v0.0.0-20250623194706-5c75c28a1d26
	github.com/rs/cors v1.11.1
	google.golang.org/grpc v1.73.0
)

require (
	github.com/gorkagg10/lovify-authentication-service v0.0.0-20250621100745-116aea9f7aa2
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6
)
