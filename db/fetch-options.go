package db

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/lib"
)

// FetchOptions exported
type FetchOptions struct {
	Table            Table
	UID              string  // ID from jwt, used to set scope in find, fetch, etc.
	Parents          []Table // Also used to set scope in find, fetch, etc.
	Column           []string
	IsNull           []string
	IsNotNull        []string
	In               map[string][]interface{}
	NotIn            map[string][]interface{}
	Equal            map[string]string
	NotEqual         map[string]string
	GreaterThan      map[string]string
	LessThan         map[string]string
	GreaterEqualThan map[string]string
	LessEqualThan    map[string]string
	Order            []string
	Offset           int
	Limit            int
	Checksum         int
	Dig              int // Flag used to include (or not) parent data in find, fetch, etc.
}

// FetchOptionsFactory exported
func FetchOptionsFactory(tbl Table, uid string, qparam url.Values, parents ...Table) *FetchOptions {

	var (
		fo   = new(FetchOptions)
		meta = GetTableMeta(tbl)
	)

	fo.setTable(tbl)
	fo.setUID(uid)
	fo.setParents(parents...)
	fo.setColumn(meta)
	fo.setIsNull(qparam, meta)
	fo.setIsNotNull(qparam, meta)
	fo.setIn(qparam, meta)
	fo.setNotIn(qparam, meta)
	fo.setEqual(qparam, meta)
	fo.setNotEqual(qparam, meta)
	fo.setGreaterThan(qparam, meta)
	fo.setLessThan(qparam, meta)
	fo.setGreaterEqualThan(qparam, meta)
	fo.setLessEqualThan(qparam, meta)
	fo.setOrder(qparam, meta)
	fo.setOffset(qparam)
	fo.setLimit(qparam)
	fo.setChecksum(qparam)
	fo.setDig(qparam)

	return fo
}

func (fo *FetchOptions) setTable(tbl Table) {

	fo.Table = tbl
}

func (fo *FetchOptions) setUID(uid string) {

	fo.UID = uid
}

func (fo *FetchOptions) setParents(parents ...Table) {

	fo.Parents = parents
}

func (fo *FetchOptions) setColumn(meta TableMeta) {

	fo.Column = meta.Fields
}

func (fo *FetchOptions) setIsNull(qparam url.Values, meta TableMeta) {

	fo.IsNull = []string{}

	if isnull, ok := qparam["isnull"]; ok {
		for _, k := range meta.Fields {
			if lib.Contains(isnull, k) {
				fo.IsNull = append(fo.IsNull, k)
			}
		}
	}
}

func (fo *FetchOptions) setIsNotNull(qparam url.Values, meta TableMeta) {

	fo.IsNotNull = []string{}

	if notnull, ok := qparam["notnull"]; ok {
		for _, k := range meta.Fields {
			if lib.Contains(notnull, k) {
				fo.IsNull = append(fo.IsNotNull, k)
			}
		}
	}
}

func (fo *FetchOptions) setIn(qparam url.Values, meta TableMeta) {

	fo.In = make(map[string][]interface{})

	if in, ok := qparam["in"]; ok {
		for _, k := range in {
			j := strings.Split(k, "|")
			if lib.Contains(meta.Fields, j[0]) {
				for _, v := range strings.Split(j[1], ",") {
					fo.In[j[0]] = append(fo.In[j[0]], v)
				}
			}
		}
	}
}

func (fo *FetchOptions) setNotIn(qparam url.Values, meta TableMeta) {

	fo.NotIn = make(map[string][]interface{})

	if notin, ok := qparam["notin"]; ok {
		for _, k := range notin {
			j := strings.Split(k, "|")
			if lib.Contains(meta.Fields, j[0]) {
				for _, v := range strings.Split(j[1], ",") {
					fo.NotIn[j[0]] = append(fo.NotIn[j[0]], v)
				}
			}
		}
	}
}

