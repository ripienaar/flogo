package nats

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	Servers     string `md:"servers,required"`
	Credentials string `md:"credentials"`
	TLSKey      string `md:"tls_key"`
	TLSCert     string `md:"tls_cert"`
	TLSCA       string `md:"tls_ca"`
	Topic       string `md:"topic"`
}

type HandlerSettings struct {
	Topic string `md:"topic,required"`
}

type Output struct {
	Message string `md:"message"`
	Topic   string `md:"topic"`
	ReplyTo string `md:"reply_to"`
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"message":  o.Message,
		"topic":    o.Topic,
		"reply_to": o.ReplyTo,
	}
}

func (o *Output) FromMap(values map[string]interface{}) (err error) {
	o.Message, err = coerce.ToString(values["message"])
	if err != nil {
		return err
	}

	o.Topic, err = coerce.ToString(values["topic"])
	if err != nil {
		return err
	}

	o.ReplyTo, err = coerce.ToString(values["reply_to"])
	if err != nil {
		return err
	}

	return nil
}
