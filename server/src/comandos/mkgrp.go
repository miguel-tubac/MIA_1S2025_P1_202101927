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
	//Obtenemos el usuario logeado
	var usuario = ObtenerUsuari()

	//Valida si el usuario es el usuario root
	if usuario.user != "root" {
		return errors.New("para crear un grupo debe se estar logeado como root")
	}

	// Obtener la partición montada
	//Tipo de retorno: (*structures.SuperBlock, *structures.PARTITION, string, error)
	partitionSuperblock, _, partitionPath, err := stores.GetMountedPartitionSuperblock(usuario.id)
	if err != nil {
		return fmt.Errorf("error al obtener el Superbloque: %w", err)
	}

	//Aca iniciamos desde el inodo numero 1
	err2 := MkgprComand(partitionPath, usuario, comando, 1, partitionSuperblock)

	//validar la salida
	if err2 != nil {
		return fmt.Errorf("error al intenter escribir en el user.txt: %w", err2)
	}

	return nil
}

// Funcion para accder al archivo de user.txt
// CrearUser: path del disco, objeto con los datos del usuario, el inicio de los inodos
func MkgprComand(path string, login *LOGIN, comando *MKGRP, inodeIndex int32, sb *structures.SuperBlock) error {
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
		var tipo = ""
		var nombre = ""
		var respaldoId = ""
		//Esto es para obtner los datos puntuales
		for _, line := range lines {
			values := strings.Split(line, ",")

			// Almacenar en variables según la cantidad de datos
			//Esto son los grupos
			if len(values) == 3 {
				id, tipo, nombre = values[0], values[1], values[2]
				//Esto quiere decir que el grupo esta elimando por lo tanto guradomos el correcto
				if id != "0" {
					respaldoId = id
				}
				if nombre == comando.name {
					return fmt.Errorf("error ya existe otro usuario: %s", nombre)
				}
				//fmt.Printf("ID: %s, Tipo: %s, Nombre: %s\n", id, tipo, nombre)
			}
		}

		//Esto solo es para comvertirlo a numero
		num, err := strconv.Atoi(respaldoId)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		//Sumamos uno al grupo ya existente
		int32Num := int32(num)
		int32Num += 1
		nuevoUsuario := strconv.Itoa(int(int32Num)) + "," + tipo + "," + comando.name + "\n"

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
