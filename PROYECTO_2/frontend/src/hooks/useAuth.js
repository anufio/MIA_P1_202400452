import { useContext } from "react";
import { AuthContext } from "../context/AuthContext";
// useAuth.js: contiene el hook personalizado que permite acceder al contexto de autenticación, proporcionando funciones para iniciar y cerrar sesión, así como el estado de la sesión actual. También maneja la verificación de que el hook se utilice dentro del AuthProvider.
export function useAuth() {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error("useAuth debe usarse dentro de AuthProvider.");
  }

  return context;
}
