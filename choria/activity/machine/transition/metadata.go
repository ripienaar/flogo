package transition

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct{}

type Input struct {
	Event string `md:"cloudevent,required"`
}

func (r *Input) FromMap(values map[string]interface{}) (err error) {
	r.Event, err = coerce.ToString(values["cloudevent"])
	if err != nil {
		return err
	}

	return nil
}

func (r *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"cloudevent": r.Event,
	}
}

type Output struct {
	Protocol   string `md:"protocol"`
	Identity   string `md:"identity"`
	ID         string `md:"id"`
	Version    string `md:"version"`
	Timestamp  int64  `md:"timestamp"`
	Machine    string `md:"machine"`
	Transition string `md:"transition"`
	FromState  string `md:"from_state"`
	ToState    string `md:"to_state"`
}

func (o *Output) FromMap(values map[string]interface{}) (err error) {
	o.Protocol, err = coerce.ToString(values["protocol"])
	if err != nil {
		return err
	}

	o.Identity, err = coerce.ToString(values["identity"])
	if err != nil {
		return err
	}

	o.ID, err = coerce.ToString(values["id"])
	if err != nil {
		return err
	}

	o.Version, err = coerce.ToString(values["version"])
	if err != nil {
		return err
	}

	o.Timestamp, err = coerce.ToInt64(values["timestamp"])
	if err != nil {
		return err
	}

	o.Machine, err = coerce.ToString(values["machine"])
	if err != nil {
		return err
	}

	o.Transition, err = coerce.ToString(values["transition"])
	if err != nil {
		return err
	}

	o.FromState, err = coerce.ToString(values["from_state"])
	if err != nil {
		return err
	}

	o.ToState, err = coerce.ToString(values["to_state"])
	if err != nil {
		return err
	}

	return err
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"protocol":   o.Protocol,
		"identity":   o.Identity,
		"id":         o.ID,
		"version":    o.Version,
		"timestamp":  o.Timestamp,
		"machine":    o.Machine,
		"transition": o.Transition,
		"from_state": o.FromState,
		"to_state":   o.ToState,
	}
}
