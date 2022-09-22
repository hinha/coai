package logger

import (
	"strings"

	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

// LogEvent logs the given event to the provided Zap logger.
func (log *Logger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		log.Console().Info("OnStart hook executing",
			zap.String("callee", e.FunctionName),
			zap.String("caller", e.CallerName),
		)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			log.Console().Info("OnStart hook failed",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.Error(e.Err),
			)
		} else {
			log.Console().Info("OnStart hook executed",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.String("runtime", e.Runtime.String()),
			)
		}
	case *fxevent.OnStopExecuting:
		log.Console().Info("OnStop hook executing",
			zap.String("callee", e.FunctionName),
			zap.String("caller", e.CallerName),
		)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			log.Console().Info("OnStop hook failed",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.Error(e.Err),
			)
		} else {
			log.Console().Info("OnStop hook executed",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.String("runtime", e.Runtime.String()),
			)
		}
	case *fxevent.Supplied:
		log.Console().Info("supplied",
			zap.String("type", e.TypeName),
			moduleField(e.ModuleName),
			zap.Error(e.Err))
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			log.Console().Info("provided",
				zap.String("constructor", e.ConstructorName),
				moduleField(e.ModuleName),
				zap.String("type", rtype),
			)
		}
		if e.Err != nil {
			log.Console().Error("error encountered while applying options",
				moduleField(e.ModuleName),
				zap.Error(e.Err))
		}
	case *fxevent.Replaced:
		for _, rtype := range e.OutputTypeNames {
			log.Console().Info("replaced",
				moduleField(e.ModuleName),
				zap.String("type", rtype),
			)
		}
		if e.Err != nil {
			log.Console().Error("error encountered while replacing",
				moduleField(e.ModuleName),
				zap.Error(e.Err))
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			log.Console().Info("decorated",
				zap.String("decorator", e.DecoratorName),
				moduleField(e.ModuleName),
				zap.String("type", rtype),
			)
		}
		if e.Err != nil {
			log.Console().Error("error encountered while applying options",
				moduleField(e.ModuleName),
				zap.Error(e.Err))
		}
	case *fxevent.Invoking:
		// Do not log stack as it will make logs hard to read.
		log.Console().Info("invoking",
			zap.String("function", e.FunctionName),
			moduleField(e.ModuleName),
		)
	case *fxevent.Invoked:
		if e.Err != nil {
			log.Console().Error("invoke failed",
				zap.Error(e.Err),
				zap.String("stack", e.Trace),
				zap.String("function", e.FunctionName),
				moduleField(e.ModuleName),
			)
		}
	case *fxevent.Stopping:
		log.Console().Info("received signal",
			zap.String("signal", strings.ToUpper(e.Signal.String())))
	case *fxevent.Stopped:
		if e.Err != nil {
			log.Console().Error("stop failed", zap.Error(e.Err))
		}
	case *fxevent.RollingBack:
		log.Console().Error("start failed, rolling back", zap.Error(e.StartErr))
	case *fxevent.RolledBack:
		if e.Err != nil {
			log.Console().Error("rollback failed", zap.Error(e.Err))
		}
	case *fxevent.Started:
		if e.Err != nil {
			log.Console().Error("start failed", zap.Error(e.Err))
		} else {
			log.Console().Info("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			log.Console().Error("custom logger initialization failed", zap.Error(e.Err))
		} else {
			log.Console().Info("initialized custom fxevent.Logger", zap.String("function", e.ConstructorName))
		}
	}
}

func moduleField(name string) zap.Field {
	if len(name) == 0 {
		return zap.Skip()
	}
	return zap.String("module", name)
}
