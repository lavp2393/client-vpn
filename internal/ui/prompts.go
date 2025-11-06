package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// PromptCallback es el tipo de función callback para los prompts
type PromptCallback func(string)

// PromptResult contiene el resultado de un prompt con opción de recordar
type PromptResult struct {
	Value    string
	Remember bool
}

// PromptCallbackWithRemember es el callback para prompts con opción de recordar
type PromptCallbackWithRemember func(PromptResult)

// ShowUsernamePrompt muestra un modal para ingresar el usuario
func ShowUsernamePrompt(window fyne.Window, callback PromptCallback) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("usuario corporativo")

	content := widget.NewForm(
		widget.NewFormItem("Usuario:", entry),
	)

	d := dialog.NewCustomConfirm(
		"Autenticación",
		"Confirmar",
		"Cancelar",
		content,
		func(submit bool) {
			if submit && entry.Text != "" {
				callback(entry.Text)
			}
		},
		window,
	)

	d.Resize(fyne.NewSize(400, 150))
	d.Show()

	// Focus en el campo de entrada
	window.Canvas().Focus(entry)
}

// ShowUsernamePromptWithRemember muestra un modal para ingresar el usuario con opción de recordar
func ShowUsernamePromptWithRemember(window fyne.Window, defaultValue string, callback PromptCallbackWithRemember) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("usuario corporativo")

	// Pre-llenar si hay un valor por defecto
	if defaultValue != "" {
		entry.SetText(defaultValue)
	}

	// Checkbox para recordar credenciales
	rememberCheck := widget.NewCheck("Recordar credenciales", nil)
	rememberCheck.Checked = defaultValue != "" // Marcar si ya hay credenciales guardadas

	form := widget.NewForm(
		widget.NewFormItem("Usuario:", entry),
	)

	content := container.NewVBox(
		form,
		rememberCheck,
	)

	d := dialog.NewCustomConfirm(
		"Autenticación",
		"Confirmar",
		"Cancelar",
		content,
		func(submit bool) {
			if submit && entry.Text != "" {
				callback(PromptResult{
					Value:    entry.Text,
					Remember: rememberCheck.Checked,
				})
			}
		},
		window,
	)

	d.Resize(fyne.NewSize(400, 180))
	d.Show()

	// Focus en el campo de entrada
	window.Canvas().Focus(entry)
}

// ShowPasswordPrompt muestra un modal para ingresar la contraseña
func ShowPasswordPrompt(window fyne.Window, callback PromptCallback) {
	ShowPasswordPromptWithDefault(window, "", callback)
}

// ShowPasswordPromptWithDefault muestra un modal para ingresar la contraseña con valor por defecto
func ShowPasswordPromptWithDefault(window fyne.Window, defaultValue string, callback PromptCallback) {
	entry := widget.NewPasswordEntry()
	entry.SetPlaceHolder("contraseña")

	// Pre-llenar si hay un valor por defecto
	if defaultValue != "" {
		entry.SetText(defaultValue)
	}

	content := widget.NewForm(
		widget.NewFormItem("Contraseña:", entry),
	)

	d := dialog.NewCustomConfirm(
		"Autenticación",
		"Confirmar",
		"Cancelar",
		content,
		func(submit bool) {
			if submit && entry.Text != "" {
				callback(entry.Text)
			}
		},
		window,
	)

	d.Resize(fyne.NewSize(400, 150))
	d.Show()

	// Focus en el campo de entrada
	window.Canvas().Focus(entry)
}

// ShowPasswordPromptWithRemember muestra un modal para ingresar la contraseña con opción de recordar
func ShowPasswordPromptWithRemember(window fyne.Window, defaultValue string, callback PromptCallbackWithRemember) {
	entry := widget.NewPasswordEntry()
	entry.SetPlaceHolder("contraseña")

	// Pre-llenar si hay un valor por defecto
	if defaultValue != "" {
		entry.SetText(defaultValue)
	}

	// Checkbox para recordar credenciales
	rememberCheck := widget.NewCheck("Recordar credenciales", nil)
	rememberCheck.Checked = defaultValue != "" // Marcar si ya hay credenciales guardadas

	form := widget.NewForm(
		widget.NewFormItem("Contraseña:", entry),
	)

	content := container.NewVBox(
		form,
		rememberCheck,
	)

	d := dialog.NewCustomConfirm(
		"Autenticación",
		"Confirmar",
		"Cancelar",
		content,
		func(submit bool) {
			if submit && entry.Text != "" {
				callback(PromptResult{
					Value:    entry.Text,
					Remember: rememberCheck.Checked,
				})
			}
		},
		window,
	)

	d.Resize(fyne.NewSize(400, 180))
	d.Show()

	// Focus en el campo de entrada
	window.Canvas().Focus(entry)
}

// ShowOTPPrompt muestra un modal para ingresar el código OTP
func ShowOTPPrompt(window fyne.Window, callback PromptCallback) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("123456")

	// Limitar a 6 dígitos
	entry.Validator = func(s string) error {
		// Validación básica: solo números, máximo 6 caracteres
		if len(s) > 6 {
			return nil // Fyne no acepta error, solo prevenir más input
		}
		return nil
	}

	hint := widget.NewLabel("Ingresa el código de 6 dígitos (se renueva cada 30s)")
	hint.Wrapping = fyne.TextWrapWord

	content := widget.NewForm(
		widget.NewFormItem("Código OTP:", entry),
	)

	d := dialog.NewCustomConfirm(
		"Código OTP",
		"Confirmar",
		"Cancelar",
		content,
		func(submit bool) {
			if submit && entry.Text != "" {
				callback(entry.Text)
			}
		},
		window,
	)

	d.Resize(fyne.NewSize(400, 150))
	d.Show()

	// Focus en el campo de entrada
	window.Canvas().Focus(entry)
}

// ShowError muestra un diálogo de error
func ShowError(window fyne.Window, title, message string) {
	dialog.ShowError(
		fmt.Errorf("%s", message),
		window,
	)
}

// ShowInfo muestra un diálogo informativo
func ShowInfo(window fyne.Window, title, message string) {
	dialog.ShowInformation(title, message, window)
}
