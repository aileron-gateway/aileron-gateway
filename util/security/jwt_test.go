package security_test

import (
	"encoding/base64"
	"fmt"

	v1 "github.com/aileron-gateway/aileron-gateway/apis/app/v1"
	"github.com/aileron-gateway/aileron-gateway/util/security"
	"github.com/golang-jwt/jwt/v5"
)

func ExampleJWTHandler_TokenWithClaims() {
	// Crate a JWT string with a common key.
	key := "this is the common key"

	spec := &v1.JWTHandlerSpec{
		PrivateKeys: []*v1.SigningKeySpec{
			{
				KeyID:     "test",
				Algorithm: v1.SigningKeyAlgorithm_HS256,
				KeyType:   v1.SigningKeyType_COMMON,
				KeyString: base64.StdEncoding.EncodeToString([]byte(key)),
			},
		},
	}

	jh, err := security.NewJWTHandler(spec, nil)
	if err != nil {
		panic(err) // handle error here.
	}

	claims := jwt.MapClaims{
		"xxx": "XXX",
		"yyy": "YYY",
		"zzz": "ZZZ",
	}

	token, err := jh.TokenWithClaims(claims)
	if err != nil {
		panic(err) // handle error here.
	}

	tokenString, err := jh.SignedString(token)
	if err != nil {
		panic(err) // handle error here.
	}

	fmt.Println(tokenString)
	// Output:
	// eyJhbGciOiJIUzI1NiIsImtpZCI6InRlc3QiLCJ0eXAiOiJKV1QifQ.eyJ4eHgiOiJYWFgiLCJ5eXkiOiJZWVkiLCJ6enoiOiJaWloifQ.11Vd7J81DhiqyCPxm5mMvI9vyIcIAdKMgSmDUZ2Ity4
}
