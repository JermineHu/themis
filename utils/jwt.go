/**
Create by jermine

Date  19-8-9-下午3:27

**/
package utils

import (
	"crypto/rsa"
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"io/ioutil"
)

func KeyFunc(token *jwtgo.Token) (interface{}, error) {
	return LoadRSAPublicKeyFromDisk(), nil
}

//// LoadJWTPublicKeys loads PEM encoded RSA public keys used to validata and decrypt the JWT.
//func LoadJWTPublicKeys() ([]jwt.Key, error) {
//	keyFiles, err := filepath.Glob("./jwtkey/*.pub")
//	if err != nil {
//		return nil, err
//	}
//	keys := make([]jwt.Key, len(keyFiles))
//	for i, keyFile := range keyFiles {
//		pem, err := ioutil.ReadFile(keyFile)
//		if err != nil {
//			return nil, err
//		}
//		key, err := jwtgo.ParseRSAPublicKeyFromPEM([]byte(pem))
//		if err != nil {
//			return nil, fmt.Errorf("failed to load key %s: %s", keyFile, err)
//		}
//		keys[i] = key
//	}
//	if len(keys) == 0 {
//		return nil, fmt.Errorf("couldn't load public keys for JWT security")
//	}
//
//	return keys, nil
//}

var (
	jwtPubK *rsa.PublicKey
	//jwtKey  []*jwt.KeyResolver
)

func LoadRSAPublicKeyFromDisk() *rsa.PublicKey {
	//location := viper.GetString("app.jwtKeyPub")
	location := "jwtkey/jwt.key.pub"
	if jwtPubK == nil {
		keyData, e := ioutil.ReadFile(location)
		if e != nil {
			panic(e.Error())
		}
		jwtPubK, e = jwtgo.ParseRSAPublicKeyFromPEM(keyData)
		if e != nil {
			panic(e.Error())
		}
	}
	return jwtPubK
}

func ConvertJWTByStr(token string) (*jwtgo.Token, error) {
	tk, err := jwtgo.Parse(token, KeyFunc)
	if err != nil {
		return nil, err
	}
	if tk == nil {
		return nil, fmt.Errorf("JWT token is missing from context") // internal error
	}
	return tk, nil
}
