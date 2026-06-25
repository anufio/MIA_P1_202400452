import axios from "axios";
import { API_BASE_URL } from "../utils/constants";
// axiosClient.js: contiene la configuración de Axios para realizar solicitudes HTTP al backend, incluyendo la URL base, el tiempo de espera y los encabezados predeterminados. También incluye una función para obtener mensajes de error de las respuestas del servidor.

const axiosClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 5000,
  headers: {
    "Content-Type": "application/json"
  }
});

axiosClient.interceptors.request.use((config) => {
  try {
    const rawSession = localStorage.getItem("mia_session");

    if (rawSession) {
      const session = JSON.parse(rawSession);

      if (session?.token) {
        config.headers.Authorization = `Bearer ${session.token}`;
      }
    }
  } catch {
    localStorage.removeItem("mia_session");
  }

  return config;
});

export function getErrorMessage(error) {
  if (error?.response?.data?.message) {
    return error.response.data.message;
  }

  if (error?.response?.data?.error) {
    return error.response.data.error;
  }

  if (error?.message) {
    return error.message;
  }

  return "Ocurrió un error inesperado.";
}

export default axiosClient;
