package auth

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/jeffprestes/build-with-celo-hackathon/podcastPOAP/backend/conf"
	"github.com/novatrixtech/cryptonx"
)

func decodeClientID(origin string) (texto string, err error) {
	err = nil
	tmp, err := hex.DecodeString(origin)
	if err != nil {
		log.Println("[decodeClientID] Error decoding clientID: ", origin, " - Error: ", err.Error())
		return
	}
	texto = string(tmp)
	return
}

func getDataFromClientID(clientIDDecoded string) (contactName string, nonce string, err error) {
	err = nil
	if !strings.Contains(clientIDDecoded, "|") {
		err = errors.New("Invalid ClientID. Não há o pipe, portanto não há como obter o nonce")
		return
	}
	tmpClientID := strings.Split(clientIDDecoded, "|")
	contactName = tmpClientID[0]
	nonce = tmpClientID[1]
	return
}

func decodeSecret(origin string, nonce string) (texto string, err error) {
	err = nil
	texto, err = cryptonx.Decrypter(conf.Cfg.Section("").Key("oauth_key").Value(), nonce, origin)
	if err != nil {
		log.Println("[decodeSecret] Error decoding the secret: ", origin, " - Error: ", err.Error())
		return
	}
	return
}

func getAndValidateDataFromSecret(secret string) (data time.Time, contatoID int, IP string, err error) {
	err = nil
	if !strings.Contains(secret, "|") {
		err = errors.New("[getAndValidateDataFromSecret] invalid secret. No pipe, then there is no way to get the nonce")
		return
	}
	tmp := strings.Split(secret, "|")
	if len(tmp) < 3 {
		err = errors.New("[getAndValidateDataFromSecret] invalid secret. Não há elementos suficientes nos dados")
		return
	}
	dataNum, err := strconv.ParseInt(tmp[0], 10, 64)
	if err != nil {
		log.Println("[getAndValidateDataFromSecret] Error parsing timestamp: ", tmp[0], " - Error: ", err.Error())
		return
	}
	threeDaysLimit, err := time.ParseDuration("36h")
	if err != nil {
		log.Println("[getAndValidateDataFromSecret] Error parsing 36h duation. Error: ", err.Error())
		return
	}
	if dataNum < int64(threeDaysLimit) {
		err = errors.New("[getAndValidateDataFromSecret] invalid secret. Data definida é menor que 2017-09-17")
		return
	}
	data, err = parseDateFromUnixTimestamp(tmp[0])
	if err != nil {
		log.Println("[getAndValidateDataFromSecret] Error parsing data: ", tmp[0], " - Error: ", err.Error())
		return
	}

	contatoID, err = strconv.Atoi(tmp[1])
	if err != nil {
		log.Println("[getAndValidateDataFromSecret] Error parsing contatoID: ", tmp[1], " - Error: ", err.Error())
		return
	}
	if contatoID < 1 {
		err = errors.New("[getAndValidateDataFromSecret] invalid contatoID")
		return
	}
	if len(tmp[2]) < 3 {
		err = errors.New("[getAndValidateDataFromSecret] IP invalid. Not enough items")
		log.Println("[getAndValidateDataFromSecret] ", tmp[2], " - Error: ", err.Error())
		return
	}
	IP = tmp[2]
	return
}

func parseDateFromUnixTimestamp(origin string) (data time.Time, err error) {
	err = nil
	i, err := strconv.ParseInt(origin, 10, 64)
	if err != nil {
		log.Println("[parseDateFromUnixTimestamp] Error parsing timestamp: ", origin, " - Error: ", err.Error())
		return
	}
	data = time.Unix(i, 0)
	return
}

func decodeClientIDAndSecret(clientID string, secret string) {
	clientIDDecoded, err := decodeClientID(clientID)
	if err != nil {
		log.Println("[GenerateCredentials] Error decoding clientID. Error: ", err.Error())
		return
	}
	log.Println("clientIDDecodado: ", clientIDDecoded)
	_, nonce, err := getDataFromClientID(clientIDDecoded)
	if err != nil {
		log.Println("[GenerateCredentials] Error getting nonce for clientID. Error: ", err.Error())
		return
	}
	secretDecoded, err := decodeSecret(secret, nonce)
	if err != nil {
		log.Println("[GenerateCredentials] Error decoding secret. Error: ", err.Error())
		return
	}
	log.Println("SecretDecodado: ", secretDecoded)
	secretData, contatoID, secretIP, err := getAndValidateDataFromSecret(secretDecoded)
	if err != nil {
		log.Println("[GenerateCredentials] Error getting secret. Error: ", err.Error())
		return
	}
	log.Println("Data: ", secretData, " - ContatoID: ", contatoID, " - IP: ", secretIP)
}

func generateUserCredentials(user User, remoteAddr string) (clientID string, secret string, err error) {
	err = nil
	nomeContatoorigin := strings.Replace(user.Name, " ", "", -1)
	dataorigin := time.Now().Unix()
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		//return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
		log.Printf("[generateUserCredentials] Error splitting host and port. userip: %q is not IP:port", remoteAddr)
	}
	ipRemotoorigin := net.ParseIP(ip)
	if ipRemotoorigin == nil {
		errStr := fmt.Sprintf("[generateUserCredentials] Error parsing userip: %q is not IP:port", ip)
		log.Println(errStr)
		err = errors.New(errStr)
		return
	}
	secretAntesCrypto := strconv.Itoa(int(dataorigin)) + "|" + strconv.Itoa(user.ID) + "|" + ipRemotoorigin.String()
	secret, nonce, err := cryptonx.Encrypter(conf.Cfg.Section("").Key("oauth_key").Value(), secretAntesCrypto)
	if err != nil {
		log.Println("[GenerateCredentials] Error encrypting text: ", err.Error())
		return
	}
	clientIDorigin := nomeContatoorigin + "|" + nonce
	clientID = hex.EncodeToString([]byte(clientIDorigin))
	return
}
