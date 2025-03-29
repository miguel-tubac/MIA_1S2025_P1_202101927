package structures

import (
	"bytes"           // Paquete para manipulación de buffers
	"encoding/binary" // Paquete para codificación y decodificación de datos binarios
	"errors"
	"fmt" // Paquete para formateo de E/S
	"os"  // Paquete para funciones del sistema operativo
	"strings"
	"time"
)

type MBR struct {
	Mbr_size           int32        // Tamaño del MBR en bytes
	Mbr_creation_date  float32      // Fecha y hora de creación del MBR
	Mbr_disk_signature int32        // Firma del disco
	Mbr_disk_fit       [1]byte      // Tipo de ajuste
	Mbr_partitions     [4]PARTITION // Particiones del MBR
}

// SerializeMBR escribe la estructura MBR al inicio de un archivo binario
func (mbr *MBR) SerializeMBR(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Serializar la estructura MBR directamente en el archivo
	err = binary.Write(file, binary.LittleEndian, mbr)
	if err != nil {
		return err
	}

	return nil
}

// DeserializeMBR lee la estructura MBR desde el inicio de un archivo binario
func (mbr *MBR) DeserializeMBR(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Obtener el tamaño de la estructura MBR
	mbrSize := binary.Size(mbr)
	if mbrSize <= 0 {
		return fmt.Errorf("invalid MBR size: %d", mbrSize)
	}

	// Leer solo la cantidad de bytes que corresponden al tamaño de la estructura MBR
	buffer := make([]byte, mbrSize)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	// Deserializar los bytes leídos en la estructura MBR
	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, mbr)
	if err != nil {
		return err
	}

	return nil
}

// Método para obtener la primera partición disponible
func (mbr *MBR) GetFirstAvailablePartition() (*PARTITION, int, int) {
	// Calcular el offset para el start de la partición
	offset := binary.Size(mbr) // Tamaño del MBR en bytes

	// Recorrer las particiones del MBR
	for i := 0; i < len(mbr.Mbr_partitions); i++ {
		// Si el start de la partición es -1, entonces está disponible
		if mbr.Mbr_partitions[i].Part_start == -1 && mbr.Mbr_size >= int32(offset) {
			// Devolver la partición, el offset y el índice
			return &mbr.Mbr_partitions[i], offset, i
		} else {
			// Calcular el nuevo offset para la siguiente partición, es decir, sumar el tamaño de la partición
			offset += int(mbr.Mbr_partitions[i].Part_size)
		}
	}
	return nil, -1, -1
}

// Funcion para revisar si existe otra particion con el mismo nombre
func (mbr *MBR) ExisteNombre(name string) bool {
	// Recorrer las particiones del MBR
	for _, partition := range mbr.Mbr_partitions {
		// Convertir Part_name a string y eliminar los caracteres nulos
		partitionName := strings.Trim(string(partition.Part_name[:]), "\x00 ")
		// Convertir el nombre de la partición a string y eliminar los caracteres nulos
		inputName := strings.Trim(name, "\x00 ")
		// Si el nombre de la partición coincide, devolver la partición y el índice
		if strings.EqualFold(partitionName, inputName) {
			fmt.Printf("El nombre %s ya existe en las particiones\n", name)
			return false
		}
	}
	return true
}

// Método para obtener una partición por nombre
func (mbr *MBR) GetPartitionByName(name string) (*PARTITION, int) {
	// Recorrer las particiones del MBR
	for i, partition := range mbr.Mbr_partitions {
		// Convertir Part_name a string y eliminar los caracteres nulos
		partitionName := strings.Trim(string(partition.Part_name[:]), "\x00 ")
		// Convertir el nombre de la partición a string y eliminar los caracteres nulos
		inputName := strings.Trim(name, "\x00 ")
		// Si el nombre de la partición coincide, devolver la partición y el índice
		if strings.EqualFold(partitionName, inputName) {
			if partition.Part_status[0] == '1' {
				return nil, -1
			}
			return &partition, i
		}
	}
	return nil, -1
}

