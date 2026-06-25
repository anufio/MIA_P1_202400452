# Manual Técnico  
## Proyecto Fase 2 - MIA

**Nombre:** Ana Lucía Nufio Roblero 
**Carnet:** 202400452 
**Curso:** Manejo e Implementación de Archivos 
**Proyecto:** Fase 2 



## 1. Introducción

El presente manual técnico describe la estructura, arquitectura, componentes principales, instalación, ejecución y funcionamiento interno del Proyecto Fase 2 de Manejo e Implementación de Archivos.

El sistema corresponde a una evolución del Proyecto Fase 1. En esta fase se agrega una interfaz gráfica web para administrar discos, particiones, reportes y navegación del sistema de archivos EXT2. También se incorporan nuevos comandos sobre el sistema de archivos, tales como `remove`, `edit`, `rename`, `copy`, `move` y nuevas opciones para `fdisk`.

El proyecto está dividido en dos aplicaciones principales:

- **Backend:** API REST desarrollada en Go.
- **Frontend:** Interfaz gráfica desarrollada con React y Vite.



## 2. Arquitectura general

La arquitectura del sistema está compuesta por tres partes principales:

1. **Frontend Web**
   - Interfaz visual para el usuario.
   - Permite crear discos, manejar particiones, iniciar sesión, explorar archivos, ejecutar comandos y visualizar reportes.
   - Se despliega en AWS S3.

2. **Backend API REST**
   - Desarrollado en Go.
   - Expone endpoints para manipular discos, particiones, reportes, autenticación y sistema de archivos.
   - Integra la lógica de la Fase 1 y las nuevas funcionalidades de la Fase 2.
   - Se despliega en AWS EC2.

3. **Sistema de archivos EXT2**
   - Manejo de estructuras como MBR, particiones, superbloque, inodos, bloques, bitmaps y archivos.
   - Las operaciones se realizan directamente sobre archivos `.dsk`.


## 4. Backend

### 4.1 Descripción

El backend es una API REST desarrollada en Go. Su función es recibir las solicitudes del frontend y ejecutar las operaciones correspondientes sobre los discos virtuales y el sistema de archivos EXT2.

El backend se encarga de:

- Crear discos.
- Eliminar discos.
- Administrar particiones.
- Montar y desmontar particiones.
- Formatear particiones.
- Iniciar y cerrar sesión.
- Ejecutar operaciones sobre archivos y carpetas.
- Generar reportes.
- Leer la estructura del sistema de archivos para mostrarla en la GUI.

### 4.2 Carpetas principales del backend

#### `internal/api/`

Contiene la configuración principal de rutas y servidor HTTP. Aquí se conectan los handlers con los servicios internos.

#### `internal/handlers/`

Contiene los controladores que reciben las peticiones HTTP. Su función es leer los datos enviados por el frontend, validar entradas básicas y llamar al servicio correspondiente.

Ejemplos de responsabilidades:

- Recibir datos para crear disco.
- Recibir datos para crear partición.
- Ejecutar login.
- Recibir comandos del sistema de archivos.
- Solicitar generación de reportes.

#### `internal/services/`

Contiene la lógica principal del sistema. Esta carpeta funciona como capa intermedia entre los handlers y las estructuras EXT2.

Aquí se implementan servicios para:

- Discos.
- Particiones.
- Montaje.
- Formateo.
- Login.
- Explorador de archivos.
- Comandos nuevos.
- Reportes.

#### `internal/disk/`

Contiene las estructuras y funciones relacionadas con discos y particiones.

Aquí se manejan:

- MBR.
- Particiones primarias.
- Partición extendida.
- Particiones lógicas.
- EBR.
- Lectura y escritura binaria sobre archivos `.dsk`.

#### `internal/ext2/`

Contiene las estructuras internas del sistema de archivos EXT2.

Aquí se manejan:

- Superbloque.
- Inodos.
- Bloques de carpeta.
- Bloques de archivo.
- Bitmaps de inodos.
- Bitmaps de bloques.
- Permisos.
- Búsqueda de rutas.
- Escritura y lectura de archivos.

#### `internal/reports/`

Contiene la generación de reportes solicitados en el proyecto.

Reportes principales:

- MBR.
- Disk.
- Tree.
- Inode.
- Block.
- Bitmap de inodos.
- Bitmap de bloques.
- Superbloque.


## 5. Frontend

### 5.1 Descripción

El frontend está desarrollado con React y Vite. Su objetivo es proporcionar una interfaz gráfica para que el usuario pueda interactuar con el sistema sin depender únicamente de comandos en consola.

El frontend permite:

