import { useEffect, useState } from "react";
import { createDisk as createDiskRequest, deleteDisk as deleteDiskRequest, getDisks } from "../api/diskApi";
import { generateReport } from "../api/reportApi";
import ReportViewer from "../components/ReportViewer";

export default function Disks() {
  const [disks, setDisks] = useState([]);
  const [selectedDisk, setSelectedDisk] = useState(null);
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);
  const [reportLoading, setReportLoading] = useState(false);
  const [report, setReport] = useState(null);

  const [form, setForm] = useState({
    name: "",
    size: "",
    unit: "M",
    fit: "FF"
  });

  async function loadDisks() {
    try {
      setLoading(true);
      const list = await getDisks();
      const normalized = Array.isArray(list) ? list : [];
      setDisks(normalized);

      if (!selectedDisk && normalized.length > 0) {
        setSelectedDisk(normalized[0]);
      }
    } catch (error) {
      setMessage("Error conectando con backend: " + error.message);
      setDisks([]);
    } finally {
      setLoading(false);
    }
  }

  async function createDisk(event) {
    event.preventDefault();

    if (!form.name.trim()) {
      setMessage("Ingrese el nombre del disco.");
      return;
    }

    if (!form.size || Number(form.size) <= 0) {
      setMessage("Ingrese un tamaño válido.");
      return;
    }

    try {
      setLoading(true);
      setMessage("");
      setReport(null);

      const disk = await createDiskRequest({
        name: form.name.trim(),
        size: Number(form.size),
        unit: form.unit,
        fit: form.fit
      });

      setForm({ name: "", size: "", unit: "M", fit: "FF" });
      setSelectedDisk(disk);
      await loadDisks();
      setMessage("Disco creado correctamente. Se puede visualizar su reporte MBR.");

      if (disk?.path) {
        await showDiskReport(disk, "mbr");
      }
    } catch (error) {
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  async function deleteDisk(path) {
    if (!window.confirm("¿Eliminar este disco?")) return;

    try {
      setLoading(true);
      setMessage("");
      setReport(null);

      await deleteDiskRequest(path);

      if (selectedDisk?.path === path) {
        setSelectedDisk(null);
      }

      setMessage("Disco eliminado correctamente.");
      await loadDisks();
    } catch (error) {
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  async function showDiskReport(disk = selectedDisk, type = "mbr") {
    if (!disk?.path) {
      setMessage("Seleccione un disco para visualizar el reporte.");
      return;
    }

    try {
      setReportLoading(true);
      setMessage("");

      const data = await generateReport({
        type,
        diskPath: disk.path,
        filePath: disk.path
      });

      setReport(data);
    } catch (error) {
      setReport(null);
      setMessage("No se pudo generar el reporte: " + error.message);
    } finally {
      setReportLoading(false);
    }
  }

  useEffect(() => {
    loadDisks();
  }, []);

  return (
    <main style={styles.page}>
      <h1 style={styles.title}>Discos</h1>

      <section style={styles.box}>
        <h2 style={styles.subtitle}>Crear nuevo disco</h2>
        

        <form onSubmit={createDisk} style={styles.form}>
          <label style={styles.label}>
            Nombre
            <input
              style={styles.input}
              value={form.name}
              placeholder="Ejemplo: disco1.dsk"
              onChange={(e) => setForm({ ...form, name: e.target.value })}
            />
          </label>

          <label style={styles.label}>
            Tamaño
            <input
              style={styles.input}
              type="number"
              value={form.size}
              placeholder="Ejemplo: 20"
              onChange={(e) => setForm({ ...form, size: e.target.value })}
            />
          </label>

          <label style={styles.label}>
            Unidad
            <select style={styles.input} value={form.unit} onChange={(e) => setForm({ ...form, unit: e.target.value })}>
              <option value="K">KB</option>
              <option value="M">MB</option>
            </select>
          </label>

          <label style={styles.label}>
            Fit
            <select style={styles.input} value={form.fit} onChange={(e) => setForm({ ...form, fit: e.target.value })}>
              <option value="FF">First Fit</option>
              <option value="BF">Best Fit</option>
              <option value="WF">Worst Fit</option>
            </select>
          </label>

          <button style={styles.button} disabled={loading}>
            Crear disco
          </button>
        </form>

        {message && <p style={styles.message}>{message}</p>}
      </section>

      <section style={styles.box}>
        <div style={styles.headerRow}>
          <div>
            <h2 style={styles.subtitle}>Discos creados</h2>
            <p style={styles.help}>Seleccione un disco para crear particiones o ver sus reportes.</p>
          </div>
          <button type="button" style={styles.secondaryButton} onClick={loadDisks} disabled={loading}>
            Actualizar
          </button>
        </div>

        {loading && <p>Cargando...</p>}
        {!loading && disks.length === 0 && <p>No hay discos creados todavía.</p>}

        <div style={styles.grid}>
          {disks.map((disk) => {
            const active = selectedDisk?.path === disk.path;

            return (
              <article
                key={disk.id || disk.path}
                style={{
                  ...styles.card,
                  border: active ? "2px solid #b85c7a" : "1px solid #e5bcc9"
                }}
              >
                <h3>{disk.name}</h3>
                <p><strong>Ruta:</strong> {disk.path}</p>
                <p><strong>Tamaño:</strong> {disk.size} bytes</p>
                <p><strong>Fit:</strong> {disk.fit}</p>

                <div style={styles.row}>
                  <button type="button" style={styles.secondaryButton} onClick={() => setSelectedDisk(disk)}>
                    Seleccionar
                  </button>
                  <button type="button" style={styles.secondaryButton} onClick={() => showDiskReport(disk, "mbr")}>
                    Ver MBR
                  </button>
                  <button type="button" style={styles.secondaryButton} onClick={() => showDiskReport(disk, "disk")}>
                    Ver Disk
                  </button>
                  <button type="button" style={styles.deleteButton} onClick={() => deleteDisk(disk.path)}>
                    Eliminar
                  </button>
                </div>
              </article>
            );
          })}
        </div>
      </section>

      <section style={styles.box}>
        <h2 style={styles.subtitle}>Reporte visual del disco</h2>
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
    marginBottom: "22px",
    maxWidth: "980px"
  },
  headerRow: {
    display: "flex",
    justifyContent: "space-between",
    gap: "14px",
    alignItems: "flex-start",
    flexWrap: "wrap"
  },
  form: {
    display: "grid",
    gap: "12px",
    maxWidth: "520px"
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
    fontSize: "15px"
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
  secondaryButton: {
    padding: "10px 12px",
    borderRadius: "10px",
    border: "1px solid #b85c7a",
    background: "white",
    color: "#b85c7a",
    fontWeight: "700",
    cursor: "pointer"
  },
  deleteButton: {
    padding: "10px 12px",
    borderRadius: "10px",
    border: "none",
    background: "#8f2f4f",
    color: "white",
    cursor: "pointer"
  },
  row: {
    display: "flex",
    gap: "10px",
    flexWrap: "wrap"
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
    gap: "14px"
  },
  card: {
    borderRadius: "14px",
    padding: "16px",
    background: "#fffafd"
  }
};