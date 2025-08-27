package assets

import (
	_ "embed"
	"encoding/json"
	"math/rand"
)

//go:embed mantras.json
var mantras []byte

type MantraLoader struct {
	mantras []string
}

func NewMantraLoader() (*MantraLoader, error) {
	var m []string
	if err := json.Unmarshal(mantras, &m); err != nil {
		return nil, err
	}
	return &MantraLoader{mantras: m}, nil
}

func (m *MantraLoader) GetMantra() string {
	return m.mantras[rand.Intn(len(m.mantras))]
}
