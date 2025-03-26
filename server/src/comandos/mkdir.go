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

// MKDIR estructura que representa el comando mkdir con sus parámetros
type MKDIR struct {
	path string // Path del directorio
	p    bool   // Opción -p (crea directorios padres si no existen)
}

/*
   mkdir -p -path=/home/user/docs/usac
   mkdir -path="/home/mis documentos/archivos clases"
*/

func ParseMkdir(tokens []string) (*MKDIR, error) {
	cmd := &MKDIR{} // Crea una nueva instancia de MKDIR

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando mkdir
	re := regexp.MustCompile(`-(?i:path=[^\s]+|path="[^"]+"|p)`)
	// Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
	matches := re.FindAllString(args, -1)

	// Itera sobre cada coincidencia encontrada
	for _, match := range matches {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		key := strings.ToLower(kv[0])

		// Switch para manejar diferentes parámetros
		switch key {
		case "-path":
			if len(kv) != 2 {
				return nil, fmt.Errorf("formato de parámetro inválido: %s", match)
			}
			value := kv[1]
			// Remove quotes from value if present
			if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
				value = strings.Trim(value, "\"")
			}
			cmd.path = value
		case "-p":
			cmd.p = true
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que el parámetro -path haya sido proporcionado
	if cmd.path == "" {
		return nil, errors.New("faltan parámetros requeridos: -path")
	}

	// Aquí se puede agregar la lógica para ejecutar el comando mkdir con los parámetros proporcionados
	err := commandMkdir(cmd)
	if err != nil {
		return nil, err
	}

	return cmd, fmt.Errorf("creado correctamente MKDIR: %+v", *cmd)
}

// Aquí debería de estar logeado un usuario, por lo cual el usuario debería tener consigo el id de la partición
func commandMkdir(mkdir *MKDIR) error {
	//Obtenemos el usuario logeado
	var usuario = ObtenerUsuari()

	// Obtener la partición montada
	partitionSuperblock, mountedPartition, partitionPath, err := stores.GetMountedPartitionSuperblock(usuario.id)
	if err != nil {
		return fmt.Errorf("error al obtener la partición montada: %w", err)
	}

	// Crear el directorio
	err = createDirectory(mkdir.path, partitionSuperblock, partitionPath, mountedPartition)
	if err != nil {
		err = fmt.Errorf("error al crear el directorio: %w", err)
	}

	return err
}

func createDirectory(dirPath string, sb *structures.SuperBlock, partitionPath string, mountedPartition *structures.PARTITION) error {
	fmt.Println("\nCreando directorio:", dirPath)

	// GetParentDirectories obtiene las carpetas padres y el directorio de destino
	parentDirs, destDir := utils.GetParentDirectories(dirPath)
	//fmt.Println("\nDirectorios padres:", parentDirs)
	//fmt.Println("Directorio destino:", destDir)

	// Crear el directorio segun el path proporcionado
	err := sb.CreateFolder(partitionPath, parentDirs, destDir)
	if err != nil {
		return fmt.Errorf("error al crear el directorio: %w", err)
	}

	// Imprimir inodos y bloques
	//sb.PrintInodes(partitionPath)
	//sb.PrintBlocks(partitionPath)

	// Serializar el superbloque
	err = sb.Serialize(partitionPath, int64(mountedPartition.Part_start))
	if err != nil {
		return fmt.Errorf("error al serializar el superbloque: %w", err)
	}

	// Imprimir inodos y bloques
	// sb.PrintInodes(partitionPath)
	// sb.PrintBlocks(partitionPath)
	return nil
}
