import Loading from "./Loading";
// ReportViewer.jsx: contiene el componente que representa un visor de reportes, mostrando la información del reporte, su tipo, archivo y contenido. También maneja el estado de carga y la visualización de un mensaje cuando no hay reporte generado todavía.

function ReportViewer({ report, loading }) {
  if (loading) {
    return <Loading text="Cargando reporte..." />;
  }

  if (!report) {
    return (
      <div className="card empty-state">
        No hay reporte generado todavía.
      </div>
    );
  }

  const url = report.imageUrl || report.url;

  if (url) {
    return (
      <div className="report-frame card">
        <p>
          <strong>Tipo:</strong> {report.type}
        </p>

        <p>
          <strong>Archivo:</strong> {report.path}
        </p>

        <div style={{ marginTop: "16px", overflow: "auto" }}>
          <img
            src={`${url}?t=${Date.now()}`}
            alt={report.type || "Reporte"}
            style={{
              maxWidth: "100%",
              display: "block",
              margin: "0 auto",
              borderRadius: "12px",
            }}
          />
        </div>

        <a
          href={url}
          target="_blank"
          rel="noreferrer"
          className="btn btn-primary"
          style={{ marginTop: "16px", display: "inline-block" }}
        >
          Abrir reporte
        </a>
      </div>
    );
  }

  return (
    <div className="report-frame card">
      <pre style={{ margin: 0, whiteSpace: "pre-wrap" }}>
        {report.content || "Reporte sin contenido."}
      </pre>
    </div>
  );
}

export default ReportViewer;
