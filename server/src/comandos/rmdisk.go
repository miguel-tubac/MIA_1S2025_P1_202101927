package analyzer

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type RMDISK struct {
	path string
}

/*
Comandos aseptados:
rmdisk -path="/home/miguel/Descargas/Archivos/Laboratorio/Proyecto1/backend/discos/Disco4.mia"

Sin comillas:
rmdisk -path=/home/miguel/Descargas/Archivos/Laboratorio/Proyecto1/backend/discos/Disco1.mia
*/

func Eliminar_Disco(tokens []string) (*RMDISK, error) {
	cmd := &RMDISK{}
	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando mkdisk
	re := regexp.MustCompile(`-(?i:path="[^"]+"|path=[^\s]+)`)
	// Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
	matches := re.FindAllString(args, -1)

	// Itera sobre cada coincidencia encontrada
	for _, match := range matches {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		// Remove comillas si estan presente
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		if key == "-path" {
			// Verifica que el path no esté vacío
			if value == "" {
				return nil, errors.New("el path no puede estar vacío")
			}

			// Intenta eliminar el archivo
			if err := os.Remove(value); err != nil {
				return nil, fmt.Errorf("error al eliminar el archivo: %v", err)
			}
			//fmt.Println("Disco eliminado correctamente")
			cmd.path = value
			return cmd, fmt.Errorf("disco eliminado: %+v", *cmd)
		}
	}
	return nil, errors.New("no se especificó un path válido")
}
