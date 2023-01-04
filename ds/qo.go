package ds

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/lib"
)

type Params map[string]string

type uparams map[string]string

type qparams map[string][]string

type ParamType int64

const (
	Primary ParamType = iota
	Url
	Qry
)

// QueryOpts exported
type QueryOptions struct {
	User             User
	DataSource       IDataSource
	Fields           []string
	WritableFields   []string
	WritableValues   []interface{}
	Checksum         int
	Dig              []string
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

/*
func (qo *QueryOptions) IsPrimary() bool {

	return len(qo.Equal[Primary]) > 0
}

func (qo *QueryOptions) SetLimit(limit *int) *QueryOptions {

	qo.Limit = limit
	return qo
}
*/

func (qo *QueryOptions) Encode(f string, crypto lib.ICrypto) bool {

	for k, v := range qo.WritableFields {
		if v == f {
			qo.WritableValues[k] = crypto.Encode(fmt.Sprint(qo.WritableValues[k]))
			return true
		}
	}
	return false
}

func (qo *QueryOptions) Copy(dsrc IDataSource, params Params) *QueryOptions {

	cqo := new(QueryOptions)

	cqo.Equal = make(map[ParamType]Params)

	cqo.User = qo.User
	cqo.DataSource = dsrc
	_, cqo.Fields, _, _, _ = Meta(dsrc)
	cqo.Equal[Primary] = params
	cqo.Dig = qo.Dig

	return cqo
}

// QueryOptsFactory exported
func QOFactory(c *gin.Context, d IDataSource) (*QueryOptions, *TagError) {

	qo := new(QueryOptions)

	keys, flds, wflds, wvals, err := Meta(d)
	if err != nil {
		return qo, err
	}

	qpar := make(qparams)
	for k, v := range c.Request.URL.Query() {
		qpar[k] = v
	}

	upar := make(uparams)
	for _, up := range c.Params {
		upar[up.Key] = up.Value
	}

	qo.setUser(c)
	qo.DataSource = d
	qo.Fields = flds
	qo.WritableFields = wflds
	qo.WritableValues = wvals
	qo.setChecksum(qpar)
	qo.setDig(qpar)

	// If Equal for Primary params is all set
	// we are done here, no more options are needed.
	if pk := qo.setEqual(upar, qpar, keys, flds); pk {
		return qo, nil
	}

	qo.setIsNull(qpar, flds)
	qo.setIsNotNull(qpar, flds)
	qo.setIn(qpar, flds)
	qo.setNotIn(qpar, flds)
	qo.setNotEqual(qpar, flds)
	qo.setGreaterThan(qpar, flds)
	qo.setLessThan(qpar, flds)
	qo.setGreaterEqualThan(qpar, flds)
	qo.setLessEqualThan(qpar, flds)
	qo.setOrder(qpar, flds)
	qo.setOffset(qpar)
	qo.setLimit(qpar)

	return qo, nil
}

func (qo *QueryOptions) setUser(c *gin.Context) {

	qo.User = User{}

	u, _ := c.Get("User")
	if u, ok := u.(User); ok {
		qo.User = u
	}
}

func (qo *QueryOptions) setChecksum(qpar qparams) {

	qo.Checksum = 0

	if checksum, ok := qpar["checksum"]; ok && (checksum[0] == "1") {
		qo.Checksum = 1
	}
}

func (qo *QueryOptions) setDig(qpar qparams) {

	qo.Dig = []string{}

	if dig, ok := qpar["dig"]; ok {
		qo.Dig = dig
	}
}

// Loads Equal params.
// Returns true if Equal for the Primary ParamType is set.
// In such case, other ParamType's are omitted from the Equal map.
func (qo *QueryOptions) setEqual(upar uparams, qpar qparams, keys []string, flds []string) (pk bool) {

	// Initialize qo.Equal
	qo.Equal = make(map[ParamType]Params)

	qo.Equal[Primary] = make(Params)
	qo.Equal[Url] = make(Params)
	qo.Equal[Qry] = make(Params)

	// Set Equal for Primary params
	for _, k := range keys {
		if v, ok := upar[k]; ok {
			qo.Equal[Primary][k] = v
		}
	}

	// Set Equal for Url params
	for _, k := range flds {
		if v, ok := upar[k]; ok {
			qo.Equal[Url][k] = v
		}
	}

	// Set Equal for Query params
	if eq, ok := qpar["eq"]; ok {
		for _, k := range eq {
			j := strings.Split(k, "|")
			if lib.Contains(flds, j[0]) {
				qo.Equal[Qry][j[0]] = j[1]
			}
		}
	}

	// If Equal for Primary params is all set remove those entries from Url and Qry params.
	// If not, reset incomplete Primary params.
	// Incomplete Primary params will be catched as Url params anyways.
	if (len(qo.Equal[Primary]) == len(keys)) && (len(keys) > 0) {
		for k := range qo.Equal[Primary] {
			delete(qo.Equal[Url], k)
			delete(qo.Equal[Qry], k)
		}
		return true
	} else {
		qo.Equal[Primary] = make(Params)
		return false
	}

}

func (qo *QueryOptions) setIsNull(qpar qparams, flds []string) {

	qo.IsNull = []string{}

	if isnull, ok := qpar["isnull"]; ok {
		for _, k := range flds {
			if lib.Contains(isnull, k) {
				qo.IsNull = append(qo.IsNull, k)
			}
		}
	}
}

func (qo *QueryOptions) setIsNotNull(qpar qparams, flds []string) {

	qo.IsNotNull = []string{}

	if notnull, ok := qpar["notnull"]; ok {
		for _, k := range flds {
			if lib.Contains(notnull, k) {
				qo.IsNotNull = append(qo.IsNotNull, k)
			}
		}
	}
}

func (qo *QueryOptions) setIn(qpar qparams, flds []string) {

	qo.In = make(map[string][]interface{})

	if in, ok := qpar["in"]; ok {
		for _, k := range in {
			j := strings.Split(k, "|")
			if lib.Contains(flds, j[0]) {
				for _, v := range strings.Split(j[1], ",") {
					qo.In[j[0]] = append(qo.In[j[0]], v)
				}
			}
		}
	}
}

func (qo *QueryOptions) setNotIn(qpar qparams, flds []string) {

	qo.NotIn = make(map[string][]interface{})

	if notin, ok := qpar["notin"]; ok {
		for _, k := range notin {
			j := strings.Split(k, "|")
			if lib.Contains(flds, j[0]) {
				for _, v := range strings.Split(j[1], ",") {
					qo.NotIn[j[0]] = append(qo.NotIn[j[0]], v)
				}
			}
		}
	}
}

func (qo *QueryOptions) setNotEqual(qpar qparams, flds []string) {

	qo.NotEqual = make(Params)

	if noteq, ok := qpar["noteq"]; ok {
		for _, k := range noteq {
			j := strings.Split(k, "|")
			if lib.Contains(flds, j[0]) {
				qo.NotEqual[j[0]] = j[1]
			}
		}
	}
}

func (qo *QueryOptions) setGreaterThan(qpar qparams, flds []string) {

	qo.GreaterThan = make(Params)

	if gt, ok := qpar["gt"]; ok {
		for _, k := range gt {
			j := strings.Split(k, "|")
			if lib.Contains(flds, j[0]) {
				qo.GreaterThan[j[0]] = j[1]
			}
		}
	}
}

func (qo *QueryOptions) setLessThan(qpar qparams, flds []string) {

	qo.LessThan = make(Params)

	if lt, ok := qpar["lt"]; ok {
		for _, k := range lt {
			j := strings.Split(k, "|")
			if lib.Contains(flds, j[0]) {
				qo.LessThan[j[0]] = j[1]
			}
		}
	}
}

func (qo *QueryOptions) setGreaterEqualThan(qpar qparams, flds []string) {

	qo.GreaterEqualThan = make(Params)

	if gteq, ok := qpar["gteq"]; ok {
		for _, k := range gteq {
			j := strings.Split(k, "|")
			if lib.Contains(flds, j[0]) {
				qo.GreaterEqualThan[j[0]] = j[1]
			}
		}
	}
}

func (qo *QueryOptions) setLessEqualThan(qpar qparams, flds []string) {

	qo.LessEqualThan = make(Params)

	if lteq, ok := qpar["lteq"]; ok {
		for _, k := range lteq {
			j := strings.Split(k, "|")
			if lib.Contains(flds, j[0]) {
				qo.LessEqualThan[j[0]] = j[1]
			}
		}
	}
}

func (qo *QueryOptions) setOrder(qpar qparams, flds []string) {

	qo.Order = []string{}

	if order, ok := qpar["order"]; ok {
		for _, k := range order {
			j := strings.Split(k, "|")
			if !lib.Contains(flds, j[0]) {
				continue
			} else if len(j) == 1 || strings.ToUpper(j[1]) == "ASC" {
				qo.Order = append(qo.Order, j[0]+" ASC")
			} else if strings.ToUpper(j[1]) == "DESC" {
				qo.Order = append(qo.Order, j[0]+" DESC")
			}
		}
	}
}

func (qo *QueryOptions) setOffset(qpar qparams) {

	qo.Offset = 0

	if offset, ok := qpar["offset"]; ok {
		qo.Offset, _ = strconv.Atoi(offset[0])
	}
}

func (qo *QueryOptions) setLimit(qpar qparams) {

	limit := config.Config().GetInt("param.icpp")

	if l, ok := qpar["limit"]; ok {
		limit, _ = strconv.Atoi(l[0])
	}

	qo.Limit = &limit
}
