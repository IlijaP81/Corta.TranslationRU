package envoy

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/pkg/errors"
)

type (
	// StoreDecoder is responsible for fetching already stored Corteza resources
	// which are then managed by envoy and imported via an encoder.
	StoreDecoder struct{}
)

// Decode returns a set of envoy nodes based on the provided params
//
// StoreDecoder expects the DecodeParam of `storer` and `dal` which conform
// to the store.Storer and dal.FullService interfaces.
func (d StoreDecoder) Decode(ctx context.Context, p envoyx.DecodeParams) (out envoyx.NodeSet, err error) {
	var (
		s  store.Storer
		dl dal.FullService
	)

	// @todo we can optionally not require them based on what we're doing
	if auxS, ok := p.Params["storer"]; ok {
		s = auxS.(store.Storer)
	}
	if auxDl, ok := p.Params["dal"]; ok {
		dl = auxDl.(dal.FullService)
	}

	return d.decode(ctx, s, dl, p)
}

func (d StoreDecoder) decode(ctx context.Context, s store.Storer, dl dal.FullService, p envoyx.DecodeParams) (out envoyx.NodeSet, err error) {
	// Transform passed filters into an ordered structure
	type (
		filterWrap struct {
			rt string
			f  envoyx.ResourceFilter
		}
	)
	wrappedFilters := make([]filterWrap, 0, len(p.Filter))
	for rt, f := range p.Filter {
		// Handle resources that don't belong to this decoder
		if !strings.HasPrefix(rt, "corteza::system") {
			continue
		}

		wrappedFilters = append(wrappedFilters, filterWrap{rt: rt, f: f})
	}

	// Get all requested scopes
	scopedNodes := make(envoyx.NodeSet, len(p.Filter))

	// @note skipping scope logic since it's currently only supported within
	//       Compose resources.

	// Get all requested references
	//
	// Keep an index for the Node and one for the reference to make our
	// lives easier.
	refNodes := make([]map[string]*envoyx.Node, len(p.Filter))
	refRefs := make([]map[string]envoyx.Ref, len(p.Filter))
	err = func() (err error) {
		for i, a := range wrappedFilters {
			if len(a.f.Refs) == 0 {
				continue
			}

			auxr := make(map[string]*envoyx.Node, len(a.f.Refs))
			auxa := make(map[string]envoyx.Ref)
			for field, ref := range a.f.Refs {
				f := ref.ResourceFilter()
				aux, err := d.decode(ctx, s, dl, envoyx.DecodeParams{
					Type:   envoyx.DecodeTypeStore,
					Filter: f,
				})
				if err != nil {
					return err
				}

				// @todo consider changing this.
				//       Currently it's required because the .decode may return some
				//       nested nodes as well.
				//       Consider a flag or a new function.
				aux = envoyx.NodesForResourceType(ref.ResourceType, aux...)
				if len(aux) == 0 {
					return fmt.Errorf("invalid reference %v", ref)
				}
				if len(aux) > 1 {
					return fmt.Errorf("ambiguous reference: too many resources returned %v", a.f)
				}

				auxr[field] = aux[0]
				auxa[field] = aux[0].ToRef()
			}

			refNodes[i] = auxr
			refRefs[i] = auxa
		}
		return
	}()
	if err != nil {
		err = errors.Wrap(err, "failed to decode node references")
		return
	}

	err = func() (err error) {
		var aux envoyx.NodeSet
		for i, wf := range wrappedFilters {
			switch wf.rt {
			case types.ApplicationResourceType:
				aux, err = d.decodeApplication(ctx, s, dl, d.makeApplicationFilter(scopedNodes[i], refNodes[i], wf.f))
				if err != nil {
					return
				}
				for _, a := range aux {
					a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
					a.References = envoyx.MergeRefs(a.References, refRefs[i])
				}
				out = append(out, aux...)

			case types.ApigwRouteResourceType:
				aux, err = d.decodeApigwRoute(ctx, s, dl, d.makeApigwRouteFilter(scopedNodes[i], refNodes[i], wf.f))
				if err != nil {
					return
				}
				for _, a := range aux {
					a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
					a.References = envoyx.MergeRefs(a.References, refRefs[i])
				}
				out = append(out, aux...)

			case types.ApigwFilterResourceType:
				aux, err = d.decodeApigwFilter(ctx, s, dl, d.makeApigwFilterFilter(scopedNodes[i], refNodes[i], wf.f))
				if err != nil {
					return
				}
				for _, a := range aux {
					a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
					a.References = envoyx.MergeRefs(a.References, refRefs[i])
				}
				out = append(out, aux...)

			case types.AuthClientResourceType:
				aux, err = d.decodeAuthClient(ctx, s, dl, d.makeAuthClientFilter(scopedNodes[i], refNodes[i], wf.f))
				if err != nil {
					return
				}
				for _, a := range aux {
					a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
					a.References = envoyx.MergeRefs(a.References, refRefs[i])
				}
				out = append(out, aux...)

			case types.QueueResourceType:
				aux, err = d.decodeQueue(ctx, s, dl, d.makeQueueFilter(scopedNodes[i], refNodes[i], wf.f))
				if err != nil {
					return
				}
				for _, a := range aux {
					a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
					a.References = envoyx.MergeRefs(a.References, refRefs[i])
				}
				out = append(out, aux...)

			case types.ReportResourceType:
				aux, err = d.decodeReport(ctx, s, dl, d.makeReportFilter(scopedNodes[i], refNodes[i], wf.f))
				if err != nil {
					return
				}
				for _, a := range aux {
					a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
					a.References = envoyx.MergeRefs(a.References, refRefs[i])
				}
				out = append(out, aux...)

			case types.RoleResourceType:
				aux, err = d.decodeRole(ctx, s, dl, d.makeRoleFilter(scopedNodes[i], refNodes[i], wf.f))
				if err != nil {
					return
				}
				for _, a := range aux {
					a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
					a.References = envoyx.MergeRefs(a.References, refRefs[i])
				}
				out = append(out, aux...)

			case types.TemplateResourceType:
				aux, err = d.decodeTemplate(ctx, s, dl, d.makeTemplateFilter(scopedNodes[i], refNodes[i], wf.f))
				if err != nil {
					return
				}
				for _, a := range aux {
					a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
					a.References = envoyx.MergeRefs(a.References, refRefs[i])
				}
				out = append(out, aux...)

			case types.UserResourceType:
				aux, err = d.decodeUser(ctx, s, dl, d.makeUserFilter(scopedNodes[i], refNodes[i], wf.f))
				if err != nil {
					return
				}
				for _, a := range aux {
					a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
					a.References = envoyx.MergeRefs(a.References, refRefs[i])
				}
				out = append(out, aux...)

			case types.DalConnectionResourceType:
				aux, err = d.decodeDalConnection(ctx, s, dl, d.makeDalConnectionFilter(scopedNodes[i], refNodes[i], wf.f))
				if err != nil {
					return
				}
				for _, a := range aux {
					a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
					a.References = envoyx.MergeRefs(a.References, refRefs[i])
				}
				out = append(out, aux...)

			case types.DalSensitivityLevelResourceType:
				aux, err = d.decodeDalSensitivityLevel(ctx, s, dl, d.makeDalSensitivityLevelFilter(scopedNodes[i], refNodes[i], wf.f))
				if err != nil {
					return
				}
				for _, a := range aux {
					a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
					a.References = envoyx.MergeRefs(a.References, refRefs[i])
				}
				out = append(out, aux...)

			default:
				aux, err = d.extendDecoder(ctx, s, dl, wf.rt, refNodes[i], wf.f)
				if err != nil {
					return
				}
				for _, a := range aux {
					a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
					a.References = envoyx.MergeRefs(a.References, refRefs[i])
				}
				out = append(out, aux...)
			}
		}
		return
	}()
	if err != nil {
		err = errors.Wrap(err, "failed to decode filters")
		return
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource application
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) decodeApplication(ctx context.Context, s store.Storer, dl dal.FullService, f types.ApplicationFilter) (out envoyx.NodeSet, err error) {
	// @todo this might need to be improved.
	//       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.SearchApplications(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		// Identifiers
		ii := envoyx.MakeIdentifiers(
			r.ID,
		)

		// Handle references
		// Omit any non-defined values
		refs := map[string]envoyx.Ref{}
		if r.OwnerID > 0 {
			refs["OwnerID"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.OwnerID),
			}
		}

		var scope envoyx.Scope

		out = append(out, &envoyx.Node{
			Resource: r,

			ResourceType: types.ApplicationResourceType,
			Identifiers:  ii,
			References:   refs,
			Scope:        scope,
		})
	}

	return
}