- Iniciar sesión.
- Crear discos.
- Eliminar discos.
- Seleccionar discos.
- Administrar particiones.
- Montar y desmontar particiones.
- Formatear particiones.
- Navegar el sistema de archivos.
- Visualizar contenido de archivos.
- Ejecutar comandos desde una consola web.
- Ver reportes generados por el backend.


### 5.2 Carpetas principales del frontend

#### `src/api/`

Contiene los archivos encargados de comunicarse con el backend.

Ejemplos:

- Cliente Axios.
- API de autenticación.
- API de discos.
- API de particiones.
- API de comandos.
- API de reportes.

#### `src/components/`

Contiene componentes reutilizables de la interfaz.

Ejemplos:

- Navbar.
- Sidebar.
- Botones.
- Tarjetas.
- Componentes visuales compartidos.

#### `src/context/`

Contiene contextos globales de React.

Ejemplos:

- Contexto de autenticación.
- Contexto de discos.
- Contexto de sistema de archivos.

#### `src/pages/`

Contiene las pantallas principales del sistema.

Pantallas principales:

- Login.
- Discos.
- Particiones.
- Explorador.
- Reportes.
- Consola de comandos.

#### `src/utils/`

Contiene constantes o funciones auxiliares utilizadas por el frontend.


## 6. Flujo de funcionamiento

### 6.1 Flujo general

1. El usuario ingresa al frontend.
2. El frontend realiza peticiones HTTP al backend.
3. El backend procesa la solicitud.
4. El backend lee o modifica el disco virtual `.dsk`.
5. El backend responde al frontend en formato JSON.
6. El frontend actualiza la información visualmente.


### 6.2 Flujo para administrar discos

1. El usuario entra a la página de discos.
2. Crea un disco indicando nombre, tamaño y unidad.
3. El frontend envía la solicitud al backend.
4. El backend crea el archivo `.dsk`.
5. El usuario puede seleccionar o eliminar el disco creado.


### 6.3 Flujo para administrar particiones

1. El usuario selecciona un disco.
2. Crea una partición indicando nombre, tamaño, tipo y ajuste.
3. El backend actualiza el MBR del disco.
4. El usuario puede montar la partición.
5. El usuario puede formatear la partición como EXT2.
6. Después del formateo se puede iniciar sesión y navegar el sistema.


### 6.4 Flujo de login

1. El usuario selecciona una partición montada.
2. Ingresa usuario y contraseña.
3. El frontend envía los datos al backend.
4. El backend busca el archivo `users.txt` dentro de la partición.
5. Si las credenciales son correctas, se crea una sesión.
6. El token de sesión se usa para validar operaciones protegidas.

Credenciales base después de formatear:

```text
Usuario: root
Contraseña: 123
```


### 6.5 Flujo del explorador

1. El usuario selecciona una partición montada y formateada.
2. El sistema abre la raíz `/`.
3. El backend lee el inodo raíz y sus bloques.
4. El frontend muestra carpetas y archivos.
5. Al abrir una carpeta, se consulta nuevamente el backend.
6. Al abrir un archivo, se muestra su contenido.


### 6.6 Flujo de reportes

1. El usuario selecciona disco o partición.
2. Solicita el reporte desde la GUI.
3. El backend genera el reporte.
4. El frontend muestra el resultado visual.
5. Los reportes sirven para validar el estado interno del sistema.


## 7. Comandos implementados de Fase 2

### 7.1 `fdisk -add`

Permite agregar o quitar espacio a una partición existente.

Ejemplo:

```bash
fdisk -add=1 -unit=M -path="/home/discos/d1.dsk" -name=Particion1
```

Si el valor de `add` es positivo, se intenta aumentar el tamaño. 
Si el valor de `add` es negativo, se intenta reducir el tamaño.

Validaciones esperadas:

- La partición debe existir.
- Debe haber espacio libre después de la partición para aumentar.
- No se debe permitir que la partición quede con tamaño negativo.



### 7.2 `fdisk -delete`

Permite eliminar una partición.

Ejemplo:

```bash
fdisk -delete=fast -path="/home/discos/d1.dsk" -name=Particion1
```

Modos:

- `fast`: marca la partición como vacía.
- `full`: marca la partición como vacía y limpia el espacio con `\0`.

Validaciones esperadas:

- La partición debe existir.
- Si se elimina una extendida, también deben eliminarse sus lógicas.
- Se debe confirmar la eliminación.



### 7.3 `remove`

Permite eliminar un archivo o carpeta.

Ejemplo:

```bash
remove -path=/home/user/docs/a.txt
```

Validaciones esperadas:

