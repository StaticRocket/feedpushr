package output

import (
	"fmt"

	"github.com/ncarlier/feedpushr/autogen/app"
	"github.com/ncarlier/feedpushr/pkg/common"
	"github.com/ncarlier/feedpushr/pkg/model"
	"github.com/ncarlier/feedpushr/pkg/plugin"
)

// GetAvailableOutputs get all available outputs
func GetAvailableOutputs() []model.Spec {
	result := []model.Spec{
		stdoutSpec,
		httpSpec,
	}
	plugin.GetRegsitry().ForEachOutputPlugin(func(plug model.OutputPlugin) error {
		result = append(result, plug.Spec())
		return nil
	})
	return result
}

// newOutputProvider creates new output provider.
func newOutputProvider(output *app.Output) (model.OutputProvider, error) {
	var provider model.OutputProvider
	switch output.Name {
	case "stdout":
		provider = newStdOutputProvider(output)
	case "http":
		provider = newHTTPOutputProvider(output)
	default:
		// Try to load plugin regarding the scheme
		plug := plugin.GetRegsitry().LookupOutputPlugin(output.Name)
		if plug == nil {
			return nil, fmt.Errorf("unsuported output provider: %s", output.Name)
		}
		var err error
		provider, err = plug.Build(output.Props, output.Tags)
		if err != nil {
			return nil, fmt.Errorf("unable to create output provider: %v", err)
		}
	}
	return provider, nil
}

// Add an output
func (m *Manager) Add(output *app.Output) (model.OutputProvider, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	m.log.Debug().Str("name", output.Name).Msg("adding output...")
	nextID := 0
	for _, _output := range m.providers {
		if nextID < _output.GetDef().ID {
			nextID = _output.GetDef().ID
		}
	}
	output.ID = nextID + 1
	provider, err := newOutputProvider(output)
	if err != nil {
		return nil, err
	}
	m.providers = append(m.providers, provider)
	m.log.Debug().Str("name", output.Name).Msg("output added")
	return provider, nil
}

// Update an output
func (m *Manager) Update(output *app.Output) (model.OutputProvider, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for idx, provider := range m.providers {
		if output.ID == provider.GetDef().ID {
			// TODO merge objects
			output.Name = provider.GetDef().Name
			m.log.Debug().Int("id", output.ID).Msg("updating output...")
			p, err := newOutputProvider(output)
			if err != nil {
				return nil, err
			}
			m.providers[idx] = p
			m.log.Debug().Int("id", output.ID).Msg("output updated")
			return p, nil
		}
	}
	return nil, common.ErrOutputNotFound
}

// Remove an output
func (m *Manager) Remove(output *app.Output) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for idx, provider := range m.providers {
		if output.ID == provider.GetDef().ID {
			m.log.Debug().Int("id", output.ID).Msg("removing output...")
			m.providers = append(m.providers[:idx], m.providers[idx+1:]...)
			m.log.Debug().Int("id", output.ID).Msg("output removed")
			return nil
		}
	}
	return common.ErrOutputNotFound
}

// Get an output
func (m *Manager) Get(id int) (model.OutputProvider, error) {
	for _, provider := range m.providers {
		if id == provider.GetDef().ID {
			return provider, nil
		}
	}
	return nil, common.ErrOutputNotFound
}