func (d StoreDecoder) makeApplicationFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.ApplicationFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.ApplicationID = ids

	if len(hh) > 0 {
		out.Name = hh[0]
	}

	// Refs
	var (
		ar *envoyx.Node
		ok bool
	)
	_ = ar
	_ = ok

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource apigwRoute
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) decodeApigwRoute(ctx context.Context, s store.Storer, dl dal.FullService, f types.ApigwRouteFilter) (out envoyx.NodeSet, err error) {
	// @todo this might need to be improved.
	//       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.SearchApigwRoutes(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		// Identifiers
		ii := envoyx.MakeIdentifiers(
			r.ID,
		)

		// Handle references
		// Omit any non-defined values
		refs := map[string]envoyx.Ref{}
		if r.CreatedBy > 0 {
			refs["CreatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.CreatedBy),
			}
		}
		if r.DeletedBy > 0 {
			refs["DeletedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.DeletedBy),
			}
		}
		if r.Group > 0 {
			refs["Group"] = envoyx.Ref{
				ResourceType: "corteza::system:apigw-group",
				Identifiers:  envoyx.MakeIdentifiers(r.Group),
			}
		}
		if r.UpdatedBy > 0 {
			refs["UpdatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.UpdatedBy),
			}
		}

		var scope envoyx.Scope

		out = append(out, &envoyx.Node{
			Resource: r,

			ResourceType: types.ApigwRouteResourceType,
			Identifiers:  ii,
			References:   refs,
			Scope:        scope,
		})
	}

	aux, err := d.extendedApigwRouteDecoder(ctx, s, dl, f, out)
	if err != nil {
		return
	}
	out = append(out, aux...)

	return
}

