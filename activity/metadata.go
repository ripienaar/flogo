package nats

import "github.com/project-flogo/core/data/coerce"

type Settings struct {
	Servers     string `md:"servers,required"`
	Credentials string `md:"credentials"`
	TLSKey      string `md:"tls_key"`
	TLSCert     string `md:"tls_cert"`
	TLSCA       string `md:"tls_ca"`
	Topic       string `md:"topic"`
}

type Input struct {
	Message string `md:"message,required"`
	Ack     bool   `md:"wait_for_ack"`
}

func (r *Input) FromMap(values map[string]interface{}) (err error) {
	r.Message, err = coerce.ToString(values["message"])
	if err != nil {
		return err
	}

	r.Ack, err = coerce.ToBool(values["wait_for_ack"])
	if err != nil {
		return err
	}

	return nil
}

func (r *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"message":      r.Message,
		"wait_for_ack": r.Ack,
	}
}

type Output struct {
	Delivered bool `md:"delvered"`
}

func (o *Output) FromMap(values map[string]interface{}) (err error) {
	o.Delivered, err = coerce.ToBool(values["delivered"])

	return err
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"delivered": o.Delivered,
	}
}
