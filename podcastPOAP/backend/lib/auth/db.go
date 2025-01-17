package auth

import (
	"database/sql"
	"log"
	"strings"

	"github.com/jeffprestes/build-with-celo-hackathon/podcastPOAP/backend/conf"
)

//AddAccessTokenAccessLog adds log when an user call the API using an Access Token
func AddAccessTokenAccessLog(accessToken string, funcName string) (err error) {
	err = nil
	db, err := conf.GetDB()
	if err != nil {
		log.Println("[AddAccessTokenAccessLog] Error getting db connection: " + err.Error())
		return
	}
	sql := "INSERT INTO logac_accesstokenacessos (logac_accesstoken, logac_funcao, logac_quando) VALUES (?,?,NOW())"
	_, err = db.Exec(sql, accessToken, funcName)
	return
}

//AddAccessTokenRequestLog logs a new access token generation
func AddAccessTokenRequestLog(accessToken string, contatoID int) (err error) {
	err = nil
	db, err := conf.GetDB()
	if err != nil {
		log.Println("[AddAccessTokenRequestLog] Error getting db connection: " + err.Error())
		return
	}
	sql := "INSERT INTO logacr_logaccesstokenrequest (logacr_accesstoken, contato_id, logacr_when) VALUES (?,?,NOW())"
	_, err = db.Exec(sql, accessToken, contatoID)
	if err != nil {
		log.Println("[AddAccessTokenRequestLog] Error inserting into the db. Query: ", sql, " parameters: [", accessToken, "] and [", contatoID, "] - Error: ", err.Error())
	}
	return
}

//GetUserRoleByContactID get active user's role based on contato_id (user ID)
func GetUserRoleByContactID(contactID int) (role string, err error) {
	err = nil
	db, err := conf.GetDB()
	if err != nil {
		log.Println("[GetUserRoleByContactID] Error getting db connection: " + err.Error())
		return
	}
	sql := `select 
				a.logcli_role
			from
				logcli_loginclient a
			where
				a.logcli_clientlegacyid = ?`
	err = db.Get(&role, sql, contactID)
	if err != nil {
		log.Println("[GetUserRoleByContactID] Error running query: ", sql, " parameter: ", contactID, " - Error: ", err.Error())
		return
	}
	return
}

//GetUserNameByContactID get active user's name based on contato_id (user ID)
func GetUserNameByContactID(contactID int) (name string, err error) {
	err = nil
	db, err := conf.GetDB()
	if err != nil {
		log.Println("[GetUserNameByContactID] Error getting db connection: " + err.Error())
		return
	}
	sql := `select 
				logcli_clientname
			FROM 
				logcli_loginclient  
			where 
				logcli_id=?`
	err = db.Get(&name, sql, contactID)
	if err != nil {
		log.Println("[GetUserNameByContactID] Error running query: ", sql, " parameter: ", contactID, " - Error: ", err.Error())
		return
	}
	return
}

//GetUserByID get active user's data based on user ID
func GetUserByID(ID int) (user User, err error) {
	err = nil
	db, err := conf.GetDB()
	if err != nil {
		log.Println("[GetUserByID] Error getting db connection: " + err.Error())
		return
	}
	sql := `select 
				COALESCE(logcli_clientlegacyid, "") AS logcli_clientlegacyid, 
				COALESCE(logcli_lastupdate, 0) AS logcli_lastupdate,
				logcli_clientname, logcli_role, logcli_clientid, logcli_secret
			FROM 
				logcli_loginclient  
			where 
				logcli_id=?`
	err = db.Get(&user, sql, ID)
	if err != nil {
		log.Println("[GetUserByID] Error running query: ", sql, " parameter: ", ID, " - Error: ", err.Error())
		return
	}
	return
}

//AddCredentialsToUser insert user's record adding her clientID and secret
func AddCredentialsToUser(user User, role string) (err error) {
	err = nil
	db, err := conf.GetDB()
	if err != nil {
		log.Println("[AddCredentialsToUser] Error getting db connection: " + err.Error())
		return
	}
	sql := "INSERT INTO logcli_loginclient (logcli_clientlegacyid, logcli_clientname, logcli_role, logcli_clientid, logcli_secret) VALUES (?, ?, ?, ?, ?)"
	_, err = db.Exec(sql, user.ID, user.Name, role, user.ClientID, user.Secret)
	if err != nil {
		log.Println("[AddCredentialsToUser] Error inserting into the db:: ", sql, " parameters: [", user.ID, "], [", user.ClientID, "] and [", user.Secret, "] - Error: ", err.Error())
		return
	}
	return
}

//UpdateUserCredentials updates user's record adding her clientID and secret
func UpdateUserCredentials(user User, clientID string, secret string) (err error) {
	err = nil
	db, err := conf.GetDB()
	if err != nil {
		log.Println("[UpdateUserCredentials] Error getting db connection: " + err.Error())
		return
	}
	sql := "UPDATE logcli_loginclient SET logcli_clientid = ?, logcli_secret = ? , logcli_lastupdate = now() WHERE igotyou_id = ?"
	_, err = db.Exec(sql, clientID, secret, user.ID)
	if err != nil {
		log.Println("[UpdateUserCredentials] Error updating th db: ", sql, " parameters: [", clientID, "], [", secret, "] and [", user.ID, "] - Error: ", err.Error())
		return
	}
	return
}

//GetUserCredentials gets user credentials
func GetUserCredentials(user User) (clientID string, secret string, err error) {
	err = nil
	db, err := conf.GetDB()
	if err != nil {
		log.Println("[GetUserCredentials] Error getting db connection: " + err.Error())
		return
	}
	strSQL := "SELECT coalesce(logcli_clientid, '') as clientID, coalesce(logcli_secret, '') as secret FROM logcli_loginclient WHERE logcli_clientlegacyid = ?"
	row := db.QueryRow(strSQL, user.ID)
	err = row.Scan(&clientID, &secret)
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}
		log.Println("[GetUserCredentials] Error running query: ", strSQL, " parameter: [", user.ID, "]  - Error: ", err.Error())
		return
	}
	return
}

//GetUserCredentialsByLogin gets user credentials by user's login
func GetUserCredentialsByLogin(clientID string, secret string) (user User, err error) {
	err = nil
	db, err := conf.GetDB()
	if err != nil {
		log.Println("[GetUserCredentialsByID] Error getting db connection: " + err.Error())
		return
	}
	strSQL := `SELECT 
							COALESCE(logcli_clientlegacyid, "") AS logcli_clientlegacyid, 
							COALESCE(logcli_lastupdate, 0) AS logcli_lastupdate,
							logcli_id, logcli_clientname, logcli_role, logcli_clientid, logcli_secret 
							FROM logcli_loginclient 
							WHERE logcli_clientid=? AND logcli_secret=?`
	err = db.Get(&user, strSQL, clientID, secret)
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}
		log.Println("[GetUserCredentialsByLogin] Error running query: ", strSQL, " parameter: [", clientID, "]  - Error: ", err.Error())
		return
	}
	return
}

//StatusUserCredentials checks if user's has role set and if she has whether credentials is defined or not
func StatusUserCredentials(user User) (hasRole bool, hasCredentials bool, err error) {
	clientID, secret, err := GetUserCredentials(user)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Println("[StatusUserCredentials] Error getting clientID and Secret: " + err.Error())
		return
	}
	hasRole = true
	if len(strings.TrimSpace(clientID)) > 5 && len(strings.TrimSpace(secret)) > 5 {
		hasCredentials = true
	}
	return
}