func (d StoreDecoder) makeApigwRouteFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.ApigwRouteFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.ApigwRouteID = ids

	if len(hh) > 0 {
		out.Endpoint = hh[0]
	}

	// Refs
	var (
		ar *envoyx.Node
		ok bool
	)
	_ = ar
	_ = ok

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource apigwFilter
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) decodeApigwFilter(ctx context.Context, s store.Storer, dl dal.FullService, f types.ApigwFilterFilter) (out envoyx.NodeSet, err error) {
	// @todo this might need to be improved.
	//       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.SearchApigwFilters(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		// Identifiers
		ii := envoyx.MakeIdentifiers(
			r.ID,
		)

		// Handle references
		// Omit any non-defined values
		refs := map[string]envoyx.Ref{}
		if r.CreatedBy > 0 {
			refs["CreatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.CreatedBy),
			}
		}
		if r.DeletedBy > 0 {
			refs["DeletedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.DeletedBy),
			}
		}
		if r.Route > 0 {
			refs["Route"] = envoyx.Ref{
				ResourceType: "corteza::system:apigw-route",
				Identifiers:  envoyx.MakeIdentifiers(r.Route),
			}
		}
		if r.UpdatedBy > 0 {
			refs["UpdatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.UpdatedBy),
			}
		}

		var scope envoyx.Scope

		out = append(out, &envoyx.Node{
			Resource: r,

			ResourceType: types.ApigwFilterResourceType,
			Identifiers:  ii,
			References:   refs,
			Scope:        scope,
		})
	}

	return
}

func (d StoreDecoder) makeApigwFilterFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.ApigwFilterFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.ApigwFilterID = ids

	// Refs
	var (
		ar *envoyx.Node
		ok bool
	)
	_ = ar
	_ = ok

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource authClient
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) decodeAuthClient(ctx context.Context, s store.Storer, dl dal.FullService, f types.AuthClientFilter) (out envoyx.NodeSet, err error) {
	// @todo this might need to be improved.
	//       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.SearchAuthClients(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		// Identifiers
		ii := envoyx.MakeIdentifiers(
			r.Handle,
			r.ID,
		)

		// Handle references
		// Omit any non-defined values
		refs := map[string]envoyx.Ref{}
		if r.CreatedBy > 0 {
			refs["CreatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.CreatedBy),
			}
		}
		if r.DeletedBy > 0 {
			refs["DeletedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.DeletedBy),
			}
		}
		if r.OwnedBy > 0 {
			refs["OwnedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.OwnedBy),
			}
		}
		if r.UpdatedBy > 0 {
			refs["UpdatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.UpdatedBy),
			}
		}

		refs = envoyx.MergeRefs(refs, d.decodeAuthClientRefs(r))

		var scope envoyx.Scope

		out = append(out, &envoyx.Node{
			Resource: r,

			ResourceType: types.AuthClientResourceType,
			Identifiers:  ii,
			References:   refs,
			Scope:        scope,
		})
	}

	return
}

