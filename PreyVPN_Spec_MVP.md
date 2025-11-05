# PreyVPN Wrapper — Especificación MVP (Ubuntu, sin .deb)

## 1) Objetivo
Binario **Go** con GUI mínima que:
- Lanza **OpenVPN** usando un perfil fijo.
- Gestiona prompts **usuario → contraseña → OTP** vía **Management Interface**.
- Permite **Conectar/Desconectar** sin terminal.
- **No** cambia nada del backend (OpenVPN + PAM/LDAP + LinOTP).

---

## 2) Supuestos del sistema
- Sistema: **Ubuntu Desktop**.
- `openvpn` instalado y disponible en `/usr/sbin/openvpn`, `/usr/bin/openvpn` o `$PATH`.
- `pkexec` disponible para elevación puntual.
- Perfil `.ovpn` provisto por la organización.

---

## 3) Perfil VPN (ubicación fija)
- Carpeta: `~/PreyVPN/`
- Archivo: **`prey-prod.ovpn`**
- Ruta esperada: `~/PreyVPN/prey-prod.ovpn`

**Regla:** El binario busca **solo** esa ruta.  
Si no existe, muestra una pantalla con instrucción y botón **Reintentar**.

---

## 4) Comportamiento del binario

### 4.1 Ciclo básico
1. **Inicio**
   - Verifica existencia de `~/PreyVPN/prey-prod.ovpn`.
   - Si no está: pantalla “perfil no encontrado” + **Reintentar**.
   - Si está: habilita **Conectar**.

2. **Conectar**
   - Selecciona **puerto de management** libre (ej. 49152–65535).
   - Lanza `openvpn` con **elevación** usando `pkexec`:
     ```bash
     openvpn --config ~/PreyVPN/prey-prod.ovpn        --management 127.0.0.1:<PORT> stdin        --auth-retry interact        --auth-nocache
     ```
   - Abre socket TCP a `127.0.0.1:<PORT>` y comienza a **parsear eventos**.

3. **Autenticación (prompts)**
   - Prompt 1 (usuario): `>PASSWORD:Need 'Auth' username`
   - Prompt 2 (contraseña): `>PASSWORD:Need 'Auth' pass`
   - Prompt 3 (OTP): cualquier `>PASSWORD:Need ...` posterior **o** presencia de `CHALLENGE/CRV1`
   - Respuestas (Management):
     - `username "Auth" <valor>`
     - `password "Auth" <valor>`
     - OTP con el **mismo tag** que solicite (normalmente `"Auth"`).

4. **Estados**
   - Mostrar **Conectando → Autenticando → Conectado**.
   - Conectado si aparece `>STATE:*,CONNECTED,SUCCESS,`.
   - Fallos:
     - `>STATE:*,AUTH_FAILED,` → mapear según etapa (pass/OTP).
     - `>FATAL:` → error de conexión.

5. **Desconectar**
   - Enviar señal para terminar el proceso `openvpn`.
   - Cerrar socket de management y volver a estado inicial.

---

## 5) UI mínima (sin imponer toolkit)
- **Ventana principal**
  - Estado: “Perfil detectado / Perfil no encontrado”
  - Botones: **Conectar / Desconectar**, **Reintentar** (solo si falta el perfil)
  - Log en vivo (últimas ~30 líneas; solo lectura)
- **Modales**
  - Usuario (placeholder: “usuario corporativo”)
  - Contraseña
  - OTP (6 dígitos; hint: “se renueva cada 30s”)
- **Mensajes**
  - Perfil no encontrado: `No encuentro ~/PreyVPN/prey-prod.ovpn. Coloca el archivo allí y presiona Reintentar.`
  - Contraseña incorrecta: `Contraseña incorrecta.`
  - OTP inválido/expirado: `OTP inválido o expirado.`
  - Conectado: `Conexión establecida ✅`

---

## 6) Seguridad
- **Obligatorio:** `--auth-nocache`.
- No persistir **contraseñas** ni **OTP**.
- (Opcional post-MVP) Recordar **solo** el usuario vía keyring del sistema.
- Logs sin secretos (no imprimir credenciales ni OTP).

---

## 7) Estructura de proyecto (sugerida)

```
/cmd/preyvpn/main.go
/internal/core/openvpn.go        // spawn/kill de proceso con pkexec; resolución de ruta openvpn
/internal/core/manager.go        // socket mgmt + parser + FSM (estados y eventos)
/internal/ui/app.go              // ventana principal, estados, log view
/internal/ui/prompts.go          // modales user/pass/otp
/internal/logs/buffer.go         // buffer de log (rotación en memoria)
```

