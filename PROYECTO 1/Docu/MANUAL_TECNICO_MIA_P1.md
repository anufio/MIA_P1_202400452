# Manual Técnico
## Proyecto 1 - Manejo e Implementación de Archivos


**Estudiante:** Ana Lucía Nufio Roblero  
**Carné:** 202400452  
**Lenguaje:** Go  
**Entorno:** Linux

---

## 1. Descripción general del proyecto

El proyecto consiste en una aplicación CLI desarrollada en Go para simular el manejo de discos, particiones y un sistema de archivos EXT2. El sistema trabaja sobre archivos binarios que representan discos virtuales. Dentro de estos discos se escriben estructuras binarias como MBR, particiones, EBR, superbloque, bitmaps, inodos y bloques.

La aplicación no usa interfaz gráfica. Todas las operaciones se realizan desde terminal. Los reportes se generan mediante Graphviz para representar visualmente las estructuras internas del disco y del sistema de archivos.

---

## 2. Objetivos técnicos

- Crear discos simulados mediante archivos binarios.
- Implementar estructuras de particionamiento como MBR, Partition y EBR.
- Aplicar teoría de particiones primarias, extendidas y lógicas.
- Implementar un sistema de archivos EXT2 sobre particiones montadas.
- Administrar usuarios y grupos mediante el archivo lógico `/users.txt`.
- Crear carpetas y archivos usando inodos, bloques y bitmaps.
- Generar reportes visuales con Graphviz.
- Procesar comandos individuales y scripts `.smia`.
- Validar parámetros y mostrar errores cuando los comandos no cumplan el formato esperado.

---

## 3. Tecnologías utilizadas

| Tecnología | Uso |
|---|---|
| Go | Lenguaje principal del proyecto. |
| Linux | Entorno de ejecución y calificación. |
| Graphviz | Generación de reportes gráficos. |
| Archivos binarios | Simulación de discos duros. |
| Terminal | Interfaz principal del usuario. |

---

## 4. Estructura del proyecto

La estructura general del proyecto es:

```text
MIA_P1_202400452/
└── PROYECTO 1/
    ├── commands/
    │   ├── cmd_cat.go
    │   ├── cmd_disk.go
    │   ├── cmd_fdisk.go
    │   ├── cmd_files.go
    │   ├── cmd_mkfs.go
    │   ├── cmd_mount.go
    │   ├── cmd_session.go
    │   └── cmd_users.go
    ├── disk/
    │   ├── disk_utils.go
    │   ├── ebr.go
    │   └── mbr.go
    ├── ext2/
    │   ├── bitmap.go
    │   ├── block.go
    │   ├── ext2_utils.go
    │   ├── inode.go
    │   └── superblock.go
    ├── parser/
    │   └── parser.go
    ├── reports/
    │   ├── rep.go
    │   ├── rep_bitmap.go
    │   ├── rep_inode_block.go
    │   ├── rep_mbr_disk.go
    │   ├── rep_sb_file_ls.go
    │   └── rep_utils.go
    ├── storage/
    │   └── storage.go
    ├── go.mod
    ├── main.go
    └── README.md
```

---

## 5. Descripción de módulos

### 5.1. `main.go`

Contiene el ciclo principal del programa. Sus responsabilidades son:

- Mostrar el encabezado de la aplicación.
- Leer comandos desde terminal.
- Detectar el comando `exit`.
- Enviar cada línea a procesamiento.
- Validar comandos y parámetros admitidos.
- Redirigir cada comando hacia su función correspondiente.

También maneja la ejecución de scripts mediante el comando `execute`.

---

### 5.2. `parser/`

Contiene el analizador de comandos.

Responsabilidades:

- Separar la línea de entrada en tokens.
- Respetar comillas dobles para valores con espacios.
- Extraer nombre de comando.
- Extraer parámetros en formato `-clave=valor`.
- Convertir nombres de comandos y parámetros a minúsculas.
- Procesar scripts línea por línea.

---

### 5.3. `commands/`

