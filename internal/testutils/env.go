package testutils

import (
	"os"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/nvsco/wallet/internal/utils"
	"github.com/subosito/gotenv"
)

var once sync.Once

func ExportEnvForTestsOnce() error {
	var err error

	once.Do(func() {
		var path string
		path, err = utils.LocateNearestFile(".env")
		if err == nil {
			err = gotenv.OverLoad(path)
			if err != nil {
				return
			}
		}

		if os.Getenv("CI") == "true" {
			// don't export local overrides on CI/CD

			spew.Dump("CI/CD Environment detected, skipping exporting overridden ENV variables")
			return
		}

		path, err = utils.LocateNearestFile(".env.local")
		if err == nil {
			err = gotenv.OverLoad(path)
			if err != nil {
				return
			}
		}
	})

	return err
}
