import { useEffect, useMemo, useState } from "react";
import {
  createPartition,
  deletePartition,
  formatPartition,
  getDisks,
  getPartitions,
  mountPartition,
  resizePartition,
  unmountPartition
} from "../api/diskApi";
import { generateReport } from "../api/reportApi";
import ReportViewer from "../components/ReportViewer";

function normalizeType(type) {
  const value = String(type || "P").toUpperCase();
  if (value === "E") return "Extendida";
  if (value === "L") return "Lógica";
  return "Primaria";
}

function formatBytes(value) {
  const number = Number(value || 0);

  if (!Number.isFinite(number) || number <= 0) return "0 bytes";
  if (number >= 1024 * 1024) return `${(number / 1024 / 1024).toFixed(2)} MB`;
  if (number >= 1024) return `${(number / 1024).toFixed(2)} KB`;

  return `${number} bytes`;
}

export default function PartitionPage() {
  const [disks, setDisks] = useState([]);
  const [selectedDiskPath, setSelectedDiskPath] = useState("");
  const [selectedPartition, setSelectedPartition] = useState(null);
  const [partitions, setPartitions] = useState([]);
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);
  const [reportLoading, setReportLoading] = useState(false);
  const [report, setReport] = useState(null);

  const [createForm, setCreateForm] = useState({
    name: "",
    size: "",
    unit: "M",
    type: "P",
    fit: "WF"
  });


  const selectedDisk = useMemo(
    () => disks.find((disk) => disk.path === selectedDiskPath),
    [disks, selectedDiskPath]
  );

  async function loadDisks() {
    try {
      const data = await getDisks();
      const list = Array.isArray(data) ? data : [];
      setDisks(list);

      if (!selectedDiskPath && list.length > 0) {
        setSelectedDiskPath(list[0].path);
      }
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
      const list = Array.isArray(data) ? data : [];

      setPartitions(list);

      if (selectedPartition) {
        const updated = list.find((partition) => partition.name === selectedPartition.name);
        setSelectedPartition(updated || null);
      }
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

  function selectPartition(partition) {
    setSelectedPartition(partition);
    setMessage(`Partición seleccionada: ${partition.name}`);
  }

  async function handleCreate(event) {
    event.preventDefault();

    if (!selectedDiskPath) {
      setMessage("Seleccione un disco.");
      return;
    }

    if (!createForm.name.trim()) {
      setMessage("Ingrese el nombre de la partición.");
      return;
    }

    if (!createForm.size || Number(createForm.size) <= 0) {
      setMessage("Ingrese un tamaño positivo para crear la partición.");
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
        fit: createForm.fit
      });

      setCreateForm({
        name: "",
        size: "",
        unit: "M",
        type: "P",
        fit: "WF"
      });

      await loadPartitions();
      setMessage("Partición creada correctamente.");
      await refreshDiskReport("disk");
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
        name: partition.name
      });

      await loadPartitions();
      setMessage(`Partición montada correctamente. ID: ${mounted.id}`);
    } catch (error) {
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleUnmount(partition) {
    if (!partition?.id) {
      setMessage("La partición no está montada.");
      return;
    }

    try {
      setLoading(true);
      setMessage("");

      await unmountPartition(partition.id);

      if (selectedPartition?.name === partition.name) {
        setSelectedPartition(null);
      }

      await loadPartitions();
      setMessage("Partición desmontada correctamente.");
    } catch (error) {
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleFormat(partition) {
    if (!partition?.mounted || !partition?.id) {
      setMessage("Primero monte la partición para poder formatearla.");
      return;
    }

    if (!window.confirm(`¿Formatear ${partition.name} con MKFS full?`)) return;

    try {
      setLoading(true);
      setMessage("");

      await formatPartition({
        id: partition.id,
        type: "full"
      });

      await loadPartitions();
      setMessage("Partición formateada correctamente. Ahora puede iniciar sesión.");
    } catch (error) {
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }


  async function refreshDiskReport(type = "disk") {
    if (!selectedDiskPath) {
      setReport(null);
      return;
    }

    try {
      setReportLoading(true);
      const data = await generateReport({
        type,
        diskPath: selectedDiskPath,
        filePath: selectedDiskPath
      });
      setReport(data);
    } catch (error) {
      setReport(null);
    } finally {
      setReportLoading(false);
    }
  }

  async function applyFdiskAdd(partition, value, unit = "M") {
    if (!partition) {
      setMessage("Seleccione una partición para cambiar su espacio.");
      return;
    }

    const add = Number(value);

    if (!Number.isFinite(add) || add === 0) {
      setMessage("Ingrese un valor válido. Puede ser positivo o negativo.");
      return;
    }

    try {
      setLoading(true);
      setMessage("");

      await resizePartition({
        diskPath: selectedDiskPath,
        path: selectedDiskPath,
        name: partition.name,
        add,
        unit
      });

      await loadPartitions();
      setMessage(`Espacio actualizado correctamente en ${partition.name}.`);
      await refreshDiskReport("disk");
    } catch (error) {
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  async function applyFdiskDelete(partition, deleteType = "fast") {
    if (!partition) {
      setMessage("Seleccione una partición para eliminar.");
      return;
    }

    const mode = String(deleteType || "fast").toLowerCase();

    if (mode !== "fast" && mode !== "full") {
      setMessage("El tipo de eliminación debe ser rápida o completa.");
      return;
    }

    if (!window.confirm(`¿Eliminar la partición ${partition.name} con modo ${mode}?`)) return;

    try {
      setLoading(true);
      setMessage("");

      await deletePartition({
        diskPath: selectedDiskPath,
        path: selectedDiskPath,
        name: partition.name,
        deleteType: mode,
        delete: mode
      });

      if (selectedPartition?.name === partition.name) {
        setSelectedPartition(null);
      }

      await loadPartitions();
      setMessage(`Partición eliminada correctamente con modo ${mode}.`);
      await refreshDiskReport("disk");
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
          onChange={(event) => setSelectedDiskPath(event.target.value)}
        >
          <option value="">Seleccione un disco</option>
          {disks.map((disk) => (
            <option key={disk.path} value={disk.path}>
              {disk.name} — {disk.path}
            </option>
          ))}
        </select>

        {selectedDisk && (
          <p style={styles.help}>
            Disco seleccionado: {selectedDisk.name} · {selectedDisk.path}
          </p>
        )}
      </section>

      <section style={styles.twoColumns}>
        <article style={styles.box}>
          <h2 style={styles.subtitle}>Crear partición</h2>
       

          <form onSubmit={handleCreate} style={styles.form}>
            <label style={styles.label}>
              Nombre
              <input
                style={styles.input}
                value={createForm.name}
                placeholder="Ejemplo: Part1"
                onChange={(event) =>
                  setCreateForm({ ...createForm, name: event.target.value })
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
                onChange={(event) =>
                  setCreateForm({ ...createForm, size: event.target.value })
                }
              />
            </label>

            <label style={styles.label}>
              Unidad
              <select
                style={styles.input}
                value={createForm.unit}
                onChange={(event) =>
                  setCreateForm({ ...createForm, unit: event.target.value })
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
                onChange={(event) =>
                  setCreateForm({ ...createForm, type: event.target.value })
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
                onChange={(event) =>
                  setCreateForm({ ...createForm, fit: event.target.value })
                }
              >
                <option value="FF">First Fit</option>
                <option value="BF">Best Fit</option>
                <option value="WF">Worst Fit</option>
              </select>
            </label>

            <button type="submit" style={styles.button} disabled={loading}>
              Crear partición
            </button>
          </form>
        </article>

      </section>

      {message && <p style={styles.message}>{message}</p>}

      <section style={styles.box}>
        <div style={styles.headerRow}>
          <div>
            <h2 style={styles.subtitle}>Particiones del disco</h2>
            
          </div>

          <button type="button" style={styles.secondaryButton} onClick={() => loadPartitions()}>
            Actualizar particiones
          </button>
        </div>

        {loading && <p>Cargando...</p>}

        {!loading && !selectedDiskPath && <p>Seleccione un disco.</p>}

        {!loading && selectedDiskPath && partitions.length === 0 && (
          <p>Este disco no tiene particiones.</p>
        )}

        <div style={styles.grid}>
          {partitions.map((partition) => {
            const isSelected = selectedPartition?.name === partition.name;

            return (
              <article
                key={`${partition.name}-${partition.start}`}
                style={{
                  ...styles.card,
                  border: isSelected ? "2px solid #b85c7a" : "1px solid #e5bcc9"
                }}
              >
                <h3 style={styles.cardTitle}>{partition.name}</h3>
                <p><strong>Tipo:</strong> {normalizeType(partition.type)}</p>
                <p><strong>Fit:</strong> {partition.fit}</p>
                <p><strong>Inicio:</strong> {partition.start}</p>
                <p><strong>Tamaño:</strong> {formatBytes(partition.size)}</p>
                <p><strong>Montada:</strong> {partition.mounted ? "Sí" : "No"}</p>
                <p><strong>Formateada:</strong> {partition.formatted ? "Sí" : "No"}</p>
                <p><strong>ID:</strong> {partition.id || "sin montar"}</p>

                <div style={styles.row}>
                  <button
                    type="button"
                    style={styles.smallButton}
                    onClick={() => selectPartition(partition)}
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
                    <button
                      type="button"
                      style={styles.smallButton}
                      onClick={() => handleUnmount(partition)}
                    >
                      Desmontar
                    </button>
                  )}

                  <button
                    type="button"
                    style={styles.smallButton}
                    onClick={() => handleFormat(partition)}
                  >
                    Formatear
                  </button>
                </div>

                <div style={styles.row}>
                  <button
                    type="button"
                    style={styles.secondarySmallButton}
                    onClick={() => applyFdiskAdd(partition, 1, "M")}
                  >
                    Aumentar +1 MB
                  </button>

                  <button
                    type="button"
                    style={styles.secondarySmallButton}
                    onClick={() => applyFdiskAdd(partition, -1, "M")}
                  >
                    Reducir -1 MB
                  </button>
                </div>

                <div style={styles.row}>
                  <button
                    type="button"
                    style={styles.dangerSmallButton}
                    onClick={() => applyFdiskDelete(partition, "fast")}
                  >
                    Eliminar rápida
                  </button>

                  <button
                    type="button"
                    style={styles.dangerSmallButton}
                    onClick={() => applyFdiskDelete(partition, "full")}
                  >
                    Eliminar completa
                  </button>
                </div>
              </article>
            );
          })}
        </div>
      </section>

      <section style={styles.box}>
        <div style={styles.headerRow}>
          <div>
            <h2 style={styles.subtitle}>Reporte del disco</h2>
            
          </div>
          <div style={styles.row}>
            <button type="button" style={styles.secondaryButton} onClick={() => refreshDiskReport("mbr")}>
              Ver MBR
            </button>
            <button type="button" style={styles.secondaryButton} onClick={() => refreshDiskReport("disk")}>
              Ver Disk
            </button>
          </div>
        </div>
        <ReportViewer report={report} loading={reportLoading} />
      </section>
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
  box: {
    background: "#ffffff",
    border: "1px solid #e5bcc9",
    borderRadius: "16px",
    padding: "22px",
    marginBottom: "22px"
  },
  twoColumns: {
    display: "grid",
    gridTemplateColumns: "repeat(auto-fit, minmax(320px, 1fr))",
    gap: "18px",
    marginBottom: "22px"
  },
  form: {
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
  secondarySmallButton: {
    padding: "9px 11px",
    borderRadius: "9px",
    border: "1px solid #b85c7a",
    background: "white",
    color: "#b85c7a",
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
  dangerSmallButton: {
    padding: "9px 11px",
    borderRadius: "9px",
    border: "none",
    background: "#8f2f4f",
    color: "white",
    cursor: "pointer"
  },
  selectedText: {
    background: "#fff7fa",
    padding: "10px",
    borderRadius: "10px",
    border: "1px solid #e5bcc9"
  },
  message: {
    background: "#f5d7e1",
    padding: "12px",
    borderRadius: "10px",
    marginBottom: "20px"
  },
  help: {
    color: "#765062",
    marginTop: "8px"
  },
  headerRow: {
    display: "flex",
    justifyContent: "space-between",
    gap: "12px",
    flexWrap: "wrap",
    alignItems: "center"
  },
  grid: {
    display: "grid",
    gap: "14px"
  },
  card: {
    borderRadius: "14px",
    padding: "16px",
    background: "#fffafd",
    display: "grid",
    gap: "6px"
  },
  cardTitle: {
    marginTop: 0,
    marginBottom: "4px"
  },
  row: {
    display: "flex",
    gap: "10px",
    flexWrap: "wrap",
    marginTop: "8px"
  }
};