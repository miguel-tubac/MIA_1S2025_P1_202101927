package reportes

import (
	structures "bakend/src/estructuras"
	utils "bakend/src/utils"
	"fmt"
	"os"
	"os/exec"
)

func ReporteLs(superblock *structures.SuperBlock, diskPath string, path string, path_file_ls string) error {
	// Crear las carpetas padre si no existen
	err := utils.CreateParentDirs(path)
	if err != nil {
		return err
	}

	// Obtener el nombre base del archivo sin la extensión
	// Estos son los nombres del archivo a crear y la ruta de la computadora
	dotFileName, outputImage := utils.GetFileNames(path)

	// GetParentDirectories obtiene las carpetas padres y el directorio de destino
	parentDirs, nombreArchivo := utils.GetParentDirectories(path_file_ls)
	// fmt.Println("Directorios padres:", parentDirs)
	// fmt.Println("Nombre archivo:", nombreArchivo)

	// Iniciar el contenido DOT
	dotContent := `digraph G {
		node [shape=plaintext]

		tabla [
        	label=<
        	<table border="1" cellborder="1" cellspacing="0">
			<tr>
                <td>Permisos</td><td>Propietario</td><td>Grupo propietario</td><td>Fecha de modificación</td><td>Tipo</td><td>Fecha de creación</td><td>Nombre</td>
            </tr>
		`

	//Se obtiene el codigo de tipo dot:
	cade, err2 := superblock.ObtenerDotLS(diskPath, parentDirs, nombreArchivo)
	if err2 != nil {
		return fmt.Errorf("error al obtener el archivo: %w", err2)
	}

	//Se agrega el codigo del .dot al codigo inicial
	dotContent += cade

	// Cerrar el contenido DOT
	dotContent += `
		</table>
        >];
	}`

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

	//fmt.Println("Archivo creado exitosamente en:", outputImage)
	return nil
}
