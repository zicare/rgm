package msg

import "fmt"

var msg map[string]Message

//Init exported
func Init(m []Message) (err error) {

	_init()

	//add client messages
	for _, v := range m {
		if _, ok := msg[v.Key]; ok {
			m := Get("1") //Invalid message key
			return &m
		}
		msg[v.Key] = v
	}

	return nil
}

//New exported
func New(k string, m string) Message {
	return Message{k, m, nil, ""}
}

//GetAll exported
func GetAll() map[string]Message {
	return msg
}

//Get exported
func Get(key string) Message {

	if message, ok := msg[key]; ok {
		return message
	}
	return New(key, fmt.Sprintf("Message %v", key))
}

func _init() {

	msg = make(map[string]Message)

	msg["1"] = New("1", "Invalid message key")
	msg["2"] = New("2", "%s tags not properly set")
	msg["3"] = New("3", "HTTP basic authentication required")
	msg["4"] = New("4", "Invalid credentials")
	msg["5"] = New("5", "We run into an unlikely error verifyng your credentials")
	msg["6"] = New("6", "Expired credentials")
	msg["7"] = New("7", "JWT authorization header malformed")
	msg["8"] = New("8", "Not enough permissions")
	msg["9"] = New("9", "Role access expired or not yet valid")
	msg["10"] = New("10", "TPS limit exceeded, access void until %s")
	msg["11"] = New("11", "Unauthorized")
	msg["12"] = New("12", "Invalid token")
	msg["13"] = New("13", "Invalid payload")
	msg["14"] = New("14", "Token tampered")
	msg["15"] = New("15", "Token expired")
	msg["16"] = New("16", "Read only model")
	msg["17"] = New("17", "Decoding Error %s")
	msg["18"] = New("18", "No found!")
	msg["19"] = New("19", "There are validation errors")
	msg["20"] = New("20", "TPS precision must be between %s and %s")
	msg["21"] = New("21", "TTPS data clean up cycle freq must be between %s and %s minutes")
	msg["22"] = New("22", "Time %s has a wrong format, required format is %s")
	msg["23"] = New("23", "Value is a %s, required type is %s")
	msg["24"] = New("24", "Value %s didn't pass %s(%s) validation")
	msg["25"] = New("25", "Server error: %s(%s)")
	msg["26"] = New("26", "Url param error")
	msg["27"] = New("27", "Couldn't retrieve Gin's default validator engine")
	msg["28"] = New("28", "Unauthorized app")
	msg["29"] = New("29", "%s row(s) deleted")
	//msg["30"] = New("30", "Auth tags no set")
	msg["31"] = New("31", "TPS penalty factor must be between %s and %s")
	msg["32"] = New("32", "Token revoked")
	msg["33"] = New("33", "PIN request accepted")
	msg["34"] = New("34", "DB Insert error")
	msg["35"] = New("35", "Created!")
	msg["36"] = New("36", "Invalid PIN")
	msg["37"] = New("37", "Updated!")
	msg["38"] = New("38", "Patched!")
	msg["39"] = New("39", "Expired PIN")
	msg["40"] = New("40", "Can't delete a parent resource, remove children and retry.")
	//msg["29"] = New("33", "CORS tags are not properly set")
}
