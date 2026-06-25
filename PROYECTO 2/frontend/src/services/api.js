const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api";
// services/api.js: contiene las funciones para manejar las solicitudes al backend, incluyendo la normalización de URLs y la comunicación con el backend a través de fetch. También incluye funciones para manejar errores y obtener mensajes de error de las respuestas del servidor.
async function request(path, options = {}) {
  const token = localStorage.getItem("mia_token");

  const headers = {
    "Content-Type": "application/json",
    ...(options.headers || {}),
  };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const response = await fetch(`${API_URL}${path}`, {
    ...options,
    headers,
  });

  const data = await response.json().catch(() => ({}));

  if (!response.ok || data.success === false) {
    throw new Error(data.error || data.message || "Error en la petición");
  }

  return data;
}

export const healthAPI = {
  check: () => request("/health"),
};

export const diskAPI = {
  list: () => request("/disks"),

  create: (payload) =>
    request("/disks", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  remove: (path) =>
    request("/disks/delete", {
      method: "POST",
      body: JSON.stringify({ path }),
    }),

  delete: (path) =>
    request("/disks/delete", {
      method: "POST",
      body: JSON.stringify({ path }),
    }),
};

export const partitionAPI = {
  list: (diskPath) =>
    request(`/partitions/list?diskPath=${encodeURIComponent(diskPath)}`),

  create: (payload) =>
    request("/partitions/create", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  add: (payload) =>
    request("/partitions/add", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  mount: (payload) =>
    request("/partitions/mount", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  unmount: (id) =>
    request("/partitions/unmount", {
      method: "POST",
      body: JSON.stringify({ id }),
    }),

  remove: (payload) =>
    request("/partitions/delete", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  delete: (payload) =>
    request("/partitions/delete", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  resize: (payload) =>
    request("/partitions/resize", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  mounted: () => request("/mounted"),
};

export const mkfsAPI = {
  format: (payload) =>
    request("/mkfs", {
      method: "POST",
      body: JSON.stringify(payload),
    }),
};

export const authAPI = {
  login: async (payload) => {
    const data = await request("/auth/login", {
      method: "POST",
      body: JSON.stringify(payload),
    });

    const token = data?.session?.token || data?.data?.token;

    if (token) {
      localStorage.setItem("mia_token", token);
      localStorage.setItem("mia_session", JSON.stringify(data.session || data.data));
    }

    return data;
  },

  logout: async () => {
    const token = localStorage.getItem("mia_token");

    try {
      if (token) {
        await request("/auth/logout", {
          method: "POST",
          body: JSON.stringify({ token }),
        });
      }
    } finally {
      localStorage.removeItem("mia_token");
      localStorage.removeItem("mia_session");
    }
  },

  me: () => request("/auth/me"),

  getSession: () => {
    const raw = localStorage.getItem("mia_session");
    if (!raw) return null;

    try {
      return JSON.parse(raw);
    } catch {
      return null;
    }
  },

  isLoggedIn: () => Boolean(localStorage.getItem("mia_token")),
};

export const fsAPI = {
  list: ({ id, path = "/" }) =>
    request(`/fs/list?id=${encodeURIComponent(id)}&path=${encodeURIComponent(path)}`),

  read: ({ id, path }) =>
    request(`/fs/read?id=${encodeURIComponent(id)}&path=${encodeURIComponent(path)}`),

  mkdir: (payload) =>
    request("/fs/mkdir", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  mkfile: (payload) =>
    request("/fs/mkfile", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  remove: (payload) =>
    request("/fs/remove", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  edit: (payload) =>
    request("/fs/edit", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  rename: (payload) =>
    request("/fs/rename", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  copy: (payload) =>
    request("/fs/copy", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  move: (payload) =>
    request("/fs/move", {
      method: "POST",
      body: JSON.stringify(payload),
    }),
};

export const reportAPI = {
  list: () => request("/reports"),

  generate: (payload) =>
    request("/reports/generate", {
      method: "POST",
      body: JSON.stringify(payload),
    }),
};

export default {
  healthAPI,
  diskAPI,
  partitionAPI,
  mkfsAPI,
  authAPI,
  fsAPI,
  reportAPI,
};
