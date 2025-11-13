//go:build darwin

package platform

import "github.com/lavp2393/navtunnel/internal/platform/darwin"

// darwinAdapter adapta darwin.DarwinPlatform a la interfaz Platform
type darwinAdapter struct {
	impl *darwin.DarwinPlatform
}

// NewDarwin retorna la implementación de Platform para macOS
func NewDarwin() Platform {
	return &darwinAdapter{
		impl: darwin.NewDarwin(),
	}
}

func (a *darwinAdapter) FindOpenVPN() (string, error) {
	return a.impl.FindOpenVPN()
}

func (a *darwinAdapter) StartOpenVPN(config StartConfig) (*Process, error) {
	darwinConfig := darwin.StartConfig{
		ConfigPath:  config.ConfigPath,
		MgmtPort:    config.MgmtPort,
		LogCallback: config.LogCallback,
	}

	darwinProc, err := a.impl.StartOpenVPN(darwinConfig)
	if err != nil {
		return nil, err
	}

	return &Process{
		Cmd:     darwinProc.Cmd,
		Running: darwinProc.Running,
		PID:     darwinProc.PID,
	}, nil
}

func (a *darwinAdapter) StopOpenVPN(proc *Process) error {
	darwinProc := &darwin.Process{
		Cmd:     proc.Cmd,
		Running: proc.Running,
		PID:     proc.PID,
	}

	err := a.impl.StopOpenVPN(darwinProc)
	proc.Running = darwinProc.Running
	return err
}

func (a *darwinAdapter) RequiresElevation() bool {
	return a.impl.RequiresElevation()
}

func (a *darwinAdapter) ElevateCommand(path string, args []string) (string, []string, error) {
	return a.impl.ElevateCommand(path, args)
}

func (a *darwinAdapter) GetConfigDir() string {
	return a.impl.GetConfigDir()
}

func (a *darwinAdapter) GetDefaultConfigPath() string {
	return a.impl.GetDefaultConfigPath()
}

func (a *darwinAdapter) GetLogPath() string {
	return a.impl.GetLogPath()
}

func (a *darwinAdapter) Name() string {
	return a.impl.Name()
}

func (a *darwinAdapter) Separator() string {
	return a.impl.Separator()
}

// NewLinux es un stub para Darwin (no se usa)
func NewLinux() Platform {
	panic("Linux platform not available on Darwin builds")
}

// NewWindows es un stub para Darwin (no se usa)
func NewWindows() Platform {
	panic("Windows platform not available on Darwin builds")
}
