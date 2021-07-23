package report

import (
	"context"
	"errors"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	model struct {
		ran         bool
		steps       []step
		datasources DatasourceSet
	}

	// M is the model interface that should be used when trying to model the datasource
	M interface {
		Add(...step) M
		Run(context.Context) error
		Load(context.Context, ...*FrameDefinition) ([]*Frame, error)
		Describe(source string) (FrameDescriptionSet, error)
	}

	stepSet []step
	step    interface {
		Name() string
		Source() []string
		Run(context.Context, ...Datasource) (Datasource, error)
		Validate() error
		Def() *StepDefinition
	}

	StepDefinitionSet []*StepDefinition
	StepDefinition    struct {
		Load  *LoadStepDefinition  `json:"load,omitempty"`
		Join  *JoinStepDefinition  `json:"join,omitempty"`
		Group *GroupStepDefinition `json:"group,omitempty"`
		// @todo Transform
	}
)

// Model initializes the model based on the provided sources and step definitions.
//
// Additional steps may be added after the model is constructed.
// Call `M.Run(context.Context)` to allow the model to be used for requesting data.
// Additional steps may not be added after the `M.Run(context.Context)` was called
func Model(ctx context.Context, sources map[string]DatasourceProvider, dd ...*StepDefinition) (M, error) {
	steps := make([]step, 0, len(dd))
	ss := make(DatasourceSet, 0, len(steps)*2)

	err := func() error {
		for _, d := range dd {
			switch {
			case d.Load != nil:
				if sources == nil {
					return errors.New("no datasources defined")
				}

				s, ok := sources[d.Load.Source]
				if !ok {
					return fmt.Errorf("unresolved datasource: %s", d.Load.Source)
				}
				ds, err := s.Datasource(ctx, d.Load)
				if err != nil {
					return err
				}

				ss = append(ss, ds)

			case d.Join != nil:
				steps = append(steps, &stepJoin{def: d.Join})

			case d.Group != nil:
				steps = append(steps, &stepGroup{def: d.Group})

			// @todo Transform

			default:
				return errors.New("malformed step definition: unsupported step kind")
			}
		}
		return nil
	}()

	if err != nil {
		return nil, fmt.Errorf("failed to create the model: %s", err.Error())
	}

	return &model{
		steps:       steps,
		datasources: ss,
	}, nil
}

// Add adds additional steps to the model
func (m *model) Add(ss ...step) M {
	m.steps = append(m.steps, ss...)
	return m
}

// Run bakes the model configuration and makes the requested data available
func (m *model) Run(ctx context.Context) (err error) {
	const errPfx = "failed to run the model"
	defer func() {
		m.ran = true
	}()

	// initial validation
	err = func() (err error) {
		if m.ran {
			return errors.New("model already ran")
		}

		if len(m.steps)+len(m.datasources) == 0 {
			return errors.New("no model steps defined")
		}

		for _, s := range m.steps {
			err = s.Validate()
			if err != nil {
				return err
			}
		}

		return nil
	}()
	if err != nil {
		return fmt.Errorf("%s: failed to validate the model: %w", errPfx, err)
	}

	// construct the step graph
	//
	// If there are no steps, there is nothing to reduce
	if len(m.steps) == 0 {
		return nil
	}
	err = func() (err error) {
		gg, err := m.buildStepGraph(m.steps, m.datasources)
		if err != nil {
			return err
		}

		m.datasources = nil
		for _, n := range gg {
			aux, err := m.reduceGraph(ctx, n)
			if err != nil {
				return err
			}
			m.datasources = append(m.datasources, aux)
		}
		return nil
	}()
	if err != nil {
		return fmt.Errorf("%s: %w", errPfx, err)
	}

	return nil
}

// Describe returns the descriptions for the requested model datasources
//
// The Run method must be called before the description can be provided.
func (m *model) Describe(source string) (out FrameDescriptionSet, err error) {
	var ds Datasource

	err = func() error {
		if !m.ran {
			return fmt.Errorf("model was not yet ran")
		}

		ds := m.datasources.Find(source)
		if ds == nil {
			return fmt.Errorf("model does not contain the datasource: %s", source)
		}

		return nil
	}()
	if err != nil {
		return nil, fmt.Errorf("unable to describe the model source: %w", err)
	}

	return ds.Describe(), nil
}

// Load returns the Frames based on the provided FrameDefinitions
//
// The Run method must be called before the frames can be provided.
func (m *model) Load(ctx context.Context, dd ...*FrameDefinition) (ff []*Frame, err error) {
	var (
		def *FrameDefinition
		ds  Datasource
	)

	// request validation
	err = func() error {
		// - all frame definitions must define the same datasource; call Load multiple times if
		//   you need to access multiple datasources
		for i, d := range dd {
			if i == 0 {
				continue
			}
			if d.Source != dd[i-1].Source {
				return fmt.Errorf("frame definition source missmatch: expected %s, got %s", dd[i-1].Source, d.Source)
			}
		}

		def = dd[0]

		ds = m.datasources.Find(def.Source)
		if ds == nil {
			return fmt.Errorf("unresolved datasource: %s", def.Source)
		}

		return nil
	}()
	if err != nil {
		return nil, fmt.Errorf("unable to load frames: invalid request: %w", err)
	}

	// apply any frame definition defaults
	aux := make([]*FrameDefinition, len(dd))
	for i, d := range dd {
		aux[i] = d.Clone()

		// assure paging is always provided so we can ignore nil checks
		if aux[i].Paging == nil {
			aux[i].Paging = &filter.Paging{
				Limit: defaultPageSize,
			}
		}

		// assure sorting is always provided so we can ignore nil checks
		if aux[i].Sorting == nil {
			aux[i].Sorting = filter.SortExprSet{}
		}
	}
	dd = aux

	// assure paging is always provided so we can ignore nil checks
	if def.Paging == nil {
		def.Paging = &filter.Paging{
			Limit: defaultPageSize,
		}
	}

	// assure sorting is always provided so we can ignore nil checks
	if def.Sorting == nil {
		def.Sorting = filter.SortExprSet{}
	}

	// load the data
	err = func() error {
		l, c, err := ds.Load(ctx, dd...)
		if err != nil {
			return err
		}
		defer c()

		ff, err = l(int(def.Paging.Limit))
		if err != nil {
			return err
		}

		return nil
	}()
	if err != nil {
		return nil, fmt.Errorf("unable to load frames: %w", err)
	}

	return ff, nil
}
