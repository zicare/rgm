package mysql

import (
	"github.com/zicare/rgm/ds"
)

// ITable defines an interface for db table access.
// Consider annonymous embedding of Table in your concrete ITable.
// Table offers default implementation for all ITable and ds.IDataSource
// methods, except Name().
// You can always overwrite the methods you need to.
// Check Table for more information.
type ITable interface {

	// ITable interfaces must fulfills ds.IDataSource
	ds.IDataSource

	// Dig offers a chance to optionally attach parent table data.
	// The included Table's Find and Fetch implementations of ITable calls Dig on each
	// found/fetched result, thus supporting nested digs. That is, it can go up in the relationship tree
	// digging parent data.
	//
	// f passes all parent relationships to dig. In a HTTP GET request, f are decoded from the
	// dig query params. ie.: https://my.site.com/clients?dig=clients.officer&dig=officers.location
	//
	// Implementation example:
	//
	// func (t *Client) Dig(f ...string) (dig []mysql.Dig) {
	//
	// 	 var p ds.Params
	//
	// 	 if lib.Contains(f, "clients.officer") && t.OfficerID != nil {
	// 		t.Officer = new(Officer)
	// 		p = make(ds.Params)
	// 		p["user_id"] = fmt.Sprint(*t.OfficerID)
	// 		dig = append(dig, mysql.NewDig(t.Officer, p))
	// 	 }
	//   return dig
	// }
	//
	// func (t *Officer) Dig(f ...string) (dig []mysql.Dig) {
	//
	// 	 var p ds.Params
	//
	// 	 if lib.Contains(f, "officers.location") {
	// 		t.Location = new(Location)
	// 		p = make(ds.Params)
	// 		p["location_id"] = fmt.Sprint(t.LocationID)
	// 		dig = append(dig, mysql.NewDig(t.Location, p))
	// 	 }
	//
	// 	 return dig
	// }
	Dig(f ...string) []Dig

	// BeforeSelect offers a chance optionally set additional constraints
	// in a per Table basis, or abort the select by returning a *ds.NotAllowedError.
	BeforeSelect(qo *ds.QueryOptions) (ds.Params, *ds.NotAllowedError)

	// BeforeInsert offers a chance to complete extra validations, alter values,
	// or abort the insert by returning an error.
	// Consider using *ds.NotAllowedError and/or validator.validationErrors, these
	// will be treated as such by ctrl.CrudController, others will be considered
	// InternalServerError's.
	BeforeInsert(qo *ds.QueryOptions) error

	// BeforeUpdate offers a chance to complete extra validations, alter values,
	// or abort the update by returning an error.
	// Consider using *ds.NotAllowedError and/or validator.validationErrors, these
	// will be treated as such by ctrl.CrudController, others will be considered
	// InternalServerError's.
	BeforeUpdate(qo *ds.QueryOptions) error

	// BeforeDelete offers a chance optionally set additional constraints
	// in a per Table basis, or even abort the select by returning a *ds.NotAllowedError.
	BeforeDelete(qo *ds.QueryOptions) (ds.Params, *ds.NotAllowedError)
}

// Table offers default implementation for all ITable and ds.IDataSource
// methods, except Name().
// Consider annonymous embedding of Table in your concrete ITable.
type Table struct{} // Dig exported

func (Table) Dig(f ...string) []Dig {

	return []Dig{}
}

func (Table) BeforeSelect(qo *ds.QueryOptions) (ds.Params, *ds.NotAllowedError) {
	return nil, nil
}

func (Table) BeforeInsert(qo *ds.QueryOptions) error {
	return nil
}

func (Table) BeforeUpdate(qo *ds.QueryOptions) error {
	return nil
}

func (Table) BeforeDelete(qo *ds.QueryOptions) (ds.Params, *ds.NotAllowedError) {
	return nil, nil
}
