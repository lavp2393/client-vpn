package core

import (
	"fmt"
	"os"

	"github.com/prey/preyvpn/internal/platform"
)

// OpenVPNProcess encapsula el proceso de OpenVPN con abstracción de plataforma
type OpenVPNProcess struct {
	proc     *platform.Process
	platform platform.Platform
}

// FindOpenVPN busca el ejecutable de OpenVPN usando la abstracción de plataforma
func FindOpenVPN() (string, error) {
	plat := platform.New()
	return plat.FindOpenVPN()
}

// GetConfigPath retorna la ruta esperada del archivo de configuración
func GetConfigPath() string {
	plat := platform.New()
	return plat.GetDefaultConfigPath()
}

// CheckConfigExists verifica si existe el archivo de configuración
func CheckConfigExists() bool {
	path := GetConfigPath()
	_, err := os.Stat(path)
	return err == nil
}

// StartOpenVPN inicia el proceso de OpenVPN con elevación de privilegios usando la abstracción de plataforma
func StartOpenVPN(configPath string, mgmtPort int, logCallback func(string)) (*OpenVPNProcess, error) {
	// Obtener la plataforma actual
	plat := platform.New()

	// Configurar el inicio de OpenVPN
	config := platform.StartConfig{
		ConfigPath:  configPath,
		MgmtPort:    mgmtPort,
		LogCallback: logCallback,
	}

	// Iniciar OpenVPN usando la abstracción de plataforma
	proc, err := plat.StartOpenVPN(config)
	if err != nil {
		return nil, fmt.Errorf("error al iniciar OpenVPN: %w", err)
	}

	return &OpenVPNProcess{
		proc:     proc,
		platform: plat,
	}, nil
}

// Stop detiene el proceso de OpenVPN limpiamente usando la abstracción de plataforma
func (p *OpenVPNProcess) Stop() error {
	if p.proc == nil {
		return nil
	}
	return p.platform.StopOpenVPN(p.proc)
}

// IsRunning retorna si el proceso está corriendo
func (p *OpenVPNProcess) IsRunning() bool {
	if p.proc == nil {
		return false
	}
	return p.proc.Running
}

// PID retorna el PID del proceso
func (p *OpenVPNProcess) PID() int {
	if p.proc == nil {
		return 0
	}
	return p.proc.PID
}
