package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type EBR struct {
	Ebr_mount [1]byte  // Indica si la partición está montada o no
	Ebr_fit   [1]byte  // Tipo de ajuste de la partición. Tendrá los valores B (Best), F(First) o W (worst)
	Ebr_start int32    // Indica en qué byte del disco inicia la partición
	Ebr_size  int32    // Contiene el tamaño total de la partición en bytes.
	Ebr_next  int32    // Byte en el que está el próximo EBR. -1 si no hay siguiente
	Ebr_name  [16]byte //Nombre de la partición
}

// SerializeMBR escribe la estructura MBR al inicio de un archivo binario
func (ebr *EBR) SerializeEBR(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Serializar la estructura MBR directamente en el archivo
	err = binary.Write(file, binary.LittleEndian, ebr)
	if err != nil {
		return err
	}

	return nil
}

// DeserializeEBR lee la estructura EBR desde el inicio de un archivo binario
func (ebr *EBR) DeserializeEBR(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Obtener el tamaño de la estructura EBR
	ebrSize := binary.Size(ebr)
	if ebrSize <= 0 {
		return fmt.Errorf("invalid EBR size: %d", ebrSize)
	}

	// Leer solo la cantidad de bytes que corresponden al tamaño de la estructura EBR
	buffer := make([]byte, ebrSize)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	// Deserializar los bytes leídos en la estructura EBR
	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, ebr)
	if err != nil {
		return err
	}

	return nil
}
