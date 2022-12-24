package ds

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/mail"
)

// IPinDataSource defines an interface to
// post pins and patch passwords using a pin.
type IPinDataSource interface {

	// Post Pin
	Post(email string) (Pin, error)

	// Patch password
	PatchPwd(patch *Patch) error
}

type Pin struct {
	Email      string    `json:"email"`
	Code       string    `json:"code"`
	Created    time.Time `json:"created"`
	Expiration time.Time `json:"expiration"`
}

func (p Pin) Send() {

	msg := new(mail.Message)
	msg.To = p.Email
	msg.Subject = "Has recibido un PIN"
	msg.Tpl = "pin.tpl"
	msg.Data = struct{ PIN string }{PIN: p.Code}
	msg.Send(1)
}

type Patch struct {
	Pin      string `json:"pin"`
	Email    string `json:"usr"`
	Password string `json:"pwd"`
}

func PatchReceiver() interface{} {

	return reflect.New(reflect.StructOf([]reflect.StructField{
		{
			Name: "Pin",
			Type: reflect.TypeOf(string("")),
			Tag:  `json:"pin" binding:"required"`,
		},
		{
			Name: "Email",
			Type: reflect.TypeOf(string("")),
			Tag:  `json:"usr" binding:"required,email"`,
		},
		{
			Name: "Password",
			Type: reflect.TypeOf(string("")),
			Tag:  reflect.StructTag(`json:"pwd" binding:"` + config.Config().GetString("account.pwd_validation") + `"`),
		},
	})).Elem().Addr().Interface()
}

func PatchDecoder(pr any) *Patch {

	/*
		    val := reflect.ValueOf(pr).Elem()
			code := val.FieldByName("Pin").Interface().(string)
			email := val.FieldByName("Email").Interface().(string)
			password := val.FieldByName("Password").Interface().(string)
	*/

	patch := new(Patch)

	data, _ := json.Marshal(pr)
	json.Unmarshal(data, patch)

	return patch
}