Contiene la lógica de los comandos del sistema.

Archivos principales:

| Archivo | Responsabilidad |
|---|---|
| `cmd_disk.go` | Creación y eliminación de discos. |
| `cmd_fdisk.go` | Creación de particiones. |
| `cmd_mount.go` | Montaje, desmontaje y tabla de particiones montadas. |
| `cmd_mkfs.go` | Formateo EXT2. |
| `cmd_session.go` | Inicio y cierre de sesión. |
| `cmd_users.go` | Administración de usuarios y grupos. |
| `cmd_files.go` | Creación de carpetas y archivos. |
| `cmd_cat.go` | Lectura de archivos. |

---

### 5.4. `disk/`

Contiene las estructuras y funciones relacionadas con el disco simulado.

Estructuras principales:

- `MBR`
- `Partition`
- `EBR`

Estas estructuras se escriben y leen directamente del archivo binario que representa el disco.

---

### 5.5. `ext2/`

Contiene las estructuras y operaciones del sistema de archivos EXT2.

Estructuras principales:

- `SuperBlock`
- `Inode`
- `FolderBlock`
- `FileBlock`
- `PointerBlock`
- Bitmaps de inodos y bloques

Funciones principales:

- Escritura y lectura de superbloque.
- Escritura y lectura de inodos.
- Escritura y lectura de bloques.
- Administración de bitmaps.
- Búsqueda de rutas dentro del sistema EXT2.
- Lectura de contenido de archivos.

---

### 5.6. `reports/`

Contiene la generación de reportes con Graphviz.

Reportes implementados:

| Reporte | Función general |
|---|---|
| `mbr` | Reporta el MBR y particiones. |
| `disk` | Reporta distribución gráfica del disco. |
| `sb` | Reporta el superbloque. |
| `tree` | Reporta árbol de inodos y bloques. |
| `inode` | Reporta inodos usados. |
| `block` | Reporta bloques usados. |
| `bm_inode` | Reporta bitmap de inodos. |
| `bm_block` | Reporta bitmap de bloques. |
| `file` | Reporta contenido de un archivo. |
| `ls` | Reporta listado de una carpeta. |

---

### 5.7. `storage/`

Contiene utilidades para resolver rutas en el sistema de archivos real. Se utiliza principalmente para ubicar archivos externos o rutas de almacenamiento cuando el comando lo necesita.

---

## 6. Flujo general de ejecución

El flujo de procesamiento de un comando es:

```text
Entrada del usuario
        ↓
Lectura en main.go
        ↓
Parser de comando
        ↓
Validación de comando y parámetros
        ↓
Ejecución de función correspondiente
        ↓
Lectura/escritura de estructuras binarias
        ↓
Mensaje de salida para el usuario
```


---

## 7. Validación de comandos y parámetros

El programa valida que cada comando reciba únicamente los parámetros permitidos.

Ejemplo correcto:

```text
mkfs -type=full -id=521A
```

Ejemplo incorrecto:

```text
mkfs -type-=full -id-=521A
```

El segundo ejemplo debe mostrar error porque `-type-` e `-id-` no son parámetros definidos.

También se validan los valores de ajuste. Los valores aceptados son:

```text
BF
FF
WF
```

---

## 8. Estructuras de disco

### 8.1. MBR

El MBR se escribe al inicio del disco. Contiene:

| Campo | Descripción |
|---|---|
| `mbr_tamano` | Tamaño total del disco en bytes. |
| `mbr_fecha_creacion` | Fecha de creación del disco. |
| `mbr_dsk_signature` | Firma aleatoria del disco. |
| `dsk_fit` | Ajuste general del disco. |
| `mbr_partitions` | Arreglo de 4 particiones. |

---

### 8.2. Partition

Representa una partición primaria o extendida dentro del MBR.

| Campo | Descripción |
|---|---|
| `part_status` | Estado de la partición. |
| `part_type` | Tipo: primaria o extendida. |
| `part_fit` | Ajuste usado. |
| `part_start` | Byte de inicio. |
| `part_size` | Tamaño en bytes. |
| `part_name` | Nombre de la partición. |
| `part_correlative` | Correlativo de la partición. |
| `part_id` | ID asignado al montar, si aplica. |

