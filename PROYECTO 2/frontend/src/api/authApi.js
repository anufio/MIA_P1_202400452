import axiosClient from "./axiosClient";
import { getErrorMessage } from "./axiosClient";
// authApi.js: contiene las funciones para manejar las solicitudes de autenticación, incluyendo login, logout y obtener información del usuario autenticado, utilizando axiosClient para comunicarse con el backend.
export async function loginRequest(credentials) {
  try {
    const response = await axiosClient.post("/auth/login", {
      id: credentials.id,
      username: credentials.username,
      user: credentials.username,
      password: credentials.password
    });

    return response.data;
  } catch (error) {
    return {
      success: false,
      message: getErrorMessage(error)
    };
  }
}

export async function logoutRequest(token) {
  try {
    const response = await axiosClient.post("/auth/logout", {
      token
    });

    return response.data;
  } catch (error) {
    return {
      success: false,
      message: getErrorMessage(error)
    };
  }
}

export async function meRequest() {
  try {
    const response = await axiosClient.get("/auth/me");
    return response.data;
  } catch (error) {
    return {
      success: false,
      message: getErrorMessage(error)
    };
  }
}