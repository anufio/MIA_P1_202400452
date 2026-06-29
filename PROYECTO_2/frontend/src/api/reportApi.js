import axiosClient from "./axiosClient";
// reportApi.js: contiene las funciones para manejar la generación y obtención de reportes, incluyendo la normalización de URLs y la comunicación con el backend a través de axiosClient.

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api";
const BACKEND_URL = API_URL.replace(/\/api\/?$/, "");

function normalizeReport(report) {
  if (!report) return null;

  const url = report.url || report.URL || "";
  const fullUrl = url.startsWith("http") ? url : `${BACKEND_URL}${url}`;

  return {
    ...report,
    url: fullUrl,
    imageUrl: fullUrl,
  };
}

function getBackendError(error) {
  return (
    error?.response?.data?.error ||
    error?.response?.data?.message ||
    error?.message ||
    "Error generando reporte"
  );
}

export async function generateReport(payload) {
  try {
    const response = await axiosClient.post("/reports/generate", payload);
    const report = response.data.report || response.data.data;
    return normalizeReport(report);
  } catch (error) {
    throw new Error(getBackendError(error));
  }
}

export async function listReports() {
  try {
    const response = await axiosClient.get("/reports");
    const reports = response.data.reports || response.data.data || [];
    return reports.map(normalizeReport);
  } catch {
    return [];
  }
}

export async function getReport(type, id) {
  return generateReport({ type, id });
}
