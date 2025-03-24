package analyzer

import (
	stores "bakend/src/almacenamiento"
	structures "bakend/src/estructuras"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type RMGRP struct {
	name string
}

func ParseRmgrp(tokens []string) (*RMGRP, error) {
	cmd := &RMGRP{}

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
			return nil, fmt.Errorf("parámetro desconocido en el rmgrp: %s", key)
		}
	}

	// Verifica que los parámetros -user, -pass y -id hayan sido proporcionados
	if cmd.name == "" {
		return nil, errors.New("faltan parámetros requeridos: -name")
	}

	// Montamos la partición
	err := commandRmgrp(cmd)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return cmd, fmt.Errorf("grupo elimando correctamnete: %+v", *cmd)
}

func commandRmgrp(comando *RMGRP) error {
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
	err2 := RmgrpComand(partitionPath, usuario, comando, 1, partitionSuperblock)

	//validar la salida
	if err2 != nil {
		return fmt.Errorf("error al intenter escribir en el user.txt: %w", err2)
	}

	return nil
}

// Funcion para accder al archivo de user.txt
// CrearUser: path del disco, objeto con los datos del usuario, el inicio de los inodos
func RmgrpComand(path string, login *LOGIN, comando *RMGRP, inodeIndex int32, sb *structures.SuperBlock) error {
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
	for _, blockIndex := range inode.I_block {
		// Si el bloque no existe, salir
		if blockIndex == -1 {
			//Esto quiere decir que no encontro al grupo
			return fmt.Errorf("no se encontro ningun grupo con el name: %s", comando.name)
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
		// var id = ""
		var tipo = ""
		var nombre = ""

		//Esta variable almacenara el resultado
		result := ""
		encontrado := false

		//Esto es para obtner los datos puntuales
		for _, line := range lines {
			values := strings.Split(line, ",")

			// Almacenar en variables según la cantidad de datos
			//Esto son los grupos
			if len(values) == 3 {
				_, tipo, nombre = values[0], values[1], values[2]
				if nombre == comando.name {
					encontrado = true
					//Se edita el id del grupo
					result += "0," + tipo + "," + nombre + "\n"
					//Continua por si hay mas datos y se agregan al contenido del fileblock
					continue
				}
				//fmt.Printf("ID: %s, Tipo: %s, Nombre: %s\n", id, tipo, nombre)
			}
			result += line + "\n"
		}

		if encontrado {

			// Copiamos el texto de usuarios en el bloque
			copy(block.B_content[:], []byte(result))
			//Se serealiza todo el contenido en el Fileblock
			err2 := block.Serialize(path, int64(sb.S_block_start+(int32(indiceFinal)*sb.S_block_size)))
			if err2 != nil {
				return err2
			}
			// fmt.Println("**********")
			// block.Print()
			//Finalizamos el bucle
			break
		}

	}

	//block.Print()

	return nil
}
