package reportes

import (
	structures "bakend/src/estructuras"
	utils "bakend/src/utils"
	"os"
	"os/exec"
)

// ReportInode genera un reporte de un inodo y lo guarda en la ruta especificada
func ReporteSB(superblock *structures.SuperBlock, diskPath string, path string) error {
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

	//Aca se imprime el superbloque
	dotContent += superblock.ObtenerDot()

	// Cerrar el contenido DOT
	dotContent += "}"

	// Crear el archivo DOT
	dotFile, err := os.Create(dotFileName)
	if err != nil {
		return err
	}
	defer dotFile.Close()

	// Escribir el contenido DOT en el archivo
	_, err = dotFile.WriteString(dotContent)
	if err != nil {
		return err
	}

	// Generar la imagen con Graphviz
	cmd := exec.Command("dot", "-Tpng", dotFileName, "-o", outputImage)
	err = cmd.Run()
	if err != nil {
		return err
	}

	//fmt.Println("Imagen de los inodos generada:", outputImage)
	//superblock.Print()
	return nil
}
