package darwin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
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

// DarwinPlatform implementa Platform para macOS
type DarwinPlatform struct{}

// NewDarwin crea una nueva instancia de DarwinPlatform
func NewDarwin() *DarwinPlatform {
	return &DarwinPlatform{}
}

// FindOpenVPN busca el ejecutable de OpenVPN en rutas comunes de macOS
func (p *DarwinPlatform) FindOpenVPN() (string, error) {
	// Rutas comunes para OpenVPN en macOS
	paths := []string{
		"/usr/local/opt/openvpn/sbin/openvpn",      // Homebrew (Intel)
		"/opt/homebrew/opt/openvpn/sbin/openvpn",   // Homebrew (Apple Silicon)
		"/usr/local/sbin/openvpn",
		"/usr/local/bin/openvpn",
		"/Applications/Tunnelblick.app/Contents/Resources/openvpn", // Tunnelblick
	}

	// Intentar rutas fijas primero
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// Buscar en PATH
	path, err := exec.LookPath("openvpn")
	if err != nil {
		return "", fmt.Errorf("OpenVPN no está instalado. Por favor instala con: brew install openvpn")
	}

	return path, nil
}

// StartOpenVPN inicia el proceso de OpenVPN con elevación de privilegios usando osascript
func (p *DarwinPlatform) StartOpenVPN(config StartConfig) (*Process, error) {
	// TODO: Implementar inicio de OpenVPN en macOS
	// Consideraciones:
	// - Usar osascript (AppleScript) para pedir privilegios de admin
	// - O usar AuthorizationExecuteWithPrivileges (deprecated pero funciona)
	// - Mejor: usar un helper tool con SMJobBless

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

	// TODO: Implementar elevación en macOS
	_ = openvpnPath
	_ = args

	return nil, fmt.Errorf("macOS support not yet implemented")
}

// StopOpenVPN detiene el proceso de OpenVPN limpiamente
func (p *DarwinPlatform) StopOpenVPN(proc *Process) error {
	if !proc.Running || proc.Cmd.Process == nil {
		return nil
	}

	// Intentar SIGTERM primero (cierre limpio)
	if err := proc.Cmd.Process.Signal(syscall.SIGTERM); err != nil {
		// Si ya está muerto, ok
		proc.Running = false
		return nil
	}

	// Esperar hasta 5 segundos a que termine
	done := make(chan error, 1)
	go func() {
		done <- proc.Cmd.Wait()
	}()

	select {
	case <-done:
		// Terminó limpiamente
		proc.Running = false
		return nil
	case <-time.After(5 * time.Second):
		// Timeout, forzar con SIGKILL
		if err := proc.Cmd.Process.Kill(); err != nil {
			return fmt.Errorf("error al forzar cierre: %w", err)
		}
		proc.Running = false
		return nil
	}
}

// RequiresElevation indica si la plataforma requiere elevación de privilegios
func (p *DarwinPlatform) RequiresElevation() bool {
	return true
}

// ElevateCommand prepara un comando para ejecutarse con privilegios elevados usando osascript
func (p *DarwinPlatform) ElevateCommand(path string, args []string) (string, []string, error) {
	// TODO: Implementar elevación con osascript
	// Ejemplo:
	// osascript -e "do shell script \"comando\" with administrator privileges"

	return "", nil, fmt.Errorf("macOS elevation not yet implemented")
}

// GetConfigDir retorna el directorio de configuración de la aplicación
func (p *DarwinPlatform) GetConfigDir() string {
	// En macOS usar ~/Library/Application Support/PreyVPN
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, "Library", "Application Support", "PreyVPN")
}

// GetDefaultConfigPath retorna la ruta esperada del archivo de configuración VPN
func (p *DarwinPlatform) GetDefaultConfigPath() string {
	return filepath.Join(p.GetConfigDir(), "prey-prod.ovpn")
}

// GetLogPath retorna la ruta del archivo de logs
func (p *DarwinPlatform) GetLogPath() string {
	// En macOS usar ~/Library/Logs/PreyVPN
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, "Library", "Logs", "PreyVPN")
}

// Name retorna el nombre de la plataforma
func (p *DarwinPlatform) Name() string {
	return "darwin"
}

// Separator retorna el separador de rutas de la plataforma
func (p *DarwinPlatform) Separator() string {
	return "/"
}
