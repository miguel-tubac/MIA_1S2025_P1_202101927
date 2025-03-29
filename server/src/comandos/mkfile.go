package analyzer

import (
	stores "bakend/src/almacenamiento"
	structures "bakend/src/estructuras"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type MKFILE struct {
	path string
	r    bool
	size int32
	cont string
}

func ParseMkfile(tokens []string) (*MKFILE, error) {
	cmd := &MKFILE{}

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando mount
	re := regexp.MustCompile(`-(?i:path="[^"]+"|path=[^\s]+|r|cont="[^"]+"|cont=[^\s]+|size=[^\s]+)`)
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
		case "-cont":
			if len(kv) != 2 {
				return nil, fmt.Errorf("formato de parámetro inválido: %s", match)
			}
			value := kv[1]
			// Remove quotes from value if present
			if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
				value = strings.Trim(value, "\"")
			}
			cmd.cont = value
		case "-size":
			if len(kv) != 2 {
				return nil, fmt.Errorf("formato de parámetro inválido: %s", match)
			}
			value := kv[1]
			//Esto para convertir el texto a numero
			num, err := strconv.Atoi(value)
			if err != nil {
				//fmt.Println("Error:", err)
				return nil, fmt.Errorf("error conversion: %s", err)
			}

			cmd.size = int32(num)
		case "-r":
			cmd.r = true
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que los parámetros -user, -pass y -id hayan sido proporcionados
	if cmd.path == "" {
		return nil, errors.New("faltan parámetros requeridos: -path")
	}

	//Verifia si el parametro no es negativo
	if cmd.size < 0 {
		return nil, errors.New("el paramtro -size no puede ser negativo")
	}

	// Agregamos al usuario
	err := commandMkfile(cmd)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return cmd, fmt.Errorf("archivo creado exitosamente: %+v", *cmd)
}

// Aquí debería de estar logeado un usuario, por lo cual el usuario debería tener consigo el id de la partición
func commandMkfile(mkfile *MKFILE) error {
	//Obtenemos el usuario logeado
	var usuario = ObtenerUsuari()

	// Obtener la partición montada
	partitionSuperblock, mountedPartition, partitionPath, err := stores.GetMountedPartitionSuperblock(usuario.id)
	if err != nil {
		return fmt.Errorf("error al obtener la partición montada: %w", err)
	}

	// Crear el directorio
	err = createFile(mkfile, partitionSuperblock, partitionPath, mountedPartition)
	if err != nil {
		err = fmt.Errorf("error al crear el directorio: %w", err)
	}

	return err
}

func createFile(mkfile *MKFILE, sb *structures.SuperBlock, partitionPath string, mountedPartition *structures.PARTITION) error {
	//fmt.Println("\nCreando directorio y archivo:", mkfile.path)

	// // GetParentDirectories obtiene las carpetas padres y el directorio de destino
	// parentDirs, nombreArchivo := utils.GetParentDirectories(mkfile.path)
	// // fmt.Println("Directorios padres:", parentDirs)
	// // fmt.Println("Nombre archivo:", nombreArchivo)

	// //Aca se deben de crear las carpetas padres
	// if mkfile.r {
	// 	//Aca se iran agregando las carpetas
	// 	destDir := ""
	// 	var nuevo []string
	// 	//Esta validacion no se avalua
	// 	validacion := true
	// 	// Iterar sobre cada inodo ya que se necesita buscar el inodo padre
	// 	for i := 0; i < len(parentDirs); i++ {
	// 		destDir = parentDirs[i]

	// 		// Asegurar que nuevo no se modifique dentro de CreateFolder
	// 		tempNuevo := append([]string{}, nuevo...)

	// 		// Crear el directorio segun el path proporcionado
	// 		err := sb.CreateFolder(partitionPath, tempNuevo, destDir, &validacion)
	// 		if err != nil {
	// 			return fmt.Errorf("error al crear el directorio: %w", err)
	// 		}

	// 		// Serializar el superbloque
	// 		err = sb.Serialize(partitionPath, int64(mountedPartition.Part_start))
	// 		if err != nil {
	// 			return fmt.Errorf("error al serializar el superbloque: %w", err)
	// 		}

	// 		// fmt.Println("********agregado***********")
	// 		// fmt.Println(nuevo)
	// 		// fmt.Println(destDir)
	// 		nuevo = append(nuevo, destDir)
	// 	}

	// 	//Aca se valida si no esta la ruta creada
	// } else {
	// 	// aca unicamente se valida si la ruta esta creada
	// 	valida := true
	// 	//Crea una copia para no editar el areglo riginal
	// 	copia := make([]string, len(parentDirs)) // Crear un slice con el mismo tamaño
	// 	copy(copia, parentDirs)                  // Copiar los elementos
	// 	//fmt.Println("---------------------")
	// 	// fmt.Println(copia)
	// 	// fmt.Println(nombreArchivo)
	// 	//fmt.Println(sb.S_inodes_count)
	// 	//Anaaliza el destino
	// 	errr := sb.ComprovarFolder(partitionPath, copia, nombreArchivo, &valida)
	// 	if errr != nil {
	// 		//aca se genero un error
	// 		return fmt.Errorf("error al comprovar si existe la ruta: %w", errr)
	// 	}

	// 	//La funcion cambia a false y por lo tanto no deberia entrar
	// 	//Si no cambia a false es porque no lo encontro
	// 	if valida {
	// 		return errors.New("no existe la carpeta padres")
	// 	}
	// }

	// contenido := ""
	// //Aca se debe de copiar el contenido del archivo al nuevo archivo
	// //Se debe de ir a buscar a mis archivos
	// if mkfile.cont != "" {
	// 	informacion, err := LeerArchivo(mkfile.cont)
	// 	if err != nil {
	// 		return fmt.Errorf("error al leer el archivo: %w", err)
	// 	}
	// 	//Se le agrega la informacion del archivo al contenido
	// 	contenido += informacion

	// 	//Aca se debe de generar el contenido del nuevo archivo, el cual debe ser numeros 0-9
	// } else if mkfile.size > 0 {
	// 	//Aca se genera la cadena numerica
	// 	contenido += GenerarCadenaNumerica(mkfile.size)
	// }

	// //Aca ya se debe de generar el archivo en el disco virtual con el -path
	// err := sb.CreateFile(partitionPath, parentDirs, nombreArchivo, contenido)
	// if err != nil {
	// 	return fmt.Errorf("error al crear el archivo: %w", err)
	// }
	// // Serializar el superbloque
	// err = sb.Serialize(partitionPath, int64(mountedPartition.Part_start))
	// if err != nil {
	// 	return fmt.Errorf("error al serializar el superbloque: %w", err)
	// }

	return nil
}

// LeerArchivo recibe una ruta y devuelve el contenido del archivo como string
func LeerArchivo(ruta string) (string, error) {
	// Leer el contenido del archivo
	contenido, err := ioutil.ReadFile(ruta)
	if err != nil {
		return "", err
	}

	// Convertir el contenido a string y retornarlo
	return string(contenido), nil
}

// GenerarCadenaNumerica genera una cadena de números del 0 al 9 repetidos hasta la longitud especificada
func GenerarCadenaNumerica(longitud int32) string {
	base := "0123456789"
	repeticiones := int(longitud) / len(base)
	resto := int(longitud) % len(base)

	// Construcción de la cadena con repeticiones
	return strings.Repeat(base, repeticiones) + base[:resto]
}
