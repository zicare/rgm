package rgm

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/ds"
	"github.com/zicare/rgm/jwt"
	"github.com/zicare/rgm/lib"
	"github.com/zicare/rgm/msg"
	"github.com/zicare/rgm/mw"
	"github.com/zicare/rgm/mysql"
	"github.com/zicare/rgm/tps"
	"github.com/zicare/rgm/validation"
)

type InitOpts struct {
	Environment  *string
	DisableAgent *bool
	Verbose      *bool
	Messages     []msg.Message
	AclDSFactory ds.AclDSFactory
	Acl          ds.IDataSource
}

// Returns a gin.HandlersChain slice loaded with
// mw.BasicAuthentication, mw.Abuse and h.
// h is the actual controller function.
func BHC(fn ds.UserDSFactory, u ds.IDataSource, crypto lib.ICrypto, h gin.HandlerFunc) gin.HandlersChain {

	dsrc, err := fn(u)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	handlersChain := gin.HandlersChain{}
	handlersChain = append(handlersChain, mw.BasicAuthentication(dsrc, crypto))
	handlersChain = append(handlersChain, mw.Abuse())
	return append(handlersChain, h)
}

// Returns a gin.HandlersChain slice loaded with
// mw.JWTAuthentication, mw.Abuse, mw.Authorization and h.
// h is the actual controller function.
func JHC(h gin.HandlerFunc) gin.HandlersChain {

	handlersChain := gin.HandlersChain{}
	handlersChain = append(handlersChain, mw.JWTAuthentication())
	handlersChain = append(handlersChain, mw.Abuse())
	handlersChain = append(handlersChain, mw.Authorization())
	return append(handlersChain, h)
}

func Init(opts InitOpts) error {

	// Check paths
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	} else if fi, err := os.Stat(dir + "/config"); err != nil || !fi.IsDir() {
		return err
	} else if fi, err := os.Stat(dir + "/tpl"); err != nil || !fi.IsDir() {
		return err
	} else if fi, err := os.Stat(dir + "/log"); err != nil || !fi.IsDir() {
		return err
	} else {
		flag.Set("log_dir", dir+"/log")
		flag.Set("stderrthreshold", "FATAL")
	}

	// Config
	if err := config.Init(*opts.Environment, dir); err != nil {
		return err
	} else if *opts.Verbose {
		fmt.Println("Config... OK")
	}

	// Timezone
	if os.Setenv("TZ", config.Config().GetString("tz")); err != nil {
		return err
	} else if *opts.Verbose {
		fmt.Println("Timezone... OK")
	}

	// MSG
	if err := msg.Init(opts.Messages); err != nil {
		return err
	} else if *opts.Verbose {
		fmt.Println("MSG... OK")
	}

	// MySQL
	if err := mysql.Init(); err != nil {
		return err
	} else if *opts.Verbose {
		fmt.Println("MySQL... OK")
	}

	// ACL
	if err := ds.Init(opts.AclDSFactory, opts.Acl); err != nil {
		return err
	} else if *opts.Verbose {
		fmt.Println("ACL... OK")
	}

	// Initialize revokedJWTMap
	jwt.Init()
	fmt.Println("JWT revokes... OK")

	// Initialize tps control
	if err = tps.Init(); err != nil {
		return err
	} else {
		fmt.Println("TPS control... OK")
	}

	// Agent
	if *opts.Verbose {
		fmt.Println("Agent enabled..." + strconv.FormatBool(!*opts.DisableAgent))
	}

	//  Custom validation
	validation.Init()

	return nil
}
