package transition

import (
	"encoding/json"

	"github.com/choria-io/go-choria/aagent/machine"
	cloudevents "github.com/cloudevents/sdk-go"
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

	return act, nil
}

type Activity struct {
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

	ce := &cloudevents.Event{}
	err = json.Unmarshal([]byte(input.Event), ce)
	if err != nil {
		return true, err
	}

	t := &machine.TransitionNotification{}
	err = ce.DataAs(t)
	if err != nil {
		return true, err
	}

	output := &Output{
		Protocol:   t.Protocol,
		Identity:   t.Identity,
		ID:         t.ID,
		Version:    t.Version,
		Timestamp:  t.Timestamp,
		Machine:    t.Machine,
		Transition: t.Transition,
		FromState:  t.FromState,
		ToState:    t.ToState,
	}

	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}
