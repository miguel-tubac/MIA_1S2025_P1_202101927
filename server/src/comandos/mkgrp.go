package analyzer

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type MKGRP struct {
	name string
}

func ParseMkgrp(tokens []string) (*MKGRP, error) {
	cmd := &MKGRP{}

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando mount
	re := regexp.MustCompile(`-(?i:name="[^"]+"|name=[^\s]+)`)
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

		// Remove comillas si estan present
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Switch para manejar diferentes parámetros
		switch key {
		case "-name":
			// Verifica que el path no esté vacío
			if value == "" {
				return nil, errors.New("el name no puede estar vacío")
			}
			cmd.name = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, fmt.Errorf("parámetro desconocido en el mkgrp: %s", key)
		}
	}

	// Verifica que los parámetros -user, -pass y -id hayan sido proporcionados
	if cmd.name == "" {
		return nil, errors.New("faltan parámetros requeridos: -name")
	}

	// Montamos la partición
	err := commandMkgrp(cmd)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return cmd, fmt.Errorf("grupo de usuarios creado: %+v", *cmd)
}

func commandMkgrp(comando *MKGRP) error {

}
