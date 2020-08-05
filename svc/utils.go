package svc

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JermineHu/themis/common"
	"github.com/JermineHu/themis/models"
	"github.com/JermineHu/themis/svc/gen/admin"
	"github.com/JermineHu/themis/svc/gen/config"
	"github.com/JermineHu/themis/utils"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/urfave/cli"
	"goa.design/goa/v3/security"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	pubKey          []byte
	privatekey      *rsa.PrivateKey
	privatekeybytes []byte
)

var (
	// ErrUnauthorized is the error returned by Login when the request credentials
	// are invalid.
	ErrUnauthorized error = config.Unauthorized("用户名或者密码错误！")

	// ErrInvalidToken is the error returned when the JWT token is invalid.
	ErrInvalidToken error = config.Unauthorized("一个无效的令牌！")

	// ErrInvalidTokenScopes is the error returned when the scopes provided in
	// the JWT token claims are invalid.
	//ErrInvalidTokenScopes error = regiment.InvalidScopes("令牌授权范围错误！")

	// Key is the key used in JWT authentication
	Key = []byte("secret")

	// ErrInvalidTokenScopes is the error returned when the scopes provided in
	// the JWT token claims are invalid.
	ErrInvalidTokenScopes error = admin.InvalidScopes("令牌授权范围错误！")
)

// 根据JWT获取用户
func GetUserIDByJWT(tokenStr string) (uID *uint64, er error) {

	claims, er := GetClaimsByTokenStr(tokenStr)
	if er != nil {
		return nil, er
	}
	// Use the claims to authorize
	usrID := uint64((*claims)["usr_id"].(float64))
	uID = &usrID
	return
}

// 根据JWT获取主机ID
func GetHostIDByJWT(tokenStr string) (uID *uint64, er error) {
	claims, er := GetClaimsByTokenStr(tokenStr)
	if er != nil {
		return nil, er
	}
	// Use the claims to authorize
	if v, ok := (*claims)["host_id"]; ok {
		hostID := uint64(v.(float64))
		uID = &hostID
	}
	return
}

//获取公钥
func GetPubKey() ([]byte, error) {
	if pubKey == nil {
		keyFiles, err := filepath.Glob("./jwtkey/*.pub")
		if len(keyFiles) == 0 || err != nil {
			return nil, errors.New("Not found the public key")
		}
		pubKey, err = ioutil.ReadFile(keyFiles[0])
		if err != nil {
			return nil, err
		}
	}
	return pubKey, nil
}

// 获取私钥
func GetPrivateKey() (pubK *rsa.PrivateKey, err error) {

	if privatekey == nil {
		//b, err := ioutil.ReadFile(viper.GetString("app.jwtKey"))
		b, err := ioutil.ReadFile("jwtkey/jwt.key")
		if err != nil {
			return nil, err
		}
		privatekey, err = jwtgo.ParseRSAPrivateKeyFromPEM(b)

		if err != nil {
			return nil, err
		}
	}
	return privatekey, nil
}

// 获取私钥
func GetPrivateKeyBytes() (pubK []byte, err error) {
	if privatekeybytes == nil {
		//b, err := ioutil.ReadFile(viper.GetString("app.jwtKey"))
		privatekeybytes, err = ioutil.ReadFile("jwtkey/jwt.key")
		if err != nil {
			return nil, err
		}
	}
	return privatekeybytes, nil
}

func GetClaimsByTokenStr(tokenStr string) (claims *jwtgo.MapClaims, err error) {
	clms := make(jwtgo.MapClaims)
	tk, err := utils.ConvertJWTByStr(tokenStr)
	if err != nil {
		return nil, ErrInvalidToken
	}
	clms = tk.Claims.(jwtgo.MapClaims)
	return &clms, nil
}

func listAllFlag() cli.BoolFlag {
	return cli.BoolFlag{
		Name:  "all,a",
		Usage: "Show stop/inactive and recently removed resources",
	}
}

func listSystemFlag() cli.BoolFlag {
	return cli.BoolFlag{
		Name:  "system,s",
		Usage: "Show system resources",
	}
}

func makeJWTWithAdmin(respon admin.Admin) (tokenStr *string, err error) {
	// Generate JWT
	token := jwtgo.New(jwtgo.SigningMethodRS512)
	in60d := time.Now().Add(time.Duration(GetTokenTimeoutTime()) * time.Second).Unix()
	uid, _ := uuid.NewV4()
	token.Claims = jwtgo.MapClaims{
		"iss":      "jermine.vdo.pub", // who creates the token and signs it
		"aud":      respon.ID,         // to whom the token is intended to be sent
		"exp":      in60d,             // time when the token will expire (60 day from now)
		"jti":      uid.String(),      // a unique identifier for the token
		"iat":      time.Now().Unix(), // when the token was issued/created (now)
		"nbf":      2,                 // time before which the token is not yet valid (2 minutes ago)
		"sub":      "themis_login",    // the subject/principal is whom the token is about
		"scopes":   "api:access",      // token scope - not a standard claim
		"usr_id":   respon.ID,         // token scope - not a standard claim
		"usr_name": respon.LoginName,  // token scope - not a standard claim
	}

	prik, err := GetPrivateKey()
	if err != nil {
		return nil, err
	}
	signedToken, err := token.SignedString(prik)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %s", err) // internal error
	}
	token_str := "Bearer " + signedToken
	return &token_str, nil
}

