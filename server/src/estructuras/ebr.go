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
func (ebr *EBR) SerializeEBR(path string, position int32) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Mover el puntero de escritura a la posición deseada
	_, err = file.Seek(int64(position), 0)
	if err != nil {
		return err
	}

	// Serializar la estructura MBR directamente en el archivo
	err = binary.Write(file, binary.LittleEndian, ebr)
	if err != nil {
		return err
	}

	return nil
}

// DeserializeEBR lee la estructura EBR desde el inicio de un archivo binario
func (ebr *EBR) DeserializeEBR(path string, position int32) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Mover el puntero de lectura a la posición deseada
	_, err = file.Seek(int64(position), 0)
	if err != nil {
		return err
	}

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

// Método para imprimir los valores del EBR
func (ebr *EBR) PrintEBR() {
	// Convertir a char cada [1]byte
	ebrMount := rune(ebr.Ebr_mount[0])
	ebrFit := rune(ebr.Ebr_fit[0])
	ebrName := string(ebr.Ebr_name[:])

	fmt.Printf("EBR mount: %c\n", ebrMount)
	fmt.Printf("EBR Fit: %c\n", ebrFit)
	fmt.Printf("EBR start: %d\n", ebr.Ebr_start)
	fmt.Printf("EBR size: %d\n", ebr.Ebr_size)
	fmt.Printf("EBR next: %d\n", ebr.Ebr_next)
	fmt.Printf("EBR name: %s\n", ebrName)
	fmt.Println()
}

// Método para obtener la primera partición disponible del ebr
func (ebr *EBR) GetFirstAvailablePartition(path string) (*EBR, int) {
	// Calcular el offset para el start de la partición
	offset := binary.Size(ebr) // Tamaño del EBR en bytes
	currentEBR := *ebr         // Copia del EBR inicial para no modificar el original

	// Recorrer las particiones del MBR
	for {
		// Si el next de la partición es -1, entonces está disponible
		if currentEBR.Ebr_next == -1 {
			// Devolver la partición, el offset y el índice
			return &currentEBR, offset
		} else {
			//Se deserealiza el sigueinte EBR
			err := currentEBR.DeserializeEBR(path, currentEBR.Ebr_next)
			if err != nil { //Validamos si no retorna error al deserealizar
				fmt.Println("Error deserializando el EBR:", err)
				return nil, -1
			}
			// Calcular el nuevo offset para la siguiente partición, es decir, sumar el tamaño de la partición
			offset += int(currentEBR.Ebr_size)
		}
	}
	//return nil, -1
}
