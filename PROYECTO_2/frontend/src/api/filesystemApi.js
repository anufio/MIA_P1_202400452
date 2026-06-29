import axiosClient from "./axiosClient";

function backendError(error) {
  return (
    error?.response?.data?.error ||
    error?.response?.data?.message ||
    error?.message ||
    "Error en sistema de archivos"
  );
}

async function safeRequest(callback) {
  try {
    return await callback();
  } catch (error) {
    throw new Error(backendError(error));
  }
}

function cleanText(value) {
  if (value === undefined || value === null) return "";
  return String(value).trim();
}

export function cleanFSPath(value) {
  let path = cleanText(value);

  if (!path) return "/";

  path = path.replaceAll('"', "").replaceAll("'", "").trim();

  if (!path.startsWith("/")) {
    path = `/${path}`;
  }

  path = path.replace(/\/+/g, "/");

  if (path.length > 1 && path.endsWith("/")) {
    path = path.slice(0, -1);
  }

  return path;
}

export function joinFSPath(parent = "/", name = "") {
  const cleanParent = cleanFSPath(parent);
  const cleanName = cleanText(name).replaceAll("/", "");

  if (!cleanName) return cleanParent;

  if (cleanParent === "/") {
    return `/${cleanName}`;
  }

  return `${cleanParent}/${cleanName}`;
}

async function getDefaultMountedId() {
  const response = await axiosClient.get("/mounted");
  const mounted = response.data.mounted || response.data.data || [];

  if (!Array.isArray(mounted) || mounted.length === 0) {
    throw new Error("No hay particiones montadas.");
  }

  return mounted[0].id;
}

async function resolveMountedId(id) {
  return id || (await getDefaultMountedId());
}

export async function listDirectory(id, path = "/") {
  return safeRequest(async () => {
    const mountedId = await resolveMountedId(id);

    const response = await axiosClient.post("/fs/list", {
      id: mountedId,
      path: cleanFSPath(path)
    });

    return (
      response.data.items ||
      response.data.files ||
      response.data.data?.items ||
      response.data.data ||
      []
    );
  });
}

export async function readFile(id, path) {
  return safeRequest(async () => {
    const mountedId = await resolveMountedId(id);

    const response = await axiosClient.post("/fs/read", {
      id: mountedId,
      path: cleanFSPath(path)
    });

    return response.data.file || response.data.data || response.data;
  });
}

export async function createFolder(id, path) {
  return safeRequest(async () => {
    const mountedId = await resolveMountedId(id);

    const response = await axiosClient.post("/fs/mkdir", {
      id: mountedId,
      path: cleanFSPath(path),
      recursive: true,
      parents: true,
      p: true
    });

    return response.data.item || response.data.data || response.data;
  });
}

export async function createFile(id, path, content = "", size = 0) {
  return safeRequest(async () => {
    const mountedId = await resolveMountedId(id);

    const response = await axiosClient.post("/fs/mkfile", {
      id: mountedId,
      path: cleanFSPath(path),
      content,
      cont: content,
      size: Number(size) || 0,
      recursive: true,
      parents: true,
      p: true
    });

    return response.data.item || response.data.data || response.data;
  });
}

export async function removeEntry(id, path) {
  return safeRequest(async () => {
    const mountedId = await resolveMountedId(id);

    const response = await axiosClient.post("/fs/remove", {
      id: mountedId,
      path: cleanFSPath(path),
      recursive: true
    });

    return response.data;
  });
}

export async function editFile(id, path, content) {
  return safeRequest(async () => {
    const mountedId = await resolveMountedId(id);

    const response = await axiosClient.post("/fs/edit", {
      id: mountedId,
      path: cleanFSPath(path),
      content,
      cont: content
    });

    return response.data.file || response.data.data || response.data;
  });
}

export async function renameEntry(id, path, name) {
  return safeRequest(async () => {
    const mountedId = await resolveMountedId(id);

    const response = await axiosClient.post("/fs/rename", {
      id: mountedId,
      path: cleanFSPath(path),
      name: cleanText(name)
    });

    return response.data.item || response.data.data || response.data;
  });
}

export async function copyEntry(id, from, to) {
  return safeRequest(async () => {
    const mountedId = await resolveMountedId(id);

    const response = await axiosClient.post("/fs/copy", {
      id: mountedId,
      from: cleanFSPath(from),
      to: cleanFSPath(to)
    });

    return response.data.item || response.data.data || response.data;
  });
}

export async function moveEntry(id, from, to) {
  return safeRequest(async () => {
    const mountedId = await resolveMountedId(id);

    const response = await axiosClient.post("/fs/move", {
      id: mountedId,
      from: cleanFSPath(from),
      to: cleanFSPath(to)
    });

    return response.data.item || response.data.data || response.data;
  });
}

/* Compatibilidad con páginas viejas del frontend. */
export async function getFileContent(path, id = "") {
  const file = await readFile(id, path);

  return {
    ...file,
    content:
      file.content ||
      file.text ||
      file.body ||
      file.value ||
      ""
  };
}

export async function getDirectoryContent(path = "/", id = "") {
  return listDirectory(id, path);
}

export async function getDirectoryItems(path = "/", id = "") {
  return listDirectory(id, path);
}