package envoy

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
)

type (
	// StoreEncoder is responsible for encoding Corteza resources into the
	// database via the Storer or the DAL interface
	//
	// @todo consider having a different encoder for the DAL resources
	StoreEncoder struct{}
)

// Prepare performs some initial processing on the resource before it can be encoded
//
// Preparation runs validation, default value initialization, matching with
// already existing instances, ...
//
// The prepare function receives a set of nodes grouped by the resource type.
// This enables some batching optimization and simplifications when it comes to
// matching with existing resources.
//
// Prepare does not receive any placeholder nodes which are used solely
// for dependency resolution.
func (e StoreEncoder) Prepare(ctx context.Context, p envoyx.EncodeParams, rt string, nn envoyx.NodeSet) (err error) {
	s, err := e.grabStorer(p)
	if err != nil {
		return
	}

	switch rt {
	case types.WorkflowResourceType:
		return e.prepareWorkflow(ctx, p, s, nn)

	case types.TriggerResourceType:
		return e.prepareTrigger(ctx, p, s, nn)
	default:
		return e.prepare(ctx, p, s, rt, nn)
	}

	return
}

// Encode encodes the given Corteza resources into the primary store
//
// Encoding should not do any additional processing apart from matching with
// dependencies and runtime validation
//
// The Encode function is called for every resource type where the resource
// appears at the root of the dependency tree.
// All of the root-level resources for that resource type are passed into the function.
// The encoding function must traverse the branches to encode all of the dependencies.
//
// This flow is used to simplify the flow of how resources are encoded into YAML
// (and other documents) as well as to simplify batching.
//
// Encode does not receive any placeholder nodes which are used solely
// for dependency resolution.
func (e StoreEncoder) Encode(ctx context.Context, p envoyx.EncodeParams, rt string, nodes envoyx.NodeSet, tree envoyx.Traverser) (err error) {
	s, err := e.grabStorer(p)
	if err != nil {
		return
	}

	switch rt {
	case types.WorkflowResourceType:
		return e.encodeWorkflows(ctx, p, s, nodes, tree)

	case types.TriggerResourceType:
		return e.encodeTriggers(ctx, p, s, nodes, tree)
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource workflow
// // // // // // // // // // // // // // // // // // // // // // // // //

// prepareWorkflow prepares the resources of the given type for encoding
func (e StoreEncoder) prepareWorkflow(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet) (err error) {
	// Grab an index of already existing resources of this type
	// @note since these resources should be fairly low-volume and existing for
	//       a short time (and because we batch by resource type); fetching them all
	//       into memory shouldn't hurt too much.
	// @todo do some benchmarks and potentially implement some smarter check such as
	//       a bloom filter or something similar.

	// Initializing the index here (and using a hashmap) so it's not escaped to the heap
	existing := make(map[int]types.Workflow, len(nn))
	err = e.matchupWorkflows(ctx, s, existing, nn)
	if err != nil {
		return
	}

	for i, n := range nn {
		if n.Resource == nil {
			panic("unexpected state: cannot call prepareWorkflow with nodes without a defined Resource")
		}

		res, ok := n.Resource.(*types.Workflow)
		if !ok {
			panic("unexpected resource type: node expecting type of workflow")
		}

		existing, hasExisting := existing[i]

		if hasExisting {
			// On existing, we don't need to re-do identifiers and references; simply
			// changing up the internal resource is enough.
			//
			// In the future, we can pass down the tree and re-do the deps like that
			switch p.Config.OnExisting {
			case envoyx.OnConflictPanic:
				err = fmt.Errorf("resource already exists")
				return

			case envoyx.OnConflictReplace:
				// Replace; simple ID change should do the trick
				res.ID = existing.ID

			case envoyx.OnConflictSkip:
				// Replace the node's resource with the fetched one
				res = &existing

				// @todo merging
			}
		} else {
			// @todo actually a bottleneck. As per sonyflake docs, it can at most
			//       generate up to 2**8 (256) IDs per 10ms in a single thread.
			//       How can we improve this?
			res.ID = id.Next()
		}

		// We can skip validation/defaults when the resource is overwritten by
		// the one already stored (the panic one errors out anyway) since it
		// should already be ok.
		if !hasExisting || p.Config.OnExisting != envoyx.OnConflictSkip {
			err = e.setWorkflowDefaults(res)
			if err != nil {
				return err
			}

			err = e.validateWorkflow(res)
			if err != nil {
				return err
			}
		}

		n.Resource = res
	}

	return
}

// encodeWorkflows encodes a set of resource into the database
func (e StoreEncoder) encodeWorkflows(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet, tree envoyx.Traverser) (err error) {
	for _, n := range nn {
		err = e.encodeWorkflow(ctx, p, s, n, tree)
		if err != nil {
			return
		}
	}

	return
}

// encodeWorkflow encodes the resource into the database
func (e StoreEncoder) encodeWorkflow(ctx context.Context, p envoyx.EncodeParams, s store.Storer, n *envoyx.Node, tree envoyx.Traverser) (err error) {
	// Grab dependency references
	var auxID uint64
	for fieldLabel, ref := range n.References {
		rn := tree.ParentForRef(n, ref)
		if rn == nil {
			err = fmt.Errorf("missing node for ref %v", ref)
			return
		}

		auxID = rn.Resource.GetID()
		if auxID == 0 {
			err = fmt.Errorf("related resource doesn't provide an ID")
			return
		}

		err = n.Resource.SetValue(fieldLabel, 0, auxID)
		if err != nil {
			return
		}
	}

	// Flush to the DB
	err = store.UpsertAutomationWorkflow(ctx, s, n.Resource.(*types.Workflow))
	if err != nil {
		return
	}

	// Handle resources nested under it
	//
	// @todo how can we remove the OmitPlaceholderNodes call the same way we did for
	//       the root function calls?

	for rt, nn := range envoyx.NodesByResourceType(tree.Children(n)...) {
		nn = envoyx.OmitPlaceholderNodes(nn...)

		switch rt {

		}
	}

	return
}

// matchupWorkflows returns an index with indicates what resources already exist
func (e StoreEncoder) matchupWorkflows(ctx context.Context, s store.Storer, uu map[int]types.Workflow, nn envoyx.NodeSet) (err error) {
	// @todo might need to do it smarter then this.
	//       Most resources won't really be that vast so this should be acceptable for now.
	aa, _, err := store.SearchAutomationWorkflows(ctx, s, types.WorkflowFilter{})
	if err != nil {
		return
	}

	idMap := make(map[uint64]*types.Workflow, len(aa))
	strMap := make(map[string]*types.Workflow, len(aa))

	for _, a := range aa {
		strMap[a.Handle] = a
		idMap[a.ID] = a

	}

	var aux *types.Workflow
	var ok bool
	for i, n := range nn {
		for _, idf := range n.Identifiers.Slice {
			if id, err := strconv.ParseUint(idf, 10, 64); err == nil {
				aux, ok = idMap[id]
				if ok {
					uu[i] = *aux
					// When any identifier matches we can end it
					break
				}
			}

			aux, ok = strMap[idf]
			if ok {
				uu[i] = *aux
				// When any identifier matches we can end it
				break
			}
		}
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource trigger
// // // // // // // // // // // // // // // // // // // // // // // // //

// prepareTrigger prepares the resources of the given type for encoding
func (e StoreEncoder) prepareTrigger(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet) (err error) {
	// Grab an index of already existing resources of this type
	// @note since these resources should be fairly low-volume and existing for
	//       a short time (and because we batch by resource type); fetching them all
	//       into memory shouldn't hurt too much.
	// @todo do some benchmarks and potentially implement some smarter check such as
	//       a bloom filter or something similar.

	// Initializing the index here (and using a hashmap) so it's not escaped to the heap
	existing := make(map[int]types.Trigger, len(nn))
	err = e.matchupTriggers(ctx, s, existing, nn)
	if err != nil {
		return
	}

	for i, n := range nn {
		if n.Resource == nil {
			panic("unexpected state: cannot call prepareTrigger with nodes without a defined Resource")
		}

		res, ok := n.Resource.(*types.Trigger)
		if !ok {
			panic("unexpected resource type: node expecting type of trigger")
		}

		existing, hasExisting := existing[i]

		if hasExisting {
			// On existing, we don't need to re-do identifiers and references; simply
			// changing up the internal resource is enough.
			//
			// In the future, we can pass down the tree and re-do the deps like that
			switch p.Config.OnExisting {
			case envoyx.OnConflictPanic:
				err = fmt.Errorf("resource already exists")
				return

			case envoyx.OnConflictReplace:
				// Replace; simple ID change should do the trick
				res.ID = existing.ID

			case envoyx.OnConflictSkip:
				// Replace the node's resource with the fetched one
				res = &existing

				// @todo merging
			}
		} else {
			// @todo actually a bottleneck. As per sonyflake docs, it can at most
			//       generate up to 2**8 (256) IDs per 10ms in a single thread.
			//       How can we improve this?
			res.ID = id.Next()
		}

		// We can skip validation/defaults when the resource is overwritten by
		// the one already stored (the panic one errors out anyway) since it
		// should already be ok.
		if !hasExisting || p.Config.OnExisting != envoyx.OnConflictSkip {
			err = e.setTriggerDefaults(res)
			if err != nil {
				return err
			}

			err = e.validateTrigger(res)
			if err != nil {
				return err
			}
		}

		n.Resource = res
	}

	return
}

// encodeTriggers encodes a set of resource into the database
func (e StoreEncoder) encodeTriggers(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet, tree envoyx.Traverser) (err error) {
	for _, n := range nn {
		err = e.encodeTrigger(ctx, p, s, n, tree)
		if err != nil {
			return
		}
	}

	return
}

// encodeTrigger encodes the resource into the database
func (e StoreEncoder) encodeTrigger(ctx context.Context, p envoyx.EncodeParams, s store.Storer, n *envoyx.Node, tree envoyx.Traverser) (err error) {
	// Grab dependency references
	var auxID uint64
	for fieldLabel, ref := range n.References {
		rn := tree.ParentForRef(n, ref)
		if rn == nil {
			err = fmt.Errorf("missing node for ref %v", ref)
			return
		}

		auxID = rn.Resource.GetID()
		if auxID == 0 {
			err = fmt.Errorf("related resource doesn't provide an ID")
			return
		}

		err = n.Resource.SetValue(fieldLabel, 0, auxID)
		if err != nil {
			return
		}
	}

	// Flush to the DB
	err = store.UpsertAutomationTrigger(ctx, s, n.Resource.(*types.Trigger))
	if err != nil {
		return
	}

	// Handle resources nested under it
	//
	// @todo how can we remove the OmitPlaceholderNodes call the same way we did for
	//       the root function calls?

	for rt, nn := range envoyx.NodesByResourceType(tree.Children(n)...) {
		nn = envoyx.OmitPlaceholderNodes(nn...)

		switch rt {

		}
	}

	return
}

// matchupTriggers returns an index with indicates what resources already exist
func (e StoreEncoder) matchupTriggers(ctx context.Context, s store.Storer, uu map[int]types.Trigger, nn envoyx.NodeSet) (err error) {
	// @todo might need to do it smarter then this.
	//       Most resources won't really be that vast so this should be acceptable for now.
	aa, _, err := store.SearchAutomationTriggers(ctx, s, types.TriggerFilter{})
	if err != nil {
		return
	}

	idMap := make(map[uint64]*types.Trigger, len(aa))
	strMap := make(map[string]*types.Trigger, len(aa))

	for _, a := range aa {
		idMap[a.ID] = a

	}

	var aux *types.Trigger
	var ok bool
	for i, n := range nn {
		for _, idf := range n.Identifiers.Slice {
			if id, err := strconv.ParseUint(idf, 10, 64); err == nil {
				aux, ok = idMap[id]
				if ok {
					uu[i] = *aux
					// When any identifier matches we can end it
					break
				}
			}

			aux, ok = strMap[idf]
			if ok {
				uu[i] = *aux
				// When any identifier matches we can end it
				break
			}
		}
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utility functions
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e *StoreEncoder) grabStorer(p envoyx.EncodeParams) (s store.Storer, err error) {
	auxs, ok := p.Params["storer"]
	if !ok {
		err = fmt.Errorf("storer not defined")
		return
	}

	s, ok = auxs.(store.Storer)
	if !ok {
		err = fmt.Errorf("invalid storer provided")
		return
	}

	return
}
