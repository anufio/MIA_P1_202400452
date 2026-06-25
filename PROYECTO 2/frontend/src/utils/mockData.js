// utils/mockData.js: contiene datos de ejemplo para discos, particiones y directorios, así como funciones para obtener directorios, contenido de archivos y reportes simulados. Estos datos se utilizan para pruebas y desarrollo de la aplicación sin necesidad de conectarse a un backend real.
export const mockDisks = [
  {
    id: "disk-1",
    name: "d1.dsk",
    path: "/home/anufio/mia/cali/d1.dsk",
    size: 78643200,
    fit: "FF"
  },
  {
    id: "disk-2",
    name: "d2.dsk",
    path: "/home/anufio/mia/cali/d2.dsk",
    size: 11534336,
    fit: "BF"
  },
  {
    id: "disk-3",
    name: "d3.dsk",
    path: "/home/anufio/mia/cali/d3.dsk",
    size: 104857600,
    fit: "WF"
  }
];

export const mockPartitions = [
  {
    id: "521A",
    diskId: "disk-1",
    name: "Part1-D1",
    type: "P",
    fit: "FF",
    start: 164,
    size: 26214400,
    mounted: true,
    formatted: true
  },
  {
    id: "522A",
    diskId: "disk-1",
    name: "Part2-D1",
    type: "P",
    fit: "WF",
    start: 26214564,
    size: 26214400,
    mounted: false,
    formatted: false
  },
  {
    id: "523A",
    diskId: "disk-3",
    name: "Part4D3",
    type: "P",
    fit: "WF",
    start: 52428800,
    size: 52428800,
    mounted: true,
    formatted: true
  }
];

const rootItems = [
  {
    name: "home",
    type: "folder",
    path: "/home",
    size: 0,
    permissions: "664"
  },
  {
    name: "users.txt",
    type: "file",
    path: "/users.txt",
    size: 128,
    permissions: "664"
  }
];

const homeItems = [
  {
    name: "archivos",
    type: "folder",
    path: "/home/archivos",
    size: 0,
    permissions: "664"
  },
  {
    name: "b1.txt",
    type: "file",
    path: "/home/b1.txt",
    size: 75,
    permissions: "664"
  },
  {
    name: "b1_cont.txt",
    type: "file",
    path: "/home/b1_cont.txt",
    size: 84,
    permissions: "664"
  }
];

const archivosItems = [
  {
    name: "mia",
    type: "folder",
    path: "/home/archivos/mia",
    size: 0,
    permissions: "664"
  }
];

const miaItems = [
  {
    name: "fase2",
    type: "folder",
    path: "/home/archivos/mia/fase2",
    size: 0,
    permissions: "664"
  },
  {
    name: "carpeta2",
    type: "folder",
    path: "/home/archivos/mia/carpeta2",
    size: 0,
    permissions: "664"
  }
];

const fase2Items = [
  {
    name: "a1",
    type: "folder",
    path: "/home/archivos/mia/fase2/a1",
    size: 0,
    permissions: "664"
  },
  {
    name: "a2",
    type: "folder",
    path: "/home/archivos/mia/fase2/a2",
    size: 0,
    permissions: "664"
  },
  {
    name: "a3",
    type: "folder",
    path: "/home/archivos/mia/fase2/a3",
    size: 0,
    permissions: "664"
  }
];

export function getMockDirectory(path) {
  if (path === "/") return rootItems;
  if (path === "/home") return homeItems;
  if (path === "/home/archivos") return archivosItems;
  if (path === "/home/archivos/mia") return miaItems;
  if (path === "/home/archivos/mia/fase2") return fase2Items;

  return [];
}

export function getMockFileContent(path) {
  if (path === "/users.txt") {
    return [
      "1,G,root",
      "1,U,root,root,123",
      "1,G,Archivos",
      "1,G,AuxEsAmor",
      "1,G,Arquizzz",
      "1,U,Archivos,user1,abc",
      "1,U,AuxEsAmor,user2,abc",
      "0,U,Arquizzz,user3,abc",
      "1,U,Arquizzz,user4,abc"
    ].join("\n");
  }

  if (path === "/home/b1.txt") {
    return "012345678901234567890123456789012345678901234567890123456789012345678901234";
  }

  if (path === "/home/b1_cont.txt") {
    return "Contenido de prueba para mkfile con parametro cont.";
  }

  return "Archivo de ejemplo sin contenido cargado desde backend.";
}

export function getMockReport(type) {
  return {
    type,
    title: `Reporte ${type}`,
    format: "text",
    content:
      `Reporte ${type}\n\n` +
      "Este visor ya está preparado para mostrar reportes generados por el backend.\n" +
      "Cuando el backend esté conectado, aquí se cargará SVG, PNG o texto según corresponda.\n"
  };
}
