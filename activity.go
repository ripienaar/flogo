package nats

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

// New optional factory method, should be used if one activity instance per configuration is desired
func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	act := &Activity{
		settings: s,
	}

	err = act.connectNats()
	if err != nil {
		return nil, err
	}

	return act, nil
}

func (a *Activity) connectNats() (err error) {
	opts := []nats.Option{
		nats.MaxReconnects(-1),
		nats.NoEcho(),
	}

	if a.settings.Credentials != "" {
		opts = append(opts, nats.UserCredentials(a.settings.Credentials))
	}

	if a.settings.TLSCA != "" && a.settings.TLSKey != "" && a.settings.TLSCert != "" {
		opts = append(opts, nats.ClientCert(a.settings.TLSCert, a.settings.TLSKey))
		opts = append(opts, nats.RootCAs(a.settings.TLSCA))
	}

	a.conn, err = nats.Connect(a.settings.Servers, opts...)
	return err
}

type Activity struct {
	conn     *nats.Conn
	settings *Settings
}

func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}

	if input.Ack {
		_, err = a.conn.Request(a.settings.Topic, []byte(input.Message), 5*time.Second)
	} else {
		err = a.conn.Publish(a.settings.Topic, []byte(input.Message))
	}
	if err != nil {
		return true, err
	}

	output := &Output{Delivered: true}
	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}