**Contratos recomendados**
- `core.Start(configPath string, mgmtPort int) (events <-chan Event, send SendFns, stop func(), err error)`
- `type Event = AskUser | AskPass | AskOTP | Connected | AuthFailed{stage} | Fatal{reason} | LogLine{text}`
- `type SendFns struct { Username(v string); Password(v string); OTP(v string) }`

---

## 8) Parsing de Management — patrones mínimos

### 8.1 Prompts (ejemplos reales orientativos)
```
>PASSWORD:Need 'Auth' username
>PASSWORD:Need 'Auth' pass
>PASSWORD:Need 'Auth' OTP
>INFO: CRV1:... <challenge-string> ...
```
**Regla simple:**  
- Si contiene `Need 'Auth'` y `username` → **AskUser**.  
- Si contiene `Need 'Auth'` y `pass` → **AskPass**.  
- Si aparece otro `Need` posterior **o** `CRV1`/`CHALLENGE` → **AskOTP**.

### 8.2 Éxito / fallo / fatal
```
>STATE:1730165123,CONNECTED,SUCCESS,10.8.0.10,xx.xx.xx.xx,,
>STATE:1730164999,AUTH_FAILED,,
>FATAL:Something bad happened
```

### 8.3 Envío de credenciales (formato)
```
username "Auth" myuser
password "Auth" mypassword
username "Auth" 123456        // si el servidor pide OTP como 'username' adicional, seguir el tag pedido
password "Auth" 123456        // o como 'password' adicional, según prompt
```
> **Nota:** usa exactamente el **tag** indicado en el prompt (normalmente `"Auth"`).

---

## 9) Elevación de privilegios
- Ejecutar OpenVPN con `pkexec` (GUI del sistema pedirá la contraseña si aplica).
- El binario debe:
  - Resolver ruta de `openvpn`.
  - Construir los argumentos.
  - Capturar PID del proceso hijo para **Desconectar** limpiamente.

---

## 10) Mapeo de errores (UX)
| Señal/Evento                         | Mensaje UI                    | Acción |
|-------------------------------------|-------------------------------|--------|
| `>STATE:*,AUTH_FAILED,` tras pass   | Contraseña incorrecta         | Re-pedir **solo** contraseña |
| `>STATE:*,AUTH_FAILED,` tras OTP    | OTP inválido o expirado       | Re-pedir **solo** OTP |
| Repetidos AUTH_FAILED al OTP        | Revisa la hora de tu equipo   | Mostrar hint de NTP |
| `>FATAL:` / timeouts                 | Error de conexión             | Permitir Reintentar |
| Falta `openvpn`                      | OpenVPN no está instalado     | Mostrar instrucción clara |

---

## 11) Criterios de aceptación (QA)
1. Con `~/PreyVPN/prey-prod.ovpn` presente:
   - **Conectar** → aparecen **3 prompts** (usuario → contraseña → OTP) y termina en **Conectado**.
2. Error de contraseña:
   - Muestra mensaje y re-pide **solo** contraseña.
3. Error de OTP:
   - Muestra mensaje y re-pide **solo** OTP.
4. **Desconectar**:
   - Mata el proceso OpenVPN y vuelve a estado inicial sin residuos.
5. `--auth-nocache`:
   - Confirmado en las líneas de arranque del log.
6. Logs:
   - Sin secretos; visor en UI muestra ~30 últimas líneas.

---

## 12) Backlog (fuera de este MVP)
- Recordar usuario (keyring).
- Soporte de múltiples perfiles en `~/PreyVPN/`.
- Auto-reconexión con backoff.
- Regla polkit por grupo (sin prompt).
- Builds Windows/macOS.

---

## 13) Notas de implementación (prácticas)
- **Selección de puerto management:** intenta N aleatorios en rango 49152–65535 hasta éxito.
- **Lectura de management:** línea-a-línea; no bloqueante; emitir `LogLine` para todo.
- **Sanitización de logs:** nunca imprimir valores enviados en `username/password`.
- **Validación de comandos:** escapar/quote argumentos al invocar `pkexec` para evitar inyección.
- **Cierre limpio:** al desconectar, enviar SIGTERM al hijo y esperar; si no termina, SIGKILL con timeout.

---

**Fin del documento.**
