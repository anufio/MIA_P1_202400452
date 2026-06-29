import axiosClient, { getErrorMessage, getStoredToken } from "./axiosClient";

export function normalizeSession(session) {
  if (!session) return null;

  const token = session.token || session.Token || "";
  const user = session.user || session.username || session.User || "";
  const partId = session.partId || session.id || session.PartID || session.partID || "";

  if (!token || !user || !partId) {
    return null;
  }

  return {
    ...session,
    token,
    user,
    username: user,
    id: partId,
    partId,
    active: session.active ?? session.Active ?? true,
    uid: session.uid ?? session.UID,
    gid: session.gid ?? session.GID,
    group: session.group ?? session.Group,
    diskPath: session.diskPath || session.DiskPath || ""
  };
}

export async function loginRequest(credentials) {
  try {
    const response = await axiosClient.post("/auth/login", {
      id: String(credentials.id || "").trim().toUpperCase(),
      username: String(credentials.username || "").trim(),
      user: String(credentials.username || "").trim(),
      password: credentials.password
    });

    const session = normalizeSession(response.data.session || response.data.data);

    if (!session) {
      return {
        success: false,
        message: "El backend no devolvió una sesión válida."
      };
    }

    return {
      ...response.data,
      success: true,
      session,
      data: session
    };
  } catch (error) {
    return {
      success: false,
      message: getErrorMessage(error)
    };
  }
}

export async function logoutRequest(token = getStoredToken()) {
  try {
    if (!token) {
      return {
        success: true,
        message: "No había sesión activa."
      };
    }

    const response = await axiosClient.post("/auth/logout", { token });
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
    const session = normalizeSession(response.data.session || response.data.data);

    if (!session) {
      return {
        success: false,
        message: "No hay sesión válida."
      };
    }

    return {
      ...response.data,
      success: true,
      session,
      data: session
    };
  } catch (error) {
    return {
      success: false,
      message: getErrorMessage(error)
    };
  }
}