// Método para imprimir los valores del MBR
func (mbr *MBR) PrintMBR() {
	// Convertir Mbr_creation_date a time.Time
	creationTime := time.Unix(int64(mbr.Mbr_creation_date), 0)

	// Convertir Mbr_disk_fit a char
	diskFit := rune(mbr.Mbr_disk_fit[0])

	fmt.Printf("MBR Size: %d\n", mbr.Mbr_size)
	fmt.Printf("Creation Date: %s\n", creationTime.Format(time.RFC3339))
	fmt.Printf("Disk Signature: %d\n", mbr.Mbr_disk_signature)
	fmt.Printf("Disk Fit: %c\n", diskFit)
}

// Método para imprimir las particiones del MBR
func (mbr *MBR) PrintPartitions() {
	for i, partition := range mbr.Mbr_partitions {
		// Convertir Part_status, Part_type y Part_fit a char
		partStatus := rune(partition.Part_status[0])
		partType := rune(partition.Part_type[0])
		partFit := rune(partition.Part_fit[0])

		// Convertir Part_name a string
		partName := string(partition.Part_name[:])
		// Convertir Part_id a string
		partID := string(partition.Part_id[:])

		fmt.Printf("Partition %d \n", i+1)
		fmt.Printf("  Status: %c\n", partStatus)
		fmt.Printf("  Type: %c\n", partType)
		fmt.Printf("  Fit: %c\n", partFit)
		fmt.Printf("  Start: %d\n", partition.Part_start)
		fmt.Printf("  Size: %d\n", partition.Part_size)
		fmt.Printf("  Name: %s\n", partName)
		fmt.Printf("  Correlative: %d\n", partition.Part_correlative)
		fmt.Printf("  ID: %s\n", partID)
		fmt.Println()
	}
}

// funcion para obtener el correlaticvo de la paricion
func (mbr *MBR) GetCorrelativo() int32 {
	var corr int32
	corr = -1

	// Recorrer las particiones del MBR
	for _, partition := range mbr.Mbr_partitions {
		corr2 := int32(partition.Part_correlative)

		if corr2 > 0 {
			corr = corr2
		}
	}
	return corr //Rernamos el ultimo valor o -1 si no hay ninguno
}

// Funcion para revisar si ya existe una particion extendida
func (mbr *MBR) GetExtendedPartition() bool {
	// Recorrer las particiones del MBR
	for _, partition := range mbr.Mbr_partitions {
		// Convertir Part_name a string y eliminar los caracteres nulos
		partitionType := strings.Trim(string(partition.Part_type[:]), "\x00 ")

		// Si el nombre de la partición coincide, devolver la partición y el índice
		if strings.EqualFold(partitionType, "E") {
			//fmt.Println("Ya existe una particion extendida")
			return true
		}
	}
	return false
}

// Funcion para obtener la particion extendida
func (mbr *MBR) GetExtendedPartition2() *PARTITION {
	// Recorrer las particiones del MBR
	for _, partition := range mbr.Mbr_partitions {
		// Convertir Part_name a string y eliminar los caracteres nulos
		partitionType := strings.Trim(string(partition.Part_type[:]), "\x00 ")
		//fmt.Println("Aqui 2:")
		//fmt.Println(partitionType)
		// Si el nombre de la partición coincide, devolver la partición y el índice
		if strings.EqualFold(partitionType, "E") {
			//fmt.Println("Ya existe una particion extendida")
			return &partition
		}
	}
	return nil
}

// Función para obtener una partición por ID
func (mbr *MBR) GetPartitionByID(id string) (*PARTITION, error) {
	for i := 0; i < len(mbr.Mbr_partitions); i++ {
		// Convertir Part_name a string y eliminar los caracteres nulos
		partitionID := strings.Trim(string(mbr.Mbr_partitions[i].Part_id[:]), "\x00 ")
		// Convertir el id a string y eliminar los caracteres nulos
		inputID := strings.Trim(id, "\x00 ")
		// Si el nombre de la partición coincide, devolver la partición
		if strings.EqualFold(partitionID, inputID) {
			return &mbr.Mbr_partitions[i], nil
		}
	}
	return nil, errors.New("partición no encontrada")
}
