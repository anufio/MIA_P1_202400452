import Loading from "./Loading";

function ReportViewer({ report, loading }) {
  if (loading) {
    return <Loading text="Cargando reporte..." />;
  }

  if (!report) {
    return <div className="card empty-state">No hay reporte generado todavía.</div>;
  }

  const url = report.imageUrl || report.url;
  const format = String(report.format || "").toLowerCase();

  if (url) {
    const cacheUrl = `${url}${url.includes("?") ? "&" : "?"}t=${Date.now()}`;
    const isImage = format === "image" || format === "svg" || /\.(png|jpg|jpeg|svg)$/i.test(url);
    const isFrame = format === "pdf" || format === "text" || !isImage;

    return (
      <div className="report-frame card">
        <p>
          <strong>Tipo:</strong> {report.type}
        </p>
        <p>
          <strong>Archivo:</strong> {report.path}
        </p>

        <div style={{ marginTop: "16px", overflow: "auto" }}>
          {isImage && (
            <img
              src={cacheUrl}
              alt={report.type || "Reporte"}
              style={{
                maxWidth: "100%",
                display: "block",
                margin: "0 auto",
                borderRadius: "12px",
                border: "1px solid #e5bcc9"
              }}
            />
          )}

          {isFrame && (
            <iframe
              title={report.type || "Reporte"}
              src={cacheUrl}
              style={{
                width: "100%",
                minHeight: "520px",
                border: "1px solid #e5bcc9",
                borderRadius: "12px",
                background: "white"
              }}
            />
          )}
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