func (fo *FetchOptions) setEqual(qparam url.Values, meta TableMeta) {

	fo.Equal = make(map[string]string)

	if eq, ok := qparam["eq"]; ok {
		for _, k := range eq {
			j := strings.Split(k, "|")
			if lib.Contains(meta.Fields, j[0]) {
				fo.Equal[j[0]] = j[1]
			}
		}
	}
}

func (fo *FetchOptions) setNotEqual(qparam url.Values, meta TableMeta) {

	fo.NotEqual = make(map[string]string)

	if noteq, ok := qparam["noteq"]; ok {
		for _, k := range noteq {
			j := strings.Split(k, "|")
			if lib.Contains(meta.Fields, j[0]) {
				fo.NotEqual[j[0]] = j[1]
			}
		}
	}
}

func (fo *FetchOptions) setGreaterThan(qparam url.Values, meta TableMeta) {

	fo.GreaterThan = make(map[string]string)

	if gt, ok := qparam["gt"]; ok {
		for _, k := range gt {
			j := strings.Split(k, "|")
			if lib.Contains(meta.Fields, j[0]) {
				fo.GreaterThan[j[0]] = j[1]
			}
		}
	}
}

func (fo *FetchOptions) setLessThan(qparam url.Values, meta TableMeta) {

	fo.LessThan = make(map[string]string)

	if lt, ok := qparam["lt"]; ok {
		for _, k := range lt {
			j := strings.Split(k, "|")
			if lib.Contains(meta.Fields, j[0]) {
				fo.LessThan[j[0]] = j[1]
			}
		}
	}
}

func (fo *FetchOptions) setGreaterEqualThan(qparam url.Values, meta TableMeta) {

	fo.GreaterEqualThan = make(map[string]string)

	if gteq, ok := qparam["gteq"]; ok {
		for _, k := range gteq {
			j := strings.Split(k, "|")
			if lib.Contains(meta.Fields, j[0]) {
				fo.GreaterEqualThan[j[0]] = j[1]
			}
		}
	}
}

func (fo *FetchOptions) setLessEqualThan(qparam url.Values, meta TableMeta) {

	fo.LessEqualThan = make(map[string]string)

	if lteq, ok := qparam["lteq"]; ok {
		for _, k := range lteq {
			j := strings.Split(k, "|")
			if lib.Contains(meta.Fields, j[0]) {
				fo.LessEqualThan[j[0]] = j[1]
			}
		}
	}
}

func (fo *FetchOptions) setOrder(qparam url.Values, meta TableMeta) {

	fo.Order = []string{}

	if order, ok := qparam["order"]; ok {
		for _, k := range order {
			j := strings.Split(k, "|")
			if !lib.Contains(meta.Fields, j[0]) {
				continue
			} else if len(j) == 1 {
				fo.Order = append(fo.Order, j[0]+" ASC")
			} else if strings.ToUpper(j[1]) == "ASC" || strings.ToUpper(j[1]) == "DESC" {
				fo.Order = append(fo.Order, j[0]+" "+strings.ToUpper(j[1]))
			}
		}
	}
}

func (fo *FetchOptions) setOffset(qparam url.Values) {

	if offset, ok := qparam["offset"]; ok {
		fo.Offset, _ = strconv.Atoi(offset[0])
	} else {
		fo.Offset = 0
	}
}

func (fo *FetchOptions) setLimit(qparam url.Values) {

	fo.Limit = config.Config().GetInt("param.icpp")

	if limit, ok := qparam["limit"]; ok {
		fo.Limit, _ = strconv.Atoi(limit[0])
	} else {
		fo.Limit = config.Config().GetInt("param.icpp")
	}
}

func (fo *FetchOptions) setChecksum(param url.Values) {

	if checksum, ok := param["checksum"]; ok && (checksum[0] == "1") {
		fo.Checksum = 1
	}
}

func (fo *FetchOptions) setDig(param url.Values) {

	if dig, ok := param["dig"]; ok && (dig[0] == "1") {
		fo.Dig = 1
	}
}
