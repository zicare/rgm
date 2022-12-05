package jwt

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/lib"
)

type JWT struct {
	header  Header
	payload Payload
	token   string
}

type Header struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
}

type Payload struct {
	UID  string    `json:"uid"`
	Type string    `json:"type"`
	Role string    `json:"role"`
	TPS  float32   `json:"tps"`
	Iat  time.Time `json:"iat"`
	Exp  time.Time `json:"exp"`
}

// Returns token and exp for an auth.User
func JWTFactory(uid string, role string, t string, tps float32, iat time.Time, exp time.Time) JWT {

	var (
		now         = time.Now()
		duration, _ = time.ParseDuration(config.Config().GetString("jwt_duration"))
	)

	// Cap exp
	if exp.After(now.Add(duration)) {
		exp = now.Add(duration)
	}

	// Cap iat
	if iat.Before(now) {
		iat = now
	}

	return jWT(Payload{
		UID:  uid,
		Type: t,
		Role: role,
		TPS:  tps,
		Iat:  iat,
		Exp:  exp,
	})
}

func (j JWT) ToString() string {

	return j.token
}

func (j JWT) GetHeader() Header {

	return j.header
}

func (j JWT) GetPayload() Payload {

	return j.payload
}

// Decode exported
func Decode(token string) (Payload, error) {

	var payload Payload

	t := strings.Split(token, ".")
	if len(t) != 3 {
		return payload, new(InvalidToken)
	}

	decodedPayload, err := lib.B64Decode(t[1])
	if err != nil {
		return payload, new(InvalidTokenPayload)
	}

	ParseErr := json.Unmarshal([]byte(decodedPayload), &payload)
	if ParseErr != nil {
		return payload, new(InvalidTokenPayload)
	}

	j := jWT(payload)

	if token != j.token {
		return payload, new(TamperedToken)
	}

	if time.Now().After(j.payload.Exp) {
		return payload, new(ExpiredToken)
	}

	return payload, nil

}

func jWT(payload Payload) JWT {

	var (
		secret    = config.Config().GetString("hmac_key")
		header    = Header{Typ: "JWT", Alg: "HS256"}
		src       = lib.B64Encode(header) + "." + lib.B64Encode(payload)
		signature = lib.Hash(src, secret)
		token     = src + "." + signature
	)

	return JWT{
		header:  header,
		payload: payload,
		token:   token,
	}
}
