import axiosClient from "./axiosClient";
//filesystemApi.js: contiene las funciones para manejar las operaciones del sistema de archivos, incluyendo la lectura de archivos y la obtención de contenido de directorios. Todas las funciones se comunican con el backend a través de axiosClient.

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

async function getDefaultMountedId() {
  const response = await axiosClient.get("/mounted");
  const mounted = response.data.mounted || response.data.data || [];

  if (!Array.isArray(mounted) || mounted.length === 0) {
    throw new Error("No hay particiones montadas.");
  }

  return mounted[0].id;
}

export async function listDirectory(id, path = "/") {
  return safeRequest(async () => {
    const mountedId = id || (await getDefaultMountedId());

    const response = await axiosClient.post("/fs/list", {
      id: mountedId,
      path,
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
    const mountedId = id || (await getDefaultMountedId());

    const response = await axiosClient.post("/fs/read", {
      id: mountedId,
      path,
    });

    return response.data.file || response.data.data || response.data;
  });
}

/*
  Compatibilidad con páginas viejas del frontend.
  FileViewerPage.jsx todavía usa getFileContent(path).
*/
export async function getFileContent(path, id = "") {
  const file = await readFile(id, path);

  return {
    ...file,
    content:
      file.content ||
      file.text ||
      file.body ||
      file.value ||
      "",
  };
}

export async function getDirectoryContent(path = "/", id = "") {
  return listDirectory(id, path);
}

export async function getDirectoryItems(path = "/", id = "") {
  return listDirectory(id, path);
}
