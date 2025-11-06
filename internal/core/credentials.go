package core

import (
	"encoding/json"

	"github.com/zalando/go-keyring"
)

const (
	// Nombre del servicio para keyring (identifica la app en gnome-keyring/KDE Wallet)
	serviceName = "PreyVPN"
	// Username para el keyring (usado como identificador único)
	keyringUser = "credentials"
)

// SavedCredentials contiene las credenciales que se guardan
type SavedCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SaveCredentials guarda las credenciales en el keyring del sistema
// En Linux usa Secret Service API (gnome-keyring o KDE Wallet)
func SaveCredentials(username, password string) error {
	creds := SavedCredentials{
		Username: username,
		Password: password,
	}

	// Serializar a JSON
	data, err := json.Marshal(creds)
	if err != nil {
		return err
	}

	// Guardar en keyring
	// service: "PreyVPN", user: "credentials", password: <json con las credenciales>
	return keyring.Set(serviceName, keyringUser, string(data))
}

// LoadCredentials carga las credenciales desde el keyring del sistema
// Retorna username, password, error
func LoadCredentials() (string, string, error) {
	// Recuperar del keyring
	data, err := keyring.Get(serviceName, keyringUser)
	if err != nil {
		// Si no existen credenciales guardadas, keyring retorna error
		// Esto es normal en el primer uso
		return "", "", err
	}

	// Deserializar
	var creds SavedCredentials
	if err := json.Unmarshal([]byte(data), &creds); err != nil {
		return "", "", err
	}

	return creds.Username, creds.Password, nil
}

// DeleteCredentials elimina las credenciales del keyring
func DeleteCredentials() error {
	return keyring.Delete(serviceName, keyringUser)
}

// HasCredentials verifica si existen credenciales guardadas
func HasCredentials() bool {
	_, _, err := LoadCredentials()
	return err == nil
}
