import { createContext, useMemo, useState } from "react";
import { loginRequest, logoutRequest } from "../api/authApi";
// AuthContext.jsx: contiene el contexto de autenticación, proporcionando funciones para iniciar y cerrar sesión, así como el estado de la sesión actual. También maneja la comunicación con el backend a través de authApi y la persistencia de la sesión en localStorage.

export const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [session, setSession] = useState(null);
  const [loading, setLoading] = useState(false);

  const login = async (credentials) => {
    setLoading(true);

    try {
      const result = await loginRequest(credentials);

      const nextSession = result.session || {
        id: credentials.id,
        username: credentials.username,
        role: credentials.username === "root" ? "root" : "user"
      };

      setSession(nextSession);

      return result;
    } finally {
      setLoading(false);
    }
  };

  const logout = async () => {
    await logoutRequest();
    setSession(null);
    localStorage.removeItem("mia_session");
  };

  const value = useMemo(
    () => ({
      session,
      loading,
      isAuthenticated: Boolean(session),
      login,
      logout
    }),
    [session, loading]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}