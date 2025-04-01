package tasker

import (
	"github.com/lcmies/go-tasker/types"
)

var mgr *types.Mgr

func Get() *types.Mgr {
	if mgr == nil {
		mgr = types.NewMgr()
	}
	return mgr
}