func (d StoreDecoder) makeAuthClientFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.AuthClientFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.AuthClientID = ids

	if len(hh) > 0 {
		out.Handle = hh[0]
	}

	// Refs
	var (
		ar *envoyx.Node
		ok bool
	)
	_ = ar
	_ = ok

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource queue
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) decodeQueue(ctx context.Context, s store.Storer, dl dal.FullService, f types.QueueFilter) (out envoyx.NodeSet, err error) {
	// @todo this might need to be improved.
	//       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.SearchQueues(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		// Identifiers
		ii := envoyx.MakeIdentifiers(
			r.ID,
			r.Queue,
		)

		// Handle references
		// Omit any non-defined values
		refs := map[string]envoyx.Ref{}
		if r.CreatedBy > 0 {
			refs["CreatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.CreatedBy),
			}
		}
		if r.DeletedBy > 0 {
			refs["DeletedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.DeletedBy),
			}
		}
		if r.UpdatedBy > 0 {
			refs["UpdatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.UpdatedBy),
			}
		}

		var scope envoyx.Scope

		out = append(out, &envoyx.Node{
			Resource: r,

			ResourceType: types.QueueResourceType,
			Identifiers:  ii,
			References:   refs,
			Scope:        scope,
		})
	}

	return
}

func (d StoreDecoder) makeQueueFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.QueueFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.QueueID = ids

	// Refs
	var (
		ar *envoyx.Node
		ok bool
	)
	_ = ar
	_ = ok

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource report
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) decodeReport(ctx context.Context, s store.Storer, dl dal.FullService, f types.ReportFilter) (out envoyx.NodeSet, err error) {
	// @todo this might need to be improved.
	//       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.SearchReports(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		// Identifiers
		ii := envoyx.MakeIdentifiers(
			r.Handle,
			r.ID,
		)

		// Handle references
		// Omit any non-defined values
		refs := map[string]envoyx.Ref{}
		if r.CreatedBy > 0 {
			refs["CreatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.CreatedBy),
			}
		}
		if r.DeletedBy > 0 {
			refs["DeletedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.DeletedBy),
			}
		}
		if r.OwnedBy > 0 {
			refs["OwnedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.OwnedBy),
			}
		}
		if r.UpdatedBy > 0 {
			refs["UpdatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.UpdatedBy),
			}
		}

		var scope envoyx.Scope

		out = append(out, &envoyx.Node{
			Resource: r,

			ResourceType: types.ReportResourceType,
			Identifiers:  ii,
			References:   refs,
			Scope:        scope,
		})
	}

	return
}

func (d StoreDecoder) makeReportFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.ReportFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.ReportID = ids

	if len(hh) > 0 {
		out.Handle = hh[0]
	}

	// Refs
	var (
		ar *envoyx.Node
		ok bool
	)
	_ = ar
	_ = ok

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource role
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) decodeRole(ctx context.Context, s store.Storer, dl dal.FullService, f types.RoleFilter) (out envoyx.NodeSet, err error) {
	// @todo this might need to be improved.
	//       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.SearchRoles(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		// Identifiers
		ii := envoyx.MakeIdentifiers(
			r.Handle,
			r.ID,
		)

		// Handle references
		// Omit any non-defined values
		refs := map[string]envoyx.Ref{}

		var scope envoyx.Scope

		out = append(out, &envoyx.Node{
			Resource: r,

			ResourceType: types.RoleResourceType,
			Identifiers:  ii,
			References:   refs,
			Scope:        scope,
		})
	}

	return
}

func (d StoreDecoder) makeRoleFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.RoleFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.RoleID = ids

	if len(hh) > 0 {
		out.Handle = hh[0]
	}

	// Refs
	var (
		ar *envoyx.Node
		ok bool
	)
	_ = ar
	_ = ok

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource template
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) decodeTemplate(ctx context.Context, s store.Storer, dl dal.FullService, f types.TemplateFilter) (out envoyx.NodeSet, err error) {
	// @todo this might need to be improved.
	//       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.SearchTemplates(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		// Identifiers
		ii := envoyx.MakeIdentifiers(
			r.Handle,
			r.ID,
		)

		// Handle references
		// Omit any non-defined values
		refs := map[string]envoyx.Ref{}
		if r.OwnerID > 0 {
			refs["OwnerID"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.OwnerID),
			}
		}

		var scope envoyx.Scope

		out = append(out, &envoyx.Node{
			Resource: r,

			ResourceType: types.TemplateResourceType,
			Identifiers:  ii,
			References:   refs,
			Scope:        scope,
		})
	}

	return
}

