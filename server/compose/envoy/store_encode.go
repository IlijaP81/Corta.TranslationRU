package envoy

import (
	"context"
	"fmt"
	"time"

	"github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/store"
)

func (e StoreEncoder) setChartDefaults(res *types.Chart) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}
	return
}

func (e StoreEncoder) validateChart(*types.Chart) (err error) {
	return
}

func (e StoreEncoder) setModuleDefaults(res *types.Module) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}
	return
}

func (e StoreEncoder) validateModule(*types.Module) (err error) {
	return
}

func (e StoreEncoder) setModuleFieldDefaults(res *types.ModuleField) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	// Update validator ID
	maxValidatorID := uint64(0)
	for _, v := range res.Expressions.Validators {
		if v.ValidatorID > maxValidatorID {
			maxValidatorID = v.ValidatorID
		}
	}

	for _, v := range res.Expressions.Validators {
		if v.ValidatorID == 0 {
			v.ValidatorID = maxValidatorID + 1
			maxValidatorID++
		}
	}

	return
}

func (e StoreEncoder) validateModuleField(*types.ModuleField) (err error) {
	return
}

func (e StoreEncoder) setNamespaceDefaults(res *types.Namespace) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}
	return
}

func (e StoreEncoder) validateNamespace(*types.Namespace) (err error) {
	return
}

func (e StoreEncoder) setPageDefaults(res *types.Page) (err error) {
	if res.CreatedAt.IsZero() {
		res.CreatedAt = time.Now()
	}

	if res.Title == "" {
		res.Title = res.Handle
	}

	// Update pageblock ID
	maxPageBlockID := uint64(0)
	for _, b := range res.Blocks {
		if b.BlockID > maxPageBlockID {
			maxPageBlockID = b.BlockID
		}
	}

	for _, b := range res.Blocks {
		if b.BlockID == 0 {
			b.BlockID = maxPageBlockID + 1
			maxPageBlockID++
		}
	}

	return
}

func (e StoreEncoder) validatePage(*types.Page) (err error) {
	return
}

func (e StoreEncoder) setRecordDefaults(*types.Record) (err error) {
	return
}

func (e StoreEncoder) validateRecord(*types.Record) (err error) {
	return
}

func (e StoreEncoder) prepare(ctx context.Context, p envoyx.EncodeParams, s store.Storer, rt string, nn envoyx.NodeSet) (err error) {
	switch rt {
	case ComposeRecordDatasourceAuxType:
		return e.prepareRecordDatasource(ctx, p, s, nn)
	}

	return
}

func (e StoreEncoder) encodeModuleExtend(ctx context.Context, p envoyx.EncodeParams, s store.Storer, n *envoyx.Node, nested envoyx.NodeSet, tree envoyx.Traverser) (err error) {

	// Push fields under mod
	mod := n.Resource.(*types.Module)
	for _, n := range nested {
		if n.ResourceType != types.ModuleFieldResourceType {
			continue
		}

		mod.Fields = append(mod.Fields, n.Resource.(*types.ModuleField))
	}

	// Register to DAL
	dl, err := e.grabDal(p)
	if err != nil {
		return
	}

	nsNode := tree.ParentForRef(n, n.References["NamespaceID"])
	ns := nsNode.Resource.(*types.Namespace)

	// @todo get connection and things from there
	models, err := service.ModulesToModelSet(dl, ns, mod)
	if err != nil {
		return err
	}

	// @note there is only one model so this is ok
	return dl.ReplaceModel(ctx, models[0])
}

func (e StoreEncoder) encodeModuleExtendSubResources(ctx context.Context, p envoyx.EncodeParams, s store.Storer, n *envoyx.Node, tree envoyx.Traverser) (err error) {
	cc := tree.ChildrenForResourceType(n, ComposeRecordDatasourceAuxType)

	dl, err := e.grabDal(p)
	if err != nil {
		return
	}

	return e.encodeRecordDatasources(ctx, p, s, dl, cc, tree)
}

func (e *StoreEncoder) grabDal(p envoyx.EncodeParams) (dl dal.FullService, err error) {
	auxs, ok := p.Params["dal"]
	if !ok {
		err = fmt.Errorf("dal not defined")
		return
	}

	dl, ok = auxs.(dal.FullService)
	if !ok {
		err = fmt.Errorf("invalid dal provided")
		return
	}

	return
}
