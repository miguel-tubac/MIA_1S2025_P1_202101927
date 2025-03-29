package structures

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type SuperBlock struct {
	S_filesystem_type   int32
	S_inodes_count      int32
	S_blocks_count      int32
	S_free_inodes_count int32
	S_free_blocks_count int32
	S_mtime             float32
	S_umtime            float32
	S_mnt_count         int32
	S_magic             int32
	S_inode_size        int32
	S_block_size        int32
	S_first_ino         int32
	S_first_blo         int32
	S_bm_inode_start    int32
	S_bm_block_start    int32
	S_inode_start       int32
	S_block_start       int32
	// Total: 68 bytes
}

// Serialize escribe la estructura SuperBlock en un archivo binario en la posición especificada
func (sb *SuperBlock) Serialize(path string, offset int64) error {
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

	// Serializar la estructura SuperBlock directamente en el archivo
	err = binary.Write(file, binary.LittleEndian, sb)
	if err != nil {
		return err
	}

	return nil
}

// Deserialize lee la estructura SuperBlock desde un archivo binario en la posición especificada
func (sb *SuperBlock) Deserialize(path string, offset int64) error {
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

	// Obtener el tamaño de la estructura SuperBlock
	sbSize := binary.Size(sb)
	if sbSize <= 0 {
		return fmt.Errorf("invalid SuperBlock size: %d", sbSize)
	}

	// Leer solo la cantidad de bytes que corresponden al tamaño de la estructura SuperBlock
	buffer := make([]byte, sbSize)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	// Deserializar los bytes leídos en la estructura SuperBlock
	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, sb)
	if err != nil {
		return err
	}

	return nil
}

// PrintSuperBlock imprime los valores de la estructura SuperBlock
func (sb *SuperBlock) Print() {
	// Convertir el tiempo de montaje a una fecha
	mountTime := time.Unix(int64(sb.S_mtime), 0)
	// Convertir el tiempo de desmontaje a una fecha
	unmountTime := time.Unix(int64(sb.S_umtime), 0)

	fmt.Printf("Filesystem Type: %d\n", sb.S_filesystem_type)
	fmt.Printf("Inodes Count: %d\n", sb.S_inodes_count)
	fmt.Printf("Blocks Count: %d\n", sb.S_blocks_count)
	fmt.Printf("Free Inodes Count: %d\n", sb.S_free_inodes_count)
	fmt.Printf("Free Blocks Count: %d\n", sb.S_free_blocks_count)
	fmt.Printf("Mount Time: %s\n", mountTime.Format(time.RFC3339))
	fmt.Printf("Unmount Time: %s\n", unmountTime.Format(time.RFC3339))
	fmt.Printf("Mount Count: %d\n", sb.S_mnt_count)
	fmt.Printf("Magic: %d\n", sb.S_magic)
	fmt.Printf("Inode Size: %d\n", sb.S_inode_size)
	fmt.Printf("Block Size: %d\n", sb.S_block_size)
	fmt.Printf("First Inode: %d\n", sb.S_first_ino)
	fmt.Printf("First Block: %d\n", sb.S_first_blo)
	fmt.Printf("Bitmap Inode Start: %d\n", sb.S_bm_inode_start)
	fmt.Printf("Bitmap Block Start: %d\n", sb.S_bm_block_start)
	fmt.Printf("Inode Start: %d\n", sb.S_inode_start)
	fmt.Printf("Block Start: %d\n", sb.S_block_start)
}

