package db

import (
	"github.com/zicare/rgm/msg"
)

// FindOptions exported
type FindOptions struct {
	Table    Table
	UID      string
	Parents  []Table // Also used to set scope in find, fetch, etc.
	Where    Params
	Checksum int
	Dig      int
}

// FindByOptionsFactory exported
func FindOptionsFactory(tbl Table, uid string, qparams QParams, params Params, pk bool, parents ...Table) (*FindOptions, *ParamError) {

	fo := new(FindOptions)

	fo.setTable(tbl)
	fo.setUID(uid)
	fo.setParents(parents...)
	fo.setChecksum(qparams)
	fo.setDig(qparams)

	err := fo.setWhere(params, pk)

	return fo, err
}

func (fo *FindOptions) setTable(tbl Table) {

	fo.Table = tbl
}

func (fo *FindOptions) setUID(uid string) {

	fo.UID = uid
}

func (fo *FindOptions) setParents(parents ...Table) {

	fo.Parents = parents
}

func (fo *FindOptions) setWhere(params Params, pk bool) *ParamError {

	var cols []string

	if pk {
		cols = Pk(fo.Table)
	} else {
		cols = Cols(fo.Table)
	}

	fo.Where = make(Params)

	for _, k := range cols {
		if v, ok := params[k]; ok {
			fo.Where[k] = v
		} else if pk {
			e := ParamError{msg.Get("26")}
			return &e
		}
	}

	return nil
}

func (fo *FindOptions) setChecksum(qparams QParams) {

	if checksum, ok := qparams["checksum"]; ok && (checksum[0] == "1") {
		fo.Checksum = 1
	}
}

func (fo *FindOptions) setDig(qparams QParams) {

	if dig, ok := qparams["dig"]; ok && (dig[0] == "1") {
		fo.Dig = 1
	}
}
