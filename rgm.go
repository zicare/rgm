package rgm

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/ds"
	"github.com/zicare/rgm/jwt"
	"github.com/zicare/rgm/lib"
	"github.com/zicare/rgm/msg"
	"github.com/zicare/rgm/mw"
	"github.com/zicare/rgm/validation"
)

type InitOpts struct {
	Environment  *string
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
	return append(handlersChain, h)
}

// Returns a gin.HandlersChain slice loaded with
// mw.JWTAuthentication, mw.Abuse, mw.Authorization and h.
// h is the actual controller function.
func JHC(h gin.HandlerFunc) gin.HandlersChain {

	handlersChain := gin.HandlersChain{}
	handlersChain = append(handlersChain, mw.JWTAuthentication())
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
	} else if *opts.Verbose {
		fmt.Println("Directories: config, tpl and log... OK")
	}
	flag.Set("log_dir", dir+"/log")
	flag.Set("stderrthreshold", "FATAL")

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

	// Load ds.Acl map in memory
	if (opts.AclDSFactory == nil) || (opts.Acl == nil) {
		fmt.Println("ACL... Not loaded")
	} else if err := ds.Init(opts.AclDSFactory, opts.Acl); err != nil {
		return err
	} else if *opts.Verbose {
		fmt.Println("ACL... OK")
	}

	// Initialize jwt.revokedJWTMap
	jwt.Init()
	fmt.Println("JWT revokes... OK")

	//  Custom validation
	validation.Init()

	//Run deferred inits
	//used to populate the backend meta registry
	ds.RunDeferredInits()

	return nil
}
