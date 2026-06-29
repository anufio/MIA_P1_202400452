import { useState } from "react";
import { resetSystem } from "../api/resetApi";
// ResetButton.jsx: contiene el componente que representa un botón para reiniciar el sistema de archivos, mostrando un mensaje de confirmación antes de realizar la acción y manejando la comunicación con el backend a través de resetApi.
export default function ResetButton() {
  const [loading, setLoading] = useState(false);

  async function handleReset() {
    const ok = window.confirm(
      "¿Seguro que querés reiniciar todo? Se borrarán discos, particiones, reportes, montajes y sesión."
    );

    if (!ok) return;

    try {
      setLoading(true);

      await resetSystem();

      localStorage.clear();
      sessionStorage.clear();

      alert("Sistema reiniciado correctamente.");
      window.location.href = "/";
    } catch (error) {
      alert(
        error?.response?.data?.error ||
          error?.response?.data?.message ||
          error?.message ||
          "Error reiniciando sistema"
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <button
      onClick={handleReset}
      disabled={loading}
      style={{
        padding: "12px 16px",
        borderRadius: "12px",
        border: "none",
        background: "#8f2f4f",
        color: "white",
        fontWeight: "700",
        cursor: "pointer",
      }}
    >
      {loading ? "Reiniciando..." : "Resetear sistema"}
    </button>
  );
}
