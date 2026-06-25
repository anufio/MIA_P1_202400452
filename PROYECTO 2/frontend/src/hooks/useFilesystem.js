import { useContext, useEffect, useState } from "react";
import { FileSystemContext } from "../context/FileSystemContext";
import { getDirectory } from "../api/filesystemApi";
// useFilesystem.js: contiene el hook personalizado que permite acceder al contexto del sistema de archivos, proporcionando funciones para cambiar la ruta actual y el estado de la ruta actual. También maneja la persistencia de la ruta en el estado del componente y la carga del contenido del directorio desde el backend a través de filesystemApi.
export function useFileSystemState() {
  const context = useContext(FileSystemContext);

  if (!context) {
    throw new Error("useFileSystemState debe usarse dentro de FileSystemProvider.");
  }

  return context;
}

export function useDirectory(path) {
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(true);

  const loadDirectory = async () => {
    setLoading(true);
    const data = await getDirectory(path);
    setItems(data);
    setLoading(false);
  };

  useEffect(() => {
    loadDirectory();
  }, [path]);

  return {
    items,
    loading,
    loadDirectory
  };
}