- La ruta debe existir.
- El usuario debe tener permiso de escritura.
- Si es carpeta, debe eliminar su contenido respetando permisos.
- Si un hijo no puede eliminarse por permisos, no deben eliminarse los padres.



### 7.4 `edit`

Permite editar el contenido de un archivo.

Ejemplo:

```bash
edit -path=/home/user/docs/a.txt -contenido=/home/anufio/contenido.txt
```

Validaciones esperadas:

- La ruta debe existir.
- La ruta debe ser un archivo.
- El usuario debe tener permiso de lectura y escritura.
- El parámetro `-contenido` indica un archivo del sistema operativo cuyo contenido será copiado al archivo EXT2.



### 7.5 `rename`

Permite cambiar el nombre de un archivo o carpeta.

Ejemplo:

```bash
rename -path=/home/user/docs/a.txt -name=b1.txt
```

Validaciones esperadas:

- La ruta debe existir.
- El usuario debe tener permiso de escritura.
- No debe existir otro archivo o carpeta con el mismo nombre en el destino.
- El nuevo nombre no debe ser inválido.



### 7.6 `copy`

Permite copiar un archivo o carpeta hacia otro destino.

Ejemplo:

```bash
copy -path="/home/user/documents" -destino="/home/images/documents"
```

Validaciones esperadas:

- La ruta origen debe existir.
- La carpeta destino debe existir.
- El usuario debe tener permisos de lectura sobre el origen.
- El usuario debe tener permisos de escritura sobre el destino.
- Si un archivo interno no tiene permiso de lectura, no se copia ese archivo.



### 7.7 `move`

Permite mover un archivo o carpeta hacia otro destino.

Ejemplo:

```bash
move -path="/home/user/documents" -destino="/home/images/documents"
```

Validaciones esperadas:

- La ruta origen debe existir.
- La carpeta destino debe existir.
- El usuario debe tener permiso de escritura sobre el origen.
- El usuario debe tener permiso de escritura sobre el destino.
- Si el origen y destino están en la misma partición, se deben actualizar referencias internas.



## 8. Reportes disponibles

### 8.1 MBR

Muestra la información principal del disco:

- Tamaño.
- Fecha de creación.
- Signature.
- Particiones.
- Estado de cada partición.
- Tipo.
- Fit.
- Inicio.
- Tamaño.
- Nombre.

### 8.2 Disk

Muestra la distribución visual del disco:

- MBR.
- Particiones primarias.
- Extendida.
- Lógicas.
- Espacios libres.

### 8.3 Tree

Muestra la estructura jerárquica del sistema de archivos EXT2:

- Superbloque.
- Inodos.
- Bloques.
- Relaciones entre carpetas y archivos.

### 8.4 Inode

Muestra información de inodos utilizados:

- UID.
- GID.
- Tamaño.
- Fechas.
- Tipo.
- Permisos.
- Apuntadores a bloques.

### 8.5 Block

Muestra bloques usados:

- Bloques de carpeta.
- Bloques de archivo.
- Contenido o referencias internas.

### 8.6 Bitmap de inodos

Muestra el estado de uso de los inodos.

### 8.7 Bitmap de bloques

Muestra el estado de uso de los bloques.

### 8.8 Superbloque

Muestra información general de la partición formateada:

- Total de inodos.
- Total de bloques.
- Inodos libres.
- Bloques libres.
- Inicio de bitmaps.
- Inicio de tabla de inodos.
- Inicio de bloques.



## 9. Endpoints principales

La API REST se organiza de forma general bajo `/api`.

Ejemplos de rutas esperadas:

```text
/api/disks
/api/partitions
/api/mount
/api/auth/login
/api/auth/logout
/api/auth/me
/api/fs/list
/api/fs/read
/api/fs/mkdir
/api/fs/mkfile
/api/fs/edit
/api/fs/rename
/api/fs/remove
/api/fs/copy
/api/fs/move
/api/reports
```



## 10. Instalación y ejecución local

### 10.1 Requisitos

Backend:

```text
Go
Linux o entorno compatible
```

Frontend:

```text
Node.js
npm
```



### 10.2 Ejecutar backend

```bash
cd "PROYECTO 2/backend"
go run .
```


### 10.3 Ejecutar frontend

```bash
cd "PROYECTO 2/frontend"
npm install
npm run dev
```



### 10.4 Compilar frontend

```bash
cd "PROYECTO 2/frontend"
npm run build
```

La carpeta generada para despliegue será:

```text
dist/
```



### 10.5 Probar backend

```bash
cd "PROYECTO 2/backend"
go test ./...
```






