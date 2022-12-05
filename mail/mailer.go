package mail

import (
	"bytes"
	"errors"
	"html/template"
	"time"

	"github.com/golang/glog"
	"github.com/zicare/rgm/config"
	"gopkg.in/mail.v2"
)

//Message exported
type Message struct {
	To      string
	Subject string
	Tpl     string
	Data    interface{}
}

//Send exported
func (msg *Message) Send(iteration int) {

	if iteration > config.Config().GetInt("smtp.retries") {
		glog.Error(errors.New("Could not send email"))
		glog.Error(errors.New("Exceeded retries"))
		return
	}

	go func(msg *Message, iteration int) {

		var (
			c    = config.Config()
			t, _ = template.New(msg.Tpl).ParseFiles("tpl/email/" + msg.Tpl)
			d    *mail.Dialer
			m    *mail.Message
			tpl  bytes.Buffer
		)

		if err := t.Execute(&tpl, msg.Data); err != nil {
			glog.Error(errors.New("Could not send  email"))
			glog.Error(err)
			return
		}

		m = mail.NewMessage()
		m.SetHeader("From", c.GetString("smtp.user"))
		m.SetHeader("To", msg.To)
		m.SetHeader("Subject", msg.Subject)
		m.SetBody("text/html", tpl.String())
		d = mail.NewDialer(c.GetString("smtp.host"), c.GetInt("smtp.port"),
			c.GetString("smtp.user"), c.GetString("smtp.password"))
		d.Timeout = c.GetDuration("smtp.timeout")
		if err := d.DialAndSend(m); err != nil {
			duration, _ := time.ParseDuration(config.Config().GetString("smtp.retry_interval"))
			time.Sleep(duration)
			msg.Send(iteration + 1)
		}
	}(msg, iteration)
}
