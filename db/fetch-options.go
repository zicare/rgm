package db

import (
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
	Equal            Params
	NotEqual         Params
	GreaterThan      Params
	LessThan         Params
	GreaterEqualThan Params
	LessEqualThan    Params
	Order            []string
	Offset           int
	Limit            int
	Checksum         int
	Dig              int // Flag used to include (or not) parent data in find, fetch, etc.
}

// FetchOptionsFactory exported
func FetchOptionsFactory(tbl Table, uid string, qparams QParams, parents ...Table) *FetchOptions {

	var (
		fo   = new(FetchOptions)
		cols = Cols(tbl)
	)

	fo.setTable(tbl)
	fo.setUID(uid)
	fo.setParents(parents...)
	fo.setColumn(cols)
	fo.setIsNull(qparams, cols)
	fo.setIsNotNull(qparams, cols)
	fo.setIn(qparams, cols)
	fo.setNotIn(qparams, cols)
	fo.setEqual(qparams, cols)
	fo.setNotEqual(qparams, cols)
	fo.setGreaterThan(qparams, cols)
	fo.setLessThan(qparams, cols)
	fo.setGreaterEqualThan(qparams, cols)
	fo.setLessEqualThan(qparams, cols)
	fo.setOrder(qparams, cols)
	fo.setOffset(qparams)
	fo.setLimit(qparams)
	fo.setChecksum(qparams)
	fo.setDig(qparams)

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

func (fo *FetchOptions) setColumn(cols []string) {

	fo.Column = cols
}

func (fo *FetchOptions) setIsNull(qparams QParams, cols []string) {

	fo.IsNull = []string{}

	if isnull, ok := qparams["isnull"]; ok {
		for _, k := range cols {
			if lib.Contains(isnull, k) {
				fo.IsNull = append(fo.IsNull, k)
			}
		}
	}
}

func (fo *FetchOptions) setIsNotNull(qparams QParams, cols []string) {

	fo.IsNotNull = []string{}

	if notnull, ok := qparams["notnull"]; ok {
		for _, k := range cols {
			if lib.Contains(notnull, k) {
				fo.IsNull = append(fo.IsNotNull, k)
			}
		}
	}
}

func (fo *FetchOptions) setIn(qparams QParams, cols []string) {

	fo.In = make(map[string][]interface{})

	if in, ok := qparams["in"]; ok {
		for _, k := range in {
			j := strings.Split(k, "|")
			if lib.Contains(cols, j[0]) {
				for _, v := range strings.Split(j[1], ",") {
					fo.In[j[0]] = append(fo.In[j[0]], v)
				}
			}
		}
	}
}

func (fo *FetchOptions) setNotIn(qparams QParams, cols []string) {

	fo.NotIn = make(map[string][]interface{})

	if notin, ok := qparams["notin"]; ok {
		for _, k := range notin {
			j := strings.Split(k, "|")
			if lib.Contains(cols, j[0]) {
				for _, v := range strings.Split(j[1], ",") {
					fo.NotIn[j[0]] = append(fo.NotIn[j[0]], v)
				}
			}
		}
	}
}

func (fo *FetchOptions) setEqual(qparams QParams, cols []string) {

	fo.Equal = make(Params)

	if eq, ok := qparams["eq"]; ok {
		for _, k := range eq {
			j := strings.Split(k, "|")
			if lib.Contains(cols, j[0]) {
				fo.Equal[j[0]] = j[1]
			}
		}
	}
}

func (fo *FetchOptions) setNotEqual(qparams QParams, cols []string) {

	fo.NotEqual = make(Params)

	if noteq, ok := qparams["noteq"]; ok {
		for _, k := range noteq {
			j := strings.Split(k, "|")
			if lib.Contains(cols, j[0]) {
				fo.NotEqual[j[0]] = j[1]
			}
		}
	}
}

func (fo *FetchOptions) setGreaterThan(qparams QParams, cols []string) {

	fo.GreaterThan = make(Params)

	if gt, ok := qparams["gt"]; ok {
		for _, k := range gt {
			j := strings.Split(k, "|")
			if lib.Contains(cols, j[0]) {
				fo.GreaterThan[j[0]] = j[1]
			}
		}
	}
}

func (fo *FetchOptions) setLessThan(qparams QParams, cols []string) {

	fo.LessThan = make(Params)

	if lt, ok := qparams["lt"]; ok {
		for _, k := range lt {
			j := strings.Split(k, "|")
			if lib.Contains(cols, j[0]) {
				fo.LessThan[j[0]] = j[1]
			}
		}
	}
}

func (fo *FetchOptions) setGreaterEqualThan(qparams QParams, cols []string) {

	fo.GreaterEqualThan = make(Params)

	if gteq, ok := qparams["gteq"]; ok {
		for _, k := range gteq {
			j := strings.Split(k, "|")
			if lib.Contains(cols, j[0]) {
				fo.GreaterEqualThan[j[0]] = j[1]
			}
		}
	}
}

func (fo *FetchOptions) setLessEqualThan(qparams QParams, cols []string) {

	fo.LessEqualThan = make(Params)

	if lteq, ok := qparams["lteq"]; ok {
		for _, k := range lteq {
			j := strings.Split(k, "|")
			if lib.Contains(cols, j[0]) {
				fo.LessEqualThan[j[0]] = j[1]
			}
		}
	}
}

func (fo *FetchOptions) setOrder(qparams QParams, cols []string) {

	fo.Order = []string{}

	if order, ok := qparams["order"]; ok {
		for _, k := range order {
			j := strings.Split(k, "|")
			if !lib.Contains(cols, j[0]) {
				continue
			} else if len(j) == 1 {
				fo.Order = append(fo.Order, j[0]+" ASC")
			} else if strings.ToUpper(j[1]) == "ASC" || strings.ToUpper(j[1]) == "DESC" {
				fo.Order = append(fo.Order, j[0]+" "+strings.ToUpper(j[1]))
			}
		}
	}
}

func (fo *FetchOptions) setOffset(qparams QParams) {

	if offset, ok := qparams["offset"]; ok {
		fo.Offset, _ = strconv.Atoi(offset[0])
	} else {
		fo.Offset = 0
	}
}

func (fo *FetchOptions) setLimit(qparams QParams) {

	fo.Limit = config.Config().GetInt("param.icpp")

	if qparams == nil {
		fo.Limit = 0
	} else if limit, ok := qparams["limit"]; ok {
		fo.Limit, _ = strconv.Atoi(limit[0])
	}
}

func (fo *FetchOptions) setChecksum(qparams QParams) {

	if checksum, ok := qparams["checksum"]; ok && (checksum[0] == "1") {
		fo.Checksum = 1
	}
}

func (fo *FetchOptions) setDig(qparams QParams) {

	if dig, ok := qparams["dig"]; ok && (dig[0] == "1") {
		fo.Dig = 1
	}
}
