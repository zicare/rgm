package db

import (
	"strconv"
	"strings"

	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/lib"
)

type ParamType int64

const (
	Primary ParamType = iota
	Url
	Query
)

// QueryOptions exported
type QueryOptions struct {
	UID              string // ID from jwt, used to set scope in find, fetch, etc.
	Table            Table
	Column           []string
	Parents          []Table // Also used to set scope in find, fetch, etc.
	Checksum         int
	Dig              []string // Used to include parent data in find, fetch, etc.
	Equal            map[ParamType]Params
	IsNull           []string
	IsNotNull        []string
	In               map[string][]interface{}
	NotIn            map[string][]interface{}
	NotEqual         Params
	GreaterThan      Params
	LessThan         Params
	GreaterEqualThan Params
	LessEqualThan    Params
	Order            []string
	Offset           int
	Limit            *int
}

func (qo *QueryOptions) IsPrimary() bool {

	return len(qo.Equal[Primary]) > 0
}

func (qo *QueryOptions) SetLimit(limit *int) *QueryOptions {

	qo.Limit = limit
	return qo
}

// QueryOptionsFactory exported
func QueryOptionsFactory(t Table, uid string, qpar QParams, upar UParams,
	p ...Table) *QueryOptions {

	var (
		cols = Cols(t)
		qo   = new(QueryOptions)
	)

	qo.setUID(uid)
	qo.setTable(t)
	qo.setColumn(cols, qpar)
	qo.setParents(p...)
	qo.setChecksum(qpar)
	qo.setDig(cols, qpar)

	// If Equal for Primary params is all set
	// we are done here, no more options are needed.
	if pk := qo.setEqual(cols, upar, qpar); pk {
		return qo
	}

	if len(qpar) > 0 {
		qo.setIsNull(cols, qpar)
		qo.setIsNotNull(cols, qpar)
		qo.setIn(cols, qpar)
		qo.setNotIn(cols, qpar)
		qo.setNotEqual(cols, qpar)
		qo.setGreaterThan(cols, qpar)
		qo.setLessThan(cols, qpar)
		qo.setGreaterEqualThan(cols, qpar)
		qo.setLessEqualThan(cols, qpar)
		qo.setOrder(cols, qpar)
		qo.setOffset(qpar)
		qo.setLimit(qpar)
	} else {
		limit := config.Config().GetInt("param.icpp_max")
		qo.SetLimit(&limit)
	}

	return qo
}

func (qo *QueryOptions) setUID(uid string) {

	qo.UID = uid
}

func (qo *QueryOptions) setTable(tbl Table) {

	qo.Table = tbl
}

func (qo *QueryOptions) setColumn(cols []string, qpar QParams) {

	qo.Column = cols
}

func (qo *QueryOptions) setParents(p ...Table) {

	qo.Parents = p
}

// Loads Equal params.
// Returns true if Equal for the Primary ParamType is set.
// In such case, other ParamType's are omitted from the Equal map.
func (qo *QueryOptions) setEqual(cols []string, upar UParams, qpar QParams) bool {

	pk := Pk(qo.Table)

	// Initialize qo.Equal
	qo.Equal = make(map[ParamType]Params)

	qo.Equal[Primary] = make(Params)
	qo.Equal[Url] = make(Params)
	qo.Equal[Query] = make(Params)

	// Set Equal for Primary params
	for _, k := range pk {
		if v, ok := upar[k]; ok {
			qo.Equal[Primary][k] = v
		}
	}

	// Set Equal for Url params
	for _, k := range cols {
		if v, ok := upar[k]; ok {
			qo.Equal[Url][k] = v
		}
	}

	// Set Equal for Query params
	if eq, ok := qpar["eq"]; ok {
		for _, k := range eq {
			j := strings.Split(k, "|")
			if lib.Contains(cols, j[0]) {
				qo.Equal[Query][j[0]] = j[1]
			}
		}
	}

	// If Equal for Primary params is all set
	// remove those entries from Url params.
	// If not, reset incomplete Primary params.
	// Incomplete Primary params will be catched
	// as Url params anyways.
	if (len(qo.Equal[Primary]) == len(pk)) && (len(pk) > 0) {
		for k := range qo.Equal[Primary] {
			delete(qo.Equal[Url], k)
		}
		return true
	} else {
		qo.Equal[Primary] = make(Params)
		return false
	}

}

func (qp *QueryOptions) setIsNull(cols []string, qpar QParams) {

	qp.IsNull = []string{}

	if isnull, ok := qpar["isnull"]; ok {
		for _, k := range cols {
			if lib.Contains(isnull, k) {
				qp.IsNull = append(qp.IsNull, k)
			}
		}
	}
}

