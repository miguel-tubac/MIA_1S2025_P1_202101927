package stores

import (
	structures "bakend/src/estructuras"
	"errors"
)

// Carnet
const Carnet string = "27" // 202101927

// Declaración de variables globales
var (
	MountedPartitions map[string]string = make(map[string]string)
)

// GetMountedPartition obtiene la partición montada con el id especificado
func GetMountedPartition(id string) (*structures.PARTITION, string, error) {
	// Obtener el path de la partición montada
	path := MountedPartitions[id]
	if path == "" {
		return nil, "", errors.New("la partición no está montada")
	}

	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.DeserializeMBR(path)
	if err != nil {
		return nil, "", err
	}

	// Buscar la partición con el id especificado
	partition, err := mbr.GetPartitionByID(id)
	if partition == nil {
		return nil, "", err
	}

	return partition, path, nil
}

//Para mostrar las particiones que el comando mounted
//Solo se recorre la lista var

// Función para obtener todas las particiones montadas
func GetPartitions() map[string]string {
	return MountedPartitions
}
