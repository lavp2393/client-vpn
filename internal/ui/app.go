package ui

import (
	"fmt"

	"github.com/prey/preyvpn/internal/core"
	"github.com/prey/preyvpn/internal/logs"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// AppState representa el estado de la aplicación
type AppState int

const (
	StateDisconnected AppState = iota
	StateConnecting
	StateAuthenticating
	StateConnected
	StateError
)

// App representa la aplicación principal
type App struct {
	fyneApp   fyne.App
	window    fyne.Window
	state     AppState
	logBuffer *logs.Buffer

	// UI elements
	statusLabel   *widget.Label
	connectBtn    *widget.Button
	disconnectBtn *widget.Button
	retryBtn      *widget.Button
	logView       *widget.Entry
	configStatus  *widget.Label

	// Core components
	manager *core.Manager
	sendFns core.SendFns

	// Credentials cache (in-memory for current session)
	savedUsername string
	savedPassword string
	rememberCreds bool
}

// NewApp crea una nueva instancia de la aplicación
func NewApp() *App {
	a := &App{
		fyneApp:   app.New(),
		logBuffer: logs.NewBuffer(30),
		state:     StateDisconnected,
	}

	// Intentar cargar credenciales guardadas
	if username, password, err := core.LoadCredentials(); err == nil {
		a.savedUsername = username
		a.savedPassword = password
		a.rememberCreds = true
	}

	a.window = a.fyneApp.NewWindow("PreyVPN")
	a.window.Resize(fyne.NewSize(700, 500))
	a.buildUI()

	return a
}

// buildUI construye la interfaz de usuario
func (a *App) buildUI() {
	// Status label
	a.statusLabel = widget.NewLabel("Estado: Desconectado")
	a.statusLabel.Wrapping = fyne.TextWrapWord

	// Config status
	a.configStatus = widget.NewLabel("")

	// Buttons
	a.connectBtn = widget.NewButton("Conectar", a.onConnect)
	a.disconnectBtn = widget.NewButton("Desconectar", a.onDisconnect)
	a.disconnectBtn.Disable()

	a.retryBtn = widget.NewButton("Reintentar", func() {
		a.updateConfigStatus()
		if core.CheckConfigExists() {
			a.connectBtn.Enable()
			a.retryBtn.Hide()
		}
	})
	a.retryBtn.Hide()

	// Actualizar estado del config después de crear todos los widgets
	a.updateConfigStatus()

	// Log view (read-only)
	a.logView = widget.NewMultiLineEntry()
	a.logView.Disable() // Read-only
	a.logView.SetPlaceHolder("Los logs aparecerán aquí...")

	// Layout
	buttonBox := container.NewHBox(
		a.connectBtn,
		a.disconnectBtn,
		a.retryBtn,
	)

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("PreyVPN - Cliente OpenVPN"),
			widget.NewSeparator(),
			a.configStatus,
			a.statusLabel,
			buttonBox,
			widget.NewSeparator(),
			widget.NewLabel("Logs:"),
		),
		nil,
		nil,
		nil,
		container.NewScroll(a.logView),
	)

	a.window.SetContent(content)
}

// updateConfigStatus actualiza el estado del archivo de configuración
func (a *App) updateConfigStatus() {
	if core.CheckConfigExists() {
		a.configStatus.SetText("✅ Perfil detectado: ~/PreyVPN/prey-prod.ovpn")
		a.connectBtn.Enable()
		a.retryBtn.Hide()
	} else {
		a.configStatus.SetText("❌ Perfil no encontrado\n\nPor favor coloca tu archivo prey-prod.ovpn en ~/PreyVPN/")
		a.connectBtn.Disable()
		a.retryBtn.Show()
	}
}

// onConnect maneja el evento de conectar
func (a *App) onConnect() {
	// Verificar que exista el config
	if !core.CheckConfigExists() {
		ShowError(a.window, "Error", "No se encontró el archivo de configuración en ~/PreyVPN/prey-prod.ovpn")
		return
	}

	a.addLog("Iniciando conexión VPN...")

	// Obtener la ruta del config
	configPath := core.GetConfigPath()

	// Buscar el binario de OpenVPN
	openvpnPath, err := core.FindOpenVPN()
	if err != nil {
		a.addLog("Error al buscar OpenVPN: " + err.Error())
		ShowError(a.window, "Error", "No se encontró OpenVPN. ¿Está instalado?")
		return
	}

	a.addLog(fmt.Sprintf("Usando OpenVPN: %s", openvpnPath))

	// Iniciar el manager directamente (sin Management Interface)
	// El manager se encarga de lanzar OpenVPN con pipes directos
	mgr, err := core.Start(configPath, openvpnPath)
	if err != nil {
		a.addLog("Error al iniciar OpenVPN: " + err.Error())
		ShowError(a.window, "Error", err.Error())
		return
	}

	a.manager = mgr
	a.sendFns = mgr.SendFunctions()

	// Actualizar UI
	a.setState(StateConnecting)
	a.connectBtn.Disable()
	a.disconnectBtn.Enable()

	a.addLog("Esperando prompts de autenticación...")

	// Iniciar procesamiento de eventos
	go a.handleEvents()
}