func (d StoreDecoder) makeTemplateFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.TemplateFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.TemplateID = ids

	if len(hh) > 0 {
		out.Handle = hh[0]
	}

	// Refs
	var (
		ar *envoyx.Node
		ok bool
	)
	_ = ar
	_ = ok

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource user
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) decodeUser(ctx context.Context, s store.Storer, dl dal.FullService, f types.UserFilter) (out envoyx.NodeSet, err error) {
	// @todo this might need to be improved.
	//       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.SearchUsers(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		// Identifiers
		ii := envoyx.MakeIdentifiers(
			r.Handle,
			r.ID,
		)

		// Handle references
		// Omit any non-defined values
		refs := map[string]envoyx.Ref{}

		var scope envoyx.Scope

		out = append(out, &envoyx.Node{
			Resource: r,

			ResourceType: types.UserResourceType,
			Identifiers:  ii,
			References:   refs,
			Scope:        scope,
		})
	}

	return
}

func (d StoreDecoder) makeUserFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.UserFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.UserID = ids

	if len(hh) > 0 {
		out.Handle = hh[0]
	}

	// Refs
	var (
		ar *envoyx.Node
		ok bool
	)
	_ = ar
	_ = ok

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource dalConnection
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) decodeDalConnection(ctx context.Context, s store.Storer, dl dal.FullService, f types.DalConnectionFilter) (out envoyx.NodeSet, err error) {
	// @todo this might need to be improved.
	//       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.SearchDalConnections(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		// Identifiers
		ii := envoyx.MakeIdentifiers(
			r.Handle,
			r.ID,
		)

		// Handle references
		// Omit any non-defined values
		refs := map[string]envoyx.Ref{}
		if r.CreatedBy > 0 {
			refs["CreatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.CreatedBy),
			}
		}
		if r.DeletedBy > 0 {
			refs["DeletedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.DeletedBy),
			}
		}
		if r.UpdatedBy > 0 {
			refs["UpdatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.UpdatedBy),
			}
		}

		refs = envoyx.MergeRefs(refs, d.decodeDalConnectionRefs(r))

		var scope envoyx.Scope

		out = append(out, &envoyx.Node{
			Resource: r,

			ResourceType: types.DalConnectionResourceType,
			Identifiers:  ii,
			References:   refs,
			Scope:        scope,
		})
	}

	return
}

func (d StoreDecoder) makeDalConnectionFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.DalConnectionFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.DalConnectionID = ids

	if len(hh) > 0 {
		out.Handle = hh[0]
	}

	// Refs
	var (
		ar *envoyx.Node
		ok bool
	)
	_ = ar
	_ = ok

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource dalSensitivityLevel
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) decodeDalSensitivityLevel(ctx context.Context, s store.Storer, dl dal.FullService, f types.DalSensitivityLevelFilter) (out envoyx.NodeSet, err error) {
	// @todo this might need to be improved.
	//       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.SearchDalSensitivityLevels(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		// Identifiers
		ii := envoyx.MakeIdentifiers(
			r.Handle,
			r.ID,
		)

		// Handle references
		// Omit any non-defined values
		refs := map[string]envoyx.Ref{}
		if r.CreatedBy > 0 {
			refs["CreatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.CreatedBy),
			}
		}
		if r.DeletedBy > 0 {
			refs["DeletedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.DeletedBy),
			}
		}
		if r.UpdatedBy > 0 {
			refs["UpdatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(r.UpdatedBy),
			}
		}

		var scope envoyx.Scope

		out = append(out, &envoyx.Node{
			Resource: r,

			ResourceType: types.DalSensitivityLevelResourceType,
			Identifiers:  ii,
			References:   refs,
			Scope:        scope,
		})
	}

	return
}

func (d StoreDecoder) makeDalSensitivityLevelFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.DalSensitivityLevelFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.DalSensitivityLevelID = ids

	if len(hh) > 0 {
		out.Handle = hh[0]
	}

	// Refs
	var (
		ar *envoyx.Node
		ok bool
	)
	_ = ar
	_ = ok

	return
}
