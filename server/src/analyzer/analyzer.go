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
			result, err := comandos.ParseMkdisk(tokens[1:])
			results = append(results, result)
			//results = append(results, "\n")
			if err != nil {
				errors = append(errors, err)
			}
		case "rmdisk":
			result, err := comandos.Eliminar_Disco(tokens[1:])
			results = append(results, result)
			//results = append(results, "\n")
			if err != nil {
				errors = append(errors, err)
			}
		case "fdisk":
			result, err := comandos.ParseFdisk(tokens[1:])
			results = append(results, result)
			//results = append(results, "\n")
			if err != nil {
				errors = append(errors, err)
			}
		case "mount":
			// Llama a la función para el mount
			result, err := comandos.ParseMount(tokens[1:])
			results = append(results, result)
			//results = append(results, "\n")
			if err != nil {
				errors = append(errors, err)
			}
		case "mounted":
			// Llama a la función para el mounted
			result, err := comandos.MountedParser()
			results = append(results, result)
			//results = append(results, "\n")
			if err != nil {
				errors = append(errors, err)
			}
		case "mkfs":
			result, err := comandos.ParseMkfs(tokens[1:])
			results = append(results, result)
			//results = append(results, "\n")
			if err != nil {
				errors = append(errors, err)
			}
		default:
			// Si el comando no es reconocido, agregamos el error
			errors = append(errors, fmt.Errorf("comando desconocido: %s", tokens[0]))
		}
	}

	// Retornamos los resultados y los errores acumulados
	return results, errors
}
