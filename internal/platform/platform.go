package platform

import (
	"os/exec"
	"runtime"
)

// Process representa un proceso de OpenVPN en ejecución
type Process struct {
	Cmd     *exec.Cmd
	Running bool
	PID     int
}

// StartConfig contiene la configuración para iniciar OpenVPN
type StartConfig struct {
	ConfigPath  string
	MgmtPort    int
	LogCallback func(string)
}

// Platform define la interfaz para operaciones específicas de cada plataforma
type Platform interface {
	// Process management
	FindOpenVPN() (string, error)
	StartOpenVPN(config StartConfig) (*Process, error)
	StopOpenVPN(proc *Process) error

	// Privilege elevation
	RequiresElevation() bool
	ElevateCommand(path string, args []string) (string, []string, error)

	// Paths
	GetConfigDir() string
	GetDefaultConfigPath() string
	GetLogPath() string

	// Platform info
	Name() string
	Separator() string
}

// New retorna la implementación de Platform para el sistema operativo actual
func New() Platform {
	switch runtime.GOOS {
	case "linux":
		return NewLinux()
	case "windows":
		return NewWindows()
	case "darwin":
		return NewDarwin()
	default:
		panic("unsupported platform: " + runtime.GOOS)
	}
}
