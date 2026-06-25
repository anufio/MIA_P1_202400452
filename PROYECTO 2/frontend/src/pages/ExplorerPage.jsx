import { useEffect, useState } from "react";
import { getMountedPartitions } from "../api/diskApi";
import { listDirectory, readFile } from "../api/filesystemApi";
// ExplorerPage.jsx: contiene la página principal del explorador de archivos, mostrando un encabezado y un formulario para seleccionar una partición montada y navegar por el sistema de archivos. También incluye una lista de archivos y carpetas en la ruta actual, así como un visor de contenido de archivos seleccionados. Utiliza los hooks useDisks y useFilesystem para manejar la carga de particiones montadas y la navegación por el sistema de archivos.
function joinPath(base, name) {
  if (base === "/") return "/" + name;
  return base + "/" + name;
}

function parentPath(path) {
  if (!path || path === "/") return "/";
  const parts = path.split("/").filter(Boolean);
  parts.pop();
  return "/" + parts.join("/");
}

function normalizeType(item) {
  const raw = String(item.type || item.kind || "").toLowerCase();

  if (raw.includes("dir") || raw.includes("folder") || raw === "carpeta") {
    return "folder";
  }

  if (item.isDir || item.directory) {
    return "folder";
  }

  return "file";
}

export default function ExplorerPage() {
  const [mounted, setMounted] = useState([]);
  const [selectedId, setSelectedId] = useState("");
  const [currentPath, setCurrentPath] = useState("/");
  const [items, setItems] = useState([]);
  const [fileContent, setFileContent] = useState(null);
  const [selectedFile, setSelectedFile] = useState("");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  async function loadMounted() {
    try {
      const data = await getMountedPartitions();
      const list = Array.isArray(data) ? data : [];

      setMounted(list);

      if (!selectedId && list.length > 0) {
        setSelectedId(list[0].id);
      }
    } catch (error) {
      setMessage(error.message);
    }
  }

  async function loadDirectory(path = currentPath, id = selectedId) {
    if (!id) {
      setItems([]);
      setMessage("Seleccione una partición montada.");
      return;
    }

    try {
      setLoading(true);
      setMessage("");
      setFileContent(null);
      setSelectedFile("");

      const data = await listDirectory(id, path);
      setItems(Array.isArray(data) ? data : []);
      setCurrentPath(path);
    } catch (error) {
      setItems([]);
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  async function openItem(item) {
    const name = item.name || item.fileName || item.filename;
    const path = item.path || joinPath(currentPath, name);
    const type = normalizeType(item);

    if (type === "folder") {
      await loadDirectory(path);
      return;
    }

    try {
      setLoading(true);
      setMessage("");

      const data = await readFile(selectedId, path);

      setSelectedFile(path);
      setFileContent(
        data.content ||
          data.text ||
          data.body ||
          data.value ||
          "Archivo sin contenido visible."
      );
    } catch (error) {
      setMessage(error.message);
      setFileContent(null);
    } finally {
      setLoading(false);
    }
  }

  function goBack() {
    loadDirectory(parentPath(currentPath));
  }

  function goRoot() {
    loadDirectory("/");
  }

  useEffect(() => {
    loadMounted();
  }, []);

  useEffect(() => {
    if (selectedId) {
      loadDirectory("/", selectedId);
    }
  }, [selectedId]);

  return (
    <main style={styles.page}>
      <h1 style={styles.title}>Explorador EXT2</h1>

      <section style={styles.box}>
        <h2 style={styles.subtitle}>Partición montada</h2>

        <div style={styles.form}>
          <select
            style={styles.input}
            value={selectedId}
            onChange={(e) => setSelectedId(e.target.value)}
          >
            <option value="">Seleccione una partición montada</option>
            {mounted.map((part) => (
              <option key={part.id} value={part.id}>
                {part.id} — {part.name}
              </option>
            ))}
          </select>

          <div style={styles.row}>
            <button style={styles.secondaryButton} onClick={loadMounted}>
              Actualizar montadas
            </button>

            <button style={styles.secondaryButton} onClick={() => loadDirectory()}>
              Actualizar carpeta
            </button>

            <button style={styles.secondaryButton} onClick={goBack}>
              Regresar
            </button>

            <button style={styles.secondaryButton} onClick={goRoot}>
              Ir a raíz
            </button>
          </div>
        </div>

        <p style={styles.path}>
          Ruta actual: <strong>{currentPath}</strong>
        </p>

        {message && <p style={styles.message}>{message}</p>}
      </section>

      <section style={styles.box}>
        <h2 style={styles.subtitle}>Contenido</h2>

        {loading && <p>Cargando...</p>}

        {!loading && items.length === 0 && (
          <p>No hay archivos o carpetas en esta ruta.</p>
        )}

        <div style={styles.grid}>
          {items.map((item, index) => {
            const name = item.name || item.fileName || item.filename || `item-${index}`;
            const type = normalizeType(item);

            return (
              <button
                key={`${name}-${index}`}
                style={styles.item}
                onClick={() => openItem(item)}
              >
                <span style={styles.itemIcon}>
                  {type === "folder" ? "□" : "▤"}
                </span>

                <span>{name}</span>

                <small style={styles.itemType}>
                  {type === "folder" ? "Carpeta" : "Archivo"}
                </small>
              </button>
            );
          })}
        </div>
      </section>

      {fileContent !== null && (
        <section style={styles.box}>
          <h2 style={styles.subtitle}>Archivo: {selectedFile}</h2>

          <pre style={styles.viewer}>{fileContent}</pre>
        </section>
      )}
    </main>
  );
}

const styles = {
  page: {
    minHeight: "100vh",
    padding: "32px",
    background: "#fff7fa",
    color: "#3b202b",
  },
  title: {
    marginTop: 0,
    marginBottom: "20px",
  },
  subtitle: {
    marginTop: 0,
  },
  box: {
    background: "#ffffff",
    border: "1px solid #e5bcc9",
    borderRadius: "16px",
    padding: "22px",
    marginBottom: "22px",
  },
  form: {
    display: "grid",
    gap: "12px",
  },
  input: {
    padding: "11px",
    borderRadius: "10px",
    border: "1px solid #d7a9b9",
    fontSize: "15px",
    width: "100%",
    boxSizing: "border-box",
  },
  row: {
    display: "flex",
    gap: "10px",
    flexWrap: "wrap",
  },
  secondaryButton: {
    padding: "10px 12px",
    borderRadius: "10px",
    border: "1px solid #b85c7a",
    background: "white",
    color: "#b85c7a",
    fontWeight: "700",
    cursor: "pointer",
  },
  path: {
    marginTop: "14px",
  },
  message: {
    background: "#f5d7e1",
    padding: "12px",
    borderRadius: "10px",
    marginTop: "14px",
  },
  grid: {
    display: "grid",
    gap: "12px",
    gridTemplateColumns: "repeat(auto-fill, minmax(180px, 1fr))",
  },
  item: {
    textAlign: "left",
    padding: "16px",
    borderRadius: "14px",
    border: "1px solid #e5bcc9",
    background: "#fffafd",
    cursor: "pointer",
    display: "grid",
    gap: "6px",
    color: "#3b202b",
  },
  itemIcon: {
    fontSize: "24px",
  },
  itemType: {
    color: "#765062",
  },
  viewer: {
    background: "#2b2025",
    color: "white",
    padding: "18px",
    borderRadius: "14px",
    whiteSpace: "pre-wrap",
    overflow: "auto",
  },
};
