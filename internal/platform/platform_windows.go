//go:build windows

package platform

import "github.com/prey/preyvpn/internal/platform/windows"

// windowsAdapter adapta windows.WindowsPlatform a la interfaz Platform
type windowsAdapter struct {
	impl *windows.WindowsPlatform
}

// NewWindows retorna la implementación de Platform para Windows
func NewWindows() Platform {
	return &windowsAdapter{
		impl: windows.NewWindows(),
	}
}

func (a *windowsAdapter) FindOpenVPN() (string, error) {
	return a.impl.FindOpenVPN()
}

func (a *windowsAdapter) StartOpenVPN(config StartConfig) (*Process, error) {
	winConfig := windows.StartConfig{
		ConfigPath:  config.ConfigPath,
		MgmtPort:    config.MgmtPort,
		LogCallback: config.LogCallback,
	}

	winProc, err := a.impl.StartOpenVPN(winConfig)
	if err != nil {
		return nil, err
	}

	return &Process{
		Cmd:     winProc.Cmd,
		Running: winProc.Running,
		PID:     winProc.PID,
	}, nil
}

func (a *windowsAdapter) StopOpenVPN(proc *Process) error {
	winProc := &windows.Process{
		Cmd:     proc.Cmd,
		Running: proc.Running,
		PID:     proc.PID,
	}

	err := a.impl.StopOpenVPN(winProc)
	proc.Running = winProc.Running
	return err
}

func (a *windowsAdapter) RequiresElevation() bool {
	return a.impl.RequiresElevation()
}

func (a *windowsAdapter) ElevateCommand(path string, args []string) (string, []string, error) {
	return a.impl.ElevateCommand(path, args)
}

func (a *windowsAdapter) GetConfigDir() string {
	return a.impl.GetConfigDir()
}

func (a *windowsAdapter) GetDefaultConfigPath() string {
	return a.impl.GetDefaultConfigPath()
}

func (a *windowsAdapter) GetLogPath() string {
	return a.impl.GetLogPath()
}

func (a *windowsAdapter) Name() string {
	return a.impl.Name()
}

func (a *windowsAdapter) Separator() string {
	return a.impl.Separator()
}

// NewLinux es un stub para Windows (no se usa)
func NewLinux() Platform {
	panic("Linux platform not available on Windows builds")
}

// NewDarwin es un stub para Windows (no se usa)
func NewDarwin() Platform {
	panic("Darwin platform not available on Windows builds")
}
