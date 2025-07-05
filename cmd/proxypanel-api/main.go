package main

import (
	fxmodules "github.com/tools4net/ezfw/backend/internal/fx"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		// Include all modules
		fxmodules.ConfigModule,
		fxmodules.LoggingModule,
		fxmodules.StoreModule,
		fxmodules.AuthModule,
		fxmodules.HandlersModule,
		fxmodules.RouterModule,
		fxmodules.ServerModule,

		// Configure fx logger to use zap
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
	)

	// Run the application
	app.Run()
}
