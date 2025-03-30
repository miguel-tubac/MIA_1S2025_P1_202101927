package structures

import (
	utils "bakend/src/utils"
	"errors"
	"strings"
	"time"
)

// Crear users.txt en nuestro sistema de archivos
func (sb *SuperBlock) CreateUsersFile(path string) error {
	// ----------- Creamos / -----------
	// Creamos el inodo raíz
	rootInode := &Inode{
		I_uid:   1,
		I_gid:   1,
		I_size:  0,
		I_atime: float32(time.Now().Unix()),
		I_ctime: float32(time.Now().Unix()),
		I_mtime: float32(time.Now().Unix()),
		I_block: [15]int32{sb.S_blocks_count, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		I_type:  [1]byte{'0'},
		I_perm:  [3]byte{'7', '7', '7'},
	}

	// Serializar el inodo raíz
	err := rootInode.Serialize(path, int64(sb.S_first_ino))
	if err != nil {
		return err
	}

	// Actualizar el bitmap de inodos
	err = sb.UpdateBitmapInode(path)
	if err != nil {
		return err
	}

	// Actualizar el superbloque
	sb.S_inodes_count++
	sb.S_free_inodes_count--
	sb.S_first_ino += sb.S_inode_size

	// Creamos el bloque del Inodo Raíz
	rootBlock := &FolderBlock{
		B_content: [4]FolderContent{
			{B_name: [12]byte{'.'}, B_inodo: 0},
			{B_name: [12]byte{'.', '.'}, B_inodo: 0},
			{B_name: [12]byte{'-'}, B_inodo: -1},
			{B_name: [12]byte{'-'}, B_inodo: -1},
		},
	}

	// Actualizar el bitmap de bloques
	err = sb.UpdateBitmapBlock(path)
	if err != nil {
		return err
	}

	// Serializar el bloque de carpeta raíz
	err = rootBlock.Serialize(path, int64(sb.S_first_blo))
	if err != nil {
		return err
	}

	// Actualizar el superbloque
	sb.S_blocks_count++
	sb.S_free_blocks_count--
	sb.S_first_blo += sb.S_block_size

	// ----------- Creamos /users.txt -----------
	usersText := "1,G,root\n1,U,root,root,123\n"

	// Deserializar el inodo raíz
	err = rootInode.Deserialize(path, int64(sb.S_inode_start+0)) // 0 porque es el inodo raíz
	if err != nil {
		return err
	}

	// Actualizamos el inodo raíz
	rootInode.I_atime = float32(time.Now().Unix())

	// Serializar el inodo raíz
	err = rootInode.Serialize(path, int64(sb.S_inode_start+0)) // 0 porque es el inodo raíz
	if err != nil {
		return err
	}

	// Deserializar el bloque de carpeta raíz
	err = rootBlock.Deserialize(path, int64(sb.S_block_start+0)) // 0 porque es el bloque de carpeta raíz
	if err != nil {
		return err
	}

	// Actualizamos el bloque de carpeta raíz
	rootBlock.B_content[2] = FolderContent{B_name: [12]byte{'u', 's', 'e', 'r', 's', '.', 't', 'x', 't'}, B_inodo: sb.S_inodes_count}

	// Serializar el bloque de carpeta raíz
	err = rootBlock.Serialize(path, int64(sb.S_block_start+0)) // 0 porque es el bloque de carpeta raíz
	if err != nil {
		return err
	}

	// Creamos el inodo users.txt
	usersInode := &Inode{
		I_uid:   1,
		I_gid:   1,
		I_size:  int32(len(usersText)),
		I_atime: float32(time.Now().Unix()),
		I_ctime: float32(time.Now().Unix()),
		I_mtime: float32(time.Now().Unix()),
		I_block: [15]int32{sb.S_blocks_count, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		I_type:  [1]byte{'1'},
		I_perm:  [3]byte{'7', '7', '7'},
	}

	// Actualizar el bitmap de inodos
	err = sb.UpdateBitmapInode(path)
	if err != nil {
		return err
	}

	// Serializar el inodo users.txt
	err = usersInode.Serialize(path, int64(sb.S_first_ino))
	if err != nil {
		return err
	}

	// Actualizamos el superbloque
	sb.S_inodes_count++
	sb.S_free_inodes_count--
	sb.S_first_ino += sb.S_inode_size

	// Creamos el bloque de users.txt
	usersBlock := &FileBlock{
		B_content: [64]byte{},
	}
	// Copiamos el texto de usuarios en el bloque
	copy(usersBlock.B_content[:], usersText)

	// Serializar el bloque de users.txt
	err = usersBlock.Serialize(path, int64(sb.S_first_blo))
	if err != nil {
		return err
	}

	// Actualizar el bitmap de bloques
	err = sb.UpdateBitmapBlock(path)
	if err != nil {
		return err
	}

	// Actualizamos el superbloque
	sb.S_blocks_count++
	sb.S_free_blocks_count--
	sb.S_first_blo += sb.S_block_size

	return nil
}

// createFolderInInode crea una carpeta en un inodo específico
func (sb *SuperBlock) createFolderInInode(path string, inodeIndex int32, destDir string) error {
	// Crear un nuevo inodo
	inode := &Inode{}
	// Deserializar el inodo
	err := inode.Deserialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
	if err != nil {
		return err
	}

	// Crear un nuevo bloque de carpeta
	block := &FolderBlock{}

	// Deserializar el bloque
	err = block.Deserialize(path, int64(sb.S_block_start+(inode.I_block[0]*sb.S_block_size))) // 64 porque es el tamaño de un bloque
	if err != nil {
		return err
	}

	// obtenemos la posicion 2 del bloque (la que contiene el apuntador al padre)
	first_content := block.B_content[1]

	// obtenemos el numero de inode del padre
	father_inode := first_content.B_inodo

	// Iterar sobre cada bloque del inodo (apuntadores)
	for i := 0; i < len(inode.I_block); i++ {
		// Si el bloque no existe, se debe de crear
		if inode.I_block[i] == -1 {
			inode.I_block[i] = sb.S_blocks_count

			new_folderblock := &FolderBlock{
				B_content: [4]FolderContent{
					{B_name: [12]byte{'.'}, B_inodo: inodeIndex},
					{B_name: [12]byte{'.', '.'}, B_inodo: father_inode},
					{B_name: [12]byte{'-'}, B_inodo: sb.S_blocks_count},
					{B_name: [12]byte{'-'}, B_inodo: -1},
				},
			}

			// creamos una instancia de foldercontent
			content_folder := &FolderContent{}

			// obtenemos el tercer campo del mismo (el primer espacio disponible)
			content_folder = &new_folderblock.B_content[2]

			// copiamos el nombre del archivo
			copy(content_folder.B_name[:], []byte(destDir))

			// Actualizar el bitmap de inodos
			err = sb.UpdateBitmapBlock(path)
			if err != nil {
				return err
			}

			// serializar el bloque
			err = new_folderblock.Serialize(path, int64(sb.S_first_blo))
			if err != nil {
				return err
			}

			// Actualizar el superbloque
			sb.S_blocks_count++
			sb.S_free_blocks_count--
			sb.S_first_blo += sb.S_block_size

			// Actualizamos el inode root
			inode.I_atime = float32(time.Now().Unix())

			// serializamos el inode root
			err = inode.Serialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
			if err != nil {
				return err
			}

			/*
				una vez creado el nuevo bloque y asignado al inode
				creamos el nuevo inode y lo guardaremos en el bloque
			*/

			new_Inode := &Inode{
				I_uid:   1,
				I_gid:   1,
				I_size:  0,
				I_atime: float32(time.Now().Unix()),
				I_ctime: float32(time.Now().Unix()),
				I_mtime: float32(time.Now().Unix()),
				I_block: [15]int32{sb.S_blocks_count, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
				I_type:  [1]byte{'0'},
				I_perm:  [3]byte{'6', '6', '4'},
			}

			//Serializa el superbloque
			err = new_Inode.Serialize(path, int64(sb.S_first_ino))
			if err != nil {
				return err
			}

			// actualizamos el bitmap de inodes
			err = sb.UpdateBitmapInode(path)
			if err != nil {
				return err
			}

			// Actualizar el superbloque
			sb.S_inodes_count++
			sb.S_free_inodes_count--
			sb.S_first_ino += sb.S_inode_size

			// creamos el bloque del inodo recien creado
			newInode_folderblock := &FolderBlock{
				B_content: [4]FolderContent{
					{B_name: [12]byte{'.'}, B_inodo: sb.S_inodes_count - 1},
					{B_name: [12]byte{'.', '.'}, B_inodo: inodeIndex},
					{B_name: [12]byte{'-'}, B_inodo: -1},
					{B_name: [12]byte{'-'}, B_inodo: -1},
				},
			}

			// serializamos el bloque de la carpeta
			err = newInode_folderblock.Serialize(path, int64(sb.S_first_blo))
			if err != nil {
				return err
			}

			// actualizar el bitmap de bloques
			err = sb.UpdateBitmapBlock(path)
			if err != nil {
				return err
			}

			// Actualizar el superbloque
			sb.S_blocks_count++
			sb.S_free_blocks_count--
			sb.S_first_blo += sb.S_block_size

			return nil
		}

		//De lo contrario continua con la creacion en la posicion dada
		// Crear un nuevo bloque de carpeta
		block := &FolderBlock{}

		// Deserializar el bloque
		err := block.Deserialize(path, int64(sb.S_block_start+(inode.I_block[0]*sb.S_block_size))) // 64 porque es el tamaño de un bloque
		if err != nil {
			return err
		}

		// Iterar sobre cada contenido del bloque, desde el index 2 porque los primeros dos son . y ..
		for indexContent := 2; indexContent < len(block.B_content); indexContent++ {
			// Obtener el contenido del bloque
			content := block.B_content[indexContent]

			//fmt.Println("---------ESTOY  CREANDO--------")

			// Si el apuntador al inodo está ocupado, continuar con el siguiente
			if content.B_inodo != -1 {
				continue
			}

			// Actualizar el contenido del bloque
			copy(content.B_name[:], destDir)
			content.B_inodo = sb.S_inodes_count

			// Actualizar el bloque
			block.B_content[indexContent] = content

			// Serializar el bloque
			err = block.Serialize(path, int64(sb.S_block_start+(inode.I_block[0]*sb.S_block_size)))
			if err != nil {
				return err
			}

			// Crear el inodo de la carpeta
			folderInode := &Inode{
				I_uid:   1,
				I_gid:   1,
				I_size:  0,
				I_atime: float32(time.Now().Unix()),
				I_ctime: float32(time.Now().Unix()),
				I_mtime: float32(time.Now().Unix()),
				I_block: [15]int32{sb.S_blocks_count, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
				I_type:  [1]byte{'0'},
				I_perm:  [3]byte{'6', '6', '4'},
			}

			// Serializar el inodo de la carpeta
			err = folderInode.Serialize(path, int64(sb.S_first_ino))
			if err != nil {
				return err
			}

			// Actualizar el bitmap de inodos
			err = sb.UpdateBitmapInode(path)
			if err != nil {
				return err
			}

			// Actualizar el superbloque
			sb.S_inodes_count++
			sb.S_free_inodes_count--
			sb.S_first_ino += sb.S_inode_size

			// Crear el bloque de la carpeta
			folderBlock := &FolderBlock{
				B_content: [4]FolderContent{
					{B_name: [12]byte{'.'}, B_inodo: content.B_inodo},
					{B_name: [12]byte{'.', '.'}, B_inodo: inodeIndex},
					{B_name: [12]byte{'-'}, B_inodo: -1},
					{B_name: [12]byte{'-'}, B_inodo: -1},
				},
			}

			// Serializar el bloque de la carpeta
			err = folderBlock.Serialize(path, int64(sb.S_first_blo))
			if err != nil {
				return err
			}

			// Actualizar el bitmap de bloques
			err = sb.UpdateBitmapBlock(path)
			if err != nil {
				return err
			}

			// Actualizar el superbloque
			sb.S_blocks_count++
			sb.S_free_blocks_count--
			sb.S_first_blo += sb.S_block_size

			return nil

		}

	}
	return errors.New("error no se pudo obtener el bloque del inodo")
}

// createFileInInode crea una archivo en un inodo específico
func (sb *SuperBlock) createFileInInode(path string, inodeIndex int32, nombreArchivo string, contenido string) error {
	// Crear un nuevo inodo
	inode := &Inode{}
	// Deserializar el inodo
	err := inode.Deserialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
	if err != nil {
		return err
	}

	// creamos un instancia de folderblock
	firstblock := &FolderBlock{}

	// deserializamos el contenido del primer bloque
	err = firstblock.Deserialize(path, int64(sb.S_block_start+(inode.I_block[0]*sb.S_block_size)))
	if err != nil {
		return err
	}

	// obtenemos la posicion 2 del bloque (la que contiene el apuntador al padre)
	first_content := firstblock.B_content[1]

	// obtenemos el numero de inode del padre
	father_inode := first_content.B_inodo

	// Iterar sobre cada bloque del inodo (apuntadores)
	for i := 0; i < len(inode.I_block); i++ {
		// Si el bloque no existe, salir
		if inode.I_block[i] == -1 {
			inode.I_block[i] = sb.S_blocks_count

			new_folderblock := &FolderBlock{
				B_content: [4]FolderContent{
					{B_name: [12]byte{'.'}, B_inodo: inodeIndex},
					{B_name: [12]byte{'.', '.'}, B_inodo: father_inode},
					{B_name: [12]byte{'-'}, B_inodo: sb.S_inodes_count},
					{B_name: [12]byte{'-'}, B_inodo: -1},
				},
			}

			// creamos una instancia de foldercontent
			content_folder := &FolderContent{}

			// obtenemos el tercer campo del mismo (el primer espacio disponible)
			content_folder = &new_folderblock.B_content[2]

			// copiamos el nombre del archivo
			copy(content_folder.B_name[:], []byte(nombreArchivo))

			// actualizamos el bitmap de bloques
			err = sb.UpdateBitmapBlock(path)
			if err != nil {
				return err
			}

			// serializar el bloque
			err = new_folderblock.Serialize(path, int64(sb.S_first_blo))
			if err != nil {
				return err
			}

			// actualizamos los campos del superbloque
			sb.S_blocks_count++
			sb.S_free_blocks_count--
			sb.S_first_blo += sb.S_block_size

			// Actualizamos el inode root
			inode.I_atime = float32(time.Now().Unix())

			// serializamos el inode root
			err = inode.Serialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
			if err != nil {
				return err
			}

			/*
				una vez creado el nuevo bloque y asignado al inode
				creamos el nuevo inode y lo guardaremos en el bloque
			*/

			new_Inode := &Inode{
				I_uid:   1,
				I_gid:   1,
				I_size:  int32(len(contenido)),
				I_atime: float32(time.Now().Unix()),
				I_ctime: float32(time.Now().Unix()),
				I_mtime: float32(time.Now().Unix()),
				I_block: [15]int32{sb.S_blocks_count, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
				I_type:  [1]byte{'1'},
				I_perm:  [3]byte{'6', '6', '4'},
			}

			// Actualizar el bitmap de inodos
			err = sb.UpdateBitmapInode(path)
			if err != nil {
				return err
			}

			//Guarda la posicion donde se serealizo el inodo
			posicionInodo := int64(sb.S_first_ino)
			// Serializar el inodo users.txt
			err = new_Inode.Serialize(path, int64(sb.S_first_ino))
			if err != nil {
				return err
			}

			// Actualizamos el superbloque
			sb.S_inodes_count++
			sb.S_free_inodes_count--
			sb.S_first_ino += sb.S_inode_size

			//Aca se valida si el contenido es demasiado grande para solo un FileBlock
			if len(contenido) > 64 {
				tamano := 64
				contador := 0
				// Recorrer el string en segmentos de 64 caracteres
				for i := 0; i < len(contenido); i += tamano {
					//fmt.Println("******************")
					fin := i + tamano
					if fin > len(contenido) {
						fin = len(contenido) // Evitar desbordamiento
					}

					//TODO: validar cuando el contenido deba de aplicar los bloques indirectos
					// Bloque Simple Indirecto: Inodo → Bloque apuntadores → bloque de datos
					// if contador == 12 {

					// 	// Bloque Doble Indirecto: Inodo → Bloque de apuntadores → Bloquede apuntadores → bloque de datos.
					// } else if contador == 13 {

					// 	// Bloque Triple Indirecto: Inodo → Bloque de apuntadores → Bloque
					// 	// de apuntadores → Bloque de apuntadores → bloque de datos.
					// } else if contador == 14 {

					// } else {
					//Se actuliza los apuntadores del inodo
					inode.I_block[i]++
					new_Inode.I_block[contador] = inode.I_block[i]
					//Serealizar el ibloque actualizado
					// Deserializar el inodo
					//err := usersInode.Serialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
					err := new_Inode.Serialize(path, posicionInodo)
					if err != nil {
						return err
					}

					parte := contenido[i:fin]
					// Creamos el bloque de users.txt
					usersBlock := &FileBlock{
						B_content: [64]byte{},
					}
					// Copiamos el texto de usuarios en el bloque
					copy(usersBlock.B_content[:], []byte(parte))

					// Serializar el bloque de users.txt
					err = usersBlock.Serialize(path, int64(sb.S_first_blo))
					if err != nil {
						return err
					}

					// Actualizar el bitmap de bloques
					err = sb.UpdateBitmapBlock(path)
					if err != nil {
						return err
					}

					// Actualizamos el superbloque
					sb.S_blocks_count++
					sb.S_free_blocks_count--
					sb.S_first_blo += sb.S_block_size
					//fmt.Println(parte)
					//}
					//Se incrementa en una unidad el contador
					contador++
				}
				//Aca es unicamente un solo FileBlock
			} else {
				//fmt.Println("aqui miguel")
				// Creamos el bloque de users.txt
				usersBlock := &FileBlock{
					B_content: [64]byte{},
				}
				// Copiamos el texto de usuarios en el bloque
				copy(usersBlock.B_content[:], []byte(contenido))

				// Serializar el bloque de users.txt
				err = usersBlock.Serialize(path, int64(sb.S_first_blo))
				if err != nil {
					return err
				}

				// Actualizar el bitmap de bloques
				err = sb.UpdateBitmapBlock(path)
				if err != nil {
					return err
				}

				// Actualizamos el superbloque
				sb.S_blocks_count++
				sb.S_free_blocks_count--
				sb.S_first_blo += sb.S_block_size
			}

			return nil
		}

		//*********************************************************Aca es cuando ya esta creado
		// Crear un nuevo bloque de carpeta
		block := &FolderBlock{}

		// Deserializar el bloque
		err := block.Deserialize(path, int64(sb.S_block_start+(inode.I_block[i]*sb.S_block_size))) // 64 porque es el tamaño de un bloque
		if err != nil {
			return err
		}

		// Iterar sobre cada contenido del bloque, desde el index 2 porque los primeros dos son . y ..
		for indexContent := 2; indexContent < len(block.B_content); indexContent++ {
			// Obtener el contenido del bloque
			content := block.B_content[indexContent]

			//fmt.Println("---------ESTOY  CREANDO--------")

			// Si el apuntador al inodo está ocupado, continuar con el siguiente
			if content.B_inodo != -1 {
				continue
			}

			// Actualizar el contenido del bloque
			copy(content.B_name[:], nombreArchivo)
			content.B_inodo = sb.S_inodes_count

			// Actualizar el bloque
			block.B_content[indexContent] = content

			// Serializar el bloque
			err = block.Serialize(path, int64(sb.S_block_start+(inode.I_block[i]*sb.S_block_size)))
			if err != nil {
				return err
			}

			// Creamos el inodo nombre.txt
			usersInode := &Inode{
				I_uid:   1,
				I_gid:   1,
				I_size:  int32(len(contenido)),
				I_atime: float32(time.Now().Unix()),
				I_ctime: float32(time.Now().Unix()),
				I_mtime: float32(time.Now().Unix()),
				I_block: [15]int32{sb.S_blocks_count, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
				I_type:  [1]byte{'1'},
				I_perm:  [3]byte{'6', '6', '4'},
			}

			// Actualizar el bitmap de inodos
			err = sb.UpdateBitmapInode(path)
			if err != nil {
				return err
			}

			//Guarda la posicion donde se serealizo el inodo
			posicionInodo := int64(sb.S_first_ino)
			// Serializar el inodo users.txt
			err = usersInode.Serialize(path, int64(sb.S_first_ino))
			if err != nil {
				return err
			}

			// Actualizamos el superbloque
			sb.S_inodes_count++
			sb.S_free_inodes_count--
			sb.S_first_ino += sb.S_inode_size

			//Aca se valida si el contenido es demasiado grande para solo un FileBlock
			if len(contenido) > 64 {
				tamano := 64
				contador := 0
				// Recorrer el string en segmentos de 64 caracteres
				for i := 0; i < len(contenido); i += tamano {
					//fmt.Println("******************")
					fin := i + tamano
					if fin > len(contenido) {
						fin = len(contenido) // Evitar desbordamiento
					}

					//TODO: validar cuando el contenido deba de aplicar los bloques indirectos
					// Bloque Simple Indirecto: Inodo → Bloque apuntadores → bloque de datos
					// if contador == 12 {

					// 	// Bloque Doble Indirecto: Inodo → Bloque de apuntadores → Bloquede apuntadores → bloque de datos.
					// } else if contador == 13 {

					// 	// Bloque Triple Indirecto: Inodo → Bloque de apuntadores → Bloque
					// 	// de apuntadores → Bloque de apuntadores → bloque de datos.
					// } else if contador == 14 {

					// } else {
					//Se actuliza los apuntadores del inodo
					inode.I_block[i]++
					usersInode.I_block[contador] = inode.I_block[i]
					//Serealizar el ibloque actualizado
					// Deserializar el inodo
					//err := usersInode.Serialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
					err := usersInode.Serialize(path, posicionInodo)
					if err != nil {
						return err
					}

					parte := contenido[i:fin]
					// Creamos el bloque de users.txt
					usersBlock := &FileBlock{
						B_content: [64]byte{},
					}
					// Copiamos el texto de usuarios en el bloque
					copy(usersBlock.B_content[:], []byte(parte))

					// Serializar el bloque de users.txt
					err = usersBlock.Serialize(path, int64(sb.S_first_blo))
					if err != nil {
						return err
					}

					// Actualizar el bitmap de bloques
					err = sb.UpdateBitmapBlock(path)
					if err != nil {
						return err
					}

					// Actualizamos el superbloque
					sb.S_blocks_count++
					sb.S_free_blocks_count--
					sb.S_first_blo += sb.S_block_size
					//fmt.Println(parte)
					//}
					//Se incrementa en una unidad el contador
					contador++
				}
				//Aca es unicamente un solo FileBlock
			} else {
				//fmt.Println("aqui miguel")
				// Creamos el bloque de users.txt
				usersBlock := &FileBlock{
					B_content: [64]byte{},
				}
				// Copiamos el texto de usuarios en el bloque
				copy(usersBlock.B_content[:], []byte(contenido))

				// Serializar el bloque de users.txt
				err = usersBlock.Serialize(path, int64(sb.S_first_blo))
				if err != nil {
					return err
				}

				// Actualizar el bitmap de bloques
				err = sb.UpdateBitmapBlock(path)
				if err != nil {
					return err
				}

				// Actualizamos el superbloque
				sb.S_blocks_count++
				sb.S_free_blocks_count--
				sb.S_first_blo += sb.S_block_size
			}

			//No retornamos nada
			return nil
		}

	}
	return nil
}

// createFolderInInode crea una carpeta en un inodo específico
func (sb *SuperBlock) getContenidoFile(path string, inodeIndex int32, parentsDir []string, destDir string) (string, error) {
	// Crear un nuevo inodo
	inode := &Inode{}
	// Deserializar el inodo
	err := inode.Deserialize(path, int64(sb.S_inode_start+(inodeIndex*sb.S_inode_size)))
	if err != nil {
		return "", err
	}
	// Verificar si el inodo es de tipo carpeta
	//fmt.Println(inodeIndex)
	if inode.I_type[0] == '1' {
		//fmt.Println("El tipo es 1 (Archivo)")
		return "", nil
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
			return "", err
		}

		// Iterar sobre cada contenido del bloque, desde el index 2 porque los primeros dos son . y ..
		for indexContent := 2; indexContent < len(block.B_content); indexContent++ {
			// Obtener el contenido del bloque
			content := block.B_content[indexContent]

			// Sí las carpetas padre no están vacías debereamos buscar la carpeta padre más cercana
			if len(parentsDir) != 0 {
				//fmt.Println("---------ESTOY  VISITANDO--------")

				// Si el contenido está vacío, salir
				if content.B_inodo == -1 {
					break
				}

				// Obtenemos la carpeta padre más cercana
				parentDir, err := utils.First(parentsDir)
				if err != nil {
					return "", err
				}

				// Convertir B_name a string y eliminar los caracteres nulos
				contentName := strings.Trim(string(content.B_name[:]), "\x00 ")
				// Convertir parentDir a string y eliminar los caracteres nulos
				parentDirName := strings.Trim(parentDir, "\x00 ")
				// fmt.Println(contentName)
				// fmt.Println(parentDirName)
				// Si el nombre del contenido coincide con el nombre de la carpeta padre
				if strings.EqualFold(contentName, parentDirName) {
					//fmt.Println("---------LA ENCONTRÉ-------")
					// Si son las mismas, entonces entramos al inodo que apunta el bloque
					cade, err := sb.getContenidoFile(path, content.B_inodo, utils.RemoveElement(parentsDir, 0), destDir)
					if err != nil {
						return "", err
					}
					return cade, nil
				}
			} else {
				//fmt.Println("---------ESTOY  OBTENIENDO--------")

				//Comprovamos que sea el mismo archivo
				contentName := strings.Trim(string(content.B_name[:]), "\x00 ")
				if contentName == destDir {
					//Aca se obtiene el inodo:
					inode2 := &Inode{}
					// Deserializar el inodo
					err := inode2.Deserialize(path, int64(sb.S_inode_start+(content.B_inodo*sb.S_inode_size)))
					if err != nil {
						return "", err
					}
					//Aca se muestra el inodo
					//inode2.Print()

					contenido := ""
					// Iterar sobre cada bloque del inodo (apuntadores)
					for _, blockIndex2 := range inode2.I_block {
						// Si el bloque no existe, salir
						if blockIndex2 == -1 {
							break
						}

						filebloque := &FileBlock{}

						// Deserializar el filebloque
						err := filebloque.Deserialize(path, int64(sb.S_block_start+(blockIndex2*sb.S_block_size))) // 64 porque es el tamaño de un bloque
						if err != nil {
							return "", err
						}

						contenido += strings.Trim(string(filebloque.B_content[:]), "\x00 ")
						// fmt.Println("****************contenido******************")
						// filebloque.Print()
					}
					//fmt.Println(contenido)

					return contenido, nil
				}

				return "", nil
			}
		}

	}
	return "", nil
}