---

### 8.3. EBR

Se utiliza para particiones lógicas dentro de una partición extendida.

| Campo | Descripción |
|---|---|
| `part_mount` | Estado de montaje. |
| `part_fit` | Ajuste de la lógica. |
| `part_start` | Byte de inicio. |
| `part_size` | Tamaño. |
| `part_next` | Apuntador al siguiente EBR. |
| `part_name` | Nombre de la lógica. |

---

## 9. Montaje de particiones

El montaje se realiza en memoria. El ID se genera con los últimos dos dígitos del carné, el número de partición y una letra del disco.

Para el carné `202400452`, los últimos dos dígitos son `52`. Al montar la primera partición del primer disco se genera `521A`.


---

## 10. Sistema de archivos EXT2

El comando `mkfs` crea las estructuras necesarias dentro de la partición montada.

El sistema EXT2 implementado contiene:

```text
Superbloque
Bitmap de inodos
Bitmap de bloques
Tabla de inodos
Tabla de bloques
```

El número de bloques se calcula como tres veces el número de inodos.

Al formatear se crean automáticamente:

- La carpeta raíz `/`.
- El archivo `/users.txt`.
- El grupo root.
- El usuario root.

Contenido inicial de `users.txt`:

```text
1,G,root
1,U,root,root,123
```

---

## 11. Superbloque

El superbloque guarda la información general del sistema EXT2.

Campos principales:

| Campo | Descripción |
|---|---|
| `s_filesystem_type` | Tipo de sistema de archivos. |
| `s_inodes_count` | Cantidad total de inodos. |
| `s_blocks_count` | Cantidad total de bloques. |
| `s_free_blocks_count` | Bloques libres. |
| `s_free_inodes_count` | Inodos libres. |
| `s_mtime` | Última fecha de montaje. |
| `s_umtime` | Última fecha de desmontaje. |
| `s_mnt_count` | Contador de montajes. |
| `s_magic` | Identificador EXT2, normalmente `0xEF53`. |
| `s_inode_s` | Tamaño del inodo. |
| `s_block_s` | Tamaño del bloque. |
| `s_bm_inode_start` | Inicio del bitmap de inodos. |
| `s_bm_block_start` | Inicio del bitmap de bloques. |
| `s_inode_start` | Inicio de la tabla de inodos. |
| `s_block_start` | Inicio de la tabla de bloques. |

---

## 12. Inodos

Cada archivo o carpeta tiene un inodo. El inodo almacena metadatos y apuntadores a bloques.

Campos principales:

| Campo | Descripción |
|---|---|
| `i_uid` | Usuario propietario. |
| `i_gid` | Grupo propietario. |
| `i_s` | Tamaño del archivo. |
| `i_atime` | Fecha de último acceso. |
| `i_ctime` | Fecha de creación. |
| `i_mtime` | Fecha de modificación. |
| `i_block` | Arreglo de 15 apuntadores. |
| `i_type` | `0` carpeta, `1` archivo. |
| `i_perm` | Permisos en formato octal. |

---

## 13. Bloques

El sistema usa bloques de 64 bytes. Existen tres tipos principales:

### 13.1. Bloque de carpeta

Guarda entradas de carpeta. Cada entrada contiene nombre y número de inodo.

Ejemplo en la raíz:

```text
.         -> 0
..        -> 0
users.txt -> 1
```

### 13.2. Bloque de archivo

Guarda contenido de archivo. Por ejemplo, el bloque de `users.txt` contiene:

```text
1,G,root
1,U,root,root,123
```

### 13.3. Bloque de apuntadores

Guarda apuntadores indirectos cuando un archivo o carpeta requiere más bloques.

---

## 14. Bitmaps

Los bitmaps indican qué inodos y bloques están ocupados.

Valor usado:

