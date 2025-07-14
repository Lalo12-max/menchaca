package utils


type ResponseCode struct {
	StatusCode int         `json:"statuscode"`
	IntCode    string      `json:"intcode"`
	Data       interface{} `json:"data"`
}


const (
	
	LOGIN_SUCCESS      = "S01"
	LOGIN_MFA_REQUIRED = "S02"

	LOGIN_INVALID_CREDENTIALS = "E01"
	LOGIN_INVALID_MFA         = "E02"
	LOGIN_PARSE_ERROR         = "E03"
	LOGIN_SERVER_ERROR        = "E04"
	LOGIN_TOKEN_ERROR         = "E05"

	LOGIN_ACCOUNT_LOCKED   = "W01"
	LOGIN_PASSWORD_EXPIRED = "W02"
)

const (
	REGISTER_SUCCESS = "S03"

	REGISTER_PARSE_ERROR   = "E06" 
	REGISTER_WEAK_PASSWORD = "E07" 
	REGISTER_HASH_ERROR    = "E08" 
	REGISTER_MFA_ERROR     = "E09"
	REGISTER_QR_ERROR      = "E10"
	REGISTER_DB_ERROR      = "E11"

	REGISTER_EMAIL_EXISTS = "W03"
)

func NewResponse(statusCode int, intCode string, data interface{}) ResponseCode {
	return ResponseCode{
		StatusCode: statusCode,
		IntCode:    intCode,
		Data:       data,
	}
}

func GetCodeDescription(code string) string {
	descriptions := map[string]string{
		LOGIN_SUCCESS:             "Login exitoso",
		LOGIN_MFA_REQUIRED:        "MFA requerido",
		LOGIN_INVALID_CREDENTIALS: "Credenciales inválidas",
		LOGIN_INVALID_MFA:         "Código MFA inválido",
		LOGIN_PARSE_ERROR:         "Error al parsear datos de login",
		LOGIN_SERVER_ERROR:        "Error del servidor en login",
		LOGIN_TOKEN_ERROR:         "Error al generar token",
		LOGIN_ACCOUNT_LOCKED:      "Cuenta bloqueada",
		LOGIN_PASSWORD_EXPIRED:    "Contraseña expirada",

		REGISTER_SUCCESS:       "Registro exitoso",
		REGISTER_PARSE_ERROR:   "Error al parsear datos de registro",
		REGISTER_WEAK_PASSWORD: "Contraseña débil",
		REGISTER_HASH_ERROR:    "Error al procesar contraseña",
		REGISTER_MFA_ERROR:     "Error al generar MFA",
		REGISTER_QR_ERROR:      "Error al generar código QR",
		REGISTER_DB_ERROR:      "Error al crear usuario",
		REGISTER_EMAIL_EXISTS:  "Email ya existe",
	}

	if desc, exists := descriptions[code]; exists {
		return desc
	}
	return "Código desconocido"
}
