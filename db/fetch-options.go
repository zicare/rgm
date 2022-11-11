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
	Dig              int
	UID              string
}

// FetchOptionsFactory exported
func FetchOptionsFactory(tbl Table, uid string, param url.Values) *FetchOptions {

	var (
		fo   = new(FetchOptions)
		meta = GetTableMeta(tbl)
	)

	fo.setTable(tbl)
	fo.setColumn(param, meta)
	fo.setIsNull(param, meta)
	fo.setIsNotNull(param, meta)
	fo.setIn(param, meta)
	fo.setNotIn(param, meta)
	fo.setEqual(param, meta)
	fo.setNotEqual(param, meta)
	fo.setGreaterThan(param, meta)
	fo.setLessThan(param, meta)
	fo.setGreaterEqualThan(param, meta)
	fo.setLessEqualThan(param, meta)
	fo.setOrder(param, meta)
	fo.setOffsetAndLimit(param)
	fo.setChecksum(param)
	fo.setDig(param)
	fo.setUID(uid)

	return fo
}

func (fo *FetchOptions) setTable(tbl Table) {

	fo.Table = tbl
}

func (fo *FetchOptions) setColumn(param url.Values, meta TableMeta) {

	fo.Column = []string{}

	if i, ok := param["cols"]; ok {
		j := strings.Split(i[0], ",")
		for _, v := range j {
			if lib.Contains(meta.Fields, v) {
				fo.Column = append(fo.Column, v)
			}
		}
	} else {
		fo.Column = meta.Fields
	}

	//xcols
	if i, ok := param["xcols"]; ok {
		fo.Column = lib.Diff(fo.Column, strings.Split(i[0], ","))
	}
}

func (fo *FetchOptions) setIsNull(param url.Values, meta TableMeta) {

	fo.IsNull = []string{}

	if i, ok := param["isnull"]; ok {
		colsAux := make(map[string]string)
		for _, v := range strings.Split(i[0], ",") {
			colsAux[v] = v
		}
		for _, k := range meta.Fields {
			if _, ok := colsAux[k]; ok {
				fo.IsNull = append(fo.IsNull, k)
			}
		}
	}
}

func (fo *FetchOptions) setIsNotNull(param url.Values, meta TableMeta) {

	fo.IsNotNull = []string{}

	if i, ok := param["notnull"]; ok {
		colsAux := make(map[string]string)
		for _, v := range strings.Split(i[0], ",") {
			colsAux[v] = v
		}
		for _, k := range meta.Fields {
			if _, ok := colsAux[k]; ok {
				fo.IsNotNull = append(fo.IsNotNull, k)
			}
		}
	}
}

func (fo *FetchOptions) setIn(param url.Values, meta TableMeta) {

	fo.In = make(map[string][]interface{})

	if i, ok := param["in"]; ok {
		for _, k := range i {
			j := strings.Split(k, "|")
			if lib.Contains(meta.Fields, j[0]) {
				for _, v := range strings.Split(j[1], ",") {
					fo.In[j[0]] = append(fo.In[j[0]], v)
				}
			}
		}
	}
}

func (fo *FetchOptions) setNotIn(param url.Values, meta TableMeta) {

	fo.NotIn = make(map[string][]interface{})

	if i, ok := param["notin"]; ok {
		for _, k := range i {
			j := strings.Split(k, "|")
			if lib.Contains(meta.Fields, j[0]) {
				for _, v := range strings.Split(j[1], ",") {
					fo.NotIn[j[0]] = append(fo.NotIn[j[0]], v)
				}
			}
		}
	}
}

func (fo *FetchOptions) setEqual(param url.Values, meta TableMeta) {

	fo.Equal = make(map[string]string)

	if i, ok := param["eq"]; ok {
		for _, k := range i {
			j := strings.Split(k, "|")
			if lib.Contains(meta.Fields, j[0]) {
				fo.Equal[j[0]] = j[1]
			}
		}
	}
}

func (fo *FetchOptions) setNotEqual(param url.Values, meta TableMeta) {

	fo.NotEqual = make(map[string]string)

	if i, ok := param["noteq"]; ok {
		for _, k := range i {
			j := strings.Split(k, "|")
			if lib.Contains(meta.Fields, j[0]) {
				fo.NotEqual[j[0]] = j[1]
			}
		}
	}
}

func (fo *FetchOptions) setGreaterThan(param url.Values, meta TableMeta) {

	fo.GreaterThan = make(map[string]string)

	if i, ok := param["gt"]; ok {
		for _, k := range i {
			j := strings.Split(k, "|")
			if lib.Contains(meta.Fields, j[0]) {
				fo.GreaterThan[j[0]] = j[1]
			}
		}
	}
}

func (fo *FetchOptions) setLessThan(param url.Values, meta TableMeta) {

	fo.LessThan = make(map[string]string)

	if i, ok := param["lt"]; ok {
		for _, k := range i {
			j := strings.Split(k, "|")
			if lib.Contains(meta.Fields, j[0]) {
				fo.LessThan[j[0]] = j[1]
			}
		}
	}
}

func (fo *FetchOptions) setGreaterEqualThan(param url.Values, meta TableMeta) {

	fo.GreaterEqualThan = make(map[string]string)

	if i, ok := param["gteq"]; ok {
		for _, k := range i {
			j := strings.Split(k, "|")
			if lib.Contains(meta.Fields, j[0]) {
				fo.GreaterEqualThan[j[0]] = j[1]
			}
		}
	}
}

func (fo *FetchOptions) setLessEqualThan(param url.Values, meta TableMeta) {

	fo.LessEqualThan = make(map[string]string)

	if i, ok := param["lteq"]; ok {
		for _, k := range i {
			j := strings.Split(k, "|")
			if lib.Contains(meta.Fields, j[0]) {
				fo.LessEqualThan[j[0]] = j[1]
			}
		}
	}
}

func (fo *FetchOptions) setOrder(param url.Values, meta TableMeta) {

	fo.Order = []string{}

	if i, ok := param["order"]; ok {
		j := strings.Split(i[0], ";")
		for _, k := range j {
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

func (fo *FetchOptions) setOffsetAndLimit(param url.Values) {

	//offset and limit
	fo.Offset = 0
	fo.Limit = config.Config().GetInt("param.icpp")

	if i, ok := param["limit"]; ok {
		j := strings.Split(i[0], ",")
		switch len(j) {
		case 1:
			fo.Offset = 0
			fo.Limit, _ = strconv.Atoi(j[0])
		case 2:
			fo.Offset, _ = strconv.Atoi(j[0])
			fo.Limit, _ = strconv.Atoi(j[1])
		}
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

func (fo *FetchOptions) setUID(uid string) {

	fo.UID = uid
}
