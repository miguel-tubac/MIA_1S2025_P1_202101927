package analyzer

import (
	"fmt"
	"strings"

	comandos "bakend/src/comandos"
)

// Analyzer analiza el comando de entrada y ejecuta la acción correspondiente
func Analyzer(input string) ([]interface{}, []error) {
	// Dividir el input en líneas
	lines := strings.Split(input, "\n")

	var results []interface{}
	var errors []error

	// Recorrer cada línea
	for _, line := range lines {
		// Ignorar líneas que comienzan con #
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		// Divide la línea en tokens usando espacios en blanco como delimitadores
		tokens := strings.Fields(line)

		// Si no se proporcionó ningún comando, continua con la siguiente línea
		if len(tokens) == 0 {
			continue
		}

		comando := strings.ToLower(tokens[0]) // Convertimos a minúsculas el comando
		// Switch para manejar diferentes comandos
		switch comando { // Toma la primera posición de la entrada
		case "mkdisk":
			//Aca se valida si ya se realizo un login
			//if comandos.ObtenerLogin() {
			result, err := comandos.ParseMkdisk(tokens[1:])
			results = append(results, result)
			//results = append(results, "\n")
			if err != nil {
				errors = append(errors, err)
			}
			//} else {
			// Si el comando no es reconocido, agregamos el error
			//errors = append(errors, fmt.Errorf("debe logearse para utilizar el comando \"mkdisk\": %s", tokens[0]))
			//}
		case "rmdisk":
			// if comandos.ObtenerLogin() {
			result, err := comandos.Eliminar_Disco(tokens[1:])
			results = append(results, result)
			//results = append(results, "\n")
			if err != nil {
				errors = append(errors, err)
			}
			// } else {
			// 	// Si el comando no es reconocido, agregamos el error
			// 	errors = append(errors, fmt.Errorf("debe logearse para utilizar el comando \"rmdisk\": %s", tokens[0]))
			// }
		case "fdisk":
			// if comandos.ObtenerLogin() {
			result, err := comandos.ParseFdisk(tokens[1:])
			results = append(results, result)
			//results = append(results, "\n")
			if err != nil {
				errors = append(errors, err)
			}
			// } else {
			// 	// Si el comando no es reconocido, agregamos el error
			// 	errors = append(errors, fmt.Errorf("debe logearse para utilizar el comando \"fdisk\": %s", tokens[0]))
			// }

		case "mount":
			// if comandos.ObtenerLogin() {
			// Llama a la función para el mount
			result, err := comandos.ParseMount(tokens[1:])
			results = append(results, result)
			//results = append(results, "\n")
			if err != nil {
				errors = append(errors, err)
			}
			// } else {
			// 	// Si el comando no es reconocido, agregamos el error
			// 	errors = append(errors, fmt.Errorf("debe logearse para utilizar el comando \"mount\": %s", tokens[0]))
			// }
		case "mounted":
			if comandos.ObtenerLogin() {
				// Llama a la función para el mounted
				result, err := comandos.MountedParser()
				results = append(results, result)
				//results = append(results, "\n")
				if err != nil {
					errors = append(errors, err)
				}
			} else {
				// Si el comando no es reconocido, agregamos el error
				errors = append(errors, fmt.Errorf("debe logearse para utilizar el comando \"mounted\": %s", tokens[0]))
			}
		case "mkfs":
			result, err := comandos.ParseMkfs(tokens[1:])
			results = append(results, result)
			//results = append(results, "\n")
			if err != nil {
				errors = append(errors, err)
			}
		case "login":
			if !comandos.ObtenerLogin() {
				result, err := comandos.ParseLogin(tokens[1:])
				results = append(results, result)
				//results = append(results, "\n")
				if err != nil {
					errors = append(errors, err)
				}
			} else {
				// Si el comando no es reconocido, agregamos el error
				errors = append(errors, fmt.Errorf("debe deslogearse para utilizar el comando \"login\": %s", tokens[0]))
			}
		case "rep":
			if comandos.ObtenerLogin() {
				result, err := comandos.ParseRep(tokens[1:])
				results = append(results, result)
				//results = append(results, "\n")
				if err != nil {
					errors = append(errors, err)
				}
			} else {
				// Si el comando no es reconocido, agregamos el error
				errors = append(errors, fmt.Errorf("debe logearse para utilizar el comando \"rep\": %s", tokens[0]))
			}
		case "logout":
			result, err := comandos.Logout(tokens[1:])
			results = append(results, result)
			if err != nil {
				errors = append(errors, err)
			}
		case "mkgrp":
			if comandos.ObtenerLogin() {
				result, err := comandos.ParseMkgrp(tokens[1:])
				results = append(results, result)
				//results = append(results, "\n")
				if err != nil {
					errors = append(errors, err)
				}
			} else {
				// Si el comando no es reconocido, agregamos el error
				errors = append(errors, fmt.Errorf("debe logearse para utilizar el comando \"mkgrp\": %s", tokens[0]))
			}
		case "rmgrp":
			if comandos.ObtenerLogin() {
				result, err := comandos.ParseRmgrp(tokens[1:])
				results = append(results, result)
				//results = append(results, "\n")
				if err != nil {
					errors = append(errors, err)
				}
			} else {
				// Si el comando no es reconocido, agregamos el error
				errors = append(errors, fmt.Errorf("debe logearse para utilizar el comando \"rmgrp\": %s", tokens[0]))
			}
		case "mkusr":
			if comandos.ObtenerLogin() {
				result, err := comandos.ParseMkusr(tokens[1:])
				results = append(results, result)
				//results = append(results, "\n")
				if err != nil {
					errors = append(errors, err)
				}
			} else {
				// Si el comando no es reconocido, agregamos el error
				errors = append(errors, fmt.Errorf("debe logearse para utilizar el comando \"mkusr\": %s", tokens[0]))
			}
		case "rmusr":
			if comandos.ObtenerLogin() {
				result, err := comandos.ParseRmusr(tokens[1:])
				results = append(results, result)
				//results = append(results, "\n")
				if err != nil {
					errors = append(errors, err)
				}
			} else {
				// Si el comando no es reconocido, agregamos el error
				errors = append(errors, fmt.Errorf("debe logearse para utilizar el comando \"rmusr\": %s", tokens[0]))
			}
		case "chgrp":
			if comandos.ObtenerLogin() {
				result, err := comandos.ParseChgrp(tokens[1:])
				results = append(results, result)
				//results = append(results, "\n")
				if err != nil {
					errors = append(errors, err)
				}
			} else {
				// Si el comando no es reconocido, agregamos el error
				errors = append(errors, fmt.Errorf("debe logearse para utilizar el comando \"chgrp\": %s", tokens[0]))
			}
		case "mkdir": //Este comando crea las carpetas es decir las rutas
			if comandos.ObtenerLogin() {
				result, err := comandos.ParseMkdir(tokens[1:])
				results = append(results, result)
				//results = append(results, "\n")
				if err != nil {
					errors = append(errors, err)
				}
			} else {
				// Si el comando no es reconocido, agregamos el error
				errors = append(errors, fmt.Errorf("debe logearse para utilizar el comando \"mkdir\": %s", tokens[0]))
			}
		case "mkfile": //Este comando crea las carpetas es decir las rutas
			if comandos.ObtenerLogin() {
				result, err := comandos.ParseMkfile(tokens[1:])
				results = append(results, result)
				//results = append(results, "\n")
				if err != nil {
					errors = append(errors, err)
				}
			} else {
				// Si el comando no es reconocido, agregamos el error
				errors = append(errors, fmt.Errorf("debe logearse para utilizar el comando \"mkfile\": %s", tokens[0]))
			}
		case "cat": //Este comando obyiene el texto del archivo
			if comandos.ObtenerLogin() {
				result, err := comandos.ParseCat(tokens[1:])
				results = append(results, result)
				//results = append(results, "\n")
				if err != nil {
					errors = append(errors, err)
				}
			} else {
				// Si el comando no es reconocido, agregamos el error
				errors = append(errors, fmt.Errorf("debe logearse para utilizar el comando \"cat\": %s", tokens[0]))
			}
		default:
			// Si el comando no es reconocido, agregamos el error
			errors = append(errors, fmt.Errorf("comando desconocido: %s", tokens[0]))
		}
	}

	// Retornamos los resultados y los errores acumulados
	return results, errors
}
