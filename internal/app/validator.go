package app

import "github.com/go-playground/validator/v10"

func (app *App) initValidator() {
	app.validator = validator.New()
}
