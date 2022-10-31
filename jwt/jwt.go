package jwt

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/zicare/rgm/lib"
	"github.com/zicare/rgm/msg"
)

//User exported
type User struct {
	Src  string
	Id   int64
	Role int64
	TPS  float32
	Usr  string
	Pwd  string
	From time.Time
	To   time.Time
}

//Payload exported
type Payload struct {
	Exp  int64   `json:"exp"`
	Iat  int64   `json:"iat"`
	Id   int64   `json:"id"`
	Role int64   `json:"role"`
	Src  string  `json:"src"`
	TPS  float32 `json:"tps"`
}

//Token exported
func Token(u User, duration time.Duration, secret string) (string, string) {

	now := time.Now()
	j := new(Payload{
		Exp:  now.Add(duration).Unix(),
		Iat:  now.Unix(),
		Id:   u.Id,
		Role: u.Role,
		Src:  u.Src,
		TPS:  u.TPS,
	}, secret)

	return j.token, j.exp

}

//Decode exported
func Decode(token string, secret string) (Payload, *msg.Message) {

	var payload Payload

	t := strings.Split(token, ".")
	if len(t) != 3 {
		//Invalid token
		return payload, msg.Get("12").M2E()
	}

	decodedPayload, PayloadErr := lib.Decode(t[1])
	if PayloadErr != nil {
		//Invalid payload
		return payload, msg.Get("13").M2E()
	}

	ParseErr := json.Unmarshal([]byte(decodedPayload), &payload)
	if ParseErr != nil {
		//Invalid payload
		return payload, msg.Get("13").M2E()
	}

	j := new(payload, secret)

	if token != j.token {
		//Token tampered
		return payload, msg.Get("14").M2E()
	}

	if j.payload.Exp < time.Now().Unix() {
		//Token expired
		return payload, msg.Get("15").M2E()
	}

	return payload, nil

}

type jwt struct {
	header  header
	payload Payload
	token   string
	exp     string
}

type header struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
}

func new(payload Payload, secret string) jwt {

	j := jwt{}
	j.header = header{
		Typ: "JWT",
		Alg: "HS256",
	}
	j.payload = payload
	j.token = j.getToken(secret)
	j.exp = time.Unix(j.payload.Exp, 0).Format(time.RFC3339)
	return j
}

func (j jwt) getToken(secret string) string {

	var (
		src       = lib.Encode(j.header) + "." + lib.Encode(j.payload)
		signature = lib.Hash(src, secret)
	)
	return src + "." + signature
}
