package windows

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Process representa un proceso en ejecución
type Process struct {
	Cmd     *exec.Cmd
	Running bool
	PID     int
}

// StartConfig contiene la configuración de inicio
type StartConfig struct {
	ConfigPath  string
	MgmtPort    int
	LogCallback func(string)
}

// WindowsPlatform implementa Platform para Windows
type WindowsPlatform struct{}

// NewWindows crea una nueva instancia de WindowsPlatform
func NewWindows() *WindowsPlatform {
	return &WindowsPlatform{}
}

// FindOpenVPN busca el ejecutable de OpenVPN en rutas comunes de Windows
func (p *WindowsPlatform) FindOpenVPN() (string, error) {
	// Rutas comunes para OpenVPN en Windows
	paths := []string{
		`C:\Program Files\OpenVPN\bin\openvpn.exe`,
		`C:\Program Files (x86)\OpenVPN\bin\openvpn.exe`,
		filepath.Join(os.Getenv("ProgramFiles"), "OpenVPN", "bin", "openvpn.exe"),
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "OpenVPN", "bin", "openvpn.exe"),
	}

	// Intentar rutas fijas primero
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// Buscar en PATH
	path, err := exec.LookPath("openvpn.exe")
	if err != nil {
		return "", fmt.Errorf("OpenVPN no está instalado. Por favor descarga e instala desde https://openvpn.net/community-downloads/")
	}

	return path, nil
}

// StartOpenVPN inicia el proceso de OpenVPN con elevación de privilegios usando UAC
func (p *WindowsPlatform) StartOpenVPN(config StartConfig) (*Process, error) {
	// TODO: Implementar inicio de OpenVPN en Windows
	// Consideraciones:
	// - Usar runas o ShellExecute con "runas" verb para elevación
	// - O usar OpenVPN GUI service si está instalado
	// - Manejar User Account Control (UAC) prompt

	openvpnPath, err := p.FindOpenVPN()
	if err != nil {
		return nil, err
	}

	// Construir argumentos base de OpenVPN
	args := []string{
		"--config", config.ConfigPath,
		"--management", "127.0.0.1", fmt.Sprintf("%d", config.MgmtPort),
		"--management-query-passwords",
		"--management-hold",
		"--auth-retry", "interact",
		"--auth-nocache",
		"--verb", "4",
	}

	// TODO: Implementar elevación en Windows
	_ = openvpnPath
	_ = args

	return nil, fmt.Errorf("Windows support not yet implemented")
}

// StopOpenVPN detiene el proceso de OpenVPN limpiamente
func (p *WindowsPlatform) StopOpenVPN(proc *Process) error {
	if !proc.Running || proc.Cmd.Process == nil {
		return nil
	}

	// En Windows, usar Process.Kill() directamente
	// TODO: Investigar si hay una forma más "limpia" de terminar OpenVPN en Windows
	if err := proc.Cmd.Process.Kill(); err != nil {
		return fmt.Errorf("error al cerrar OpenVPN: %w", err)
	}

	// Esperar a que termine
	done := make(chan error, 1)
	go func() {
		done <- proc.Cmd.Wait()
	}()

	select {
	case <-done:
		proc.Running = false
		return nil
	case <-time.After(5 * time.Second):
		// Ya debería estar muerto, pero marcar como no corriendo
		proc.Running = false
		return nil
	}
}

// RequiresElevation indica si la plataforma requiere elevación de privilegios
func (p *WindowsPlatform) RequiresElevation() bool {
	return true
}

// ElevateCommand prepara un comando para ejecutarse con privilegios elevados usando UAC
func (p *WindowsPlatform) ElevateCommand(path string, args []string) (string, []string, error) {
	// TODO: Implementar elevación con UAC
	// Opciones:
	// 1. Usar runas: runas /user:Administrator "comando"
	// 2. Usar ShellExecute con verb "runas" (requiere syscall)
	// 3. Usar PowerShell: Start-Process -Verb RunAs

	return "", nil, fmt.Errorf("Windows elevation not yet implemented")
}

// GetConfigDir retorna el directorio de configuración de la aplicación
func (p *WindowsPlatform) GetConfigDir() string {
	// En Windows usar %APPDATA%\PreyVPN
	appData := os.Getenv("APPDATA")
	if appData == "" {
		homeDir, _ := os.UserHomeDir()
		appData = filepath.Join(homeDir, "AppData", "Roaming")
	}
	return filepath.Join(appData, "PreyVPN")
}

// GetDefaultConfigPath retorna la ruta esperada del archivo de configuración VPN
func (p *WindowsPlatform) GetDefaultConfigPath() string {
	return filepath.Join(p.GetConfigDir(), "prey-prod.ovpn")
}

// GetLogPath retorna la ruta del archivo de logs
func (p *WindowsPlatform) GetLogPath() string {
	// En Windows usar %LOCALAPPDATA%\PreyVPN\logs
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		homeDir, _ := os.UserHomeDir()
		localAppData = filepath.Join(homeDir, "AppData", "Local")
	}
	return filepath.Join(localAppData, "PreyVPN", "logs")
}

// Name retorna el nombre de la plataforma
func (p *WindowsPlatform) Name() string {
	return "windows"
}

// Separator retorna el separador de rutas de la plataforma
func (p *WindowsPlatform) Separator() string {
	return "\\"
}
