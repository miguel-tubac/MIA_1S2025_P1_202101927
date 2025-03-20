package analyzer

import (
	stores "bakend/src/almacenamiento"
	structures "bakend/src/estructuras"

	//utils "bakend/src/utils"
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

var logeado = false

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

	return cmd, fmt.Errorf("login realizado: %+v", *cmd) // Devuelve el comando LOGIN creado
}

// Fincion para validar el usuario logeado
func commandLogear(login *LOGIN) error {
	// Obtener la partición montada
	//Tipo de retorno: (*structures.SuperBlock, *structures.PARTITION, string, error)
	partitionSuperblock, _, partitionPath, err := stores.GetMountedPartitionSuperblock(login.id)
	if err != nil {
		return fmt.Errorf("error al obtener la partición montada: %w", err)
	}

	//Aca iniciamos desde el inodo numero 1
	err2 := Login(partitionPath, login, 1, partitionSuperblock)

	//TODO: validar la salida
	if err2 != nil {
		return fmt.Errorf("error al obtener el suuario y contraseña: %w", err2)
	}

	return nil
}

// Funcion para accder al archivo de user.txt
// Login: path del disco, objeto con los datos del usuario, el inicio de los inodos
func Login(path string, login *LOGIN, inodeIndex int32, sb *structures.SuperBlock) error {
	//Se crea una instancia de un objeto de tipo Inode
	inode := &structures.Inode{}

	// Deserializar el inodo
	err := inode.Deserialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
	if err != nil {
		return err
	}

	//Aca se almacenara el contenido de user.txt
	data := ""
	// Iterar sobre cada bloque del inodo (apuntadores)
	for _, blockIndex := range inode.I_block {
		// Si el bloque no existe, salir
		if blockIndex == -1 {
			break
		}

		// Crear un nuevo bloque de archivo
		block := &structures.FileBlock{}

		// Deserializar el bloque desde el incio de los blokes  + posicion por el peso de los bloques  que es 64
		err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size))) // 64 porque es el tamaño de un bloque
		if err != nil {
			return err
		}

		//Obtengo el texto del archivo uset.txt
		data = strings.Trim(string(block.B_content[:]), "\x00 ")

	}

	/*
		1,G,root
		1,U,root,123
	*/

	// Dividir por salto de línea
	lines := strings.Split(data, "\n")
	// Recorrer cada línea y dividir por comas
	for _, line := range lines {
		values := strings.Split(line, ",")

		// Almacenar en variables según la cantidad de datos
		//Esto son los grupos
		if len(values) == 3 {
			//id, tipo, nombre := values[0], values[1], values[2]
			//fmt.Printf("ID: %s, Tipo: %s, Nombre: %s\n", id, tipo, nombre)
		} else if len(values) == 4 {
			//Estos son los usuarios
			_, _, nombre, extra := values[0], values[1], values[2], values[3]
			if nombre == login.user && extra == login.pass {
				logeado = true
				//fmt.Println("Logeado")
			} else {
				logeado = false
				return fmt.Errorf("error con el suario: %s ó contraseña: %s", nombre, extra)
			}

			//fmt.Printf("ID: %s, Tipo: %s, Nombre: %s, Extra: %s\n", id, tipo, nombre, extra)
		}
	}

	return nil
}

func ObtenerLogin() bool {
	return logeado
}
