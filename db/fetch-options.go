package db

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/lib"
)

//FetchOptions exported
type FetchOptions struct {
	Column           []string
	IsNull           []string
	IsNotNull        []string
	Equal            map[string]string
	GreaterThan      map[string]string
	LessThan         map[string]string
	GreaterEqualThan map[string]string
	LessEqualThan    map[string]string
	Order            []string
	Offset           int
	Limit            int
	Checksum         int
}

func GetFetchOptions(c *gin.Context, tbl Table) *FetchOptions {

	var (
		meta  = GetTableMeta(tbl)
		param = c.Request.URL.Query()
	)

	fo := new(FetchOptions)

	fo.setColumn(param, meta)
	fo.setIsNull(param, meta)
	fo.setIsNotNull(param, meta)
	fo.setEqual(param, meta)
	fo.setGreaterThan(param, meta)
	fo.setLessThan(param, meta)
	fo.setGreaterEqualThan(param, meta)
	fo.setLessEqualThan(param, meta)
	fo.setChecksum(param)
	fo.setOrder(param, meta)
	fo.setOffsetAndLimit(param)

	return fo
}

func (so *FetchOptions) setColumn(param url.Values, meta TableMeta) {

	so.Column = []string{}

	if i, ok := param["cols"]; ok {
		j := strings.Split(i[0], ",")
		for _, v := range j {
			if lib.Contains(meta.Fields, v) {
				so.Column = append(so.Column, v)
			}
		}
	} else {
		so.Column = meta.Fields
	}

	//xcols
	if i, ok := param["xcols"]; ok {
		so.Column = lib.Diff(so.Column, strings.Split(i[0], ","))
	}
}

func (so *FetchOptions) setIsNull(param url.Values, meta TableMeta) {

	so.IsNull = []string{}

	if i, ok := param["isnull"]; ok {
		colsAux := make(map[string]string)
		for _, v := range strings.Split(i[0], ",") {
			colsAux[v] = v
		}
		for _, k := range meta.Fields {
			if _, ok := colsAux[k]; ok {
				so.IsNull = append(so.IsNull, k)
			}
		}
	}
}

func (so *FetchOptions) setIsNotNull(param url.Values, meta TableMeta) {

	so.IsNotNull = []string{}

	if i, ok := param["notnull"]; ok {
		colsAux := make(map[string]string)
		for _, v := range strings.Split(i[0], ",") {
			colsAux[v] = v
		}
		for _, k := range meta.Fields {
			if _, ok := colsAux[k]; ok {
				so.IsNotNull = append(so.IsNotNull, k)
			}
		}
	}
}

func (so *FetchOptions) setEqual(param url.Values, meta TableMeta) {

	so.Equal = make(map[string]string)

	if i, ok := param["eq"]; ok {
		for _, k := range i {
			j := strings.Split(k, "|")
			if lib.Contains(meta.Fields, j[0]) {
				so.Equal[j[0]] = j[1]
			}
		}
	}
}

func (so *FetchOptions) setGreaterThan(param url.Values, meta TableMeta) {

	so.GreaterThan = make(map[string]string)

	if i, ok := param["gt"]; ok {
		for _, k := range i {
			j := strings.Split(k, "|")
			if lib.Contains(meta.Fields, j[0]) {
				so.GreaterThan[j[0]] = j[1]
			}
		}
	}
}

func (so *FetchOptions) setLessThan(param url.Values, meta TableMeta) {

	so.LessThan = make(map[string]string)

	if i, ok := param["lt"]; ok {
		for _, k := range i {
			j := strings.Split(k, "|")
			if lib.Contains(meta.Fields, j[0]) {
				so.LessThan[j[0]] = j[1]
			}
		}
	}
}

func (so *FetchOptions) setGreaterEqualThan(param url.Values, meta TableMeta) {

	so.GreaterEqualThan = make(map[string]string)

	if i, ok := param["gteq"]; ok {
		for _, k := range i {
			j := strings.Split(k, "|")
			if lib.Contains(meta.Fields, j[0]) {
				so.GreaterEqualThan[j[0]] = j[1]
			}
		}
	}
}

func (so *FetchOptions) setLessEqualThan(param url.Values, meta TableMeta) {

	so.LessEqualThan = make(map[string]string)

	if i, ok := param["lteq"]; ok {
		for _, k := range i {
			j := strings.Split(k, "|")
			if lib.Contains(meta.Fields, j[0]) {
				so.LessEqualThan[j[0]] = j[1]
			}
		}
	}
}

func (so *FetchOptions) setChecksum(param url.Values) { //checksum

	if checksum, ok := param["checksum"]; ok && (checksum[0] == "1") {
		so.Checksum = 1
	}
}

func (so *FetchOptions) setOrder(param url.Values, meta TableMeta) {

	so.Order = []string{}

	if i, ok := param["order"]; ok {
		j := strings.Split(i[0], ";")
		for _, k := range j {
			j := strings.Split(k, "|")
			if !lib.Contains(meta.Fields, j[0]) {
				continue
			} else if len(j) == 1 {
				so.Order = append(so.Order, j[0]+" ASC")
			} else if strings.ToUpper(j[1]) == "ASC" || strings.ToUpper(j[1]) == "DESC" {
				so.Order = append(so.Order, j[0]+" "+strings.ToUpper(j[1]))
			}
		}
	}
}

func (so *FetchOptions) setOffsetAndLimit(param url.Values) {

	//offset and limit
	so.Offset = 0
	so.Limit = config.Config().GetInt("param.icpp")

	if i, ok := param["limit"]; ok {
		j := strings.Split(i[0], ",")
		switch len(j) {
		case 1:
			so.Offset = 0
			so.Limit, _ = strconv.Atoi(j[0])
		case 2:
			so.Offset, _ = strconv.Atoi(j[0])
			so.Limit, _ = strconv.Atoi(j[1])
		}
	}
}
