package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{})

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
}

// Factory is a NATS trigger factory
type Factory struct{}

// Trigger is a NATS trigger
type Trigger struct {
	settings *Settings
	conn     *nats.Conn
	log      log.Logger
	handlers []*Handler
}

// Handler is a NATS topic handler
type Handler struct {
	shutdown chan struct{}
	handler  trigger.Handler
	topic    string
	conn     *nats.Conn
	log      log.Logger
}

// Metadata implements trigger.Factory.Metadata
func (*Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// New implements trigger.Factory.New
func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}

	return &Trigger{settings: s}, nil
}

// Initialize initializes the trigger
func (t *Trigger) Initialize(ctx trigger.InitContext) (err error) {
	t.log = ctx.Logger()

	err = t.connectNats()
	if err != nil {
		return err
	}

	for _, h := range ctx.GetHandlers() {
		nh, err := newNATSHandler(ctx, h, t.conn)
		if err != nil {
			return err
		}

		t.handlers = append(t.handlers, nh)
	}

	return nil
}

func (t *Trigger) connectNats() (err error) {
	opts := []nats.Option{
		nats.MaxReconnects(-1),
		nats.NoEcho(),
	}

	if t.settings.Credentials != "" {
		opts = append(opts, nats.UserCredentials(t.settings.Credentials))
	}

	if t.settings.TLSCA != "" && t.settings.TLSKey != "" && t.settings.TLSCert != "" {
		opts = append(opts, nats.ClientCert(t.settings.TLSCert, t.settings.TLSKey))
		opts = append(opts, nats.RootCAs(t.settings.TLSCA))
	}

	t.conn, err = nats.Connect(t.settings.Servers, opts...)

	return err
}

// Start starts the NATS trigger
func (t *Trigger) Start() error {
	for _, h := range t.handlers {
		h.Start()
	}

	return nil
}

// Stop implements ext.Trigger.Stop
func (t *Trigger) Stop() error {
	for _, h := range t.handlers {
		h.Stop()
	}

	t.conn.Close()

	return nil
}

func newNATSHandler(ctx trigger.InitContext, handler trigger.Handler, nc *nats.Conn) (h *Handler, err error) {
	h = &Handler{
		shutdown: make(chan struct{}),
		handler:  handler,
		conn:     nc,
		log:      ctx.Logger(),
	}

	handlerSetting := &HandlerSettings{}
	err = metadata.MapToStruct(handler.Settings(), handlerSetting, true)
	if err != nil {
		return nil, err
	}

	if handlerSetting.Topic == "" {
		return nil, fmt.Errorf("topic string was not provided for handler: [%s]", handler)
	}

	h.topic = handlerSetting.Topic

	return h, nil
}

func (h *Handler) natsSub() {
	msgs := make(chan *nats.Msg, 100)

	sub, err := h.conn.Subscribe(h.topic, func(m *nats.Msg) {
		msgs <- m
	})

	if err != nil {
		h.log.Errorf("subscribe to %s failed: %s", h.topic, err)
		return
	}

	for {
		select {
		case m := <-msgs:
			h.handler.Handle(context.Background(), &Output{string(m.Data)})
		case <-h.shutdown:
			sub.Unsubscribe()
			return
		}
	}
}

func (h *Handler) Start() error {
	go h.natsSub()
}

func (h *Handler) Stop() error {
	close(h.shutdown)
}
