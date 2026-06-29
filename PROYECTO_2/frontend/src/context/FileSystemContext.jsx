import { createContext, useMemo, useState } from "react";
// FileSystemContext.jsx: contiene el contexto del sistema de archivos, proporcionando funciones para cambiar la ruta actual y el estado de la ruta actual. También maneja la persistencia de la ruta en el estado del componente.
export const FileSystemContext = createContext(null);

export function FileSystemProvider({ children }) {
  const [currentPath, setCurrentPath] = useState("/");

  const value = useMemo(
    () => ({
      currentPath,
      setCurrentPath
    }),
    [currentPath]
  );

  return (
    <FileSystemContext.Provider value={value}>
      {children}
    </FileSystemContext.Provider>
  );
}
