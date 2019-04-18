package authorization

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"


	"constant"
	"db"

	"github.com/dgrijalva/jwt-go"
	//"github.com/fantasy/routes"
)

// Params :
type Params struct {
	endpoint   func(http.ResponseWriter, *http.Request)
	userAccess int
}

var mySigningKey = "thisisnotaneasyphasetocrack"

//var mySigningKey = os.Getenv("LOCAL_JWT_TOKEN")

// GenerateJWT :
func GenerateJWT(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user"] = "Elliot Forbes"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	key := mySigningKey + username
	tokenString, err := token.SignedString([]byte(key))

	if err != nil {
		fmt.Println("Something went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

// IsAuthorized :
func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request), userAccess int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				key := mySigningKey + r.Header["User-Context"][0]
				return []byte(key), nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				if userAccess != constant.UserAny {
					if isUserRoleAuthorized(userAccess, r.Header["User-Context"][0]) {
						endpoint(w, r)
					} else {
						fmt.Fprintf(w, "Not an authorized user")
					}
				} else {
					endpoint(w, r)
				}
			}
		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

func isUserRoleAuthorized(accessRole int, username string) bool {
	var role int
	err := database.Db.QueryRow("SELECT Role FROM users WHERE Username=@Username", sql.Named("Username", username)).Scan(&role)
	if err != nil {
		return false
	} else if role == accessRole {
		return true
	} else {
		return false
	}
}
