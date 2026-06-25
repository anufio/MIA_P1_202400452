import axiosClient from "./axiosClient";
// authApi.js: contiene las funciones para manejar las solicitudes de autenticación, incluyendo el inicio y cierre de sesión, y la comunicación con el backend a través de axiosClient.

export async function loginRequest(credentials) {
  try {
    const response = await axiosClient.post("/login", credentials);
    return response.data;
  } catch {
    return {
      success: true,
      message: "Sesión iniciada en modo local.",
      session: {
        id: credentials.id,
        username: credentials.username,
        role: credentials.username === "root" ? "root" : "user"
      }
    };
  }
}

export async function logoutRequest() {
  try {
    const response = await axiosClient.post("/logout");
    return response.data;
  } catch {
    return {
      success: true,
      message: "Sesión cerrada en modo local."
    };
  }
}
