
package gnmi_server

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/metadata"
	"time"
	jwt "github.com/dgrijalva/jwt-go"
	"crypto/rand"
	spb "proto/gnoi"
	"common_utils"
)

var (
	JwtRefreshInt time.Duration
	JwtValidInt   time.Duration
	hmacSampleSecret = make([]byte, 16)
)
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}


type Claims struct {
	Username string `json:"username"`
	Roles []string `json:"roles"`
	jwt.StandardClaims
}



func generateJWT(username string, roles []string, expire_dt time.Time) string {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	claims := &Claims{
		Username: username,
		Roles: roles,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expire_dt.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString(hmacSampleSecret)

	return tokenString
}
func GenerateJwtSecretKey() {
	rand.Read(hmacSampleSecret)
}

func tokenResp(username string, roles []string) *spb.JwtToken {
	exp_tm := time.Now().Add(JwtValidInt)
	token := spb.JwtToken{AccessToken: generateJWT(username, roles, exp_tm), Type: "Bearer", ExpiresIn: int64(JwtValidInt/time.Second)}
	return &token
}

func JwtAuthenAndAuthor(ctx context.Context, admin_required bool) (*spb.JwtToken, context.Context, error) {
	rc, ctx := common_utils.GetContext(ctx)
	var token spb.JwtToken
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ctx, status.Errorf(codes.Unknown, "Invalid context")
	}


	if token_str, ok := md["access_token"]; ok {
		token.AccessToken = token_str[0]
	}else {
		return nil, ctx, status.Errorf(codes.Unauthenticated, "No JWT Token Provided")
	}

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token.AccessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return hmacSampleSecret, nil
	})
	if err != nil {
		return &token, ctx, status.Errorf(codes.Unauthenticated, err.Error())
	}
	if !tkn.Valid {
		return &token, ctx, status.Errorf(codes.Unauthenticated, "Invalid JWT Token")
	}
	// if err := PopulateAuthStruct(claims.Username, &rc.Auth); err != nil {
	// 	glog.Infof("[%s] Failed to retrieve authentication information; %v", rc.ID, err)
	// 	return &token, ctx, status.Errorf(codes.Unauthenticated, "")	
	// }
	rc.Auth.User = claims.Username
	rc.Auth.Roles = claims.Roles
	return &token, ctx, nil
}

