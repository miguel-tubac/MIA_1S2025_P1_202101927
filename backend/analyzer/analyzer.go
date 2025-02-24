package analyzer

import (
	"errors"
	"fmt"
	"strings"

	comandos "bakend/comandos"
)

// Analyzer analiza el comando de entrada y ejecuta la acción correspondiente
func Analyzer(input string) (interface{}, error) {
	// Divide la entrada en tokens usando espacios en blanco como delimitadores
	tokens := strings.Fields(input) //Solo crea un areglo de las palabras ["hola", "mundo"]

	// Si no se proporcionó ningún comando, devuelve un error
	if len(tokens) == 0 {
		return nil, errors.New("no se proporcionó ningún comando")
	}

	comando := strings.ToLower(tokens[0]) //Convertimos a minisculas el comando
	// Switch para manejar diferentes comandos
	switch comando { //Toma la primera posicion de la entrada
	case "mkdisk":
		// Llama a la función Mkdir del paquete commands con los argumentos restantes
		return comandos.ParseMkdisk(tokens[1:])
	case "rmdisk":
		// Llama a la función de eliminar disco
		return comandos.Eliminar_Disco(tokens[1:])
	case "fdisk":
		return comandos.ParseFdisk(tokens[1:])
	default:
		// Si el comando no es reconocido, devuelve un error
		return nil, fmt.Errorf("comando desconocido: %s", tokens[0])
	}
}
