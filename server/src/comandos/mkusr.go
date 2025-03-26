package analyzer

import (
	stores "bakend/src/almacenamiento"
	structures "bakend/src/estructuras"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type MKUSR struct {
	user string
	pass string
	grp  string
}

func ParseMkusr(tokens []string) (*MKUSR, error) {
	cmd := &MKUSR{}

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando mount
	re := regexp.MustCompile(`-(?i:user="[^"]+"|user=[^\s]+|pass="[^"]+"|pass=[^\s]+|grp="[^"]+"|grp=[^\s]+)`)
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
			// Verifica que el path no esté vacío
			if value == "" {
				return nil, errors.New("el pass no puede estar vacío")
			}
			cmd.pass = value
		case "-grp":
			// Verifica que el path no esté vacío
			if value == "" {
				return nil, errors.New("el grp no puede estar vacío")
			}
			cmd.grp = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, fmt.Errorf("parámetro desconocido en el mkusr: %s", key)
		}
	}

	// Verifica que los parámetros -user, -pass y -id hayan sido proporcionados
	if cmd.user == "" {
		return nil, errors.New("faltan parámetros requeridos: -user")
	}
	if cmd.pass == "" {
		return nil, errors.New("faltan parámetros requeridos: -pass")
	}
	if cmd.grp == "" {
		return nil, errors.New("faltan parámetros requeridos: -grp")
	}

	// Agregamos al usuario
	err := commandMkusr(cmd)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return cmd, fmt.Errorf("usuario creado exitosamente: %+v", *cmd)
}

// Para respaldar el identificador
var respaldoId = ""
var tipo = ""

// Esto es para obtener el superbloque
func commandMkusr(comando *MKUSR) error {
	respaldoId = ""
	tipo = ""
	//Obtenemos el usuario logeado
	var usuario = ObtenerUsuari()

	//Valida si el usuario es el usuario root
	if usuario.user != "root" {
		return errors.New("para crear un usuario debe se estar logeado como root")
	}

	// Obtener la partición montada
	//Tipo de retorno: (*structures.SuperBlock, *structures.PARTITION, string, error)
	partitionSuperblock, particion, partitionPath, err := stores.GetMountedPartitionSuperblock(usuario.id)
	if err != nil {
		return fmt.Errorf("error al obtener el Superbloque: %w", err)
	}

	//Aca iniciamos desde el inodo numero 1
	err2 := MkusrComand(partitionPath, usuario, comando, 1, partitionSuperblock, particion)

	//validar la salida
	if err2 != nil {
		return fmt.Errorf("error al intenter escribir en el user.txt: %w", err2)
	}

	return nil
}

