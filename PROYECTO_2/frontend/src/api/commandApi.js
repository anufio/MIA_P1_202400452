import axiosClient from "./axiosClient";

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

  if (!clean) {
    return "";
  }

  if (!clean.startsWith("/")) {
    clean = "/" + clean;
  }

  return clean;
}

function cleanDiskPath(value) {
  return cleanText(value).replaceAll('"', "").replaceAll("'", "").trim();
}

function parseCommand(input) {
  const clean = input.trim();
  const firstSpace = clean.search(/\s/);

  const command =
    firstSpace === -1 ? clean.toLowerCase() : clean.slice(0, firstSpace).toLowerCase();

  const rest = firstSpace === -1 ? "" : clean.slice(firstSpace + 1);

  const flags = {};
  const regex = /-([a-zA-Z_][\w-]*)=("([^"]*)"|'([^']*)'|[^\s]+)/g;

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

function requireNumber(value, name) {
  const number = Number(value);

  if (!Number.isFinite(number) || number === 0) {
    throw new Error(`El parámetro ${name} debe ser un número diferente de cero.`);
  }

  return number;
}

async function request(endpoint, payload) {
  try {
    const response = await axiosClient.post(endpoint, payload);
    return response.data;
  } catch (error) {
    throw new Error(backendError(error));
  }
}

function buildOutput(result, fallback) {
  return result?.message || result?.output || fallback;
}

async function executeFdisk(flags) {
  const diskPath = cleanDiskPath(flags.path || flags.diskpath);
  const name = cleanText(flags.name);

  requireValue(diskPath, "-path");
  requireValue(name, "-name");

  if (flags.add !== undefined) {
    const add = requireNumber(flags.add, "-add");
    const unit = cleanText(flags.unit || "K").toUpperCase();

    const result = await request("/partitions/resize", {
      path: diskPath,
      diskPath,
      name,
      add,
      unit
    });

    return {
      success: true,
      output: buildOutput(result, `FDISK ADD aplicado en ${name}.`)
    };
  }

  if (flags.delete !== undefined) {
    const deleteType = cleanText(flags.delete || "fast").toLowerCase();

    if (deleteType !== "fast" && deleteType !== "full") {
      throw new Error("El parámetro -delete debe ser fast o full.");
    }

    const result = await request("/partitions/delete", {
      path: diskPath,
      diskPath,
      name,
      deleteType,
      delete: deleteType
    });

    return {
      success: true,
      output: buildOutput(result, `FDISK DELETE aplicado en ${name}.`)
    };
  }

  throw new Error("Para fdisk use -add o -delete.");
}

export async function executeCommand(commandText, mountedId) {
  if (!commandText.trim()) {
    throw new Error("Ingrese un comando.");
  }

  const { command, flags } = parseCommand(commandText);

  if (command === "fdisk") {
    return executeFdisk(flags);
  }

  if (!mountedId) {
    throw new Error("Seleccione una partición montada antes de ejecutar comandos de sistema de archivos.");
  }

  switch (command) {
    case "mkdir": {
      const path = cleanPath(flags.path);

      const result = await request("/fs/mkdir", {
        id: mountedId,
        path: requireValue(path, "-path"),
        recursive: true,
        p: true
      });

      return {
        success: true,
        output: buildOutput(result, `Carpeta creada: ${path}`)
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
        p: true
      });

      return {
        success: true,
        output: buildOutput(result, `Archivo creado: ${path}`)
      };
    }

    case "edit": {
      const path = cleanPath(flags.path);
      const contenido = cleanText(flags.contenido || "");
      const content = cleanText(flags.cont || flags.content || "");

      const payload = {
        id: mountedId,
        path: requireValue(path, "-path")
      };

      if (contenido) {
        payload.contenido = contenido;
      } else {
        payload.content = requireValue(content, "-contenido o -cont");
        payload.cont = content;
      }

      const result = await request("/fs/edit", payload);

      return {
        success: true,
        output: buildOutput(result, `Archivo editado: ${path}`)
      };
    }

    case "rename": {
      const path = cleanPath(flags.path);
      const name = cleanText(flags.name || flags.newname);

      const result = await request("/fs/rename", {
        id: mountedId,
        path: requireValue(path, "-path"),
        name: requireValue(name, "-name"),
        newName: name
      });

      return {
        success: true,
        output: buildOutput(result, `Renombrado: ${path} -> ${name}`)
      };
    }

    case "copy": {
      const path = cleanPath(flags.path);
      const dest = cleanPath(flags.dest || flags.destino || flags.destination);

      const result = await request("/fs/copy", {
        id: mountedId,
        from: requireValue(path, "-path"),
        to: requireValue(dest, "-destino")
      });

      return {
        success: true,
        output: buildOutput(result, `Copiado: ${path} -> ${dest}`)
      };
    }

    case "move": {
      const path = cleanPath(flags.path);
      const dest = cleanPath(flags.dest || flags.destino || flags.destination);

      const result = await request("/fs/move", {
        id: mountedId,
        from: requireValue(path, "-path"),
        to: requireValue(dest, "-destino")
      });

      return {
        success: true,
        output: buildOutput(result, `Movido: ${path} -> ${dest}`)
      };
    }

    case "remove": {
      const path = cleanPath(flags.path);

      const result = await request("/fs/remove", {
        id: mountedId,
        path: requireValue(path, "-path"),
        recursive: true
      });

      return {
        success: true,
        output: buildOutput(result, `Eliminado: ${path}`)
      };
    }

    default:
      throw new Error("Comando no reconocido. Use: fdisk, mkdir, mkfile, edit, rename, copy, move o remove.");
  }
}