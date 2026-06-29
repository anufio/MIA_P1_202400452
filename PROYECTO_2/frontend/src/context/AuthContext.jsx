import { createContext, useEffect, useMemo, useState } from "react";
import { loginRequest, logoutRequest, meRequest, normalizeSession } from "../api/authApi";

const SESSION_KEY = "mia_session";
const TOKEN_KEY = "mia_token";

export const AuthContext = createContext(null);

function readSavedSession() {
  const raw = localStorage.getItem(SESSION_KEY);

  if (!raw) return null;

  try {
    return normalizeSession(JSON.parse(raw));
  } catch {
    localStorage.removeItem(SESSION_KEY);
    localStorage.removeItem(TOKEN_KEY);
    return null;
  }
}

function saveSession(session) {
  localStorage.setItem(SESSION_KEY, JSON.stringify(session));
  localStorage.setItem(TOKEN_KEY, session.token);
}

function clearSession() {
  localStorage.removeItem(SESSION_KEY);
  localStorage.removeItem(TOKEN_KEY);
}

export function AuthProvider({ children }) {
  const [session, setSession] = useState(() => readSavedSession());
  const [loading, setLoading] = useState(false);
  const [checkingSession, setCheckingSession] = useState(Boolean(readSavedSession()));

  useEffect(() => {
    let alive = true;

    async function verifySession() {
      const saved = readSavedSession();

      if (!saved?.token) {
        clearSession();
        if (alive) {
          setSession(null);
          setCheckingSession(false);
        }
        return;
      }

      const result = await meRequest();

      if (!alive) return;

      if (result.success && result.session) {
        saveSession(result.session);
        setSession(result.session);
      } else {
        clearSession();
        setSession(null);
      }

      setCheckingSession(false);
    }

    verifySession();

    return () => {
      alive = false;
    };
  }, []);

  const login = async (credentials) => {
    setLoading(true);

    try {
      const result = await loginRequest(credentials);

      if (!result.success || !result.session) {
        throw new Error(result.message || "No se pudo iniciar sesión.");
      }

      saveSession(result.session);
      setSession(result.session);

      return result;
    } finally {
      setLoading(false);
    }
  };

  const logout = async () => {
    const token = session?.token || readSavedSession()?.token;

    await logoutRequest(token);
    clearSession();
    setSession(null);
  };

  const value = useMemo(
    () => ({
      session,
      loading,
      checkingSession,
      isAuthenticated: Boolean(session?.token),
      login,
      logout
    }),
    [session, loading, checkingSession]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}