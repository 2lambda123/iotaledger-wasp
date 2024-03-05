package vmimpl

import (
	"github.com/iotaledger/wasp/packages/isc"
)

var _ isc.LogInterface = &requestContext{}

func (reqctx *requestContext) LogInfof(format string, params ...interface{}) {
	reqctx.vm.task.Log.LogInfof(format, params...)
}

func (reqctx *requestContext) LogDebugf(format string, params ...interface{}) {
	reqctx.vm.task.Log.LogDebugf(format, params...)
}

func (reqctx *requestContext) LogPanicf(format string, params ...interface{}) {
	reqctx.vm.task.Log.LogPanicf(format, params...)
}

func (reqctx *requestContext) LogWarnf(format string, params ...interface{}) {
	reqctx.vm.task.Log.LogWarnf(format, params...)
}