// onDisconnect maneja el evento de desconectar
func (a *App) onDisconnect() {
	a.addLog("Desconectando...")

	// El manager se encarga de matar el proceso OpenVPN cuando se llama Stop()
	if a.manager != nil {
		a.manager.Stop()
		a.addLog("Proceso OpenVPN detenido")
		a.manager = nil
	}

	a.setState(StateDisconnected)
	a.connectBtn.Enable()
	a.disconnectBtn.Disable()
}

// handleEvents procesa los eventos del manager
func (a *App) handleEvents() {
	for event := range a.manager.Events() {
		switch event.Type {
		case core.EventLogLine:
			a.addLog(event.Message)

		case core.EventAskUser:
			a.setState(StateAuthenticating)
			ShowUsernamePromptWithRemember(a.window, a.savedUsername, func(result PromptResult) {
				a.savedUsername = result.Value
				a.rememberCreds = result.Remember

				// Enviar username a OpenVPN
				if err := a.sendFns.Username(result.Value); err != nil {
					a.addLog("Error al enviar usuario: " + err.Error())
				}
			})

		case core.EventAskPass:
			a.setState(StateAuthenticating)
			// Usar prompt con valor por defecto pero sin checkbox (la decisión ya se tomó en el modal de usuario)
			ShowPasswordPromptWithDefault(a.window, a.savedPassword, func(password string) {
				a.savedPassword = password

				// Guardar o eliminar credenciales según la decisión tomada en el modal de usuario
				if a.rememberCreds {
					if err := core.SaveCredentials(a.savedUsername, a.savedPassword); err != nil {
						a.addLog("Advertencia: No se pudieron guardar las credenciales: " + err.Error())
					} else {
						a.addLog("✓ Credenciales guardadas")
					}
				} else {
					// Usuario no marcó recordar, eliminar credenciales guardadas si existían
					if err := core.DeleteCredentials(); err != nil {
						// Ignorar error si no existían
					}
					a.savedUsername = ""
					a.savedPassword = ""
				}

				// Enviar password a OpenVPN
				if err := a.sendFns.Password(password); err != nil {
					a.addLog("Error al enviar contraseña: " + err.Error())
				}
			})

		case core.EventAskOTP:
			a.setState(StateAuthenticating)
			ShowOTPPrompt(a.window, func(otp string) {
				if err := a.sendFns.OTP(otp); err != nil {
					a.addLog("Error al enviar OTP: " + err.Error())
				}
			})

		case core.EventConnected:
			a.setState(StateConnected)
			a.addLog(event.Message)
			ShowInfo(a.window, "Conectado", "Conexión VPN establecida exitosamente")

		case core.EventAuthFailed:
			a.setState(StateAuthenticating)
			a.addLog("Error: " + event.Message)
			ShowError(a.window, "Error de autenticación", event.Message)

			// Re-pedir solo el campo que falló
			if event.Stage == "password" {
				// Re-pedir password sin checkbox (usa la decisión ya tomada)
				ShowPasswordPromptWithDefault(a.window, a.savedPassword, func(password string) {
					a.savedPassword = password

					// Actualizar credenciales guardadas si estaba marcado recordar
					if a.rememberCreds {
						if err := core.SaveCredentials(a.savedUsername, a.savedPassword); err != nil {
							a.addLog("Advertencia: No se pudieron guardar las credenciales: " + err.Error())
						}
					}

					if err := a.sendFns.Password(password); err != nil {
						a.addLog("Error al enviar contraseña: " + err.Error())
					}
				})
			} else if event.Stage == "otp" {
				ShowOTPPrompt(a.window, func(otp string) {
					if err := a.sendFns.OTP(otp); err != nil {
						a.addLog("Error al enviar OTP: " + err.Error())
					}
				})
			}

		case core.EventFatal:
			a.setState(StateError)
			a.addLog("Error fatal: " + event.Message)
			ShowError(a.window, "Error Fatal", event.Message)
			a.onDisconnect()

		case core.EventDisconnected:
			a.addLog("Conexión cerrada")
			a.onDisconnect()
		}
	}
}

// setState actualiza el estado de la aplicación
func (a *App) setState(state AppState) {
	a.state = state

	switch state {
	case StateDisconnected:
		a.statusLabel.SetText("Estado: Desconectado")
	case StateConnecting:
		a.statusLabel.SetText("Estado: Conectando...")
	case StateAuthenticating:
		a.statusLabel.SetText("Estado: Autenticando...")
	case StateConnected:
		a.statusLabel.SetText("Estado: Conectado ✅")
	case StateError:
		a.statusLabel.SetText("Estado: Error ❌")
	}
}

// addLog agrega una línea al buffer de logs y actualiza la UI
func (a *App) addLog(line string) {
	a.logBuffer.Add(line)
	a.logView.SetText(a.logBuffer.GetText())

	// Auto-scroll al final
	a.logView.CursorRow = len(a.logBuffer.GetAll())
}

// Run inicia la aplicación
func (a *App) Run() {
	a.window.ShowAndRun()
}