func (qo *QueryOptions) setIsNotNull(cols []string, qpar QParams) {

	qo.IsNotNull = []string{}

	if notnull, ok := qpar["notnull"]; ok {
		for _, k := range cols {
			if lib.Contains(notnull, k) {
				qo.IsNotNull = append(qo.IsNotNull, k)
			}
		}
	}
}

func (qo *QueryOptions) setIn(cols []string, qpar QParams) {

	qo.In = make(map[string][]interface{})

	if in, ok := qpar["in"]; ok {
		for _, k := range in {
			j := strings.Split(k, "|")
			if lib.Contains(cols, j[0]) {
				for _, v := range strings.Split(j[1], ",") {
					qo.In[j[0]] = append(qo.In[j[0]], v)
				}
			}
		}
	}
}

func (qo *QueryOptions) setNotIn(cols []string, qpar QParams) {

	qo.NotIn = make(map[string][]interface{})

	if notin, ok := qpar["notin"]; ok {
		for _, k := range notin {
			j := strings.Split(k, "|")
			if lib.Contains(cols, j[0]) {
				for _, v := range strings.Split(j[1], ",") {
					qo.NotIn[j[0]] = append(qo.NotIn[j[0]], v)
				}
			}
		}
	}
}

func (qo *QueryOptions) setNotEqual(cols []string, qpar QParams) {

	qo.NotEqual = make(Params)

	if noteq, ok := qpar["noteq"]; ok {
		for _, k := range noteq {
			j := strings.Split(k, "|")
			if lib.Contains(cols, j[0]) {
				qo.NotEqual[j[0]] = j[1]
			}
		}
	}
}

func (qo *QueryOptions) setGreaterThan(cols []string, qpar QParams) {

	qo.GreaterThan = make(Params)

	if gt, ok := qpar["gt"]; ok {
		for _, k := range gt {
			j := strings.Split(k, "|")
			if lib.Contains(cols, j[0]) {
				qo.GreaterThan[j[0]] = j[1]
			}
		}
	}
}

func (qo *QueryOptions) setLessThan(cols []string, qpar QParams) {

	qo.LessThan = make(Params)

	if lt, ok := qpar["lt"]; ok {
		for _, k := range lt {
			j := strings.Split(k, "|")
			if lib.Contains(cols, j[0]) {
				qo.LessThan[j[0]] = j[1]
			}
		}
	}
}

func (qo *QueryOptions) setGreaterEqualThan(cols []string, qpar QParams) {

	qo.GreaterEqualThan = make(Params)

	if gteq, ok := qpar["gteq"]; ok {
		for _, k := range gteq {
			j := strings.Split(k, "|")
			if lib.Contains(cols, j[0]) {
				qo.GreaterEqualThan[j[0]] = j[1]
			}
		}
	}
}

func (qo *QueryOptions) setLessEqualThan(cols []string, qpar QParams) {

	qo.LessEqualThan = make(Params)

	if lteq, ok := qpar["lteq"]; ok {
		for _, k := range lteq {
			j := strings.Split(k, "|")
			if lib.Contains(cols, j[0]) {
				qo.LessEqualThan[j[0]] = j[1]
			}
		}
	}
}

func (qo *QueryOptions) setOrder(cols []string, qpar QParams) {

	qo.Order = []string{}

	if order, ok := qpar["order"]; ok {
		for _, k := range order {
			j := strings.Split(k, "|")
			if !lib.Contains(cols, j[0]) {
				continue
			} else if len(j) == 1 {
				qo.Order = append(qo.Order, j[0]+" ASC")
			} else if strings.ToUpper(j[1]) == "ASC" || strings.ToUpper(j[1]) == "DESC" {
				qo.Order = append(qo.Order, j[0]+" "+strings.ToUpper(j[1]))
			}
		}
	}
}

func (qo *QueryOptions) setOffset(qpar QParams) {

	if offset, ok := qpar["offset"]; ok {
		qo.Offset, _ = strconv.Atoi(offset[0])
	} else {
		qo.Offset = 0
	}
}

func (qo *QueryOptions) setLimit(qpar QParams) {

	limit := config.Config().GetInt("param.icpp")

	if qpar == nil {
		limit = 0
	} else if l, ok := qpar["limit"]; ok {
		limit, _ = strconv.Atoi(l[0])
	}

	qo.Limit = &limit
}

func (qo *QueryOptions) setChecksum(qpar QParams) {

	if checksum, ok := qpar["checksum"]; ok && (checksum[0] == "1") {
		qo.Checksum = 1
	}
}

func (qo *QueryOptions) setDig(cols []string, qpar QParams) {

	if dig, ok := qpar["dig"]; ok {
		qo.Dig = dig
	} else {
		qo.Dig = []string{}
	}

	/*
		qo.Dig = []string{}

		if dig, ok := qpar["dig"]; ok {
			for _, k := range cols {
				if lib.Contains(dig, k) {
					qo.Dig = append(qo.Dig, k)
				}
			}
		}
	*/
}
