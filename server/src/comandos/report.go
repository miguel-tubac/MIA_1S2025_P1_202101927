package analyzer

import (
	stores "bakend/src/almacenamiento"
	reports "bakend/src/reportes"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// REP estructura que representa el comando rep con sus parámetros
type REP struct {
	id           string // ID del disco
	path         string // Ruta del archivo del disco
	name         string // Nombre del reporte
	path_file_ls string // Ruta del archivo ls (opcional)
}

// ParserRep parsea el comando rep y devuelve una instancia de REP
func ParseRep(tokens []string) (*REP, error) {
	cmd := &REP{} // Crea una nueva instancia de REP

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando rep
	re := regexp.MustCompile(`-(?i:id=[^\s]+|path="[^"]+"|path=[^\s]+|name=[^\s]+|path_file_ls="[^"]+"|path_file_ls=[^\s]+)`)
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

		// Remove quotes from value if present
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Switch para manejar diferentes parámetros
		switch key {
		case "-id":
			// Verifica que el id no esté vacío
			if value == "" {
				return nil, errors.New("el id no puede estar vacío")
			}
			cmd.id = value
		case "-path":
			// Verifica que el path no esté vacío
			if value == "" {
				return nil, errors.New("el path no puede estar vacío")
			}
			cmd.path = value
		case "-name":
			//Convertimos todo a minuscula
			value = strings.ToLower(value)
			// Verifica que el nombre sea uno de los valores permitidos
			validNames := []string{"mbr", "disk", "inode", "block", "bm_inode", "bm_block", "sb", "file", "ls"}
			if !contains(validNames, value) {
				return nil, errors.New("nombre inválido, debe ser uno de los siguientes: mbr, disk, inode, block, bm_inode, bm_block, sb, file, ls")
			}
			cmd.name = value
		case "-path_file_ls":
			cmd.path_file_ls = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que los parámetros obligatorios hayan sido proporcionados
	if cmd.id == "" {
		return nil, errors.New("faltan parámetros requeridos: -id")
	}
	if cmd.path == "" {
		return nil, errors.New("faltan parámetros requeridos: -path")
	}
	if cmd.name == "" {
		return nil, errors.New("faltan parámetros requeridos: -name")
	}

	// Aquí se puede agregar la lógica para ejecutar el comando rep con los parámetros proporcionados
	err := commandRep(cmd)
	if err != nil {
		//fmt.Println("Error:", err)
		return nil, err
	}

	// Devuelve el comando REP creado
	return cmd, fmt.Errorf("reporte creado correctamnete: %+v", *cmd)
}

// Función auxiliar para verificar si un valor está en una lista
func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// Ejemplo de función commandRep (debe ser implementada)
func commandRep(rep *REP) error {
	// Obtener la partición montada
	//*structures.MBR, *structures.SuperBlock, path(particion) string, error
	mountedMbr, mountedSb, mountedDiskPath, err := stores.GetMountedPartitionRep(rep.id)
	if err != nil {
		return err
	}

	// Switch para manejar diferentes tipos de reportes
	switch rep.name {
	case "mbr":
		err = reports.ReportMBR(mountedMbr, rep.path, mountedDiskPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}
		//Rerornamos el mensaje de satisfacion
		return fmt.Errorf("imagen del MBR generada: %s", rep.path)
	case "inode":
		err = reports.ReportInode(mountedSb, mountedDiskPath, rep.path)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}
		//Rerornamos el mensaje de satisfacion
		return fmt.Errorf("imagen del INODE generada: %s", rep.path)
	case "block":
		err = reports.ReportBlock(mountedSb, mountedDiskPath, rep.path)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}
		//Rerornamos el mensaje de satisfacion
		return fmt.Errorf("imagen del BLOCK generada: %s", rep.path)
	case "bm_inode":
		err = reports.ReportBMInode(mountedSb, mountedDiskPath, rep.path)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}
		//Rerornamos el mensaje de satisfacion
		return fmt.Errorf("reporte (txt) del BM_INODE generado: %s", rep.path)
	case "bm_block":
		err = reports.ReportBMBloc(mountedSb, mountedDiskPath, rep.path)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}
		//Rerornamos el mensaje de satisfacion
		return fmt.Errorf("reporte (txt) del BM_Block generado: %s", rep.path)
	case "sb":
		err = reports.ReporteSB(mountedSb, mountedDiskPath, rep.path)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}
		//Rerornamos el mensaje de satisfacion
		return fmt.Errorf("imagen del Superbloque generado: %s", rep.path)
	case "disk":
		err = reports.ReporteDisk(mountedMbr, mountedSb, mountedDiskPath, rep.path)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}
		//Rerornamos el mensaje de satisfacion
		return fmt.Errorf("imagen del Disco generado: %s", rep.path)
	case "file":
		err = reports.ReporteFile(mountedSb, mountedDiskPath, rep.path, rep.path_file_ls)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}
		//Rerornamos el mensaje de satisfacion
		return fmt.Errorf("reporte (txt) del file generado: %s", rep.path)
	case "ls":
		err = reports.ReporteLs(mountedSb, mountedDiskPath, rep.path, rep.path_file_ls)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}
		//Rerornamos el mensaje de satisfacion
		return fmt.Errorf("imagen del LS generado: %s", rep.path)
	}

	return nil
}
