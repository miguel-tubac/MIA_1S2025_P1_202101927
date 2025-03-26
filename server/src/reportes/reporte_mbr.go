package reportes

import (
	structures "bakend/src/estructuras"
	utils "bakend/src/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ReportMBR genera un reporte del MBR y lo guarda en la ruta especificada
func ReportMBR(mbr *structures.MBR, path string, pathDisco string) error {
	// Crear las carpetas padre si no existen
	err := utils.CreateParentDirs(path)
	if err != nil {
		return err
	}

	// Obtener el nombre base del archivo sin la extensión
	dotFileName, outputImage := utils.GetFileNames(path)

	// Definir el contenido DOT con una tabla
	dotContent := fmt.Sprintf(`digraph G {
        node [shape=plaintext]
        tabla [label=<
            <table border="0" cellborder="1" cellspacing="0">
                <tr><td colspan="2" bgcolor="#800080"><font color="white"> REPORTE MBR </font></td></tr>
                <tr><td>mbr_tamano</td><td>%d</td></tr>
                <tr><td>mrb_fecha_creacion</td><td>%s</td></tr>
                <tr><td>mbr_disk_signature</td><td>%d</td></tr>
            `, mbr.Mbr_size, time.Unix(int64(mbr.Mbr_creation_date), 0), mbr.Mbr_disk_signature)

	// Agregar las particiones a la tabla
	for i, part := range mbr.Mbr_partitions {

		// Continuar si el tamaño de la partición es -1 (o sea, no está asignada)
		if part.Part_size == -1 {
			continue
		}

		// Convertir Part_name a string y eliminar los caracteres nulos
		partName := strings.TrimRight(string(part.Part_name[:]), "\x00")
		// Convertir Part_status, Part_type y Part_fit a char
		partStatus := rune(part.Part_status[0])
		partType := rune(part.Part_type[0])
		partFit := rune(part.Part_fit[0])

		// Agregar la partición a la tabla
		dotContent += fmt.Sprintf(`
			<tr><td colspan="2" bgcolor="#800080"><font color="white"> PARTICIÓN %d </font></td></tr>
			<tr><td>part_status</td><td>%c</td></tr>
			<tr><td>part_type</td><td>%c</td></tr>
			<tr><td>part_fit</td><td>%c</td></tr>
			<tr><td>part_start</td><td>%d</td></tr>
			<tr><td>part_size</td><td>%d</td></tr>
			<tr><td>part_name</td><td>%s</td></tr>
		`, i+1, partStatus, partType, partFit, part.Part_start, part.Part_size, partName)

		if partType == 'E' {
			contador := 1
			posicion := part.Part_start
			var ebr structures.EBR
			for {
				err2 := ebr.DeserializeEBR(pathDisco, posicion)
				if err2 != nil {
					fmt.Println("Error deserializando el EBR:", err2)
					return err2
				}

				if ebr.Ebr_next == -1 {
					break
				}

				EBRmount := rune(ebr.Ebr_mount[0])
				EBRfit := rune(ebr.Ebr_fit[0])
				EBRname := strings.TrimRight(string(ebr.Ebr_name[:]), "\x00")
				// Agregar la partición a la tabla
				dotContent += fmt.Sprintf(`
					<tr><td colspan="2" bgcolor="#FF1493"><font color="white"> PARTICIÓN LOGICA %d </font></td></tr>
					<tr><td>part_mount</td><td>%c</td></tr>
					<tr><td>part_fit</td><td>%c</td></tr>
					<tr><td>part_start</td><td>%d</td></tr>
					<tr><td>part_size</td><td>%d</td></tr>
					<tr><td>part_next</td><td>%d</td></tr>
					<tr><td>part_name</td><td>%s</td></tr>
				`, contador, EBRmount, EBRfit, ebr.Ebr_start, ebr.Ebr_size, ebr.Ebr_next, EBRname)

				contador++
				posicion = ebr.Ebr_next
			}
		}
	}

	// Cerrar la tabla y el contenido DOT
	dotContent += "</table>>] }"

	// Guardar el contenido DOT en un archivo
	file, err := os.Create(dotFileName)
	if err != nil {
		return fmt.Errorf("error al crear el archivo: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(dotContent)
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo: %v", err)
	}

	// Ejecutar el comando Graphviz para generar la imagen
	cmd := exec.Command("dot", "-Tpng", dotFileName, "-o", outputImage)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error al ejecutar el comando Graphviz: %v", err)
	}

	//fmt.Println("Imagen de la tabla generada:", outputImage)
	return nil
}
