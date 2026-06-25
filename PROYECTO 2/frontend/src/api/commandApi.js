import axiosClient from "./axiosClient";
// commandApi.js: contiene las funciones para manejar la ejecución de comandos en el sistema de archivos, incluyendo la limpieza y validación de rutas y nombres, el análisis de comandos y banderas, y la comunicación con el backend a través de axiosClient.

function backendError(error) {
  return (
    error?.response?.data?.error ||
    error?.response?.data?.message ||
    error?.message ||
    "Error ejecutando comando"
  );
}

function cleanText(value) {
  if (value === undefined || value === null) return "";

  let clean = String(value).trim();

  if (
    (clean.startsWith('"') && clean.endsWith('"')) ||
    (clean.startsWith("'") && clean.endsWith("'"))
  ) {
    clean = clean.slice(1, -1);
  }

  return clean;
}

function cleanPath(value) {
  let clean = cleanText(value);

  clean = clean.replaceAll('"', "");
  clean = clean.replaceAll("'", "");
  clean = clean.trim();

  if (!clean.startsWith("/")) {
    clean = "/" + clean;
  }

  return clean;
}

function parseCommand(input) {
  const clean = input.trim();
  const firstSpace = clean.search(/\s/);

  const command =
    firstSpace === -1 ? clean.toLowerCase() : clean.slice(0, firstSpace).toLowerCase();

  const rest = firstSpace === -1 ? "" : clean.slice(firstSpace + 1);

  const flags = {};
  const regex = /-(\w+)=("([^"]*)"|'([^']*)'|[^\s]+)/g;

  let match;

  while ((match = regex.exec(rest)) !== null) {
    const key = match[1].toLowerCase();
    const rawValue = match[3] ?? match[4] ?? match[2] ?? "";
    flags[key] = cleanText(rawValue);
  }

  return { command, flags };
}

function requireValue(value, name) {
  if (!value || String(value).trim() === "") {
    throw new Error(`Falta el parámetro ${name}`);
  }

  return value;
}

async function request(endpoint, payload) {
  try {
    const response = await axiosClient.post(endpoint, payload);
    return response.data;
  } catch (error) {
    throw new Error(backendError(error));
  }
}

export async function executeCommand(commandText, mountedId) {
  if (!commandText.trim()) {
    throw new Error("Ingrese un comando.");
  }

  if (!mountedId) {
    throw new Error("Seleccione una partición montada antes de ejecutar comandos.");
  }

  const { command, flags } = parseCommand(commandText);

  switch (command) {
    case "mkdir": {
      const path = cleanPath(flags.path);

      const result = await request("/fs/mkdir", {
        id: mountedId,
        path: requireValue(path, "-path"),
        recursive: true,
        p: true,
      });

      return {
        success: true,
        output: result.message || `Carpeta creada: ${path}`,
      };
    }

    case "mkfile": {
      const path = cleanPath(flags.path);
      const content = cleanText(flags.cont || flags.content || "");
      const size = Number(flags.size || 0);

      const result = await request("/fs/mkfile", {
        id: mountedId,
        path: requireValue(path, "-path"),
        content,
        cont: content,
        size,
        recursive: true,
        p: true,
      });

      return {
        success: true,
        output: result.message || `Archivo creado: ${path}`,
      };
    }

    case "edit": {
      const path = cleanPath(flags.path);
      const content = cleanText(flags.cont || flags.content);

      const result = await request("/fs/edit", {
        id: mountedId,
        path: requireValue(path, "-path"),
        content: requireValue(content, "-cont"),
        cont: content,
      });

      return {
        success: true,
        output: result.message || `Archivo editado: ${path}`,
      };
    }

    case "rename": {
      const path = cleanPath(flags.path);
      const name = cleanText(flags.name || flags.newname);

      const result = await request("/fs/rename", {
        id: mountedId,
        path: requireValue(path, "-path"),
        name: requireValue(name, "-name"),
        newName: name,
      });

      return {
        success: true,
        output: result.message || `Renombrado: ${path} -> ${name}`,
      };
    }

    case "copy": {
      const path = cleanPath(flags.path);
      const dest = cleanPath(flags.dest || flags.destino || flags.destination);

      const result = await request("/fs/copy", {
        id: mountedId,
        path: requireValue(path, "-path"),
        destPath: requireValue(dest, "-dest"),
        destination: dest,
        targetPath: dest,
      });

      return {
        success: true,
        output: result.message || `Copiado: ${path} -> ${dest}`,
      };
    }

    case "move": {
      const path = cleanPath(flags.path);
      const dest = cleanPath(flags.dest || flags.destino || flags.destination);

      const result = await request("/fs/move", {
        id: mountedId,
        path: requireValue(path, "-path"),
        destPath: requireValue(dest, "-dest"),
        destination: dest,
        targetPath: dest,
      });

      return {
        success: true,
        output: result.message || `Movido: ${path} -> ${dest}`,
      };
    }

    case "remove": {
      const path = cleanPath(flags.path);

      const result = await request("/fs/remove", {
        id: mountedId,
        path: requireValue(path, "-path"),
        recursive: true,
      });

      return {
        success: true,
        output: result.message || `Eliminado: ${path}`,
      };
    }

    default:
      throw new Error("Comando no reconocido. Use: mkdir, mkfile, edit, rename, copy, move o remove.");
  }
}
