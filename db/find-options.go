package db

import (
	"net/url"

	"github.com/zicare/rgm/msg"
)

// FindOptions exported
type FindOptions struct {
	Table    Table
	Where    map[string]string
	Checksum int
	Dig      int
	UID      string
}

// FindOptionsFactory exported
func FindOptionsFactory(tbl Table, uid string, param url.Values, idv []string) (*FindOptions, *ParamError) {

	var (
		fo  = new(FindOptions)
		idk = Pk(tbl)
	)

	if err := fo.setWhere(idk, idv); err != nil {
		return fo, err
	}

	fo.setTable(tbl)
	fo.setChecksum(param)
	fo.setDig(param)
	fo.setUID(uid)

	return fo, nil
}

func (fo *FindOptions) setTable(tbl Table) {

	fo.Table = tbl
}

func (fo *FindOptions) setWhere(idk []string, idv []string) *ParamError {

	fo.Where = make(map[string]string)

	if len(idk) != len(idv) {
		e := ParamError{msg.Get("26")}
		return &e
	}

	for i, j := range idk {
		fo.Where[j] = idv[i]
	}

	return nil
}

func (fo *FindOptions) setChecksum(param url.Values) {

	if checksum, ok := param["checksum"]; ok && (checksum[0] == "1") {
		fo.Checksum = 1
	}
}

func (fo *FindOptions) setDig(param url.Values) {

	if dig, ok := param["dig"]; ok && (dig[0] == "1") {
		fo.Dig = 1
	}
}

func (fo *FindOptions) setUID(uid string) {

	fo.UID = uid
}
