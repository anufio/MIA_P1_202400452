// utils/formatters.js: contiene funciones para formatear bytes, normalizar rutas, obtener la ruta padre, unir rutas y obtener el nombre de archivo desde una ruta. Estas funciones son útiles para manejar y mostrar información relacionada con archivos y directorios en la aplicación.
export function formatBytes(bytes) {
  const value = Number(bytes || 0);

  if (value >= 1024 * 1024) {
    return `${(value / (1024 * 1024)).toFixed(2)} MB`;
  }

  if (value >= 1024) {
    return `${(value / 1024).toFixed(2)} KB`;
  }

  return `${value} B`;
}

export function normalizePath(path) {
  if (!path || path.trim() === "") {
    return "/";
  }

  let result = path.trim();

  if (!result.startsWith("/")) {
    result = `/${result}`;
  }

  if (result.length > 1 && result.endsWith("/")) {
    result = result.slice(0, -1);
  }

  return result;
}

export function parentPath(path) {
  const normalized = normalizePath(path);

  if (normalized === "/") {
    return "/";
  }

  const parts = normalized.split("/").filter(Boolean);
  parts.pop();

  if (parts.length === 0) {
    return "/";
  }

  return `/${parts.join("/")}`;
}

export function joinPath(base, name) {
  const cleanBase = normalizePath(base);
  const cleanName = String(name || "").replace(/^\/+/, "");

  if (cleanBase === "/") {
    return `/${cleanName}`;
  }

  return `${cleanBase}/${cleanName}`;
}

export function getFileNameFromPath(path) {
  const normalized = normalizePath(path);

  if (normalized === "/") {
    return "/";
  }

  const parts = normalized.split("/").filter(Boolean);
  return parts[parts.length - 1] || "/";
}
