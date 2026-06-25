// services/index.js: exporta todas las funciones de la API desde el archivo api.js, permitiendo que otros módulos importen estas funciones de manera centralizada.
export {
  healthAPI,
  diskAPI,
  partitionAPI,
  mkfsAPI,
  authAPI,
  fsAPI,
  reportAPI,
} from "./api";
