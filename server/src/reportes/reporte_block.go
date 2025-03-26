package reportes

import (
	structures "bakend/src/estructuras"
	utils "bakend/src/utils"
	"fmt"
	"os"
	"os/exec"
)

// ReportInode genera un reporte de un inodo y lo guarda en la ruta especificada
func ReportBlock(superblock *structures.SuperBlock, diskPath string, path string) error {
	// Crear las carpetas padre si no existen
	err := utils.CreateParentDirs(path)
	if err != nil {
		return err
	}

	// Obtener el nombre base del archivo sin la extensión
	dotFileName, outputImage := utils.GetFileNames(path)

	// Iniciar el contenido DOT
	dotContent := `digraph G {
	node [shape=plaintext]
	`

	// Iterar sobre cada inodo
	for i := int32(0); i < superblock.S_inodes_count; i++ {
		inode := &structures.Inode{}

		// Deserializar el inodo
		err := inode.Deserialize(diskPath, int64(superblock.S_inode_start+(i*superblock.S_inode_size)))
		if err != nil {
			return err
		}

		// Iterar sobre cada bloque del inodo (apuntadores)
		for indice, blockIndex := range inode.I_block {
			// Si el bloque no existe, salir
			if blockIndex == -1 {
				break
			}
			//TODO: validar los apuntadores 13,14,15

			// Definir el contenido DOT para el inodo actual
			dotContent += fmt.Sprintf(`bloque%d [label=<
			<table border="0" cellborder="1" cellspacing="0">`, blockIndex)
			// Si el inodo es de tipo carpeta
			if inode.I_type[0] == '0' {
				block := &structures.FolderBlock{}
				// Deserializar el bloque
				err := block.Deserialize(diskPath, int64(superblock.S_block_start+(blockIndex*superblock.S_block_size))) // 64 porque es el tamaño de un bloque
				if err != nil {
					return err
				}
				dotContent += fmt.Sprintf(`
				<tr><td colspan="2" bgcolor="#0000FF"><font color="white"> REPORTE BLOQUE %d </font></td></tr>`, blockIndex)
				//Obtinen el dot de un bloque de carpeta
				dotContent += block.ObtenerDot()
				//continue

				// Si el inodo es de tipo archivo
			} else if inode.I_type[0] == '1' {
				block := &structures.FileBlock{}
				// Deserializar el bloque
				err := block.Deserialize(diskPath, int64(superblock.S_block_start+(blockIndex*superblock.S_block_size))) // 64 porque es el tamaño de un bloque
				if err != nil {
					return err
				}
				// Obtiene el bloque
				dotContent += block.ObtenerDot()
				//continue
			}
			dotContent += "	</table>>];\n"

			// Agregar enlace al siguiente bloque si no es el último
			//fmt.Println(inode.I_block[indice+1])
			if inode.I_block[indice] != -1 {
				dotContent += fmt.Sprintf("\t\tbloque%d -> bloque%d;\n", blockIndex, blockIndex+1)
			}
		}
		//Aca cambia de inodo por lo tanto no debe de estar enlazado
	}

	// Cerrar el contenido DOT
	dotContent += "}"

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
