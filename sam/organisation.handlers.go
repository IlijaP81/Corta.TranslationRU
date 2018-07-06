package sam

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `organisation.go`, `organisation.util.go` or `organisation_test.go` to
	implement your API calls, helper functions and tests. The file `organisation.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"net/http"

	"github.com/titpetric/factory/resputil"
)

func (oh *OrganisationHandlers) Edit(w http.ResponseWriter, r *http.Request) {
	params := organisationEditRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return oh.Organisation.Edit(params) })
}
func (oh *OrganisationHandlers) Remove(w http.ResponseWriter, r *http.Request) {
	params := organisationRemoveRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return oh.Organisation.Remove(params) })
}
func (oh *OrganisationHandlers) Read(w http.ResponseWriter, r *http.Request) {
	params := organisationReadRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return oh.Organisation.Read(params) })
}
func (oh *OrganisationHandlers) Search(w http.ResponseWriter, r *http.Request) {
	params := organisationSearchRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return oh.Organisation.Search(params) })
}
func (oh *OrganisationHandlers) Archive(w http.ResponseWriter, r *http.Request) {
	params := organisationArchiveRequest{}.new()
	resputil.JSON(w, params.Fill(r), func() (interface{}, error) { return oh.Organisation.Archive(params) })
}
