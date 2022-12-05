package rgm

import (
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/zicare/rgm/auth"
	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/db"
	"github.com/zicare/rgm/jwt"
	"github.com/zicare/rgm/msg"
	"github.com/zicare/rgm/mw"
	"github.com/zicare/rgm/tps"
)

// Returns a gin.HandlersChain slice loaded
// with mw.BasicAuthentication, mw.Abuse and h.
// mw.BasicAuthentication must be passed an auth.UserDS param,
// BHC relies for this on auth.TUserDS, a default implementation
// of auth.UserDS that uses t as the user data store.
// t must be propperly annotated with auth tags,
// otherwise a 500 http response code will be issued
// when calling mw.BasicAuthentication.
func BHC(t db.Table, h gin.HandlerFunc) gin.HandlersChain {

	handlersChain := gin.HandlersChain{}
	handlersChain = append(handlersChain, mw.BasicAuthentication(t))
	handlersChain = append(handlersChain, mw.Abuse())
	return append(handlersChain, h)
}

// Returns a gin.HandlersChain slice with
// mw.JWTAuthentication, mw.Abuse, mw.Authorization and h.
func JHC(h gin.HandlerFunc) gin.HandlersChain {

	handlersChain := gin.HandlersChain{}
	handlersChain = append(handlersChain, mw.JWTAuthentication())
	handlersChain = append(handlersChain, mw.Abuse())
	handlersChain = append(handlersChain, mw.Authorization())
	return append(handlersChain, h)
}

func Init(environment string, acl db.Table, messages []msg.Message) (verbose []string, err error) {

	// Initialize config
	if dir, err := filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		return verbose, err
	} else if err = config.Init(environment, dir); err != nil {
		return verbose, err
	} else if info, err := os.Stat(dir + "/log"); err != nil || !info.IsDir() {
		return verbose, err
	} else {
		flag.Set("log_dir", dir+"/log")
		flag.Set("stderrthreshold", "FATAL")
		verbose = append(verbose, "Binary path... ok")
		verbose = append(verbose, "Config file... ok")
		verbose = append(verbose, "Log directory... ok")
	}

	// Set timezone
	if os.Setenv("TZ", config.Config().GetString("tz")); err != nil {
		return verbose, err
	} else {
		verbose = append(verbose, "Timezone set... ok")
	}

	// Initialize msg
	if err = msg.Init(messages); err != nil {
		return verbose, err
	} else {
		verbose = append(verbose, "System messages... ok")
	}

	// Initialize db
	if err = db.Init(); err != nil {
		return verbose, err
	} else {
		verbose = append(verbose, "DB connection... ok")
	}

	// Initialize auth
	if acl == nil {
		// Do nothing.
		// Looks like client app doesn't has acl in database.
		// Client app can always call auth.Init by itself passing
		// custom implementation of auth.AclDS
	} else if aclDS, err := auth.TAclDSFactory(acl); err != nil {
		return verbose, err
	} else if err := auth.Init(aclDS); err != nil {
		return verbose, err
	} else {
		verbose = append(verbose, "Access control list... ok")
	}

	// Initialize revokedJWTMap
	jwt.Init()

	// Initialize tps control
	if err = tps.Init(); err != nil {
		return verbose, err
	} else {
		verbose = append(verbose, "TPS control... ok")
	}

	// Validation setup
	// This is a workaround for FieldError.Field() bug
	// in validation v10, that returns the actual struct field name
	// instead of the json name, which is needed for custom error messages.
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	return verbose, nil
}
