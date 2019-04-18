module falcons-fantasy/src/fantasy

go 1.12

require (
	authorization v0.0.0-00010101000000-000000000000
	constant v0.0.0-00010101000000-000000000000
	db v0.0.0-00010101000000-000000000000
	github.com/denisenkom/go-mssqldb v0.0.0-20190412130859-3b1d194e553a
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/mux v1.7.1
	golang.org/x/crypto v0.0.0-20190417174047-f416ebab96af
	gopkg.in/go-playground/validator.v8 v8.18.2
	hashing v0.0.0-00010101000000-000000000000 // indirect
	models v0.0.0-00010101000000-000000000000 // indirect
	routes v0.0.0-00010101000000-000000000000
)

replace (
	authorization => ./falcons-fantasy/src/fantasy/authorization
	constant => ./falcons-fantasy/src/fantasy/constant
	db => ./falcons-fantasy/src/fantasy/database
	hashing => ./falcons-fantasy/src/fantasy/hashing
	models => ./falcons-fantasy/src/fantasy/models
	routes => ./falcons-fantasy/src/fantasy/routes
)
