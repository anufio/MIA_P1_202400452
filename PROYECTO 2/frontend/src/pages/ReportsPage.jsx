import { useEffect, useState } from "react";
import ReportViewer from "../components/ReportViewer";
import { generateReport, listReports } from "../api/reportApi";
import { getMountedPartitions } from "../api/diskApi";
import { REPORT_TYPES } from "../utils/constants";
// ReportsPage.jsx: contiene la página principal de reportes, mostrando un encabezado y un formulario para generar reportes de particiones montadas. También incluye una lista de reportes generados y un visor de reportes para mostrar el contenido del reporte seleccionado. Utiliza los hooks useState y useEffect para manejar la carga de particiones montadas, la generación de reportes y la visualización de reportes.
export default function ReportsPage() {
  const [type, setType] = useState("disk");
  const [mounted, setMounted] = useState([]);
  const [selectedId, setSelectedId] = useState("");
  const [report, setReport] = useState(null);
  const [reports, setReports] = useState([]);
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
      setMessage("Error cargando particiones montadas: " + error.message);
    }
  }

  async function loadReports() {
    const data = await listReports();
    setReports(Array.isArray(data) ? data : []);
  }

  async function handleGenerate() {
    if (!selectedId) {
      setMessage("Primero montá una partición.");
      setReport(null);
      return;
    }

    try {
      setLoading(true);
      setMessage("");
      setReport(null);

      const data = await generateReport({
        id: selectedId,
        type,
      });

      setReport(data);
      setMessage("Reporte generado correctamente.");
      await loadReports();
    } catch (error) {
      setReport(null);
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadMounted();
    loadReports();
  }, []);

  return (
    <main style={styles.page}>
      <h1 style={styles.title}>Reportes</h1>

      <section style={styles.box}>
        <h2 style={styles.subtitle}>Generar reporte</h2>

        <div style={styles.form}>
          <label style={styles.label}>
            Partición montada
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
          </label>

          <label style={styles.label}>
            Tipo de reporte
            <select
              style={styles.input}
              value={type}
              onChange={(e) => setType(e.target.value)}
            >
              {REPORT_TYPES.map((reportType) => (
                <option key={reportType.value} value={reportType.value}>
                  {reportType.label}
                </option>
              ))}
            </select>
          </label>

          <button style={styles.button} onClick={handleGenerate} disabled={loading}>
            Generar reporte
          </button>

          <button style={styles.secondaryButton} onClick={loadMounted} disabled={loading}>
            Actualizar montadas
          </button>
        </div>

        {message && <p style={styles.message}>{message}</p>}
      </section>

      <ReportViewer report={report} loading={loading} />

      <section style={styles.box}>
        <h2 style={styles.subtitle}>Reportes generados</h2>

        {reports.length === 0 && <p>No hay reportes generados todavía.</p>}

        <div style={styles.grid}>
          {reports.map((item) => (
            <article key={item.path} style={styles.card}>
              <h3>{item.type}</h3>
              <p>{item.path}</p>
              <a href={item.url} target="_blank" rel="noreferrer">
                Abrir
              </a>
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
    maxWidth: "620px",
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
  secondaryButton: {
    padding: "12px",
    borderRadius: "10px",
    border: "1px solid #b85c7a",
    background: "white",
    color: "#b85c7a",
    fontWeight: "700",
    cursor: "pointer",
  },
  message: {
    background: "#f5d7e1",
    padding: "12px",
    borderRadius: "10px",
    marginTop: "14px",
  },
  grid: {
    display: "grid",
    gap: "14px",
  },
  card: {
    border: "1px solid #e5bcc9",
    borderRadius: "14px",
    padding: "16px",
    background: "#fffafd",
  },
};
