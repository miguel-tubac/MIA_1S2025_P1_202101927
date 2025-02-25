package analyzer

import (
	"errors"  // Paquete para manejar errores y crear nuevos errores con mensajes personalizados
	"fmt"     // Paquete para formatear cadenas y realizar operaciones de entrada/salida
	"regexp"  // Paquete para trabajar con expresiones regulares, útil para encontrar y manipular patrones en cadenas
	"strconv" // Paquete para convertir cadenas a otros tipos de datos, como enteros
	"strings" // Paquete para manipular cadenas, como unir, dividir, y modificar contenido de cadenas

	structures "bakend/src/estructuras"
	utils "bakend/src/utils"
)

// FDISK estructura que representa el comando fdisk con sus parámetros
type FDISK struct {
	size int    // Tamaño de la partición
	unit string // Unidad de medida del tamaño (K o M)
	fit  string // Tipo de ajuste (BF, FF, WF)
	path string // Ruta del archivo del disco
	typ  string // Tipo de partición (P, E, L)
	name string // Nombre de la partición
}

/*
	fdisk -size=1 -type=L -unit=M -fit=BF -name="Particion3" -path="/home/keviin/University/PRACTICAS/MIA_LAB_S2_2024/CLASEEXTRA/disks/Disco1.mia"
	fdisk -size=300 -path=/home/Disco1.mia -name=Particion1
	fdisk -type=E -path=/home/Disco2.mia -Unit=K -name=Particion2 -size=300
*/

