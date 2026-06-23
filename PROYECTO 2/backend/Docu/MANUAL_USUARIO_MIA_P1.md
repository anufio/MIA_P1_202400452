# Manual de Usuario
## Proyecto 1 - Manejo e Implementación de Archivos

**Estudiante:** Ana Lucía Nufio Roblero  
**Carné:** 202400452  
**Repositorio sugerido:** `MIA_P1_202400452`  
**Lenguaje:** Go

---

## 1. Introducción

Este manual explica cómo utilizar la aplicación desarrollada para el Proyecto 1 de Manejo e Implementación de Archivos. El programa funciona como una aplicación de línea de comandos, también conocida como CLI, que permite crear discos simulados, crear particiones, montar particiones, formatearlas como EXT2, administrar usuarios y grupos, crear carpetas y archivos, leer contenido y generar reportes visuales mediante Graphviz.

El sistema trabaja de forma local y no utiliza una interfaz gráfica propia. Todas las acciones se realizan desde la terminal mediante comandos escritos por el usuario o cargados desde un archivo de entrada con extensión `.smia`.

---

## 2. Requisitos para ejecutar el programa

Para utilizar el programa se necesita una instalación GNU/Linux con Go y Graphviz instalados.

### 2.1. Requisitos mínimos

- Sistema operativo GNU/Linux.
- Go instalado.
- Graphviz instalado para la generación de reportes.
- Terminal con permisos para crear archivos en las rutas usadas por los comandos.

### 2.2. Instalación de dependencias

En distribuciones basadas en Ubuntu, Debian o Linux Mint, se pueden instalar las dependencias con:

```bash
sudo apt update
sudo apt install golang-go graphviz
```

---

## 3. Compilación del proyecto

Desde la carpeta principal del proyecto, ejecutar:

```bash
go fmt ./...
go build -o mia
```

El comando `go fmt ./...` ordena el formato del código fuente.  
El comando `go build -o mia` genera el ejecutable llamado `mia`.

---

## 4. Ejecución del programa

Para iniciar la aplicación:

```bash
./mia
```

Al iniciar, el programa muestra un encabezado y la lista de comandos disponibles:

```text
mkdisk, rmdisk, fdisk, mount, unmount, mounted, mkfs,
login, logout, mkgrp, rmgrp, mkusr, rmusr, chgrp,
mkfile, mkdir, cat, rep, execute, pause
```

Para salir del programa se escribe:

```text
exit
```

---

## 5. Reglas generales de uso

Los comandos pueden escribirse en mayúsculas o minúsculas. Por ejemplo, `MKDISK`, `mkdisk` y `Mkdisk` son interpretados como el mismo comando.

Los parámetros deben escribirse con guion y signo igual:

```text
-parametro=valor
```

Si un valor contiene espacios, debe encerrarse entre comillas dobles:

```text
-path="/home/mis discos/disco1.mia"
```

Si se usa un parámetro no reconocido, el programa muestra un mensaje de error. Esto es importante porque el sistema solo debe aceptar los comandos y parámetros definidos.

---

## 6. Limpieza antes de una prueba

Antes de repetir una prueba desde cero, se recomienda limpiar discos y reportes anteriores. Desde la terminal, fuera del programa:

```bash
rm -f /tmp/d1.dsk /tmp/d2.dsk /tmp/d3.dsk
rm -f /tmp/eliminar_1.dsk /tmp/eliminar_2.dsk
rm -f /tmp/test.mia /tmp/testp1.mia
sudo rm -rf /home/parte1 /home/parte2
```

Como se uso una carpeta local llamada `reportes`, esta puede eliminarse con:

```bash
rm -rf reportes
```

---

## 7. Uso de archivos de entrada `.smia`

El comando `execute` permite cargar comandos desde un archivo de texto con extensión `.smia`.

Ejemplo:

```text
execute -path="/home/anufio/Descargas/calificacion_mia.smia"
```

El archivo puede contener:

- Comandos válidos.
- Comentarios iniciados con `#`.
- Líneas en blanco.
- Comandos `pause` para detener temporalmente la ejecución.

Cuando el programa encuentra `pause`, se detiene y espera que el usuario presione Enter.

---

## 8. Comandos de administración de discos

### 8.1. `mkdisk`

Crea un archivo binario que simula un disco.

Sintaxis:

```text
mkdisk -size=<tamaño> -unit=<K|M> -fit=<BF|FF|WF> -path=<ruta>
```