// Esta funcion retorna el codigo de dot del superbloque
func (sb *SuperBlock) ObtenerDot() string {
	// Convertir el tiempo de montaje a una fecha
	mountTime := time.Unix(int64(sb.S_mtime), 0)
	// Convertir el tiempo de desmontaje a una fecha
	unmountTime := time.Unix(int64(sb.S_umtime), 0)

	// Agregar los bloques indirectos a la tabla
	dotContent := fmt.Sprintf(`tabla [label=<
        <table border="0" cellborder="1" cellspacing="0">
			<tr><td colspan="2" bgcolor="#006400"><font color="white"> REPORTE SUPERBLOQUE</font></td></tr>
			<tr><td>S_filesystem_type</td><td>%d</td></tr>
			<tr><td bgcolor="#32CD32"><font color="white">S_inodes_count</font></td><td bgcolor="#32CD32"><font color="white">%d</font></td></tr>
			<tr><td>S_blocks_count</td><td>%d</td></tr>
			<tr><td bgcolor="#32CD32"><font color="white">S_free_blocks_count</font></td><td bgcolor="#32CD32"><font color="white">%d</font></td></tr>
			<tr><td>S_free_inodes_count</td><td>%d</td></tr>
			<tr><td bgcolor="#32CD32"><font color="white">S_mtime</font></td><td bgcolor="#32CD32"><font color="white">%s</font></td></tr>
			<tr><td>S_umtime</td><td>%s</td></tr>
			<tr><td bgcolor="#32CD32"><font color="white">S_mnt_count</font></td><td bgcolor="#32CD32"><font color="white">%d</font></td></tr>
			<tr><td>S_magic</td><td>%d</td></tr>
			<tr><td bgcolor="#32CD32"><font color="white">S_inode_s</font></td><td bgcolor="#32CD32"><font color="white">%d</font></td></tr>
			<tr><td>S_block_s</td><td>%d</td></tr>
			<tr><td bgcolor="#32CD32"><font color="white">S_firts_ino</font></td><td bgcolor="#32CD32"><font color="white">%d</font></td></tr>
			<tr><td>S_first_blo</td><td>%d</td></tr>
			<tr><td bgcolor="#32CD32"><font color="white">S_bm_inode_start</font></td><td bgcolor="#32CD32"><font color="white">%d</font></td></tr>
			<tr><td>S_bm_block_start</td><td>%d</td></tr>
			<tr><td bgcolor="#32CD32"><font color="white">S_inode_start</font></td><td bgcolor="#32CD32"><font color="white">%d</font></td></tr>
			<tr><td>S_block_start</td><td>%d</td></tr>
		</table>>];
		`, sb.S_filesystem_type, sb.S_inodes_count, sb.S_blocks_count,
		sb.S_free_blocks_count, sb.S_free_inodes_count, mountTime.Format(time.RFC3339),
		unmountTime.Format(time.RFC3339), sb.S_mnt_count, sb.S_magic, sb.S_inode_size, sb.S_block_size,
		sb.S_first_ino, sb.S_first_blo, sb.S_bm_inode_start, sb.S_bm_block_start, sb.S_inode_start, sb.S_block_start)

	return dotContent
}

// Imprimir inodos
func (sb *SuperBlock) PrintInodes(path string) error {
	// Imprimir inodos
	fmt.Println("\nInodos\n----------------")
	// Iterar sobre cada inodo
	for i := int32(0); i < sb.S_inodes_count; i++ {
		inode := &Inode{}
		// Deserializar el inodo
		err := inode.Deserialize(path, int64(sb.S_inode_start+(i*sb.S_inode_size)))
		if err != nil {
			return err
		}
		// Imprimir el inodo
		fmt.Printf("\nInodo %d:\n", i)
		inode.Print()
	}

	return nil
}

// Impriir bloques
func (sb *SuperBlock) PrintBlocks(path string) error {
	// Imprimir bloques
	fmt.Println("\nBloques\n----------------")
	// Iterar sobre cada inodo
	for i := int32(0); i < sb.S_inodes_count; i++ {
		inode := &Inode{}
		// Deserializar el inodo
		err := inode.Deserialize(path, int64(sb.S_inode_start+(i*sb.S_inode_size)))
		if err != nil {
			return err
		}
		// Iterar sobre cada bloque del inodo (apuntadores)
		for _, blockIndex := range inode.I_block {
			// Si el bloque no existe, salir
			if blockIndex == -1 {
				break
			}
			// Si el inodo es de tipo carpeta
			if inode.I_type[0] == '0' {
				block := &FolderBlock{}
				// Deserializar el bloque
				err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size))) // 64 porque es el tamaño de un bloque
				if err != nil {
					return err
				}
				// Imprimir el bloque
				fmt.Printf("\nBloque %d:\n", blockIndex)
				block.Print()
				continue

				// Si el inodo es de tipo archivo
			} else if inode.I_type[0] == '1' {
				block := &FileBlock{}
				// Deserializar el bloque
				err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size))) // 64 porque es el tamaño de un bloque
				if err != nil {
					return err
				}
				// Imprimir el bloque
				fmt.Printf("\nBloque %d:\n", blockIndex)
				block.Print()
				continue
			}

		}
	}

	return nil
}

// CreateFolder crea una carpeta en el sistema de archivos
func (sb *SuperBlock) CreateFolder(crear_padres bool, path string, parentsDir []string, destDir string) error {
	// Si parentsDir está vacío, solo trabajar con el primer inodo que sería el raíz "/"
	if len(parentsDir) == 0 {
		// enviamos a buscar el directorio a crear, para saber si existe
		next_inode, err := sb.Encontrar_Directorio(path, 0, destDir)
		// si hay un error, lo devolvemos
		if err != nil {
			return err
		}
		// si el valor de la variable es un -1, significa que el directorio no existe, por ende, hay que crearlo
		if next_inode == int32(-1) {
			//Aca se genera el indo y el fileblock
			err := sb.createFolderInInode(path, 0, destDir)

			if err != nil {
				return err
			}

		}

		return nil
	}

	// Iterar sobre cada inodo ya que se necesita buscar el inodo padre
	Posicion := int32(0)
	for i := 0; i < len(parentsDir); i++ {
		// enviamos a buscar el directorio en la posicion i, para saber si existe
		next_inode, err := sb.Encontrar_Directorio(path, Posicion, parentsDir[i])
		// si hay un error, lo devolvemos
		if err != nil {
			return err
		}

		//fmt.Println("VALOR NEXT_INODE: ", next_inode)

		// si el valor de la variable es un -1, significa que el directorio no existe, por ende, hay que crearlo
		if next_inode == int32(-1) {
			if crear_padres {
				err := sb.createFolderInInode(path, Posicion, parentsDir[i])

				if err != nil {
					return err
				}

				/*
					ya que creamos una carpeta (inode) y es el inmediato siguiente al que teniamos, simplemente sumamos 1 a Inode_destino y mantenemos la continuidad
				*/
				Posicion = sb.S_inodes_count - 1

			} else {
				return errors.New("error los directorios padres de la ruta no existe")
			}
		} else {
			/*
				si no devuelve un -1, significa que encontro el inode de la siguiente carpeta, por ende se lo asignamos a Inode_destino y mantenemos la continuidad
			*/
			Posicion = next_inode
		}
	}
	//sb.Print()

	return nil
}

