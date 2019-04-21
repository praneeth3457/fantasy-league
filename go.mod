module src/fantasy

go 1.12

require (
	authorization v0.0.0-00010101000000-000000000000
	constant v0.0.0-00010101000000-000000000000
	db v0.0.0-00010101000000-000000000000
	github.com/0xAX/notificator v0.0.0-20181105090803-d81462e38c21 // indirect
	github.com/codegangsta/envy v0.0.0-20141216192214-4b78388c8ce4 // indirect
	github.com/codegangsta/gin v0.0.0-20171026143024-cafe2ce98974 // indirect
	github.com/denisenkom/go-mssqldb v0.0.0-20190412130859-3b1d194e553a
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/mux v1.7.1
	github.com/mattn/go-shellwords v1.0.5 // indirect
	github.com/rs/cors v1.6.0
	golang.org/x/crypto v0.0.0-20190417174047-f416ebab96af
	gopkg.in/go-playground/validator.v8 v8.18.2
	gopkg.in/urfave/cli.v1 v1.20.0 // indirect
	hashing v0.0.0-00010101000000-000000000000 // indirect
	models v0.0.0-00010101000000-000000000000 // indirect
	routes v0.0.0-00010101000000-000000000000
)

replace (
	authorization => ./src/fantasy/authorization
	constant => ./src/fantasy/constant
	db => ./src/fantasy/database
	hashing => ./src/fantasy/hashing
	models => ./src/fantasy/models
	routes => ./src/fantasy/routes
)
