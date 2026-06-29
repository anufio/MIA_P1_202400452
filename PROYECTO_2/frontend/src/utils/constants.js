// utils/constants.js: contiene constantes utilizadas en la aplicación, incluyendo información de la aplicación, rutas, símbolos y tipos de reportes. También define la URL base para las solicitudes a la API del backend.
export const APP_INFO = {
  name: "MIA Proyecto 2",
  student: "Ana Lucía Nufio Roblero",
  carnet: "202400452",
  course: "Manejo e Implementación de Archivos"
};

export const ROUTES = {
  home: "/",
  login: "/login",
  disks: "/disks",
  partitions: "/partitions",
  explorer: "/explorer",
  fileViewer: "/file-viewer",
  reports: "/reports",
  commands: "/commands"
};

export const SYMBOLS = {
  home: "⌂",
  login: "⌁",
  disk: "◉",
  partition: "◌",
  folder: "□",
  file: "▤",
  report: "▧",
  console: "›",
  edit: "✎",
  delete: "×",
  save: "✓",
  back: "←",
  next: "→",
  user: "◇",
  warning: "!",
  dot: "•"
};

export const REPORT_TYPES = [
  { value: "mbr", label: "MBR" },
  { value: "disk", label: "Disk" },
  { value: "tree", label: "Tree" },
  { value: "inode", label: "Inodos" },
  { value: "block", label: "Bloques" },
  { value: "bm_inode", label: "Bitmap Inodos" },
  { value: "bm_block", label: "Bitmap Bloques" },
  { value: "sb", label: "Superbloque" }
];

export const API_BASE_URL =
  import.meta.env.VITE_API_URL || "http://localhost:8080/api";
