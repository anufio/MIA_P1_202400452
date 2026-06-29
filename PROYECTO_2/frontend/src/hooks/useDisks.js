import { useContext, useEffect, useState } from "react";
import { DiskContext } from "../context/DiskContext";
import { getDisks } from "../api/diskApi";
// useDisks.js: contiene el hook personalizado que permite acceder al contexto de disco, proporcionando funciones para seleccionar un disco y una partición, así como el estado de la selección actual. También maneja la persistencia de la selección en el estado del componente y la carga de los discos desde el backend a través de diskApi.
export function useDiskSelection() {
  const context = useContext(DiskContext);

  if (!context) {
    throw new Error("useDiskSelection debe usarse dentro de DiskProvider.");
  }

  return context;
}

export function useDisks() {
  const [disks, setDisks] = useState([]);
  const [loading, setLoading] = useState(true);

  const loadDisks = async () => {
    setLoading(true);
    const data = await getDisks();
    setDisks(data);
    setLoading(false);
  };

  useEffect(() => {
    loadDisks();
  }, []);

  return {
    disks,
    loading,
    loadDisks
  };
}