// Funcion para accder al archivo de user.txt
// CrearUser: path del disco, objeto con los datos del usuario, el inicio de los inodos
func MkusrComand(path string, login *LOGIN, comando *MKUSR, inodeIndex int32, sb *structures.SuperBlock, mountedPartition *structures.PARTITION) error {
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
	// Iterar sobre cada bloque del inodo (apuntadores)
	for indiceList, blockIndex := range inode.I_block {
		// Si el bloque no existe, salir
		if blockIndex == -1 {
			break
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

		/*
			1,G,root
			1,U,root,root,123
		*/
		// Dividir por salto de línea
		lines := strings.Split(data, "\n")
		// Recorrer cada línea y dividir por comas
		var id = ""
		// var grupo = ""
		var usuario = ""
		// var password = ""

		// //Para validar si existe el grupo
		// var existeGrupo = true
		//Esto es para obtner los datos puntuales
		for _, line := range lines {
			values := strings.Split(line, ",")

			// Almacenar en variables según la cantidad de datos
			//Esto son los grupos
			if len(values) == 5 {
				//Estos son los usuarios
				id, tipo, _, usuario, _ = values[0], values[1], values[2], values[3], values[4]
				//Esto quiere decir que el grupo esta elimando por lo tanto guradomos el correcto
				if id != "0" {
					respaldoId = id
				}
				// //Esto valida que exista el grupo
				// if grupo == comando.grp {
				// 	existeGrupo = false
				// }
				//Aca se valida si el usuario ya esta pero tambien si no esta eliminado
				if usuario == comando.user && id != "0" {
					return fmt.Errorf("error ya existe otro usuario: %s", usuario)
				}
				//fmt.Printf("ID: %s, Tipo: %s, Nombre: %s\n", id, tipo, nombre)
			}

		}

		// if existeGrupo {
		// 	return fmt.Errorf("error no existe el grupo: %s", comando.grp)
		// }

		//Esto solo es para comvertirlo a numero
		num, err := strconv.Atoi(respaldoId)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		//Sumamos uno al grupo ya existente
		int32Num := int32(num)
		int32Num += 1
		nuevoUsuario := strconv.Itoa(int(int32Num)) + "," + tipo + "," + comando.grp + "," + comando.user + "," + comando.pass + "\n"

		//Aca debo de validar que el indice no se salga de 14
		if indiceList+1 == 15 {
			//Aca si llegamos al ultimo inodo y esta llena el FileBlock
			if len(data) >= 64 {
				return errors.New("ya no existe espacio para crear mas grupos")
				//Aca valido si la suma del nuevo grupo supera los 64 bytes
			} else if (len(data) + len(nuevoUsuario)) >= 64 {
				return errors.New("ya no existe espacio en el ultimo bloque para crear mas grupos")
				//Aca se agrega el nuevo grupo
			} else {
				//Unicamente se debe de agergar el grupo y detener el programa
				//Se agrega a la cadena principal
				data += nuevoUsuario

				//Sustituir los datos anteriores por los nuevos
				// Copiamos el texto de usuarios en el bloque
				copy(block.B_content[:], data)

				//Se serealiza todo el contenido en el Fileblock
				err2 := block.Serialize(path, int64(sb.S_block_start+(int32(indiceFinal)*sb.S_block_size)))
				if err2 != nil {
					return err2
				}

				break
			}
		}

		//validar si el bloque esta lleno ó si la suma del nuevo texto supera el espacio de 64 bytes
		if len(data) >= 64 {
			//Aca el sigueinte bloque esta basillo
			if inode.I_block[indiceList+1] == -1 {
				//Aca se debe de generar un nuevo FileBlock
				//TODO: pendiente
				//Primero se actualiza el ibloque
				inode.I_block[indiceList+1] = blockIndex + 1
				//Serealizar el ibloque actualizado
				// Deserializar el inodo
				err := inode.Serialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
				if err != nil {
					return err
				}
				//Generamos el nuevo Filebloque
				nuevoFilebloque := &structures.FileBlock{}
				// Copiamos el texto de usuarios en el bloque
				copy(nuevoFilebloque.B_content[:], nuevoUsuario)
				//Se serealiza todo el contenido en el Fileblock
				//					inicio de la tabla de bloques + (indice +1* temaño del bloque)
				err2 := nuevoFilebloque.Serialize(path, int64(sb.S_block_start+(int32(indiceFinal+1)*sb.S_block_size)))
				if err2 != nil {
					return err2
				}
				// Actualizar el bitmap de bloques
				err = sb.UpdateBitmapBlock(path)
				if err != nil {
					return err
				}
				//Se actualiza el superbloque
				sb.S_blocks_count++
				sb.S_free_blocks_count--
				sb.S_first_blo += sb.S_block_size

				// Serializar el superbloque
				err = sb.Serialize(path, int64(mountedPartition.Part_start))
				if err != nil {
					return fmt.Errorf("error al serializar el superbloque: %w", err)
				}
				// fmt.Println("**********")
				// nuevoFilebloque.Print()
				break
			} else {
				//Aca necesito avanzar al sigueinte bloque
			}
		} else if (len(data) + len(nuevoUsuario)) >= 64 {
			//Validamos si el nuevo usuario mas el texto anterior no se pase de 64 bytes
			//Validamos que el sigueinte iblock esta basillo y si si lo creamos
			if inode.I_block[indiceList+1] == -1 {
				//Aca se debe de generar un nuevo FileBlock
				//TODO: pendiente
				//Primero se actualiza el ibloque
				inode.I_block[indiceList+1] = blockIndex + 1
				//Serealizar el ibloque actualizado
				// Deserializar el inodo
				err := inode.Serialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
				if err != nil {
					return err
				}
				//Generamos el nuevo Filebloque
				nuevoFilebloque := &structures.FileBlock{}
				// Copiamos el texto de usuarios en el bloque
				copy(nuevoFilebloque.B_content[:], nuevoUsuario)
				//Se serealiza todo el contenido en el Fileblock
				//					inicio de la tabla de bloques + (indice +1* temaño del bloque)
				err2 := nuevoFilebloque.Serialize(path, int64(sb.S_block_start+(int32(indiceFinal+1)*sb.S_block_size)))
				if err2 != nil {
					return err2
				}
				// Actualizar el bitmap de bloques
				err = sb.UpdateBitmapBlock(path)
				if err != nil {
					return err
				}
				//Se actualiza el superbloque
				sb.S_blocks_count++
				sb.S_free_blocks_count--
				sb.S_first_blo += sb.S_block_size

				// Serializar el superbloque
				err = sb.Serialize(path, int64(mountedPartition.Part_start))
				if err != nil {
					return fmt.Errorf("error al serializar el superbloque: %w", err)
				}
				// fmt.Println("**********")
				// nuevoFilebloque.Print()
			} else {
				//Aca necesito avanzar al sigueinte bloque
			}
		} else {
			//Aca quiere decir que ya puedo agregar el nuevo grupo al fileblock y finalizo el bucle
			//Se agrega a la cadena principal
			data += nuevoUsuario

			//Sustituir los datos anteriores por los nuevos
			// Copiamos el texto de usuarios en el bloque
			copy(block.B_content[:], data)

			//Se serealiza todo el contenido en el Fileblock
			err2 := block.Serialize(path, int64(sb.S_block_start+(int32(indiceFinal)*sb.S_block_size)))
			if err2 != nil {
				return err2
			}
			// fmt.Println("**********")
			// block.Print()
			break
		}

	}

	//block.Print()

	return nil
}
