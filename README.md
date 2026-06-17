# MIA_P1_202400452

# Proyecto 1 - Sistema de Archivos EXT2

## Información General

**Curso:** Manejo e Implementación de Archivos (MIA)  
**Vacaciones Junio 2026**  
**Facultad de Ingeniería - Escuela de Ciencias y Sistemas**

**Estudiante:** Ana Lucía Nufio Roblero  
**Carnet:** 202400452

---

## Descripción

Este proyecto consiste en el desarrollo de una aplicación de línea de comandos (CLI) para la administración y simulación de un sistema de archivos EXT2 utilizando archivos binarios como discos virtuales.

La aplicación permite la creación y administración de discos, particiones, usuarios, grupos, archivos y carpetas, implementando estructuras internas similares a las utilizadas por el sistema de archivos EXT2.

Todo el manejo del almacenamiento se realiza mediante archivos binarios con extensión `.mia`, simulando el comportamiento de un disco duro real.

---

## Objetivos

- Comprender el funcionamiento interno del sistema de archivos EXT2.
- Implementar estructuras de almacenamiento utilizando lenguaje Go.
- Aplicar conceptos de particionamiento y administración de discos.
- Implementar gestión de usuarios, grupos y permisos.
- Simular operaciones sobre archivos y directorios.
- Generar reportes gráficos utilizando Graphviz.

---

## Tecnologías Utilizadas

- Go
- Graphviz
- GNU/Linux
- Git y GitHub

---

## Funcionalidades Implementadas

### Administración de Discos

- MKDISK
- RMDISK
- FDISK
- MOUNT

### Sistema de Archivos

- MKFS
- CAT

### Administración de Usuarios y Grupos

- LOGIN
- LOGOUT
- MKGRP
- RMGRP
- MKUSR
- RMUSR
- CHGRP

### Administración de Archivos y Carpetas

- MKFILE
- MKDIR

### Generación de Reportes

- MBR
- DISK
- INODE
- BLOCK
- BM_INODE
- BM_BLOCK
- SB
- FILE
- LS
- TREE

---

## Estructuras Implementadas

### Discos y Particiones

- MBR (Master Boot Record)
- Particiones Primarias
- Particiones Extendidas
- EBR (Extended Boot Record)
- Particiones Lógicas

### Sistema de Archivos EXT2

- SuperBloque
- Bitmap de Inodos
- Bitmap de Bloques
- Tabla de Inodos
- Bloques de Carpetas
- Bloques de Archivos
- Bloques de Apuntadores

### Gestión de Usuarios

- Archivo users.txt
- Usuarios
- Grupos
- Permisos UGO

---

## Reportes

El sistema genera reportes mediante Graphviz para visualizar las estructuras internas del sistema de archivos.

Entre los reportes implementados se encuentran:

- Reporte MBR
- Reporte DISK
- Reporte INODE
- Reporte BLOCK
- Reporte BM_INODE
- Reporte BM_BLOCK
- Reporte SUPERBLOCK
- Reporte FILE
- Reporte LS
- Reporte TREE


---