Parámetros:

| Parámetro | Obligatorio | Descripción |
|---|---:|---|
| `-size` | Sí | Tamaño del disco. Debe ser mayor que cero. |
| `-unit` | No | Unidad del tamaño. Puede ser `K` o `M`. Si no se indica, se usa `M`. |
| `-fit` | No | Ajuste del disco. Puede ser `BF`, `FF` o `WF`. Si no se indica, se usa `FF`. |
| `-path` | Sí | Ruta donde se creará el disco. |

Ejemplo:

```text
mkdisk -size=75 -unit=M -path=/tmp/d1.dsk
```

---

### 8.2. `rmdisk`

Elimina un disco existente. Antes de eliminar, solicita confirmación.

Sintaxis:

```text
rmdisk -path=<ruta>
```

Ejemplo:

```text
rmdisk -path="/tmp/eliminar_1.dsk"
```

Si el archivo no existe, el programa muestra un error.

---

### 8.3. `fdisk`

Crea particiones primarias, extendidas o lógicas dentro de un disco.

Sintaxis:

```text
fdisk -size=<tamaño> -unit=<B|K|M> -path=<ruta> -name=<nombre> -type=<P|E|L> -fit=<BF|FF|WF>
```

Parámetros:

| Parámetro | Obligatorio | Descripción |
|---|---:|---|
| `-size` | Sí | Tamaño de la partición. |
| `-unit` | No | Unidad: `B`, `K` o `M`. |
| `-path` | Sí | Ruta del disco. |
| `-name` | Sí | Nombre de la partición. |
| `-type` | No | Tipo de partición: `P`, `E` o `L`. Por defecto es `P`. |
| `-fit` | No | Ajuste: `BF`, `FF` o `WF`. |

Ejemplo:

```text
fdisk -type=P -unit=M -name=Part1 -size=25 -path=/tmp/d1.dsk
```

Errores comunes:

```text
fdisk -type=E -unit=M -name=Part1 -size=10 -path=/tmp/d2.dsk -fit=FirstFit
```

Ese comando debe mostrar error porque el ajuste correcto es `FF`, no `FirstFit`.

---

## 9. Comandos de montaje

### 9.1. `mount`

Monta una partición primaria del disco y le asigna un ID.

Sintaxis:

```text
mount -path=<ruta_disco> -name=<nombre_particion>
```

Ejemplo:

```text
mount -path=/tmp/d1.dsk -name=Part1
```

Salida esperada:

```text
Partición 'Part1' montada exitosamente con ID: '521A'

Particiones montadas:
 ID: 521A | Disco: /tmp/d1.dsk | Partición: Part1
```

El ID se forma usando los últimos dos dígitos del carné, el número correlativo de la partición y una letra asignada al disco.

---

### 9.2. `mounted`

Muestra las particiones montadas actualmente.

Sintaxis:

```text
mounted
```

---

### 9.3. `unmount`

Desmonta una partición usando su ID.

Sintaxis:

```text
unmount -id=<id>
```

Ejemplo:

```text
unmount -id=521A
```

---

## 10. Comando de formateo EXT2

### 10.1. `mkfs`

Formatea una partición montada como EXT2.

Sintaxis:

```text
mkfs -type=full -id=<id>
```

Ejemplo:

```text
mkfs -type=full -id=521A
```

Salida esperada:

```text
Partición '521A' formateada como EXT2 exitosamente.
Inodos: 79678 | Bloques: 239034
```

Al formatear, el sistema crea la carpeta raíz `/` y el archivo `/users.txt` con el usuario root inicial.

Importante: el ID debe estar montado. Si se ejecuta `unmount` antes de `mkfs`, se debe volver a montar la partición.

---

## 11. Comandos de sesión

### 11.1. `login`

Inicia sesión con un usuario existente.

Sintaxis:

```text
login -usr=<usuario> -pass=<contraseña> -id=<id>
```

Ejemplo:

```text
login -usr=root -pass=123 -id=521A
```

Salida esperada:

```text
Login exitoso. Bienvenido, root (UID=1, GID=1)
```

---

### 11.2. `logout`

Cierra la sesión activa.

Sintaxis:

```text
logout
```

Si no hay sesión activa, se muestra un error.

---

## 12. Comandos de usuarios y grupos

Estos comandos requieren sesión activa. Para crear o eliminar grupos y usuarios, se recomienda iniciar sesión como root.

### 12.1. `mkgrp`

Crea un grupo.

