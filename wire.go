//go:build wireinject

package main

import (
	"github.com/google/wire"
	"stori/helpers"
	"stori/internal/core/usecases"
	"stori/internal/handlers"
	"stori/internal/repositories"
)

func InitializeApp() *App {
	wire.Build(
		ProvideApp,
		adapters.ProvideTransactionRepository,
		usecases.ProvideTransactionUseCase,
		helpers.ProvideDatabaseParameters,
		helpers.ProvideDatabase,
		handlers.ProvideHttpService,
		handlers.ProvideHttpServiceParams,
		handlers.ProvideRouter,
		handlers.ProvideTransactionHttpHandler,
	)
	return &App{}
}
