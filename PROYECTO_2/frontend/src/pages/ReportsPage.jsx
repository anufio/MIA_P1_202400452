import { useEffect, useMemo, useState } from "react";
import ReportViewer from "../components/ReportViewer";
import { generateReport, listReports } from "../api/reportApi";
import { getDisks, getMountedPartitions } from "../api/diskApi";
import { REPORT_TYPES } from "../utils/constants";

const DISK_REPORTS = new Set(["mbr", "disk"]);

export default function ReportsPage() {
  const [type, setType] = useState("disk");
  const [mounted, setMounted] = useState([]);
  const [disks, setDisks] = useState([]);
  const [selectedId, setSelectedId] = useState("");
  const [selectedDiskPath, setSelectedDiskPath] = useState("");
  const [report, setReport] = useState(null);
  const [reports, setReports] = useState([]);
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  const selectedReportIsDisk = useMemo(() => DISK_REPORTS.has(type), [type]);

  async function loadData() {
    try {
      const [mountedData, diskData, reportsData] = await Promise.all([
        getMountedPartitions().catch(() => []),
        getDisks().catch(() => []),
        listReports().catch(() => [])
      ]);

      const mountedList = Array.isArray(mountedData) ? mountedData : [];
      const diskList = Array.isArray(diskData) ? diskData : [];

      setMounted(mountedList);
      setDisks(diskList);
      setReports(Array.isArray(reportsData) ? reportsData : []);

      if (!selectedId && mountedList.length > 0) {
        setSelectedId(mountedList[0].id);
      }

      if (!selectedDiskPath && diskList.length > 0) {
        setSelectedDiskPath(diskList[0].path);
      }
    } catch (error) {
      setMessage("Error cargando datos de reportes: " + error.message);
    }
  }

  async function handleGenerate(nextType = type) {
    const isDiskReport = DISK_REPORTS.has(nextType);

    if (isDiskReport && !selectedDiskPath) {
      setMessage("Seleccione un disco para generar el reporte.");
      setReport(null);
      return;
    }

    if (!isDiskReport && !selectedId) {
      setMessage("Primero monte y formatee una partición para generar este reporte.");
      setReport(null);
      return;
    }

    try {
      setLoading(true);
      setMessage("");
      setReport(null);

      const data = await generateReport(
        isDiskReport
          ? { type: nextType, diskPath: selectedDiskPath, filePath: selectedDiskPath }
          : { type: nextType, id: selectedId }
      );

      setType(nextType);
      setReport(data);
      setMessage("Reporte generado correctamente.");
      await loadData();
    } catch (error) {
      setReport(null);
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadData();
  }, []);

  return (
    <main style={styles.page}>
      <h1 style={styles.title}>Reportes</h1>

      <section style={styles.box}>
        <h2 style={styles.subtitle}>Generar reporte visual</h2>
        
        <div style={styles.form}>
          <label style={styles.label}>
            Tipo de reporte
            <select style={styles.input} value={type} onChange={(e) => setType(e.target.value)}>
              {REPORT_TYPES.map((reportType) => (
                <option key={reportType.value} value={reportType.value}>
                  {reportType.label}
                </option>
              ))}
            </select>
          </label>

          {selectedReportIsDisk ? (
            <label style={styles.label}>
              Disco
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
            </label>
          ) : (
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
          )}

          <div style={styles.row}>
            <button style={styles.button} onClick={() => handleGenerate()} disabled={loading}>
              Generar reporte
            </button>
            <button style={styles.secondaryButton} onClick={loadData} disabled={loading}>
              Actualizar listas
            </button>
          </div>
        </div>

        <div style={styles.quickGrid}>
          <button style={styles.quickButton} onClick={() => handleGenerate("mbr")} disabled={loading}>
            Ver MBR
          </button>
          <button style={styles.quickButton} onClick={() => handleGenerate("disk")} disabled={loading}>
            Ver Disk
          </button>
          <button style={styles.quickButton} onClick={() => handleGenerate("tree")} disabled={loading}>
            Ver Tree
          </button>
          <button style={styles.quickButton} onClick={() => handleGenerate("inode")} disabled={loading}>
            Ver Inodos
          </button>
          <button style={styles.quickButton} onClick={() => handleGenerate("block")} disabled={loading}>
            Ver Bloques
          </button>
          <button style={styles.quickButton} onClick={() => handleGenerate("bm_inode")} disabled={loading}>
            Ver Bitmap Inodos
          </button>
          <button style={styles.quickButton} onClick={() => handleGenerate("bm_block")} disabled={loading}>
            Ver Bitmap Bloques
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
  form: {
    display: "grid",
    gap: "12px",
    maxWidth: "720px"
  },
  row: {
    display: "flex",
    gap: "10px",
    flexWrap: "wrap"
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
    padding: "12px 16px",
    borderRadius: "10px",
    border: "none",
    background: "#b85c7a",
    color: "white",
    fontWeight: "700",
    cursor: "pointer"
  },
  secondaryButton: {
    padding: "12px 16px",
    borderRadius: "10px",
    border: "1px solid #b85c7a",
    background: "white",
    color: "#b85c7a",
    fontWeight: "700",
    cursor: "pointer"
  },
  quickGrid: {
    display: "flex",
    flexWrap: "wrap",
    gap: "10px",
    marginTop: "18px"
  },
  quickButton: {
    padding: "10px 12px",
    borderRadius: "10px",
    border: "1px solid #e5bcc9",
    background: "#fffafd",
    color: "#3b202b",
    fontWeight: "700",
    cursor: "pointer"
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
    border: "1px solid #e5bcc9",
    borderRadius: "14px",
    padding: "16px",
    background: "#fffafd"
  }
};