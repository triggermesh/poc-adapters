package fixedwithtojson

import pkgadapter "knative.dev/eventing/pkg/adapter/v2"

type FixedWithJSONRepresentation struct {
	Fields []Field `json:"fields"`
}

type Field struct {
	Value      string `json:"value"`
	Spaceleft  int    `json:"spaceLeft"`
	LineNumber int    `json:"lineNumber"`
}

type envAccessor struct {
	pkgadapter.EnvConfig
	BridgeIdentifier string `envconfig:"EVENTS_BRIDGE_IDENTIFIER"`
	// CloudEvents responses parametrization
	CloudEventPayloadPolicy string `envconfig:"EVENTS_PAYLOAD_POLICY" default:"error"`
	// Sink defines the target sink for the events. If no Sink is defined the
	// events are replied back to the sender.
	Sink string `envconfig:"K_SINK"`
}
