package analyzer

import (
	stores "bakend/src/almacenamiento"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type LOGIN struct {
	user string //Almacenara el usuario
	pass string // Alamcenara la contraseña
	id   string //Almacenara el id de la particion
}

/*
	Comandos aseptados:

	login -user=root -pass=123 -id=062A

	login -user="mi usuario" -pass="mi pwd" -id=062A
*/

// Commando para validar el login
func ParseLogin(tokens []string) (*LOGIN, error) {
	cmd := &LOGIN{} // Crea una nueva instancia de LOGIN

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando mount
	re := regexp.MustCompile(`-(?i:user="[^"]+"|user=[^\s]+|pass="[^"]+"|pass=[^\s]+|id=[^\s]+)`)
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
		case "-user":
			// Verifica que el path no esté vacío
			if value == "" {
				return nil, errors.New("el user no puede estar vacío")
			}
			cmd.user = value
		case "-pass":
			// Verifica que el nombre no esté vacío
			if value == "" {
				return nil, errors.New("la contraseña (pass) no puede estar vacío")
			}
			cmd.pass = value
		case "-id":
			// Verifica que el nombre no esté vacío
			if value == "" {
				return nil, errors.New("el id no puede estar vacío")
			}
			cmd.id = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, fmt.Errorf("parámetro desconocido en el login: %s", key)
		}
	}

	// Verifica que los parámetros -user, -pass y -id hayan sido proporcionados
	if cmd.user == "" {
		return nil, errors.New("faltan parámetros requeridos: -user")
	}
	if cmd.pass == "" {
		return nil, errors.New("faltan parámetros requeridos: -pass")
	}
	if cmd.id == "" {
		return nil, errors.New("faltan parámetros requeridos: -id")
	}

	// Montamos la partición
	err := commandLogear(cmd)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return cmd, fmt.Errorf("login realizado: %+v", *cmd) // Devuelve el comando MOUNT creado
}

// Fincion para validar el usuario logeado
func commandLogear(login *LOGIN) error {
	// Obtener la partición montada
	//Tipo de retorno: (particion) structures.PARTITION, (path Disco) string,(error por si algo salia mal) error
	mountedPartition, partitionPath, err := stores.GetMountedPartition(login.id)
	if err != nil {
		return err
	}

	//1. Obtener el superbloque a partir del start de la particion
	posision := mountedPartition.Part_start
	//2. Deserealizar para obtner el superbloque
	//3. OBtnener el primer inodo
	//4. Recorrer el bloque de carpetas hasta encontrar el bloque de archivos y deserealizar el contenido de user.txt
	//5. Comparar la contraseña y agregar una variable global
	//6. desde el analyzer obter la variable y validar si ya esta logeado
}