// CommandFdisk parsea el comando fdisk y devuelve una instancia de FDISK
func ParseFdisk(tokens []string) (*FDISK, error) {
	cmd := &FDISK{} // Crea una nueva instancia de FDISK

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando fdisk
	re := regexp.MustCompile(`-(?i:size=\d+|unit=[kKmMbB]|fit=[bBfFwW]{2}|path="[^"]+"|path=[^\s]+|type=[pPeElL]|name="[^"]+"|name=[^\s]+)`)
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
		switch key {
		case "-size":
			// Convierte el valor del tamaño a un entero
			size, err := strconv.Atoi(value)
			if err != nil || size <= 0 {
				return nil, errors.New("el tamaño debe ser un número entero positivo")
			}
			cmd.size = size
		case "-unit":
			value = strings.ToUpper(value)
			// Verifica que la unidad sea "K" o "M"
			if value != "K" && value != "M" && value != "B" {
				return nil, errors.New("la unidad debe ser K, M o B")
			}
			cmd.unit = value
		case "-fit":
			// Verifica que el ajuste sea "BF", "FF" o "WF"
			value = strings.ToUpper(value)
			if value != "BF" && value != "FF" && value != "WF" {
				return nil, errors.New("el ajuste debe ser BF, FF o WF")
			}
			cmd.fit = value
		case "-path":
			// Verifica que el path no esté vacío
			if value == "" {
				return nil, errors.New("el path no puede estar vacío")
			}
			cmd.path = value
		case "-type":
			// Verifica que el tipo sea "P", "E" o "L"
			value = strings.ToUpper(value)
			if value != "P" && value != "E" && value != "L" {
				return nil, errors.New("el tipo debe ser P, E o L")
			}
			cmd.typ = value
		case "-name":
			// Verifica que el nombre no esté vacío
			if value == "" {
				return nil, errors.New("el nombre no puede estar vacío")
			}
			cmd.name = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que los parámetros -size, -path y -name hayan sido proporcionados
	if cmd.size == 0 {
		return nil, errors.New("faltan parámetros requeridos: -size")
	}
	if cmd.path == "" {
		return nil, errors.New("faltan parámetros requeridos: -path")
	}
	if cmd.name == "" {
		return nil, errors.New("faltan parámetros requeridos: -name")
	}

	// Si no se proporcionó la unidad, se establece por defecto a "M"
	if cmd.unit == "" {
		cmd.unit = "K"
	}

	// Si no se proporcionó el ajuste, se establece por defecto a "FF"
	if cmd.fit == "" {
		cmd.fit = "WF"
	}

	// Si no se proporcionó el tipo, se establece por defecto a "P"
	if cmd.typ == "" {
		cmd.typ = "P" //Es una particion primaria por defecto
	}

	// Crear la partición con los parámetros proporcionados
	err := commandFdisk(cmd)
	if err != nil {
		fmt.Println("Error:", err)
		return cmd, err
	}

	return cmd, nil // Devuelve el comando FDISK creado
}

func commandFdisk(fdisk *FDISK) error {
	// Convertir el tamaño a bytes
	sizeBytes, err := utils.ConvertToBytes(fdisk.size, fdisk.unit)
	if err != nil {
		fmt.Println("Error al convertir las unidades de size:", err)
		return err
	}

	if fdisk.typ == "P" {
		// Crear partición primaria
		err = createPrimaryPartition(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error creando partición primaria:", err)
			return err
		}
	} else if fdisk.typ == "E" {
		fmt.Println("Creando partición extendida...") // Les toca a ustedes implementar la partición extendida
		err = createExtendidaPartition(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error creando partición extendida:", err)
			return err
		}
		err2 := createEBR(fdisk, sizeBytes)
		if err2 != nil {
			fmt.Println("Error creando el EBR:", err2)
			return err2
		}
	} else if fdisk.typ == "L" {
		fmt.Println("Creando partición lógica...") // Les toca a ustedes implementar la partición lógica
		err = createLogicPartition(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error agregando la particion Logica:", err)
			return err
		}
	}

	return nil
}

// -------------------------------------------------------------Particion Primaria--------------------------------------------------------------
// Creacion de particiones primarias
func createPrimaryPartition(fdisk *FDISK, sizeBytes int) error {
	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.DeserializeMBR(fdisk.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return err
	}

	/* SOLO PARA VERIFICACIÓN */
	// Imprimir MBR
	fmt.Println("\n--MBR original--")
	mbr.PrintMBR()

	// Obtener la primera partición disponible
	/*PARTITION:
	Part_status      [1]byte  // Estado de la partición
	Part_type        [1]byte  // Tipo de partición
	Part_fit         [1]byte  // Ajuste de la partición
	Part_start       int32    // Byte de inicio de la partición
	Part_size        int32    // Tamaño de la partición
	Part_name        [16]byte // Nombre de la partición
	Part_correlative int32    // Correlativo de la partición
	Part_id          [4]byte  // ID de la partición */
	availablePartition, startPartition, indexPartition := mbr.GetFirstAvailablePartition() //*PARTITION, int, int   (Retornos)
	var err2 error
	err2 = nil
	if availablePartition == nil {
		err2 = errors.New("no hay particiones disponibles o la particion es mas grande que el disco")
	}

	//Se comprueba si existe otra particion con el mismo nombre
	if mbr.ExisteNombre(fdisk.name) && availablePartition != nil && startPartition != -1 {
		/* SOLO PARA VERIFICACIÓN */
		// Print para verificar que la partición esté disponible
		//fmt.Println("\n--Partición disponible--")
		//availablePartition.PrintPartition()

		var corre = mbr.GetCorrelativo()

		// Crear la partición con los parámetros proporcionados
		availablePartition.CreatePartition(startPartition, sizeBytes, fdisk.typ, fdisk.fit, fdisk.name, corre)

		// Print para verificar que la partición se haya creado correctamente
		//fmt.Println("\n--Partición creada (modificada)--")
		//availablePartition.PrintPartition()

		// Se cambio el puntero de las particiones a las actuales
		mbr.Mbr_partitions[indexPartition] = *availablePartition

		// Imprimir las particiones del MBR
		fmt.Println("\n--Particiones del MBR--")
		mbr.PrintPartitions()
	}

	// Serializar el MBR en el archivo binario
	err = mbr.SerializeMBR(fdisk.path)
	if err != nil {
		fmt.Println("Error:", err)
		err2 = err //En caso de que ocura un error se retorna
	}

	return err2 //Si no ocurrio ningun error se retorna nil
}

// -------------------------------------------------------------Particion Extendida--------------------------------------------------------------
// Creacion de particiones primarias
func createExtendidaPartition(fdisk *FDISK, sizeBytes int) error {
	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.DeserializeMBR(fdisk.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return err
	}

	/* SOLO PARA VERIFICACIÓN */
	// Imprimir MBR
	fmt.Println("\n--MBR original--")
	mbr.PrintMBR()

	// Obtener la primera partición disponible
	/*PARTITION:
	  Part_status      [1]byte  // Estado de la partición
	  Part_type        [1]byte  // Tipo de partición
	  Part_fit         [1]byte  // Ajuste de la partición
	  Part_start       int32    // Byte de inicio de la partición
	  Part_size        int32    // Tamaño de la partición
	  Part_name        [16]byte // Nombre de la partición
	  Part_correlative int32    // Correlativo de la partición
	  Part_id          [4]byte  // ID de la partición */
	availablePartition, startPartition, indexPartition := mbr.GetFirstAvailablePartition() //*PARTITION, int, int   (Retornos)
	var err2 error
	err2 = nil
	if availablePartition == nil {
		err2 = errors.New("no hay particiones disponibles o la particion es mas grande que el disco")
	}

	//Validar si existe otra paricion extendida(Solo puede aver una)
	validacion_extendida := mbr.GetExtendedPartition()
	if validacion_extendida {
		err2 = errors.New("solo puede existir una partricion extendida dentro del disco")
	}

	// Se comprueba si existe otra particion con el mismo nombre
	if mbr.ExisteNombre(fdisk.name) && availablePartition != nil && startPartition != -1 && !validacion_extendida {
		/* SOLO PARA VERIFICACIÓN */
		// Print para verificar que la partición esté disponible
		//fmt.Println("\n--Partición disponible--")
		//availablePartition.PrintPartition()

		var corre = mbr.GetCorrelativo()

		// Crear la partición con los parámetros proporcionados
		availablePartition.CreatePartition(startPartition, sizeBytes, fdisk.typ, fdisk.fit, fdisk.name, corre)

		// Print para verificar que la partición se haya creado correctamente
		//fmt.Println("\n--Partición creada (modificada)--")
		//availablePartition.PrintPartition()

		// Se cambio el puntero de las particiones a las actuales
		mbr.Mbr_partitions[indexPartition] = *availablePartition

		// Imprimir las particiones del MBR
		fmt.Println("\n--Particiones del MBR--")
		mbr.PrintPartitions()
	}

	// Serializar el MBR en el archivo binario
	err = mbr.SerializeMBR(fdisk.path)
	if err != nil {
		fmt.Println("Error:", err)
		err2 = err //En caso de que ocura un error se retorna
	}

	return err2 //Si no ocurrio ningun error se retorna nil
}

// Se crea el EBR dentro de la particion
func createEBR(fdisk *FDISK, sizeBytes int) error {
	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.DeserializeMBR(fdisk.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return err
	}

	paricion, _ := mbr.GetPartitionByName(fdisk.name)

	// Crear el EBR con los valores proporcionados
	ebr := &structures.EBR{
		Ebr_mount: [1]byte{'N'},
		Ebr_fit:   [1]byte{'N'},
		Ebr_start: int32(-1),
		Ebr_size:  int32(-1),
		Ebr_next:  int32(-1),
		Ebr_name:  [16]byte{'N'},
	}

	// Serializar el EBR en el archivo
	err2 := ebr.SerializeEBR(fdisk.path, paricion.Part_start) //Se le manda en que posicion se debe de agregar el EBR
	if err2 != nil {
		fmt.Println("Error:", err2)
		return err2
	}

	return nil
}

// -------------------------------------------------------------Particion Logica--------------------------------------------------------------
//  1. Obtenermos la particion extendida
//  2. Obtenemos el start de la extendida
//  3. Deserealizamos el EBR con el ofset del start
//     Funcion Recursuiva:
//  4. Comprobamos que la propiedad NEXT del EBR sea -1 (se agrega la informacion del EBR logico)
//  5. SI NEXT != -1 se accede al sigueinte EBR y se vuelve a evaluar
func createLogicPartition(fdisk *FDISK, sizeBytes int) error {
	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.DeserializeMBR(fdisk.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return err
	}

	/* SOLO PARA VERIFICACIÓN */
	// Imprimir MBR
	fmt.Println("\n--MBR original--")
	mbr.PrintMBR()

	//Validar si existe la particion extendida
	var err2 error
	err2 = nil
	validacion_extendida := mbr.GetExtendedPartition()
	//fmt.Println(validacion_extendida)
	if validacion_extendida {
		//obtenemos la particion
		particion := mbr.GetExtendedPartition2()

		//Obtenemos el objeto EBR deserealizandolo del dico
		var ebr structures.EBR
		//fmt.Println("Aqui Miguel:")
		//fmt.Println(particion.Part_start)                           //Objeto donde se deserealizara el ebr
		err2 = ebr.DeserializeEBR(fdisk.path, particion.Part_start) //Se deserealiza el primer EBR de la extendida
		if err2 != nil {
			fmt.Println("Error deserializando el EBR:", err2)
			return err2
		}

		//Imprimimos el EBR para verificar
		fmt.Println("\n--EBR original--")
		ebr.PrintEBR()

		//Devolver el start tomando en cuenta el valor de Bytes del EBR
		starPartition := ebr.BuscarEBRDisponible(fdisk.path, particion.Part_start)

		// Seleccionar el tipo de ajuste
		var fitByte byte
		switch fdisk.fit {
		case "FF":
			fitByte = 'F'
		case "BF":
			fitByte = 'B'
		case "WF":
			fitByte = 'W'
		default:
			fmt.Println("Invalid fit type")
			return nil
		}
		ebr.CreateLogical(fitByte, int32(starPartition), fdisk.size, fdisk.name)

		/*
			Ebr_mount [1]byte  // Indica si la partición está montada o no
			Ebr_fit   [1]byte  // Tipo de ajuste de la partición. Tendrá los valores B (Best), F(First) o W (worst)
			Ebr_start int32    // Indica en qué byte del disco inicia la partición
			Ebr_size  int32    // Contiene el tamaño total de la partición en bytes.
			Ebr_next  int32    // Byte en el que está el próximo EBR. -1 si no hay siguiente
			Ebr_name  [16]byte //Nombre de la partición*/
	}

	// Se comprueba si existe otra particion con el mismo nombre
	/*if mbr.ExisteNombre(fdisk.name) && availablePartition != nil && startPartition != -1 && !validacion_extendida {
		// SOLO PARA VERIFICACIÓN
		// Print para verificar que la partición esté disponible
		//fmt.Println("\n--Partición disponible--")
		//availablePartition.PrintPartition()

		var corre = mbr.GetCorrelativo()

		// Crear la partición con los parámetros proporcionados
		availablePartition.CreatePartition(startPartition, sizeBytes, fdisk.typ, fdisk.fit, fdisk.name, corre)

		// Print para verificar que la partición se haya creado correctamente
		//fmt.Println("\n--Partición creada (modificada)--")
		//availablePartition.PrintPartition()

		// Se cambio el puntero de las particiones a las actuales
		mbr.Mbr_partitions[indexPartition] = *availablePartition

		// Imprimir las particiones del MBR
		fmt.Println("\n--Particiones del MBR--")
		mbr.PrintPartitions()
	}*/

	/*// Serializar el MBR en el archivo binario
	err = mbr.SerializeMBR(fdisk.path)
	if err != nil {
		fmt.Println("Error:", err)
		err2 = err //En caso de que ocura un error se retorna
	}*/

	return err2 //Si no ocurrio ningun error se retorna nil
}
