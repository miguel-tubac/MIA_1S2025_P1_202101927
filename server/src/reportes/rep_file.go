package reportes

import (
	structures "bakend/src/estructuras"
	utils "bakend/src/utils"
	"fmt"
	"os"
)

func ReporteFile(superblock *structures.SuperBlock, diskPath string, path string, path_file_ls string) error {
	// Crear las carpetas padre si no existen
	err := utils.CreateParentDirs(path)
	if err != nil {
		return err
	}

	// Obtener el nombre base del archivo sin la extensión
	// Estos son los nombres del archivo a crear y la ruta de la computadora
	_, outputImage := utils.GetFileNames(path)

	// GetParentDirectories obtiene las carpetas padres y el directorio de destino
	parentDirs, nombreArchivo := utils.GetParentDirectories(path_file_ls)
	// fmt.Println("Directorios padres:", parentDirs)
	// fmt.Println("Nombre archivo:", nombreArchivo)

	//Aca ya se debe de obtner el archivo en el disco virtual con el -path
	cade, err2 := superblock.GetFileContent(diskPath, parentDirs, nombreArchivo)
	if err2 != nil {
		return fmt.Errorf("error al obtener el archivo: %w", err2)
	}

	// Crear y escribir el archivo en la ruta especificada
	err = os.WriteFile(outputImage, []byte(cade), 0644)
	if err != nil {
		return fmt.Errorf("error al escribir el archivo: %w", err)
	}

	//fmt.Println("Archivo creado exitosamente en:", outputImage)
	return nil
}
