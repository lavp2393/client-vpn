package linux

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

// LinuxPlatform implementa Platform para Linux
type LinuxPlatform struct{}

// NewLinux crea una nueva instancia de LinuxPlatform
func NewLinux() *LinuxPlatform {
	return &LinuxPlatform{}
}

// FindOpenVPN busca el ejecutable de OpenVPN en rutas comunes de Linux
func (p *LinuxPlatform) FindOpenVPN() (string, error) {
	// Rutas comunes para OpenVPN en Linux
	paths := []string{
		"/usr/sbin/openvpn",
		"/usr/bin/openvpn",
		"/usr/local/sbin/openvpn",
		"/usr/local/bin/openvpn",
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
		return "", fmt.Errorf("OpenVPN no está instalado. Por favor instala openvpn: sudo apt install openvpn")
	}

	return path, nil
}

// StartOpenVPN inicia el proceso de OpenVPN con elevación de privilegios usando pkexec
func (p *LinuxPlatform) StartOpenVPN(config StartConfig) (*Process, error) {
	// Buscar OpenVPN
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
		"--verb", "4", // Más verbosidad para debugging
	}

	// Elevar comando con pkexec
	elevatedCmd, elevatedArgs, err := p.ElevateCommand(openvpnPath, args)
	if err != nil {
		return nil, err
	}

	// Crear comando
	cmd := exec.Command(elevatedCmd, elevatedArgs...)

	// Capturar stdout y stderr para debugging
	if config.LogCallback != nil {
		stdout, err := cmd.StdoutPipe()
		if err == nil {
			go func() {
				buf := make([]byte, 1024)
				for {
					n, err := stdout.Read(buf)
					if n > 0 {
						config.LogCallback("[OpenVPN stdout] " + string(buf[:n]))
					}
					if err != nil {
						break
					}
				}
			}()
		}

		stderr, err := cmd.StderrPipe()
		if err == nil {
			go func() {
				buf := make([]byte, 1024)
				for {
					n, err := stderr.Read(buf)
					if n > 0 {
						config.LogCallback("[OpenVPN stderr] " + string(buf[:n]))
					}
					if err != nil {
						break
					}
				}
			}()
		}
	}

	// Configurar para poder matar el proceso
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// Iniciar el proceso
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("error al iniciar OpenVPN: %w", err)
	}

	proc := &Process{
		Cmd:     cmd,
		Running: true,
		PID:     cmd.Process.Pid,
	}

	return proc, nil
}

// StopOpenVPN detiene el proceso de OpenVPN limpiamente
func (p *LinuxPlatform) StopOpenVPN(proc *Process) error {
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
func (p *LinuxPlatform) RequiresElevation() bool {
	return true
}

// ElevateCommand prepara un comando para ejecutarse con privilegios elevados usando pkexec
func (p *LinuxPlatform) ElevateCommand(path string, args []string) (string, []string, error) {
	// Verificar que pkexec esté disponible
	if _, err := exec.LookPath("pkexec"); err != nil {
		return "", nil, fmt.Errorf("pkexec no está disponible. Instala con: sudo apt install policykit-1")
	}

	// pkexec requiere el path completo como primer argumento, seguido de los args
	elevatedArgs := append([]string{path}, args...)
	return "pkexec", elevatedArgs, nil
}

// GetConfigDir retorna el directorio de configuración de la aplicación
func (p *LinuxPlatform) GetConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	// Seguir XDG Base Directory Specification
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		configHome = filepath.Join(homeDir, ".config")
	}

	return filepath.Join(configHome, "PreyVPN")
}

// GetDefaultConfigPath retorna la ruta esperada del archivo de configuración VPN
func (p *LinuxPlatform) GetDefaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	// Por ahora mantenemos ~/PreyVPN para MVP (compatibilidad)
	return filepath.Join(homeDir, "PreyVPN", "prey-prod.ovpn")
}

// GetLogPath retorna la ruta del archivo de logs
func (p *LinuxPlatform) GetLogPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	// Seguir XDG Base Directory Specification
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		cacheHome = filepath.Join(homeDir, ".cache")
	}

	return filepath.Join(cacheHome, "PreyVPN", "logs")
}

// Name retorna el nombre de la plataforma
func (p *LinuxPlatform) Name() string {
	return "linux"
}

// Separator retorna el separador de rutas de la plataforma
func (p *LinuxPlatform) Separator() string {
	return "/"
}