func makeJWTWithHost(respon models.Host) (tokenStr *string, err error) {
	// Generate JWT
	token := jwtgo.New(jwtgo.SigningMethodRS512)
	in60d := time.Now().Add(time.Duration(GetTokenTimeoutTime()) * time.Second).Unix()
	uid, _ := uuid.NewV4()
	token.Claims = jwtgo.MapClaims{
		"iss":     "jermine.vdo.pub", // who creates the token and signs it
		"aud":     respon.ID,         // to whom the token is intended to be sent
		"exp":     in60d,             // time when the token will expire (60 day from now)
		"jti":     uid.String(),      // a unique identifier for the token
		"iat":     time.Now().Unix(), // when the token was issued/created (now)
		"nbf":     2,                 // time before which the token is not yet valid (2 minutes ago)
		"sub":     "themis_login",    // the subject/principal is whom the token is about
		"scopes":  "api:access",      // token scope - not a standard claim
		"host_id": respon.ID,         // token scope - not a standard claim
	}

	prik, err := GetPrivateKey()
	if err != nil {
		return nil, err
	}
	signedToken, err := token.SignedString(prik)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %s", err) // internal error
	}
	token_str := "Bearer " + signedToken
	return &token_str, nil
}

func makeJWT(respon admin.Admin) (tokenStr *string, err error) {
	// Generate JWT
	token := jwtgo.New(jwtgo.SigningMethodRS512)
	in60d := time.Now().Add(time.Duration(GetTokenTimeoutTime()) * time.Second).Unix()
	uid, _ := uuid.NewV4()
	token.Claims = jwtgo.MapClaims{
		"iss":      "jermine.vdo.pub", // who creates the token and signs it
		"aud":      respon.ID,         // to whom the token is intended to be sent
		"exp":      in60d,             // time when the token will expire (60 day from now)
		"jti":      uid.String(),      // a unique identifier for the token
		"iat":      time.Now().Unix(), // when the token was issued/created (now)
		"nbf":      2,                 // time before which the token is not yet valid (2 minutes ago)
		"sub":      "themis_login",    // the subject/principal is whom the token is about
		"scopes":   "api:access",      // token scope - not a standard claim
		"usr_id":   respon.ID,         // token scope - not a standard claim
		"usr_name": respon.LoginName,  // token scope - not a standard claim
	}

	prik, err := GetPrivateKey()
	if err != nil {
		return nil, err
	}
	signedToken, err := token.SignedString(prik)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %s", err) // internal error
	}
	token_str := "Bearer " + signedToken
	return &token_str, nil
}

func GetTokenTimeoutTime() int64 {
	tm := os.Getenv(common.TOKEN_TIMEOUT)
	if len(tm) > 0 {
		t, _ := strconv.ParseInt(tm, 10, 64)
		return t
	}
	return 30
}

// 将任何数据类型转化为二进制数据
func M2BytesData(anyType interface{}) []byte {
	data, _ := json.Marshal(anyType)
	return data
}

func JWTCheck(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	claims := make(jwtgo.MapClaims)

	// authorize request
	// 1. parse JWT token, token key is hardcoded to "secret" in this example
	_, err := jwtgo.ParseWithClaims(token, claims, func(_ *jwtgo.Token) (interface{}, error) { return utils.LoadRSAPublicKeyFromDisk(), nil })

	//
	//ret, b := pares(str, pubkey) //直接使用公钥字符串
	//if b {
	//	fmt.Printf("111 pares ok,value:%+v", ret)
	//} else {
	//	fmt.Println("pares error")
	//}
	//
	if err != nil {
		return ctx, ErrInvalidToken
	}
	// 2. validate provided "scopes" claim
	if claims["scopes"] == nil {
		return ctx, ErrInvalidTokenScopes
	}
	//scopes, ok := claims["scopes"].([]interface{})
	//if !ok {
	//	return ctx, ErrInvalidTokenScopes
	//}
	//scopesInToken := make([]string, len(scopes))
	//for _, scp := range scopes {
	//	scopesInToken = append(scopesInToken, scp.(string))
	//}
	//if err := scheme.Validate(scopesInToken); err != nil {
	//	return ctx, ErrInvalidTokenScopes
	//}
	//return ctx, nil
	return ctx, nil

}

func JWTCheckForHost(ctx context.Context, token string) (context.Context, error) {
	claims := make(jwtgo.MapClaims)

	// authorize request
	// 1. parse JWT token, token key is hardcoded to "secret" in this example
	_, err := jwtgo.ParseWithClaims(token, claims, func(_ *jwtgo.Token) (interface{}, error) { return utils.LoadRSAPublicKeyFromDisk(), nil })

	//
	//ret, b := pares(str, pubkey) //直接使用公钥字符串
	//if b {
	//	fmt.Printf("111 pares ok,value:%+v", ret)
	//} else {
	//	fmt.Println("pares error")
	//}
	//
	if err != nil {
		return ctx, ErrInvalidToken
	}
	// 2. validate provided "scopes" claim
	if claims["scopes"] == nil {
		return ctx, ErrInvalidTokenScopes
	}
	if v, ok := claims["host_id"]; !ok || v == nil {
		return ctx, ErrInvalidTokenScopes
	}
	//scopes, ok := claims["scopes"].([]interface{})
	//if !ok {
	//	return ctx, ErrInvalidTokenScopes
	//}
	//scopesInToken := make([]string, len(scopes))
	//for _, scp := range scopes {
	//	scopesInToken = append(scopesInToken, scp.(string))
	//}
	//if err := scheme.Validate(scopesInToken); err != nil {
	//	return ctx, ErrInvalidTokenScopes
	//}
	//return ctx, nil
	return ctx, nil

}