// Esta funcion para buscar el directorio en donde se debe de crear el fileblok
func (sb *SuperBlock) Encontrar_Directorio(path string, inodeIndex int32, destDir string) (int32, error) {
	// Crear un nuevo inodo
	inode := &Inode{}
	// Deserializar el inodo
	err := inode.Deserialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
	if err != nil {
		return int32(-1), err
	}
	// Verificar si el inodo es de tipo carpeta
	//fmt.Println(inodeIndex)
	if inode.I_type[0] == '1' {
		// fmt.Println("aqui miguel")
		return int32(-1), errors.New("error los directorios de la ruta es un archivo")
	}

	// Iterar sobre cada bloque del inodo (apuntadores)
	for _, blockIndex := range inode.I_block {
		// Si el bloque no existe, salir
		if blockIndex == -1 {
			break
		}

		// Crear un nuevo bloque de carpeta
		block := &FolderBlock{}

		// Deserializar el bloque
		err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size))) // 64 porque es el tamaño de un bloque
		if err != nil {
			return int32(-1), err
		}

		// Iterar sobre cada contenido del bloque, desde el index 2 porque los primeros dos son . y ..
		for indexContent := 2; indexContent < len(block.B_content); indexContent++ {
			// Obtener el contenido del bloque
			content := block.B_content[indexContent]

			// Sí las carpetas padre no están vacías debereamos buscar la carpeta padre más cercana
			//fmt.Println("---------ESTOY  VISITANDO--------")

			// Si el contenido está vacío, salir
			if content.B_inodo == -1 {
				break
			}

			// Obtenemos la carpeta padre más cercana
			// parentDir, err := utils.First(parentsDir)
			// if err != nil {
			// 	return err
			// }

			// Convertir B_name a string y eliminar los caracteres nulos
			contentName := strings.Trim(string(content.B_name[:]), "\x00 ")
			// Convertir parentDir a string y eliminar los caracteres nulos
			parentDirName := strings.Trim(destDir, "\x00 ")
			fmt.Println(contentName)
			fmt.Println(parentDirName)
			// Si el nombre del contenido coincide con el nombre de la carpeta padre
			if strings.EqualFold(contentName, parentDirName) {
				return int32(content.B_inodo), nil
				//return nil
			}

		}

	}
	return int32(-1), nil
}

// CreateFile crea una archivo en el sistema de archivos
func (sb *SuperBlock) CreateFile(path string, parentsDir []string, nombreArchivo string, contenido string) error {
	// Si parentsDir está vacío, solo trabajar con el primer inodo que sería el raíz "/"
	if len(parentsDir) == 0 {
		return sb.createFileInInode(path, 0, parentsDir, nombreArchivo, contenido)
	}

	// Iterar sobre cada inodo ya que se necesita buscar el inodo padre
	for i := int32(0); i < sb.S_inodes_count; i++ {
		err := sb.createFileInInode(path, i, parentsDir, nombreArchivo, contenido)
		if err != nil {
			return err
		}
		//fmt.Println(parentsDir)
	}

	return nil
}

// CreateFolder crea una carpeta en el sistema de archivos
func (sb *SuperBlock) GetFileContent(path string, parentsDir []string, destDir string) (string, error) {
	// Si parentsDir está vacío, solo trabajar con el primer inodo que sería el raíz "/"
	if len(parentsDir) == 0 {
		return sb.getContenidoFile(path, 0, parentsDir, destDir)
	}

	// Iterar sobre cada inodo ya que se necesita buscar el inodo padre
	for i := int32(0); i < sb.S_inodes_count; i++ {
		cade, err := sb.getContenidoFile(path, i, parentsDir, destDir)
		if err != nil {
			return "", err
		}
		//Si la cadnea es distinta de cadena vasilla
		if cade != "" {
			return cade, nil
		}
	}
	//sb.Print()

	return "", nil
}
