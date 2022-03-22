package controllers

import (
	"fmt"
	"goChallenge/chat/config"
	"goChallenge/chat/db"
	"goChallenge/chat/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type claims struct {
	Username string
}

type tokenResponse struct {
	Token string `json:"token"`
}

// swagger:route POST /auth/login Auth login
// Signup an user
//
// responses:
//	200:
//  400: ErrorResponse
//  500: Error

// Signup users
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	var err error
	DecodeRequest(r, &user)

	var userData models.User

	db.DB.Where(&models.User{Email: user.Email}).First(&userData)

	if (models.User{}) == userData {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}

	userPwd := []byte(user.Password)
	dbPwd := []byte(userData.Password)

	passErr := bcrypt.CompareHashAndPassword(dbPwd, userPwd)

	if passErr != nil {
		log.Println(passErr)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"response":"Wrong Password!"}`))
		return
	}

	jwtToken, err := generateJWT(&claims{Username: user.Email})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	EncodeResponse(w, jwtToken)
}

func HandleSignup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	var err error
	DecodeRequest(r, &user)

	user.Password, err = getHash([]byte(user.Password))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("password hash failed"))
		return
	}

	data := db.DB.Create(&user)

	if data.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("User creation failed"))
		return
	}
	EncodeResponse(w, data)
}

func getHash(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func generateJWT(c *claims) (tokenResponse, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user"] = c.Username
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString([]byte(config.JwtConfig.SecretKey))
	if err != nil {
		log.Println("Error in JWT token generation")
		return tokenResponse{}, err
	}

	return tokenResponse{Token: tokenString}, nil
}

func ExtractTokenMetadata(tokenString string) (*claims, error) {

	fmt.Println("Claims:")
	token, err := VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	claimsData, _ := token.Claims.(jwt.MapClaims)
	return &claims{
		Username: fmt.Sprintf("%v", claimsData["user"]),
	}, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv(config.JwtConfig.SecretKey)), nil
	})
	return token, nil
}
