package rgm

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/zicare/rgm/acl"
	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/db"
	"github.com/zicare/rgm/msg"
	"github.com/zicare/rgm/tps"
)

func Init(environment string, grants db.Table, messages []msg.Message) (verbose []string, err error) {

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

	// Initialize msg
	if err = msg.Init(messages); err != nil {
		return verbose, err
	} else {
		verbose = append(verbose, "System messages... ok")
	}

	//initialize db
	if err = db.Init(); err != nil {
		return verbose, err
	} else {
		verbose = append(verbose, "DB connection... ok")
	}

	//initialize acl
	if err = acl.Init(grants); err != nil {
		return verbose, err
	} else {
		verbose = append(verbose, "Access control list... ok")
	}

	//Initialize tps control
	if err = tps.Init(); err != nil {
		return verbose, err
	} else {
		verbose = append(verbose, "TPS control... ok")
	}

	return verbose, nil
}
