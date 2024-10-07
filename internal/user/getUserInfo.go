package user

/*
import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	// 1) Get the cookie off the req
	tokenString, err := r.Cookie("Authorization")
	if err != nil {
		http.Error(w, "cookie not found", http.StatusBadRequest)
		return
	}
	// 2) Decode/validate cookie
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		http.Error(w, "Unauthorized", http.StatusBadRequest)
	}
	// 3) Check the expiration

	// 4) Find the user with the token sub

	// 5) Attatch to req

	// 6) Continue ------ this is where you would do .Next() and do any logic related to
	// the actual aplication, such as get UserCourse
}
*/
