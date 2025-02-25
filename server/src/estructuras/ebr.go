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

// Crea la particion logica
func (ebr *EBR) CreateLogical(fit byte, start int32, size int, name string) {
	/*
		Ebr_mount [1]byte  // Indica si la partición está montada o no
		Ebr_fit   [1]byte  // Tipo de ajuste de la partición. Tendrá los valores B (Best), F(First) o W (worst)
		Ebr_start int32    // Indica en qué byte del disco inicia la partición
		Ebr_size  int32    // Contiene el tamaño total de la partición en bytes.
		Ebr_next  int32    // Byte en el que está el próximo EBR. -1 si no hay siguiente
		Ebr_name  [16]byte //Nombre de la partición*/

	// Asignar status de la partición
	ebr.Ebr_fit[0] = fit
	ebr.Ebr_start = int32(start)
	ebr.Ebr_size = int32(size)
	ebr.Ebr_next = int32(size) + int32(start)
	copy(ebr.Ebr_name[:], name)

}

// Esta funcion retorna el EBR disponible es desir que NEXT = -1
func (ebr *EBR) BuscarEBRDisponible(path string, position int32) int {
	offset := int(position) + binary.Size(ebr) // posiscion mas el tamaño del EBR en bytes
	for ebr.Ebr_next != -1 {
		err := ebr.DeserializeEBR(path, ebr.Ebr_next) //Accedemos al sigueinte ebr
		if err != nil {
			return -1 //Retornamos -1 si ay error
		}

		offset += int(ebr.Ebr_size) //Sumamos el nuevo size de la particion logica
		offset += binary.Size(ebr)  //Sumamos el nuevo EBR

	}
	return offset
}

func (ebr *EBR) BuscarEBRDisponible(path string, position int32) int32 {
	// Moverse al primer EBR en la posición indicada
	err := ebr.DeserializeEBR(path, position)
	if err != nil {
		return -1 // Retornamos -1 en caso de error
	}

	// Recorremos los EBR hasta encontrar uno donde Ebr_next == -1
	for ebr.Ebr_next != -1 {
		err := ebr.DeserializeEBR(path, ebr.Ebr_next) // Accedemos al siguiente EBR
		if err != nil {
			return -1 // Retornamos -1 si hay error
		}
	}

	// Cuando Ebr_next == -1, hemos encontrado el último EBR disponible
	return ebr.Ebr_start
}
