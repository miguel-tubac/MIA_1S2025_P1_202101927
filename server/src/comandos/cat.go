package analyzer

import (
	stores "bakend/src/almacenamiento"
	structures "bakend/src/estructuras"
	utils "bakend/src/utils"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type CAT struct {
	file     []string
	textObte string
}

func ParseCat(tokens []string) (*CAT, error) {
	cmd := &CAT{}

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando fdisk
	re := regexp.MustCompile(`-(?i:file[0-9]+="[^"]+"|file[0-9]+=[^\s]+)`)
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

		// Remove las comillas si estan presentes
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Switch para manejar diferentes parámetros
		filePattern := regexp.MustCompile(`^-file[0-9]+$`)
		if filePattern.MatchString(key) {
			// Verifica que el path no esté vacío
			if value == "" {
				return nil, errors.New("el parámetro file no puede estar vacío")
			}
			cmd.file = append(cmd.file, value)
		} else {
			return nil, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Obtiene la informacion con los parámetros proporcionados
	err := commandFile(cmd)
	if err != nil {
		fmt.Println("Error:", err)
		return cmd, err
	}

	return cmd, fmt.Errorf("%+v", cmd.textObte) // Devuelve el comando FDISK creado
}

// Esto es para obtener el superbloque
func commandFile(comando *CAT) error {
	//Obtenemos el usuario logeado
	var usuario = ObtenerUsuari()

	// Obtener la partición montada
	//Tipo de retorno: (*structures.SuperBlock, *structures.PARTITION, path del disco string, error)
	partitionSuperblock, _, partitionPath, err := stores.GetMountedPartitionSuperblock(usuario.id)
	if err != nil {
		return fmt.Errorf("error al obtener el Superbloque: %w", err)
	}

	comando.textObte += "***************** CAT ********************"
	for _, direccion := range comando.file {
		//Aca se recorre y se obtiene
		err2 := ObtnerFileDisco(comando, partitionSuperblock, partitionPath, direccion)
		if err2 != nil {
			return fmt.Errorf("error al intenter obtener el archivo .txt: %w", err2)
		}
	}

	return nil
}

func ObtnerFileDisco(comando *CAT, superblock *structures.SuperBlock, diskPath string, path_file_ls string) error {

	// GetParentDirectories obtiene las carpetas padres y el directorio de destino
	parentDirs, nombreArchivo := utils.GetParentDirectories(path_file_ls)
	// fmt.Println("Directorios padres:", parentDirs)
	// fmt.Println("Nombre archivo:", nombreArchivo)

	//Aca ya se debe de obtner el archivo en el disco virtual con el -path
	cade, err2 := superblock.GetFileContent(diskPath, parentDirs, nombreArchivo)
	if err2 != nil {
		return fmt.Errorf("error al obtener el archivo: %w", err2)
	}

	//Aca se agrega el encabezado de cada archivo
	comando.textObte += "\n-------------- " + nombreArchivo + " --------------\n"
	//Se concatena la informacion
	comando.textObte += cade
	//fmt.Println("Archivo creado exitosamente en:", outputImage)
	return nil
}
