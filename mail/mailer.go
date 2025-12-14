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

// Message exported
type Message struct {
	To      string
	ToN     []string
	Cc      []string
	Bcc     []string
	Subject string
	Tpl     string
	Data    interface{}
}

// Send exported
func (msg *Message) Send(iteration int) {

	if iteration > config.Config().GetInt("smtp.retries") {
		glog.Error(errors.New("could not send email"))
		glog.Error(errors.New("exceeded retries"))
		glog.Flush()
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
			glog.Error(errors.New("could not send email"))
			glog.Error(err)
			glog.Flush()
			return
		}

		m = mail.NewMessage()
		m.SetHeader("From", c.GetString("smtp.user"))
		m.SetHeader("To", msg.to()...)
		m.SetHeader("Cc", msg.cc()...)
		m.SetHeader("Bcc", msg.bcc()...)
		m.SetHeader("Subject", msg.Subject)
		m.SetBody("text/html", tpl.String())
		d = mail.NewDialer(c.GetString("smtp.host"), c.GetInt("smtp.port"),
			c.GetString("smtp.user"), c.GetString("smtp.password"))
		d.Timeout = c.GetDuration("smtp.timeout")
		if err := d.DialAndSend(m); err != nil {
			glog.Error(errors.New("*mail.Dialer.DialAndSend error"))
			glog.Error(err)
			glog.Flush()
			duration, _ := time.ParseDuration(config.Config().GetString("smtp.retry_interval"))
			time.Sleep(duration)
			msg.Send(iteration + 1)
		}
	}(msg, iteration)
}

func (msg *Message) to() []string {

	var (
		seen = map[string]bool{}
		rcpt []string
	)

	// Prefer ToN if present
	if len(msg.ToN) > 0 {
		for _, to := range msg.ToN {
			if to == "" {
				continue
			}
			if seen[to] {
				continue
			}
			seen[to] = true
			rcpt = append(rcpt, to)
		}
		return rcpt
	}

	// Fallback to To
	if msg.To != "" {
		rcpt = append(rcpt, msg.To)
	}

	return rcpt
}

func (msg *Message) cc() []string {

	var (
		seen = map[string]bool{}
		rcpt []string
	)

	// Seed seen with To recipients so we don't duplicate addresses across headers.
	for _, to := range msg.to() {
		if to == "" {
			continue
		}
		seen[to] = true
	}

	for _, cc := range msg.Cc {
		if cc == "" {
			continue
		}
		if seen[cc] {
			continue
		}
		seen[cc] = true
		rcpt = append(rcpt, cc)
	}

	return rcpt
}

func (msg *Message) bcc() []string {

	var (
		seen = map[string]bool{}
		rcpt []string
	)

	// Seed seen with To + Cc so we don't duplicate addresses across headers.
	for _, to := range msg.to() {
		if to == "" {
			continue
		}
		seen[to] = true
	}
	for _, cc := range msg.cc() {
		if cc == "" {
			continue
		}
		seen[cc] = true
	}

	for _, bcc := range msg.Bcc {
		if bcc == "" {
			continue
		}
		if seen[bcc] {
			continue
		}
		seen[bcc] = true
		rcpt = append(rcpt, bcc)
	}

	return rcpt
}
