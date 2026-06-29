import { useEffect, useMemo, useState } from "react";
import { getMountedPartitions } from "../api/diskApi";
import {
  cleanFSPath,
  copyEntry,
  createFile,
  createFolder,
  editFile,
  joinFSPath,
  listDirectory,
  moveEntry,
  readFile,
  removeEntry,
  renameEntry
} from "../api/filesystemApi";
import { useAuth } from "../hooks/useAuth";

function parentPath(path) {
  if (!path || path === "/") return "/";
  const parts = path.split("/").filter(Boolean);
  parts.pop();
  return `/${parts.join("/")}` || "/";
}

function baseName(path) {
  if (!path || path === "/") return "";
  const parts = path.split("/").filter(Boolean);
  return parts[parts.length - 1] || "";
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

function normalizeItem(item, currentPath, index = 0) {
  const name = item.name || item.fileName || item.filename || `item-${index}`;
  const path = item.path || joinFSPath(currentPath, name);
  const type = normalizeType(item);

  return {
    ...item,
    name,
    path,
    type
  };
}

function readableSize(size) {
  const number = Number(size || 0);

  if (!Number.isFinite(number) || number <= 0) return "0 bytes";

  if (number >= 1024 * 1024) return `${(number / 1024 / 1024).toFixed(2)} MB`;
  if (number >= 1024) return `${(number / 1024).toFixed(2)} KB`;

  return `${number} bytes`;
}

export default function ExplorerPage() {
  const { session } = useAuth();
  const [mounted, setMounted] = useState([]);
  const [selectedId, setSelectedId] = useState(session?.partId || session?.id || "");
  const [currentPath, setCurrentPath] = useState("/");
  const [items, setItems] = useState([]);
  const [selectedItem, setSelectedItem] = useState(null);
  const [fileContent, setFileContent] = useState("");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  const [folderName, setFolderName] = useState("");
  const [fileName, setFileName] = useState("");
  const [newFileContent, setNewFileContent] = useState("");
  const [newName, setNewName] = useState("");
  const [destinationPath, setDestinationPath] = useState("");
  const [editContent, setEditContent] = useState("");

  const selectedLabel = useMemo(() => {
    if (!selectedItem) return "No hay archivo o carpeta seleccionada.";
    return `${selectedItem.name} (${selectedItem.type === "folder" ? "carpeta" : "archivo"})`;
  }, [selectedItem]);

  function resetSelected() {
    setSelectedItem(null);
    setFileContent("");
    setEditContent("");
    setNewName("");
    setDestinationPath("");
  }

  function selectItem(item) {
    setSelectedItem(item);
    setNewName(item.name || "");
    setDestinationPath(joinFSPath(currentPath, item.name || baseName(item.path)));

    if (item.type !== "file") {
      setFileContent("");
      setEditContent("");
    }
  }

  async function loadMounted() {
    try {
      const data = await getMountedPartitions();
      const list = Array.isArray(data) ? data : [];
      const sessionId = session?.partId || session?.id || "";

      setMounted(list);

      if (sessionId) {
        setSelectedId(sessionId);
      } else if (!selectedId && list.length > 0) {
        setSelectedId(list[0].id);
      }
    } catch (error) {
      setMessage(error.message);
    }
  }

  async function loadDirectory(path = currentPath, id = selectedId) {
    if (!id) {
      setItems([]);
      setMessage("Seleccione una partición montada o inicie sesión en una partición formateada.");
      return;
    }

    try {
      setLoading(true);
      setMessage("");
      resetSelected();

      const cleanPath = cleanFSPath(path);
      const data = await listDirectory(id, cleanPath);
      const list = Array.isArray(data) ? data : [];

      setItems(list.map((item, index) => normalizeItem(item, cleanPath, index)));
      setCurrentPath(cleanPath);
    } catch (error) {
      setItems([]);
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  async function openItem(item) {
    const normalized = normalizeItem(item, currentPath);
    selectItem(normalized);

    if (normalized.type === "folder") {
      await loadDirectory(normalized.path);
      return;
    }

    try {
      setLoading(true);
      setMessage("");

      const data = await readFile(selectedId, normalized.path);
      const content =
        data.content ||
        data.text ||
        data.body ||
        data.value ||
        "";

      setFileContent(content || "Archivo sin contenido visible.");
      setEditContent(content || "");
    } catch (error) {
      setMessage(error.message);
      setFileContent("");
    } finally {
      setLoading(false);
    }
  }

  async function handleCreateFolder(event) {
    event.preventDefault();

    if (!folderName.trim()) {
      setMessage("Ingrese el nombre de la carpeta.");
      return;
    }

    try {
      setLoading(true);
      setMessage("");
      await createFolder(selectedId, joinFSPath(currentPath, folderName));
      setFolderName("");
      await loadDirectory(currentPath);
      setMessage("Carpeta creada correctamente.");
    } catch (error) {
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleCreateFile(event) {
    event.preventDefault();

    if (!fileName.trim()) {
      setMessage("Ingrese el nombre del archivo.");
      return;
    }

    try {
      setLoading(true);
      setMessage("");
      await createFile(selectedId, joinFSPath(currentPath, fileName), newFileContent);
      setFileName("");
      setNewFileContent("");
      await loadDirectory(currentPath);
      setMessage("Archivo creado correctamente.");
    } catch (error) {
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleRemove() {
    if (!selectedItem) {
      setMessage("Seleccione una entrada para eliminar.");
      return;
    }

    const confirmText = `¿Eliminar ${selectedItem.path}?`;
    if (!window.confirm(confirmText)) return;

    try {
      setLoading(true);
      setMessage("");
      await removeEntry(selectedId, selectedItem.path);
      await loadDirectory(currentPath);
      setMessage("Entrada eliminada correctamente.");
    } catch (error) {
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleRename(event) {
    event.preventDefault();

    if (!selectedItem) {
      setMessage("Seleccione una entrada para renombrar.");
      return;
    }

    if (!newName.trim()) {
      setMessage("Ingrese el nuevo nombre.");
      return;
    }

    try {
      setLoading(true);
      setMessage("");
      await renameEntry(selectedId, selectedItem.path, newName);
      await loadDirectory(currentPath);
      setMessage("Entrada renombrada correctamente.");
    } catch (error) {
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleEdit(event) {
    event.preventDefault();

    if (!selectedItem || selectedItem.type !== "file") {
      setMessage("Seleccione un archivo para editar.");
      return;
    }

    try {
      setLoading(true);
      setMessage("");
      await editFile(selectedId, selectedItem.path, editContent);
      const data = await readFile(selectedId, selectedItem.path);
      const content = data.content || data.text || data.body || data.value || "";
      setFileContent(content || "Archivo sin contenido visible.");
      setEditContent(content || "");
      await loadDirectory(currentPath);
      setMessage("Archivo editado correctamente.");
    } catch (error) {
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleCopy(event) {
    event.preventDefault();

    if (!selectedItem) {
      setMessage("Seleccione una entrada para copiar.");
      return;
    }

    if (!destinationPath.trim()) {
      setMessage("Ingrese la ruta destino completa.");
      return;
    }

    try {
      setLoading(true);
      setMessage("");
      await copyEntry(selectedId, selectedItem.path, destinationPath);
      await loadDirectory(currentPath);
      setMessage("Entrada copiada correctamente.");
    } catch (error) {
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleMove(event) {
    event.preventDefault();

    if (!selectedItem) {
      setMessage("Seleccione una entrada para mover.");
      return;
    }

    if (!destinationPath.trim()) {
      setMessage("Ingrese la ruta destino completa.");
      return;
    }

    try {
      setLoading(true);
      setMessage("");
      await moveEntry(selectedId, selectedItem.path, destinationPath);
      await loadDirectory(currentPath);
      setMessage("Entrada movida correctamente.");
    } catch (error) {
      setMessage(error.message);
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

        <div style={styles.formWide}>
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
            <button type="button" style={styles.secondaryButton} onClick={loadMounted}>
              Actualizar montadas
            </button>

            <button type="button" style={styles.secondaryButton} onClick={() => loadDirectory()}>
              Actualizar carpeta
            </button>

            <button type="button" style={styles.secondaryButton} onClick={goBack}>
              Regresar
            </button>

            <button type="button" style={styles.secondaryButton} onClick={goRoot}>
              Ir a raíz
            </button>
          </div>
        </div>

        <p style={styles.path}>
          Ruta actual: <strong>{currentPath}</strong>
        </p>

        {message && <p style={styles.message}>{message}</p>}
      </section>

      <section style={styles.twoColumns}>
        <article style={styles.box}>
          <h2 style={styles.subtitle}>Crear carpeta</h2>
          <p style={styles.help}>Crear una carpeta dentro de la ruta actual.</p>

          <form onSubmit={handleCreateFolder} style={styles.form}>
            <label style={styles.label}>
              Nombre de carpeta
              <input
                style={styles.input}
                value={folderName}
                placeholder="Ejemplo: docs"
                onChange={(e) => setFolderName(e.target.value)}
              />
            </label>

            <button type="submit" style={styles.button} disabled={loading}>
              Crear carpeta
            </button>
          </form>
        </article>

        <article style={styles.box}>
          <h2 style={styles.subtitle}>Crear archivo</h2>
          <p style={styles.help}>Crear un archivo dentro de la ruta actual.</p>

          <form onSubmit={handleCreateFile} style={styles.form}>
            <label style={styles.label}>
              Nombre de archivo
              <input
                style={styles.input}
                value={fileName}
                placeholder="Ejemplo: nota.txt"
                onChange={(e) => setFileName(e.target.value)}
              />
            </label>

            <label style={styles.label}>
              Contenido
              <textarea
                style={styles.textarea}
                value={newFileContent}
                placeholder="Contenido inicial del archivo"
                onChange={(e) => setNewFileContent(e.target.value)}
              />
            </label>

            <button type="submit" style={styles.button} disabled={loading}>
              Crear archivo
            </button>
          </form>
        </article>
      </section>

      <section style={styles.box}>
        <h2 style={styles.subtitle}>Contenido</h2>

        {loading && <p>Cargando...</p>}

        {!loading && items.length === 0 && (
          <p>No hay archivos o carpetas en esta ruta.</p>
        )}

        <div style={styles.grid}>
          {items.map((item, index) => {
            const normalized = normalizeItem(item, currentPath, index);
            const isSelected = selectedItem?.path === normalized.path;

            return (
              <article
                key={`${normalized.path}-${index}`}
                style={{
                  ...styles.item,
                  border: isSelected ? "2px solid #b85c7a" : "1px solid #e5bcc9"
                }}
              >
                <button
                  type="button"
                  style={styles.itemMainButton}
                  onClick={() => selectItem(normalized)}
                >
                  <span style={styles.itemIcon}>
                    {normalized.type === "folder" ? "□" : "▤"}
                  </span>
                  <span style={styles.itemTitle}>{normalized.name}</span>
                  <small style={styles.itemType}>
                    {normalized.type === "folder" ? "Carpeta" : "Archivo"} · {readableSize(normalized.size)}
                  </small>
                  <small style={styles.itemType}>Permisos: {normalized.permissions || "---"}</small>
                </button>

                <div style={styles.row}>
                  <button type="button" style={styles.smallButton} onClick={() => selectItem(normalized)}>
                    Seleccionar
                  </button>
                  <button type="button" style={styles.smallButton} onClick={() => openItem(normalized)}>
                    {normalized.type === "folder" ? "Abrir" : "Ver archivo"}
                  </button>
                </div>
              </article>
            );
          })}
        </div>
      </section>

      <section style={styles.box}>
        <h2 style={styles.subtitle}>Acciones sobre archivo o carpeta</h2>
        <p style={styles.path}>Seleccionado: <strong>{selectedLabel}</strong></p>
        {selectedItem && <p style={styles.help}>Ruta seleccionada: {selectedItem.path}</p>}

        <div style={styles.actionGrid}>
          <form onSubmit={handleRename} style={styles.form}>
            <h3 style={styles.sectionTitle}>Renombrar</h3>
            <label style={styles.label}>
              Nuevo nombre
              <input
                style={styles.input}
                value={newName}
                placeholder="nuevo.txt"
                onChange={(e) => setNewName(e.target.value)}
                disabled={!selectedItem}
              />
            </label>
            <button type="submit" style={styles.button} disabled={loading || !selectedItem}>
              Renombrar
            </button>
          </form>

          <form onSubmit={handleCopy} style={styles.form}>
            <h3 style={styles.sectionTitle}>Copiar</h3>
            <label style={styles.label}>
              Ruta destino completa
              <input
                style={styles.input}
                value={destinationPath}
                placeholder="/docs/copia.txt"
                onChange={(e) => setDestinationPath(e.target.value)}
                disabled={!selectedItem}
              />
            </label>
            <button type="submit" style={styles.button} disabled={loading || !selectedItem}>
              Copiar
            </button>
          </form>

          <form onSubmit={handleMove} style={styles.form}>
            <h3 style={styles.sectionTitle}>Mover</h3>
            <label style={styles.label}>
              Ruta destino completa
              <input
                style={styles.input}
                value={destinationPath}
                placeholder="/carpeta/nombre.txt"
                onChange={(e) => setDestinationPath(e.target.value)}
                disabled={!selectedItem}
              />
            </label>
            <button type="submit" style={styles.button} disabled={loading || !selectedItem}>
              Mover
            </button>
          </form>

          <div style={styles.form}>
            <h3 style={styles.sectionTitle}>Eliminar</h3>
            <p style={styles.help}>Elimina archivos o carpetas desde la interfaz.</p>
            <button type="button" style={styles.dangerButton} onClick={handleRemove} disabled={loading || !selectedItem}>
              Eliminar seleccionado
            </button>
          </div>
        </div>
      </section>

      {selectedItem?.type === "file" && (
        <section style={styles.box}>
          <h2 style={styles.subtitle}>Ver y editar archivo</h2>
          <p style={styles.help}>Archivo: {selectedItem.path}</p>

          <div style={styles.twoColumns}>
            <article>
              <h3 style={styles.sectionTitle}>Contenido actual</h3>
              <pre style={styles.viewer}>{fileContent || "Presione Ver archivo para cargar el contenido."}</pre>
            </article>

            <form onSubmit={handleEdit} style={styles.form}>
              <h3 style={styles.sectionTitle}>Editar contenido</h3>
              <textarea
                style={styles.editArea}
                value={editContent}
                placeholder="Nuevo contenido del archivo"
                onChange={(e) => setEditContent(e.target.value)}
              />
              <button type="submit" style={styles.button} disabled={loading}>
                Guardar edición
              </button>
            </form>
          </div>
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
    color: "#3b202b"
  },
  title: {
    marginTop: 0,
    marginBottom: "20px"
  },
  subtitle: {
    marginTop: 0,
    marginBottom: "8px"
  },
  sectionTitle: {
    margin: "0 0 10px 0",
    fontSize: "16px"
  },
  box: {
    background: "#ffffff",
    border: "1px solid #e5bcc9",
    borderRadius: "16px",
    padding: "22px",
    marginBottom: "22px"
  },
  twoColumns: {
    display: "grid",
    gridTemplateColumns: "repeat(auto-fit, minmax(280px, 1fr))",
    gap: "18px",
    marginBottom: "22px"
  },
  actionGrid: {
    display: "grid",
    gridTemplateColumns: "repeat(auto-fit, minmax(230px, 1fr))",
    gap: "16px"
  },
  form: {
    display: "grid",
    gap: "12px"
  },
  formWide: {
    display: "grid",
    gap: "12px"
  },
  label: {
    display: "grid",
    gap: "6px",
    fontWeight: "600"
  },
  input: {
    padding: "11px",
    borderRadius: "10px",
    border: "1px solid #d7a9b9",
    fontSize: "15px",
    width: "100%",
    boxSizing: "border-box"
  },
  textarea: {
    minHeight: "96px",
    padding: "11px",
    borderRadius: "10px",
    border: "1px solid #d7a9b9",
    fontSize: "15px",
    width: "100%",
    boxSizing: "border-box",
    resize: "vertical"
  },
  editArea: {
    minHeight: "220px",
    padding: "11px",
    borderRadius: "10px",
    border: "1px solid #d7a9b9",
    fontSize: "15px",
    width: "100%",
    boxSizing: "border-box",
    resize: "vertical"
  },
  row: {
    display: "flex",
    gap: "10px",
    flexWrap: "wrap"
  },
  button: {
    padding: "12px",
    borderRadius: "10px",
    border: "none",
    background: "#b85c7a",
    color: "white",
    fontWeight: "700",
    cursor: "pointer"
  },
  smallButton: {
    padding: "9px 11px",
    borderRadius: "9px",
    border: "none",
    background: "#b85c7a",
    color: "white",
    cursor: "pointer"
  },
  secondaryButton: {
    padding: "10px 12px",
    borderRadius: "10px",
    border: "1px solid #b85c7a",
    background: "white",
    color: "#b85c7a",
    fontWeight: "700",
    cursor: "pointer"
  },
  dangerButton: {
    padding: "12px",
    borderRadius: "10px",
    border: "none",
    background: "#8f2f4f",
    color: "white",
    fontWeight: "700",
    cursor: "pointer"
  },
  path: {
    marginTop: "14px"
  },
  help: {
    color: "#765062",
    marginTop: 0
  },
  message: {
    background: "#f5d7e1",
    padding: "12px",
    borderRadius: "10px",
    marginTop: "14px"
  },
  grid: {
    display: "grid",
    gap: "12px",
    gridTemplateColumns: "repeat(auto-fill, minmax(220px, 1fr))"
  },
  item: {
    textAlign: "left",
    padding: "14px",
    borderRadius: "14px",
    background: "#fffafd",
    display: "grid",
    gap: "10px",
    color: "#3b202b"
  },
  itemMainButton: {
    border: "none",
    background: "transparent",
    padding: 0,
    textAlign: "left",
    display: "grid",
    gap: "5px",
    color: "#3b202b",
    cursor: "pointer"
  },
  itemIcon: {
    fontSize: "24px"
  },
  itemTitle: {
    fontWeight: "700"
  },
  itemType: {
    color: "#765062"
  },
  viewer: {
    background: "#2b2025",
    color: "white",
    padding: "18px",
    borderRadius: "14px",
    whiteSpace: "pre-wrap",
    overflow: "auto",
    minHeight: "180px"
  }
};