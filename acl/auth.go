package acl

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/lib"
	"github.com/zicare/rgm/msg"
	"github.com/zicare/rgm/tps"

	sqlbuilder "github.com/huandu/go-sqlbuilder"
	"golang.org/x/crypto/bcrypt"

	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/db"
	"github.com/zicare/rgm/jwt"
)

// BasicAuth executes HTTP basic authentication.
// If passed, a new key/value pair is stored in the request context.
// key: "User"
// value: acl.User
// In order to pass validation the username and password must be correct,
// user's system access date range must be valid
// and her role can't be nil
func BasicAuth(m db.Table) gin.HandlerFunc {

	return func(c *gin.Context) {

		var (
			u      jwt.User
			now    = time.Now()
			ms     = sqlbuilder.NewStruct(m).For(sqlbuilder.MySQL)
			pepper = config.Config().GetString("pepper")
		)

		//get usr and pwd from http request headers
		username, password, ok := c.Request.BasicAuth()
		if ok == false {
			abort(c, 401, msg.Get("3")) //HTTP basic authentication required
			return
		}

		//get m table's auth-related field names
		f, err := lib.TaggedFields(m, "auth", []string{"id", "role", "tps", "usr", "pwd", "from", "to"})
		if err != nil {
			abort(c, 500, msg.Get("30")) //Auth tags are not properly set
			return
		}

		//find user
		sb := ms.SelectFrom(m.View())
		sb.Select(f...)
		sb.Where(sb.Equal(f[3], username), sb.IsNotNull(f[1]))
		sql, args := sb.Build()
		if err := db.Db().QueryRow(sql, args...).Scan(&u.Id, &u.Role, &u.TPS, &u.Usr, &u.Pwd, &u.From, &u.To); err != nil {
			abort(c, 401, msg.Get("4")) //Invalid credentials
			return
		}
		u.Src = m.Table()

		//pwd creation
		//hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(password+pepper), bcrypt.DefaultCost)
		//hpwd := string(hashedBytes)
		//fmt.Println(hpwd)

		err = bcrypt.CompareHashAndPassword([]byte(u.Pwd), []byte(password+pepper))
		switch err {
		case nil:
			//pwd okay
		case bcrypt.ErrMismatchedHashAndPassword:
			//wrong password
			abort(c, 401, msg.Get("4")) //Invalid credentials
			return
		default:
			//something went wrong
			//possibly a mal formed hashed password in database
			abort(c, 401, msg.Get("4")) //Invalid credentials
			return
		}

		if now.Before(u.From) || now.After(u.To) {
			abort(c, 401, msg.Get("6")) //Credentials expired or not yet valid
			return
		}

		c.Set("User", u)
		c.Next()
	}
}

// Auth exported.
// JWT must be correct and not expired.
// Authorization towards the ACL must be passed.
// User can't exceed her TPS rate.
// JWT can't be found in acl.RevokedJWTMap registry.
func Auth() gin.HandlerFunc {

	return func(c *gin.Context) {

		var (
			payload jwt.Payload
			secret  = config.Config().GetString("hmac_key")
		)

		//verify JWT authorization header is properly set
		token := strings.Split(c.GetHeader("Authorization"), " ")
		if (len(token) != 2) || (token[0] != "JWT") {
			abort(c, 401, msg.Get("7")) //JWT authorization header malformed
			return
		}

		//validate jwt
		var e *msg.Message
		payload, e = jwt.Decode(token[1], secret)
		if e != nil {
			abort(c, 401, *e) //invalid JWT, maybe tampered, etc.
			return
		}

		//verify revoked JWTs
		if revoked := jwt.IsRevoked(payload); revoked {
			abort(c, 401, msg.Get("11")) //Unauthorized
			return
		}

		//verify TPS
		if tps.IsEnabled() {
			if tps.Transaction(payload.Src, payload.Id, payload.TPS) != nil {
				abort(c, 401, msg.Get("10")) //TPS limit exceeded
				return
			}
		}

		//verify authorization in ACL
		a := ACL()
		g := Grant{Role: payload.Role, Route: c.FullPath(), Method: c.Request.Method}
		r, ok := a[g]
		if !ok {
			abort(c, 401, msg.Get("8")) //Not enough permissions
			return
		}

		//verify grant time range
		now := time.Now()
		if now.Before(r.From) || now.After(r.To) {
			abort(c, 401, msg.Get("9")) //Role access expired or not yet valid
			return
		}

		//all good
		c.Set("Auth", payload)
		c.Next()
	}
}

func abort(c *gin.Context, code int, msg msg.Message) {

	c.JSON(
		code,
		gin.H{"message": msg},
	)
	c.Abort()
}

/*
//TsAndUserID exported
func TsAndUserID(c *gin.Context) (*time.Time, *int64) {

	ts := time.Now()
	if jp, exists := c.Get("Auth"); !exists {
		return &ts, nil
	} else if py, ok := jp.(JwtPayload); ok {
		return &ts, &py.UserID
	}
	return &ts, nil
}

//UserID exported
func UserID(c *gin.Context) int64 {

	if jp, exists := c.Get("Auth"); !exists {
		return 0
	} else if py, ok := jp.(JwtPayload); ok {
		return py.UserID
	}
	return 0
}

//RoleID exported
func RoleID(c *gin.Context) *int64 {

	if jp, exists := c.Get("Auth"); !exists {
		return nil
	} else if py, ok := jp.(JwtPayload); ok {
		return py.RoleID
	}
	return nil
}

//Session exported
func Session(c *gin.Context) (JwtPayload, bool) {

	if jp, exists := c.Get("Auth"); !exists {
		return JwtPayload{}, false
	} else if session, ok := jp.(JwtPayload); ok {
		return session, ok
	}
	return JwtPayload{}, false
}
*/
