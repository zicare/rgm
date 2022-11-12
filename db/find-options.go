package db

import (
	"net/url"

	"github.com/zicare/rgm/msg"
)

// FindOptions exported
type FindOptions struct {
	Table    Table
	UID      string
	Parents  []Table // Also used to set scope in find, fetch, etc.
	Where    map[string]string
	Checksum int
	Dig      int
}

// FindOptionsFactory exported
func FindOptionsFactory(tbl Table, uid string, qparam url.Values, uparam map[string]string, parents ...Table) (*FindOptions, *ParamError) {

	var (
		fo = new(FindOptions)
		pk = Pk(tbl)
	)

	if err := fo.setWhere(pk, uparam); err != nil {
		return fo, err
	}

	fo.setTable(tbl)
	fo.setUID(uid)
	fo.setParents(parents...)
	fo.setChecksum(qparam)
	fo.setDig(qparam)

	return fo, nil
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

func (fo *FindOptions) setWhere(pk []string, uparam map[string]string) *ParamError {

	fo.Where = make(map[string]string)

	for _, k := range pk {
		if _, ok := uparam[k]; ok {
			fo.Where[k] = uparam[k]
		} else {
			e := ParamError{msg.Get("26")}
			return &e
		}
	}

	return nil
}

func (fo *FindOptions) setChecksum(qparam url.Values) {

	if checksum, ok := qparam["checksum"]; ok && (checksum[0] == "1") {
		fo.Checksum = 1
	}
}

func (fo *FindOptions) setDig(qparam url.Values) {

	if dig, ok := qparam["dig"]; ok && (dig[0] == "1") {
		fo.Dig = 1
	}
}
