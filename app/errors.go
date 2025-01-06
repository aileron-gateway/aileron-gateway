package app

import "github.com/aileron-gateway/aileron-gateway/kernel/errorutil"

var (
	// app: E1001 - E1999
	ErrAppMain              = errorutil.NewKind("E1001", "AppMain", "main function exit with error")
	ErrAppMainRun           = errorutil.NewKind("E1002", "AppMainRun", "running service failed")
	ErrAppMainService       = errorutil.NewKind("E1003", "AppMainService", "service operation failed")
	ErrAppMainLoadEnv       = errorutil.NewKind("E1004", "AppMainLoadEnv", "failed to load environmental variables. {{path}}")
	ErrAppMainLoadConfigs   = errorutil.NewKind("E1005", "AppMainLoadConfigs", "failed to load configs. {{path}} {{content}}")
	ErrAppMainGetEntrypoint = errorutil.NewKind("E1006", "AppMainGetEntrypoint", "failed to get entrypoint resource")
)
