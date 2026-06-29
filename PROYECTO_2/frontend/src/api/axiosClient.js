import axios from "axios";
import { API_BASE_URL } from "../utils/constants";

const axiosClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    "Content-Type": "application/json"
  }
});

export function getStoredSession() {
  const rawSession = localStorage.getItem("mia_session");

  if (!rawSession) return null;

  try {
    return JSON.parse(rawSession);
  } catch {
    localStorage.removeItem("mia_session");
    localStorage.removeItem("mia_token");
    return null;
  }
}

export function getStoredToken() {
  const session = getStoredSession();
  return session?.token || localStorage.getItem("mia_token") || "";
}

axiosClient.interceptors.request.use((config) => {
  const token = getStoredToken();

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});

axiosClient.interceptors.response.use(
  (response) => {
    if (response?.data?.success === false) {
      const message = response.data.message || response.data.error || "Error en la petición";
      const error = new Error(message);
      error.response = response;
      throw error;
    }

    return response;
  },
  (error) => Promise.reject(error)
);

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