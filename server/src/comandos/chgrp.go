package analyzer

import (
	stores "bakend/src/almacenamiento"
	structures "bakend/src/estructuras"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type CHGRP struct {
	user string
	grp  string
}

func ParseChgrp(tokens []string) (*CHGRP, error) {
	cmd := &CHGRP{}

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando mount
	re := regexp.MustCompile(`-(?i:user="[^"]+"|user=[^\s]+|grp="[^"]+"|grp=[^\s]+)`)
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
		case "-grp":
			// Verifica que el path no esté vacío
			if value == "" {
				return nil, errors.New("el grp no puede estar vacío")
			}
			cmd.grp = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, fmt.Errorf("parámetro desconocido en el chgrp: %s", key)
		}
	}

	// Verifica que los parámetros -user, -pass y -id hayan sido proporcionados
	if cmd.user == "" {
		return nil, errors.New("faltan parámetros requeridos: -user")
	}
	if cmd.grp == "" {
		return nil, errors.New("faltan parámetros requeridos: -grp")
	}

	// Agregamos al usuario
	err := commandChgrp(cmd)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return cmd, fmt.Errorf("cambio de grupo realizado: %+v", *cmd)
}

func commandChgrp(comando *CHGRP) error {
	//Obtenemos el usuario logeado
	var usuario = ObtenerUsuari()

	//Valida si el usuario es el usuario root
	if usuario.user != "root" {
		return errors.New("para cambiar de grupo debe se estar logeado como root")
	}

	// Obtener la partición montada
	//Tipo de retorno: (*structures.SuperBlock, *structures.PARTITION, string, error)
	partitionSuperblock, _, partitionPath, err := stores.GetMountedPartitionSuperblock(usuario.id)
	if err != nil {
		return fmt.Errorf("error al obtener el Superbloque: %w", err)
	}

	//Aca iniciamos desde el inodo numero 1
	err2 := ChgrpComand(partitionPath, usuario, comando, 1, partitionSuperblock)

	//validar la salida
	if err2 != nil {
		return fmt.Errorf("error al intenter escribir en el user.txt: %w", err2)
	}

	return nil
}

// Funcion para accder al archivo de user.txt
// CrearUser: path del disco, objeto con los datos del usuario, el inicio de los inodos
func ChgrpComand(path string, login *LOGIN, comando *CHGRP, inodeIndex int32, sb *structures.SuperBlock) error {
	//Se crea una instancia de un objeto de tipo Inode
	inode := &structures.Inode{}

	// Deserializar el inodo
	err := inode.Deserialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
	if err != nil {
		return err
	}

	//Aca se almacenara el contenido de user.txt
	data := ""
	// Crear un nuevo bloque de archivo
	block := &structures.FileBlock{}
	indiceFinal := 0

	encontrado := false
	encontradoGrupo := false
	// Iterar sobre cada bloque del inodo (apuntadores)
	for _, blockIndex := range inode.I_block {
		// Si el bloque no existe, salir
		if blockIndex == -1 {
			//Esto quiere decir que no encontro al grupo
			return fmt.Errorf("no se encontro ningun usuario con el user: %s o el grupo: %s", comando.user, comando.grp)
		}
		indiceFinal = int(blockIndex)
		// Crear un nuevo bloque de archivo
		block = &structures.FileBlock{}

		// Deserializar el bloque desde el incio de los blokes  + posicion por el peso de los bloques  que es 64
		err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size))) // 64 porque es el tamaño de un bloque
		if err != nil {
			return err
		}

		//Obtengo el texto del archivo uset.txt, de un fileblock
		data = strings.Trim(string(block.B_content[:]), "\x00 ")
		// fmt.Println("*******************")
		// fmt.Println(data)
		// fmt.Println("*******************")
		/*
			1,G,root
			1,U,root,root,123
		*/
		// Dividir por salto de línea
		lines := strings.Split(data, "\n")
		// Recorrer cada línea y dividir por comas
		var id = ""
		var tipo = ""
		var usuario = ""
		//var grupo = ""
		var password = ""

		//Esta variable almacenara el resultado
		result := ""

		//Esto es para obtner los datos puntuales
		for _, line := range lines {
			values := strings.Split(line, ",")

			// Almacenar en variables según la cantidad de datos
			//Esto son los grupos
			if len(values) == 5 {
				//Estos son los usuarios
				id, tipo, _, usuario, password = values[0], values[1], values[2], values[3], values[4]
				if usuario == comando.user {
					//Esto quiere decir que ya esta borrado
					if id == "0" {
						return fmt.Errorf("el usuario con el name: %s ya fue borrado", comando.user)
					}
					encontrado = true
					//Se edita el id del grupo
					result += id + "," + tipo + "," + comando.grp + "," + usuario + "," + password //+ "\n"
					//Continua por si hay mas datos y se agregan al contenido del fileblock
					continue
				}
				//fmt.Printf("ID: %s, Tipo: %s, Nombre: %s\n", id, tipo, nombre)
			} else if len(values) == 3 {
				//Aca es para validar si el grupo existe ya que se cambio de grupo
				id2, _, nombre2 := values[0], values[1], values[2]
				if nombre2 == comando.grp {
					//Esto quiere decir que el grupo esta elimando por lo tanto guradomos el correcto
					if id2 == "0" {
						return fmt.Errorf("el grupo con el name: %s ya fue borrado", comando.grp)
					}
					encontradoGrupo = true
				}
			}
			result += line + "\n"
		}

		if encontrado && encontradoGrupo {
			//Aca eliminamos caracteres nulos si es que existen
			//fmt.Println(result)
			result = strings.Trim(string(result), "\x00 ")
			// Limpiar el array antes de copiar
			block.B_content = [64]byte{}
			copy(block.B_content[:], []byte(result))
			//Se serealiza todo el contenido en el Fileblock
			err2 := block.Serialize(path, int64(sb.S_block_start+(int32(indiceFinal)*sb.S_block_size)))
			if err2 != nil {
				return err2
			}
			// fmt.Println("******cambio grupo****")
			// block.Print()
			//Finalizamos el bucle
			break
		}

	}

	//block.Print()

	return nil
}
