package filter

import (
	"fmt"
	"sync"

	"github.com/ncarlier/feedpushr/autogen/app"
	"github.com/ncarlier/feedpushr/pkg/common"
	"github.com/ncarlier/feedpushr/pkg/model"
	"github.com/ncarlier/feedpushr/pkg/plugin"
	"github.com/rs/zerolog/log"
)

// Chain contains filter chain
type Chain struct {
	filters []model.Filter
	lock    sync.RWMutex
}

// NewChainFilter create new chain filter
func NewChainFilter() *Chain {
	return &Chain{}
}

func newFilter(filter *app.Filter) (model.Filter, error) {
	var _filter model.Filter
	switch filter.Name {
	case "title":
		_filter = newTitleFilter(filter)
	case "fetch":
		_filter = newFetchFilter(filter)
	case "minify":
		_filter = newMinifyFilter(filter)
	default:
		// Try to load plugin regarding the name
		plug := plugin.GetRegsitry().LookupFilterPlugin(filter.Name)
		if plug == nil {
			return nil, fmt.Errorf("unsuported filter: %s", filter.Name)
		}
		var err error
		_filter, err = plug.Build(filter.Props, filter.Tags)
		if err != nil {
			return nil, fmt.Errorf("unable to create filter: %v", err)
		}
	}
	return _filter, nil
}

// GetAvailableFilters get all available filters
func GetAvailableFilters() []model.Spec {
	result := []model.Spec{
		titleSpec,
		fetchSpec,
		minifySpec,
	}
	plugin.GetRegsitry().ForEachFilterPlugin(func(plug model.FilterPlugin) error {
		result = append(result, plug.Spec())
		return nil
	})
	return result
}

// Add a filter to the chain
func (chain *Chain) Add(filter *app.Filter) (model.Filter, error) {
	chain.lock.RLock()
	defer chain.lock.RUnlock()

	log.Debug().Str("name", filter.Name).Msg("creating filter...")
	nextID := 0
	for _, _filter := range chain.filters {
		if nextID < _filter.GetDef().ID {
			nextID = _filter.GetDef().ID
		}
	}
	filter.ID = nextID + 1
	_filter, err := newFilter(filter)
	if err != nil {
		return nil, err
	}

	chain.filters = append(chain.filters, _filter)
	log.Debug().Int("id", filter.ID).Str("name", filter.Name).Msg("filter created")
	return _filter, nil
}

// Update a filter of the chain
func (chain *Chain) Update(filter *app.Filter) (model.Filter, error) {
	chain.lock.RLock()
	defer chain.lock.RUnlock()

	for idx, _filter := range chain.filters {
		if filter.ID == _filter.GetDef().ID {
			log.Debug().Int("id", filter.ID).Msg("updating filter...")
			// TODO merge objects
			filter.Name = _filter.GetDef().Name
			f, err := newFilter(filter)
			if err != nil {
				return nil, err
			}
			chain.filters[idx] = f
			log.Debug().Int("id", filter.ID).Msg("filter updated")
			return f, nil
		}
	}
	return nil, common.ErrFilterNotFound
}

// Remove a filter from the chain
func (chain *Chain) Remove(filter *app.Filter) error {
	chain.lock.RLock()
	defer chain.lock.RUnlock()

	for idx, _filter := range chain.filters {
		if filter.ID == _filter.GetDef().ID {
			log.Debug().Int("id", filter.ID).Msg("removing filter...")
			chain.filters = append(chain.filters[:idx], chain.filters[idx+1:]...)
			log.Debug().Int("id", filter.ID).Msg("filter removed")
			return nil
		}
	}
	return common.ErrFilterNotFound
}

// Apply applies filter chain on an article
func (chain *Chain) Apply(article *model.Article) error {
	for idx, filter := range chain.filters {
		tags := filter.GetDef().Tags
		if !filter.GetDef().Enabled || !article.Match(tags) {
			// Ignore disabled filters or that do not match the article tags
			continue
		}
		err := filter.DoFilter(article)
		if err != nil {
			return fmt.Errorf("error while applying filter #%d: %v", idx, err)
		}
	}
	return nil
}

// Get a filter from the chain
func (chain *Chain) Get(id int) (model.Filter, error) {
	for _, _filter := range chain.filters {
		if id == _filter.GetDef().ID {
			return _filter, nil
		}
	}
	return nil, common.ErrFilterNotFound
}

// GetFilterDefs return definitions of the chain filter
func (chain *Chain) GetFilterDefs() []model.FilterDef {
	result := make([]model.FilterDef, len(chain.filters))
	for idx, filter := range chain.filters {
		result[idx] = filter.GetDef()
	}
	return result
}
