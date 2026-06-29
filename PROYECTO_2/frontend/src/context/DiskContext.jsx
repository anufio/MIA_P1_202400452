import { createContext, useMemo, useState } from "react";
// DiskContext.jsx: contiene el contexto de disco, proporcionando funciones para seleccionar un disco y una partición, así como el estado de la selección actual. También maneja la persistencia de la selección en el estado del componente.

export const DiskContext = createContext(null);

export function DiskProvider({ children }) {
  const [selectedDisk, setSelectedDisk] = useState(null);
  const [selectedPartition, setSelectedPartition] = useState(null);

  const value = useMemo(
    () => ({
      selectedDisk,
      selectedPartition,
      setSelectedDisk,
      setSelectedPartition
    }),
    [selectedDisk, selectedPartition]
  );

  return <DiskContext.Provider value={value}>{children}</DiskContext.Provider>;
}
