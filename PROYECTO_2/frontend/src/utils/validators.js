// utils/validators.js: contiene funciones de validación para formularios y entradas de usuario, incluyendo validaciones para campos requeridos, inicio de sesión, comandos y formularios de disco. Estas funciones devuelven mensajes de error apropiados si las entradas no cumplen con los criterios establecidos.
export function required(value) {
  return value !== undefined && value !== null && String(value).trim() !== "";
}

export function validateLogin(values) {
  const errors = {};

  if (!required(values.id)) {
    errors.id = "El ID de partición es obligatorio.";
  }

  if (!required(values.username)) {
    errors.username = "El usuario es obligatorio.";
  }

  if (!required(values.password)) {
    errors.password = "La contraseña es obligatoria.";
  }

  return errors;
}

export function validateCommand(command) {
  if (!required(command)) {
    return "Ingrese un comando antes de ejecutar.";
  }

  return "";
}

export function validateDiskForm(values) {
  const errors = {};

  if (!required(values.name)) {
    errors.name = "El nombre del disco es obligatorio.";
  }

  if (!required(values.path)) {
    errors.path = "La ruta del disco es obligatoria.";
  }

  if (!required(values.size) || Number(values.size) <= 0) {
    errors.size = "El tamaño debe ser mayor a cero.";
  }

  return errors;
}
