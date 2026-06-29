import axiosClient from "./axiosClient";
// resetApi.js: contiene la función para manejar la solicitud de reinicio del sistema de archivos, incluyendo la comunicación con el backend a través de axiosClient.

export async function resetSystem() {
  const response = await axiosClient.post("/reset");
  return response.data;
}
