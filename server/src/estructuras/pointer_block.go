package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type PointerBlock struct {
	P_pointers [16]int32 // 16 * 4 = 64 bytes
	// Total: 64 bytes
}

// Serialize escribe la estructura FileBlock en un archivo binario en la posición especificada
func (fb *PointerBlock) Serialize(path string, offset int64) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Mover el puntero del archivo a la posición especificada
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Serializar la estructura PointerBlock directamente en el archivo
	err = binary.Write(file, binary.LittleEndian, fb)
	if err != nil {
		return err
	}

	return nil
}

// Deserialize lee la estructura FileBlock desde un archivo binario en la posición especificada
func (fb *PointerBlock) Deserialize(path string, offset int64) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Mover el puntero del archivo a la posición especificada
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Obtener el tamaño de la estructura FileBlock
	fbSize := binary.Size(fb)
	if fbSize <= 0 {
		return fmt.Errorf("invalid PointerBlock size: %d", fbSize)
	}

	// Leer solo la cantidad de bytes que corresponden al tamaño de la estructura FileBlock
	buffer := make([]byte, fbSize)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	// Deserializar los bytes leídos en la estructura FileBlock
	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, fb)
	if err != nil {
		return err
	}

	return nil
}

// PrintContent prints the content of P_pointers as a string
func (fb *PointerBlock) Print() {
	fmt.Printf("%v\n", fb.P_pointers) // Imprime como un array de enteros
}
