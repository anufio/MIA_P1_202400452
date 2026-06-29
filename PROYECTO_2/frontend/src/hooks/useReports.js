import { useEffect, useState } from "react";
import { getReport } from "../api/reportApi";
// useReports.js: contiene el hook personalizado que permite acceder a los reportes generados, proporcionando funciones para cargar un reporte específico y el estado de carga del reporte. También maneja la persistencia del reporte en el estado del componente y la carga del reporte desde el backend a través de reportApi.
export function useReport(type) {
  const [report, setReport] = useState(null);
  const [loading, setLoading] = useState(true);

  const loadReport = async () => {
    setLoading(true);
    const data = await getReport(type);
    setReport(data);
    setLoading(false);
  };

  useEffect(() => {
    if (type) {
      loadReport();
    }
  }, [type]);

  return {
    report,
    loading,
    loadReport
  };
}
