package jwt

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/lib"
	"github.com/zicare/rgm/msg"
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

	return new(Payload{
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

// Decode exported
func Decode(token string) (Payload, error) {

	var payload Payload

	t := strings.Split(token, ".")
	if len(t) != 3 {
		//Invalid token
		e := InvalidToken{Message: msg.Get("12")}
		return payload, &e
	}

	decodedPayload, err := lib.B64Decode(t[1])
	if err != nil {
		//Invalid payload
		e := InvalidTokenPayload{Message: msg.Get("13")}
		return payload, &e
	}

	ParseErr := json.Unmarshal([]byte(decodedPayload), &payload)
	if ParseErr != nil {
		//Token tampered
		e := InvalidTokenPayload{Message: msg.Get("14")}
		return payload, &e
	}

	j := new(payload)

	if token != j.token {
		//Token tampered
		e := TamperedToken{Message: msg.Get("14")}
		return payload, &e
	}

	if time.Now().After(j.payload.Exp) {
		//Token expired
		e := ExpiredToken{Message: msg.Get("15")}
		return payload, &e
	}

	return payload, nil

}

func new(payload Payload) JWT {

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
