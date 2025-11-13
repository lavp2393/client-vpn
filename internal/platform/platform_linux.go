//go:build linux

package platform

import "github.com/lavp2393/navtunnel/internal/platform/linux"

// linuxAdapter adapta linux.LinuxPlatform a la interfaz Platform
type linuxAdapter struct {
	impl *linux.LinuxPlatform
}

// NewLinux retorna la implementación de Platform para Linux
func NewLinux() Platform {
	return &linuxAdapter{
		impl: linux.NewLinux(),
	}
}

func (a *linuxAdapter) FindOpenVPN() (string, error) {
	return a.impl.FindOpenVPN()
}

func (a *linuxAdapter) StartOpenVPN(config StartConfig) (*Process, error) {
	// Convertir StartConfig de platform a linux
	linuxConfig := linux.StartConfig{
		ConfigPath:  config.ConfigPath,
		MgmtPort:    config.MgmtPort,
		LogCallback: config.LogCallback,
	}

	// Llamar a la implementación de linux
	linuxProc, err := a.impl.StartOpenVPN(linuxConfig)
	if err != nil {
		return nil, err
	}

	// Convertir linux.Process a platform.Process
	return &Process{
		Cmd:     linuxProc.Cmd,
		Running: linuxProc.Running,
		PID:     linuxProc.PID,
	}, nil
}

func (a *linuxAdapter) StopOpenVPN(proc *Process) error {
	// Convertir platform.Process a linux.Process
	linuxProc := &linux.Process{
		Cmd:     proc.Cmd,
		Running: proc.Running,
		PID:     proc.PID,
	}

	err := a.impl.StopOpenVPN(linuxProc)

	// Actualizar el estado de vuelta
	proc.Running = linuxProc.Running

	return err
}

func (a *linuxAdapter) RequiresElevation() bool {
	return a.impl.RequiresElevation()
}

func (a *linuxAdapter) ElevateCommand(path string, args []string) (string, []string, error) {
	return a.impl.ElevateCommand(path, args)
}

func (a *linuxAdapter) GetConfigDir() string {
	return a.impl.GetConfigDir()
}

func (a *linuxAdapter) GetDefaultConfigPath() string {
	return a.impl.GetDefaultConfigPath()
}

func (a *linuxAdapter) GetLogPath() string {
	return a.impl.GetLogPath()
}

func (a *linuxAdapter) Name() string {
	return a.impl.Name()
}

func (a *linuxAdapter) Separator() string {
	return a.impl.Separator()
}

// NewWindows es un stub para Linux (no se usa)
func NewWindows() Platform {
	panic("Windows platform not available on Linux builds")
}

// NewDarwin es un stub para Linux (no se usa)
func NewDarwin() Platform {
	panic("Darwin platform not available on Linux builds")
}
