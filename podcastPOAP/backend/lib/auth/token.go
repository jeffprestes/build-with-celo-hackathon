package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jeffprestes/build-with-celo-hackathon/podcastPOAP/backend/conf"
	"github.com/jeffprestes/build-with-celo-hackathon/podcastPOAP/backend/lib/contx"
	"github.com/novatrixtech/cryptonx"
)

// CreateJWTCookie create cookie with jwt token
func CreateJWTCookie(ID string, issuer string, expiration int, ctx *contx.Context) (err error) {
	ip := ctx.RemoteAddr()
	expireCookie := time.Now().Add(time.Second * time.Duration(expiration))
	signedToken, err := generateJWTToken(ID, ip, issuer, expiration)
	if err != nil {
		log.Println("CreateJWTCookie error generating JWT: ", err.Error())
		return
	}
	cookie := http.Cookie{Name: cookieName, Value: signedToken, Expires: expireCookie, HttpOnly: true}
	http.SetCookie(ctx.Resp, &cookie)
	return
}

// InvalidateJWTToken invalidate jwt token
func InvalidateJWTToken(ctx *contx.Context) {
	deleteCookie := http.Cookie{Name: cookieName, Value: "none", Expires: time.Now()}
	http.SetCookie(ctx.Resp, &deleteCookie)
}

func generateJWTToken(ID string, ip string, issuer string, expiration int) (signedToken string, err error) {
	expireToken := time.Now().Add(time.Second * time.Duration(expiration)).Unix()

	if issuer == "" {
		issuer = "localhost:8080"
	}
	claims := Claims{
		IP: ip,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    issuer,
			Id:        ID,
		},
	}
	log.Printf("generateJWTToken Claims: %+v\n", claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(conf.Cfg.Section("").Key("oauth_key").Value()))
	if err != nil {
		log.Println("generateJWTToken error signing: ", err.Error())
		return
	}
	return
}

// ClientDecrypter decrypt client token
func ClientDecrypter(key, clientID, clientSecret string) (name, id string, err error) {
	text, err := cryptonx.Decrypter(key, clientSecret, clientID)
	if err != nil {
		return "", "", err
	}
	values := strings.Split(string(text), "|")
	name = values[0]
	id = values[1]
	return
}

//ClientEncrypter encrypts new client
func ClientEncrypter(key, appName, appID string) (clientID, clientSecret string, err error) {
	clientID, clientSecret, err = cryptonx.Encrypter(key, appName+"|"+appID)
	if err != nil {
		log.Println("[ClientEncrypter] Error encrypting text: ", err.Error())
		return
	}
	return
}

func parse(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("[parse] unexpected Signing method")
	}
	return []byte(conf.Cfg.Section("").Key("oauth_key").Value()), nil
}
