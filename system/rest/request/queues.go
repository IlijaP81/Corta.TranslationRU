package request

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/messagebus"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/go-chi/chi"
	sqlxTypes "github.com/jmoiron/sqlx/types"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

// dummy vars to prevent
// unused imports complain
var (
	_ = chi.URLParam
	_ = multipart.ErrMessageTooLarge
	_ = payload.ParseUint64s
	_ = strings.ToLower
	_ = io.EOF
	_ = fmt.Errorf
	_ = json.NewEncoder
)

type (
	// Internal API interface
	QueuesList struct {
		// Query GET parameter
		//
		// Search query
		Query string

		// Limit GET parameter
		//
		// Limit
		Limit uint

		// PageCursor GET parameter
		//
		// Page cursor
		PageCursor string

		// Sort GET parameter
		//
		// Sort items
		Sort string

		// Deleted GET parameter
		//
		// Exclude (0, default), include (1) or return only (2) deleted queues
		Deleted uint
	}

	QueuesCreate struct {
		// Queue POST parameter
		//
		// Name of queue
		Queue string

		// Handler POST parameter
		//
		// Queue handler
		Handler string

		// Meta POST parameter
		//
		// Meta data for queue
		Meta sqlxTypes.JSONText
	}

	QueuesRead struct {
		// QueueID PATH parameter
		//
		// Queue ID
		QueueID uint64 `json:",string"`
	}

	QueuesUpdate struct {
		// QueueID PATH parameter
		//
		// Queue ID
		QueueID uint64 `json:",string"`

		// Queue POST parameter
		//
		// Name of queue
		Queue string

		// Handler POST parameter
		//
		// Queue handler
		Handler string

		// Meta POST parameter
		//
		// Meta data for queue
		Meta messagebus.QueueSettingsMeta
	}

	QueuesDelete struct {
		// QueueID PATH parameter
		//
		// Queue ID
		QueueID uint64 `json:",string"`
	}
)

// NewQueuesList request
func NewQueuesList() *QueuesList {
	return &QueuesList{}
}

// Auditable returns all auditable/loggable parameters
func (r QueuesList) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"query":      r.Query,
		"limit":      r.Limit,
		"pageCursor": r.PageCursor,
		"sort":       r.Sort,
		"deleted":    r.Deleted,
	}
}

// Auditable returns all auditable/loggable parameters
func (r QueuesList) GetQuery() string {
	return r.Query
}

// Auditable returns all auditable/loggable parameters
func (r QueuesList) GetLimit() uint {
	return r.Limit
}

// Auditable returns all auditable/loggable parameters
func (r QueuesList) GetPageCursor() string {
	return r.PageCursor
}

// Auditable returns all auditable/loggable parameters
func (r QueuesList) GetSort() string {
	return r.Sort
}

// Auditable returns all auditable/loggable parameters
func (r QueuesList) GetDeleted() uint {
	return r.Deleted
}

// Fill processes request and fills internal variables
func (r *QueuesList) Fill(req *http.Request) (err error) {

	{
		// GET params
		tmp := req.URL.Query()

		if val, ok := tmp["query"]; ok && len(val) > 0 {
			r.Query, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["limit"]; ok && len(val) > 0 {
			r.Limit, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["pageCursor"]; ok && len(val) > 0 {
			r.PageCursor, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["sort"]; ok && len(val) > 0 {
			r.Sort, err = val[0], nil
			if err != nil {
				return err
			}
		}
		if val, ok := tmp["deleted"]; ok && len(val) > 0 {
			r.Deleted, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewQueuesCreate request
func NewQueuesCreate() *QueuesCreate {
	return &QueuesCreate{}
}

// Auditable returns all auditable/loggable parameters
func (r QueuesCreate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"queue":   r.Queue,
		"handler": r.Handler,
		"meta":    r.Meta,
	}
}

// Auditable returns all auditable/loggable parameters
func (r QueuesCreate) GetQueue() string {
	return r.Queue
}

// Auditable returns all auditable/loggable parameters
func (r QueuesCreate) GetHandler() string {
	return r.Handler
}

// Auditable returns all auditable/loggable parameters
func (r QueuesCreate) GetMeta() sqlxTypes.JSONText {
	return r.Meta
}

// Fill processes request and fills internal variables
func (r *QueuesCreate) Fill(req *http.Request) (err error) {

	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["queue"]; ok && len(val) > 0 {
			r.Queue, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["handler"]; ok && len(val) > 0 {
			r.Handler, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta"]; ok && len(val) > 0 {
			r.Meta, err = payload.ParseJSONTextWithErr(val[0])
			if err != nil {
				return err
			}
		}
	}

	return err
}

// NewQueuesRead request
func NewQueuesRead() *QueuesRead {
	return &QueuesRead{}
}

// Auditable returns all auditable/loggable parameters
func (r QueuesRead) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"queueID": r.QueueID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r QueuesRead) GetQueueID() uint64 {
	return r.QueueID
}

// Fill processes request and fills internal variables
func (r *QueuesRead) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "queueID")
		r.QueueID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewQueuesUpdate request
func NewQueuesUpdate() *QueuesUpdate {
	return &QueuesUpdate{}
}

// Auditable returns all auditable/loggable parameters
func (r QueuesUpdate) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"queueID": r.QueueID,
		"queue":   r.Queue,
		"handler": r.Handler,
		"meta":    r.Meta,
	}
}

// Auditable returns all auditable/loggable parameters
func (r QueuesUpdate) GetQueueID() uint64 {
	return r.QueueID
}

// Auditable returns all auditable/loggable parameters
func (r QueuesUpdate) GetQueue() string {
	return r.Queue
}

// Auditable returns all auditable/loggable parameters
func (r QueuesUpdate) GetHandler() string {
	return r.Handler
}

// Auditable returns all auditable/loggable parameters
func (r QueuesUpdate) GetMeta() messagebus.QueueSettingsMeta {
	return r.Meta
}

// Fill processes request and fills internal variables
func (r *QueuesUpdate) Fill(req *http.Request) (err error) {

	if strings.ToLower(req.Header.Get("content-type")) == "application/json" {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["queue"]; ok && len(val) > 0 {
			r.Queue, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["handler"]; ok && len(val) > 0 {
			r.Handler, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["meta[]"]; ok {
			r.Meta, err = messagebus.ParseQueueSettingsMeta(val)
			if err != nil {
				return err
			}
		} else if val, ok := req.Form["meta"]; ok {
			r.Meta, err = messagebus.ParseQueueSettingsMeta(val)
			if err != nil {
				return err
			}
		}
	}

	{
		var val string
		// path params

		val = chi.URLParam(req, "queueID")
		r.QueueID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}

// NewQueuesDelete request
func NewQueuesDelete() *QueuesDelete {
	return &QueuesDelete{}
}

// Auditable returns all auditable/loggable parameters
func (r QueuesDelete) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"queueID": r.QueueID,
	}
}

// Auditable returns all auditable/loggable parameters
func (r QueuesDelete) GetQueueID() uint64 {
	return r.QueueID
}

// Fill processes request and fills internal variables
func (r *QueuesDelete) Fill(req *http.Request) (err error) {

	{
		var val string
		// path params

		val = chi.URLParam(req, "queueID")
		r.QueueID, err = payload.ParseUint64(val), nil
		if err != nil {
			return err
		}

	}

	return err
}
