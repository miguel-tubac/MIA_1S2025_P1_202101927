package reportes

import (
	structures "bakend/src/estructuras"
	utils "bakend/src/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ReporteDisk(mbrparticion *structures.MBR, superblock *structures.SuperBlock, diskPath string, path string) error {
	// Crear las carpetas padre si no existen
	err := utils.CreateParentDirs(path)
	if err != nil {
		return err
	}

	// Obtener el nombre base del archivo sin la extensión
	dotFileName, outputImage := utils.GetFileNames(path)

	// Leer y procesar los EBRs si hay particiones extendidas
	var ebrs []structures.EBR
	contador := 0
	for i := 0; i < 4; i++ {
		if string(mbrparticion.Mbr_partitions[i].Part_type[:]) == "E" { // Partición extendida
			ebrPosition := mbrparticion.Mbr_partitions[i].Part_start
			for ebrPosition != -1 {
				var tempEBR structures.EBR
				err := tempEBR.DeserializeEBR(diskPath, int32(ebrPosition))
				if err != nil {
					break
				}
				ebrs = append(ebrs, tempEBR)
				ebrPosition = tempEBR.Ebr_next
				contador++
			}
		}
	}

	contador = contador * len(ebrs)
	// Iniciar el contenido del archivo en formato Graphviz (.dot)
	content := "digraph G {\n"
	content += "\tnode [shape=none];\n"
	content += "\tgraph [splines=false];\n"
	content += "\tsubgraph cluster_disk {\n"
	content += "\t\tlabel=\"Disco1.dsk\";\n"
	content += "\t\tstyle=rounded;\n"
	content += "\t\tcolor=black;\n"

	// Iniciar tabla para las particiones
	content += "\t\ttable [label=<\n\t\t\t<TABLE BORDER=\"1\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"10\">\n"
	content += "\t\t\t<TR>\n"
	content += "\t\t\t<TD>MBR (159 bytes)</TD>\n"

	// Variables para el porcentaje y espacio libre
	totalDiskSize := mbrparticion.Mbr_size
	var usedSpace int32 = 159 // Tamaño del MBR en bytes
	var freeSpace int32 = totalDiskSize - usedSpace

	for i := 0; i < 4; i++ {
		part := mbrparticion.Mbr_partitions[i]
		if part.Part_size > 0 { // Si la partición tiene un tamaño valido
			percentage := float64(part.Part_size) / float64(totalDiskSize) * 100
			partName := strings.TrimRight(string(part.Part_name[:]), "\x00") // Limpiar el nombre de la partición

			if string(part.Part_type[:]) == "P" { // Partición primaria
				content += fmt.Sprintf("\t\t\t<TD>Primaria<br/>%s<br/>%.2f%% del disco</TD>\n", partName, percentage)
				usedSpace += part.Part_size
			} else if string(part.Part_type[:]) == "E" { // Partición extendida
				content += "\t\t\t<TD>\n"
				content += "\t\t\t\t<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\">\n"
				content += fmt.Sprintf("\t\t\t\t<TR><TD COLSPAN=\"%d\">Extendida</TD></TR>\n", contador)
				//fmt.Println(contador)
				// Leer los EBRs y agregar las particiones lógicas
				content += "\t\t\t\t<TR>\n"
				for _, ebr := range ebrs {
					logicalPercentage := float64(ebr.Ebr_size) / float64(totalDiskSize) * 100
					content += fmt.Sprintf("\t\t\t\t<TD>EBR (32 bytes)</TD>\n\t\t\t\t<TD>Lógica<br/>%.2f%% del disco</TD>\n", logicalPercentage)
					usedSpace += ebr.Ebr_size + 32 // Añadir el tamaño de la partición lógica y el EBR
				}
				content += "\t\t\t\t</TR>\n"
				content += "\t\t\t\t</TABLE>\n"
				content += "\t\t\t</TD>\n"
			}
		}
	}

	// Recalcular el espacio libre
	freeSpace = totalDiskSize - usedSpace
	freePercentage := float64(freeSpace) / float64(totalDiskSize) * 100

	// Agregar el espacio libre restante
	content += fmt.Sprintf("\t\t\t<TD>Libre<br/>%.2f%% del disco</TD>\n", freePercentage)
	content += "\t\t\t</TR>\n"
	content += "\t\t\t</TABLE>\n>];\n"
	content += "\t}\n"
	content += "}\n"

	// Guardar el contenido DOT en un archivo
	file, err := os.Create(dotFileName)
	if err != nil {
		return fmt.Errorf("error al crear el archivo: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
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