| Valor | Significado |
|---|---|
| `1` | Ocupado. |
| `0` | Libre. |

Después de un `mkfs` inicial, normalmente los bitmaps comienzan con:

```text
1100000...
```

Eso indica que están ocupados el inodo raíz, el inodo de `users.txt`, el bloque raíz y el bloque de archivo de `users.txt`.

---

## 15. Manejo de sesiones

La sesión activa se guarda en memoria. Para ejecutar comandos de usuarios, grupos, carpetas, archivos y lectura, debe existir una sesión activa.

El usuario inicial es:

```text
usuario: root
contraseña: 123
```

Comando de inicio de sesión:

```text
login -usr=root -pass=123 -id=521A
```

---

## 16. Archivo `users.txt`

El archivo `/users.txt` se crea durante el formateo EXT2.

Formato de grupo:

```text
GID,G,nombre_grupo
```

Formato de usuario:

```text
UID,U,grupo,usuario,contraseña
```

Ejemplo inicial:

```text
1,G,root
1,U,root,root,123
```

Cuando se elimina un usuario o grupo, su ID se marca como `0` en lugar de borrar físicamente la línea.

---

## 17. Generación de reportes

Los reportes se generan desde el paquete `reports`. El comando principal es:

```text
rep -id=<id> -Path=<ruta_salida> -name=<reporte>
```

El programa obtiene la partición montada a partir del ID y luego lee las estructuras correspondientes en el disco.

### 17.1. Reporte MBR

Muestra el contenido del MBR y sus particiones. Sirve para validar tamaño del disco, fecha, firma, ajuste y datos de las particiones.

### 17.2. Reporte DISK

Muestra visualmente la distribución del disco: MBR, particiones y espacio libre.

### 17.3. Reporte SB

Muestra el superbloque del sistema EXT2.

### 17.4. Reporte TREE

Muestra la relación entre inodos y bloques. Es útil para validar la estructura jerárquica del sistema de archivos.

### 17.5. Reporte INODE

Muestra los inodos utilizados y sus campos.

### 17.6. Reporte BLOCK

Muestra los bloques ocupados, ya sean bloques de carpeta, archivo o apuntadores.

### 17.7. Reporte BM_INODE

Muestra el bitmap de inodos.

### 17.8. Reporte BM_BLOCK

Muestra el bitmap de bloques.

### 17.9. Reporte FILE

Muestra el contenido de un archivo específico.

### 17.10. Reporte LS

Muestra el listado de una ruta, incluyendo permisos, propietario, grupo, tamaño, fechas, tipo y nombre.

---

## 18. Flujo técnico de una prueba completa

```text
mkdisk -size=75 -unit=M -path=/tmp/d1.dsk
fdisk -type=P -unit=M -name=Part1 -size=25 -path=/tmp/d1.dsk
mount -path=/tmp/d1.dsk -name=Part1
mkfs -type=full -id=521A
login -usr=root -pass=123 -id=521A
mkdir -P -path=/home/archivos/mia/fase2
mkfile -r -path="/home/AhoraSiExiste/b2.txt" -s=175
cat -file=/users.txt
rep -id=521A -Path="/home/parte2/inicial/ext2_sb_1.jpg" -name=sb
rep -id=521A -Path="/home/parte2/inicial/ext2_tree_1.jpg" -name=tree
rep -id=521A -Path="/home/parte2/avanzado/reporte_inode2.pdf" -name=inode
rep -id=521A -Path="/home/parte2/avanzado/reporte_block2.pdf" -name=block
rep -id=521A -Path="/home/parte2/avanzado/reporte_bm_inode2.pdf" -name=bm_inode
rep -id=521A -Path="/home/parte2/avanzado/reporte_bm_block2.pdf" -name=bm_block
```

---

## 19. Compilación y ejecución

Compilar:

```bash
go fmt ./...
go build -o mia
```

Ejecutar:

```bash
./mia
```

Ejecutar archivo de entrada:

```text
execute -path="/home/anufio/Descargas/calificacion_mia.smia"
```

---


