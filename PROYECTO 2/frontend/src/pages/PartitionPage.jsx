import { useEffect, useState } from "react";
import {
  createPartition,
  deletePartition,
  formatPartition,
  getDisks,
  getPartitions,
  mountPartition,
  resizePartition,
  unmountPartition,
} from "../api/diskApi";
// PartitionPage.jsx: contiene la página principal de administración de particiones, mostrando un encabezado y un formulario para crear, eliminar, montar, desmontar, formatear y redimensionar particiones. También incluye una lista de particiones disponibles en el disco seleccionado y el estado de la partición seleccionada actualmente. Utiliza los hooks useState y useEffect para manejar la carga de discos y particiones, así como las acciones del usuario.
export default function PartitionPage() {
  const [disks, setDisks] = useState([]);
  const [selectedDiskPath, setSelectedDiskPath] = useState("");
  const [selectedPartition, setSelectedPartition] = useState(null);
  const [partitions, setPartitions] = useState([]);
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  const [createForm, setCreateForm] = useState({
    name: "",
    size: "",
    unit: "M",
    type: "P",
    fit: "WF",
  });

  const [resizeForm, setResizeForm] = useState({
    add: "",
    unit: "M",
    deleteType: "fast",
  });

  async function loadDisks() {
    try {
      const data = await getDisks();
      setDisks(Array.isArray(data) ? data : []);
    } catch (error) {
      setMessage("Error cargando discos: " + error.message);
    }
  }

  async function loadPartitions(path = selectedDiskPath) {
    if (!path) {
      setPartitions([]);
      return;
    }

    try {
      setLoading(true);
      const data = await getPartitions(path);
      setPartitions(Array.isArray(data) ? data : []);
    } catch (error) {
      setMessage("Error cargando particiones: " + error.message);
      setPartitions([]);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadDisks();
  }, []);

  useEffect(() => {
    setSelectedPartition(null);
    loadPartitions(selectedDiskPath);
  }, [selectedDiskPath]);

  async function handleCreate(e) {
    e.preventDefault();

    if (!selectedDiskPath) {
      setMessage("Seleccioná un disco.");
      return;
    }

    if (!createForm.name.trim()) {
      setMessage("Ingresá el nombre de la partición.");
      return;
    }

    if (!createForm.size || Number(createForm.size) <= 0) {
      setMessage("Ingresá un tamaño válido.");
      return;
    }

    try {
      setLoading(true);
      setMessage("");

      await createPartition({
        diskPath: selectedDiskPath,
        path: selectedDiskPath,
        name: createForm.name.trim(),
        size: Number(createForm.size),
        unit: createForm.unit,
        type: createForm.type,
        fit: createForm.fit,
      });

      setCreateForm({
        name: "",
        size: "",
        unit: "M",
        type: "P",
        fit: "WF",
      });

      setMessage("Partición creada correctamente.");
      await loadPartitions();
    } catch (error) {
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleMount(partition) {
    try {
      setLoading(true);
      setMessage("");

      const mounted = await mountPartition({
        diskPath: selectedDiskPath,
        path: selectedDiskPath,
        name: partition.name,
      });

      setMessage(`Partición montada correctamente. ID: ${mounted.id}`);
      await loadPartitions();
    } catch (error) {
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleUnmount(partition) {
    try {
      setLoading(true);
      setMessage("");

      await unmountPartition(partition.id);

      setMessage("Partición desmontada correctamente.");
      setSelectedPartition(null);
      await loadPartitions();
    } catch (error) {
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleFormat(partition) {
    try {
      setLoading(true);
      setMessage("");

      await formatPartition({
        id: partition.id,
        type: "full",
      });

      setMessage("Partición formateada correctamente.");
      await loadPartitions();
    } catch (error) {
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleDelete() {
    if (!selectedPartition) {
      setMessage("Seleccioná una partición.");
      return;
    }

    if (!window.confirm("¿Eliminar esta partición?")) return;

    try {
      setLoading(true);
      setMessage("");

      await deletePartition({
        diskPath: selectedDiskPath,
        path: selectedDiskPath,
        name: selectedPartition.name,
        deleteType: resizeForm.deleteType,
        delete: resizeForm.deleteType,
      });

      setMessage("Partición eliminada correctamente.");
      setSelectedPartition(null);
      await loadPartitions();
    } catch (error) {
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleResize() {
    if (!selectedPartition) {
      setMessage("Seleccioná una partición.");
      return;
    }

    if (!resizeForm.add || Number(resizeForm.add) === 0) {
      setMessage("Ingresá un valor add válido.");
      return;
    }

    try {
      setLoading(true);
      setMessage("");

      await resizePartition({
        diskPath: selectedDiskPath,
        path: selectedDiskPath,
        name: selectedPartition.name,
        add: Number(resizeForm.add),
        unit: resizeForm.unit,
      });

      setMessage("Partición redimensionada correctamente.");
      await loadPartitions();
    } catch (error) {
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <main style={styles.page}>
      <h1 style={styles.title}>Particiones</h1>

      <section style={styles.box}>
        <h2 style={styles.subtitle}>Seleccionar disco</h2>

        <select
          style={styles.input}
          value={selectedDiskPath}
          onChange={(e) => setSelectedDiskPath(e.target.value)}
        >
          <option value="">Seleccione un disco</option>
          {disks.map((disk) => (
            <option key={disk.path} value={disk.path}>
              {disk.name} — {disk.path}
            </option>
          ))}
        </select>
      </section>

      <section style={styles.box}>
        <h2 style={styles.subtitle}>Crear partición</h2>

        <form onSubmit={handleCreate} style={styles.form}>
          <label style={styles.label}>
            Nombre
            <input
              style={styles.input}
              value={createForm.name}
              placeholder="Ejemplo: Part1"
              onChange={(e) =>
                setCreateForm({ ...createForm, name: e.target.value })
              }
            />
          </label>

          <label style={styles.label}>
            Tamaño
            <input
              style={styles.input}
              type="number"
              value={createForm.size}
              placeholder="Ejemplo: 5"
              onChange={(e) =>
                setCreateForm({ ...createForm, size: e.target.value })
              }
            />
          </label>

          <label style={styles.label}>
            Unidad
            <select
              style={styles.input}
              value={createForm.unit}
              onChange={(e) =>
                setCreateForm({ ...createForm, unit: e.target.value })
              }
            >
              <option value="B">B</option>
              <option value="K">K</option>
              <option value="M">M</option>
            </select>
          </label>

          <label style={styles.label}>
            Tipo
            <select
              style={styles.input}
              value={createForm.type}
              onChange={(e) =>
                setCreateForm({ ...createForm, type: e.target.value })
              }
            >
              <option value="P">Primaria</option>
              <option value="E">Extendida</option>
              <option value="L">Lógica</option>
            </select>
          </label>

          <label style={styles.label}>
            Fit
            <select
              style={styles.input}
              value={createForm.fit}
              onChange={(e) =>
                setCreateForm({ ...createForm, fit: e.target.value })
              }
            >
              <option value="FF">First Fit</option>
              <option value="BF">Best Fit</option>
              <option value="WF">Worst Fit</option>
            </select>
          </label>

          <button style={styles.button} disabled={loading}>
            Crear partición
          </button>
        </form>
      </section>

      <section style={styles.box}>
        <h2 style={styles.subtitle}>Modificar partición seleccionada</h2>

        <p>
          {selectedPartition
            ? `Seleccionada: ${selectedPartition.name}`
            : "No hay partición seleccionada."}
        </p>

        <div style={styles.form}>
          <label style={styles.label}>
            Add
            <input
              style={styles.input}
              type="number"
              value={resizeForm.add}
              placeholder="Ejemplo: 1 o -1"
              onChange={(e) =>
                setResizeForm({ ...resizeForm, add: e.target.value })
              }
            />
          </label>

          <label style={styles.label}>
            Unidad
            <select
              style={styles.input}
              value={resizeForm.unit}
              onChange={(e) =>
                setResizeForm({ ...resizeForm, unit: e.target.value })
              }
            >
              <option value="B">B</option>
              <option value="K">K</option>
              <option value="M">M</option>
            </select>
          </label>

          <label style={styles.label}>
            Tipo delete
            <select
              style={styles.input}
              value={resizeForm.deleteType}
              onChange={(e) =>
                setResizeForm({ ...resizeForm, deleteType: e.target.value })
              }
            >
              <option value="fast">Fast</option>
              <option value="full">Full</option>
            </select>
          </label>

          <div style={styles.row}>
            <button type="button" style={styles.button} onClick={handleResize}>
              Aplicar add
            </button>

            <button type="button" style={styles.dangerButton} onClick={handleDelete}>
              Eliminar partición
            </button>
          </div>
        </div>
      </section>

      {message && <p style={styles.message}>{message}</p>}

      <section style={styles.box}>
        <h2 style={styles.subtitle}>Particiones del disco</h2>

        {loading && <p>Cargando...</p>}

        {!loading && !selectedDiskPath && <p>Seleccione un disco.</p>}

        {!loading && selectedDiskPath && partitions.length === 0 && (
          <p>Este disco no tiene particiones.</p>
        )}

        <div style={styles.grid}>
          {partitions.map((partition) => (
            <article
              key={`${partition.name}-${partition.start}`}
              style={{
                ...styles.card,
                border:
                  selectedPartition?.name === partition.name
                    ? "2px solid #b85c7a"
                    : "1px solid #e5bcc9",
              }}
            >
              <h3>{partition.name}</h3>
              <p><strong>Tipo:</strong> {partition.type}</p>
              <p><strong>Fit:</strong> {partition.fit}</p>
              <p><strong>Inicio:</strong> {partition.start}</p>
              <p><strong>Tamaño:</strong> {partition.size} bytes</p>
              <p><strong>Montada:</strong> {partition.mounted ? "Sí" : "No"}</p>
              <p><strong>Formateada:</strong> {partition.formatted ? "Sí" : "No"}</p>
              <p><strong>ID:</strong> {partition.id}</p>

              <div style={styles.row}>
                <button
                  type="button"
                  style={styles.smallButton}
                  onClick={() => setSelectedPartition(partition)}
                >
                  Seleccionar
                </button>

                {!partition.mounted ? (
                  <button
                    type="button"
                    style={styles.smallButton}
                    onClick={() => handleMount(partition)}
                  >
                    Montar
                  </button>
                ) : (
                  <>
                    <button
                      type="button"
                      style={styles.smallButton}
                      onClick={() => handleUnmount(partition)}
                    >
                      Desmontar
                    </button>

                    <button
                      type="button"
                      style={styles.smallButton}
                      onClick={() => handleFormat(partition)}
                    >
                      Formatear
                    </button>
                  </>
                )}
              </div>
            </article>
          ))}
        </div>
      </section>
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
    maxWidth: "520px",
  },
  label: {
    display: "grid",
    gap: "6px",
    fontWeight: "600",
  },
  input: {
    padding: "11px",
    borderRadius: "10px",
    border: "1px solid #d7a9b9",
    fontSize: "15px",
    width: "100%",
    boxSizing: "border-box",
  },
  button: {
    padding: "12px",
    borderRadius: "10px",
    border: "none",
    background: "#b85c7a",
    color: "white",
    fontWeight: "700",
    cursor: "pointer",
  },
  smallButton: {
    padding: "9px 11px",
    borderRadius: "9px",
    border: "none",
    background: "#b85c7a",
    color: "white",
    cursor: "pointer",
  },
  dangerButton: {
    padding: "12px",
    borderRadius: "10px",
    border: "none",
    background: "#8f2f4f",
    color: "white",
    fontWeight: "700",
    cursor: "pointer",
  },
  message: {
    background: "#f5d7e1",
    padding: "12px",
    borderRadius: "10px",
    marginBottom: "20px",
  },
  grid: {
    display: "grid",
    gap: "14px",
  },
  card: {
    borderRadius: "14px",
    padding: "16px",
    background: "#fffafd",
  },
  row: {
    display: "flex",
    gap: "10px",
    flexWrap: "wrap",
  },
};