```text
mkgrp -name=Archivos
```

### 12.2. `rmgrp`

Marca un grupo como eliminado en `users.txt`.

```text
rmgrp -name=Archivos
```

### 12.3. `mkusr`

Crea un usuario dentro de un grupo existente.

```text
mkusr -usr="user1" -grp=root -pass=user1
```

### 12.4. `rmusr`

Marca un usuario como eliminado.

```text
rmusr -usr="user1"
```

### 12.5. `chgrp`

Cambia de grupo a un usuario.

```text
chgrp -usr="user1" -grp=Archivos
```

---

## 13. Comandos de carpetas y archivos

### 13.1. `mkdir`

Crea carpetas dentro del sistema EXT2.

```text
mkdir -P -path=/home/archivos/mia/fase2
```

El parámetro `-P` permite crear carpetas padre si no existen.

### 13.2. `mkfile`

Crea archivos dentro del sistema EXT2.

Crear archivo con tamaño generado automáticamente:

```text
mkfile -path="/home/b1.txt" -s=75
```

Crear archivo y crear carpetas padre si no existen:

```text
mkfile -r -path="/home/AhoraSiExiste/b2.txt" -s=175
```

Crear archivo con contenido desde un archivo del sistema real:

```text
mkfile -path="/home/b3.txt" -cont=/home/anufio/Documentos/algo.txt
```

### 13.3. `cat`

Muestra el contenido de un archivo del sistema EXT2.

```text
cat -file=/users.txt
```

Salida inicial esperada:

```text
1,G,root
1,U,root,root,123
```

---

## 14. Comando de reportes

### 14.1. `rep`

Genera reportes visuales o de texto a partir de las estructuras del disco o del sistema EXT2.

Sintaxis general:

```text
rep -id=<id> -Path=<ruta_salida> -name=<tipo_reporte>
```

Algunos reportes necesitan además `-path_file_ls`.

Tipos de reporte disponibles:

| Reporte | Descripción |
|---|---|
| `mbr` | Muestra el MBR y las particiones del disco. |
| `disk` | Muestra gráficamente la distribución del disco. |
| `sb` | Muestra el superbloque de EXT2. |
| `tree` | Muestra el árbol de inodos y bloques. |
| `inode` | Muestra los inodos utilizados. |
| `block` | Muestra los bloques utilizados. |
| `bm_inode` | Muestra el bitmap de inodos. |
| `bm_block` | Muestra el bitmap de bloques. |
| `file` | Muestra el contenido de un archivo específico. |
| `ls` | Lista el contenido de una ruta. |

Ejemplos:

```text
rep -id=521A -Path=/home/parte1/particiones/d1.jpg -name=disk
rep -id=521A -Path=/home/parte1/mbr1.jpg -name=mbr
rep -id=521A -Path="/home/parte2/inicial/ext2_sb_1.jpg" -name=sb
rep -id=521A -Path="/home/parte2/inicial/ext2_tree_1.jpg" -name=tree
rep -id=521A -Path="/home/parte2/avanzado/reporte_inode2.pdf" -name=inode
rep -id=521A -Path="/home/parte2/avanzado/reporte_block2.pdf" -name=block
rep -id=521A -Path="/home/parte2/avanzado/reporte_bm_inode2.pdf" -name=bm_inode
rep -id=521A -Path="/home/parte2/avanzado/reporte_bm_block2.pdf" -name=bm_block
rep -id=521A -Path="/home/parte2/extra/users_report.txt" -name=file -path_file_ls=/users.txt
rep -id=521A -Path="/home/parte2/extra/ls_raiz.jpg" -name=ls -path_file_ls=/
```

---

## 15. Errores comunes y solución

### Error: el ID no está montado

Mensaje:

```text
Error: el ID '521A' no está montado
```

Solución:

```text
mount -path=/tmp/d1.dsk -name=Part1
```

### Error: parámetro no reconocido

Ejemplo:

```text
mkfs -type-=full -id-=521A
```

Ese comando está mal. El correcto es:

```text
mkfs -type=full -id=521A
```

### Error: ajuste inválido

Ejemplo:

```text
fdisk -fit=BestFit
```

Debe usarse:

```text
-fit=BF
```

Los valores válidos son `BF`, `FF` y `WF`.

### Error: debe iniciar sesión

Mensaje:

```text
Error: debe iniciar sesión para ejecutar este comando
```

Solución:

```text
login -usr=root -pass=123 -id=521A
```

