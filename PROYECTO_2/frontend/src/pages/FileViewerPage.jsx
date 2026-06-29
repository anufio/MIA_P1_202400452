import { useEffect, useState } from "react";
import { useSearchParams } from "react-router-dom";
import Loading from "../components/Loading";
import PageHeader from "../components/PageHeader";
import { getFileContent } from "../api/filesystemApi";
// FileViewerPage.jsx: contiene la página principal del visor de archivos, mostrando un encabezado y el contenido del archivo seleccionado. Utiliza el hook useSearchParams para obtener la ruta del archivo desde los parámetros de búsqueda de la URL y carga el contenido del archivo desde el backend a través de filesystemApi.
function FileViewerPage() {
  const [searchParams] = useSearchParams();
  const path = searchParams.get("path") || "/users.txt";

  const [content, setContent] = useState("");
  const [loading, setLoading] = useState(true);

  const loadContent = async () => {
    setLoading(true);
    const data = await getFileContent(path);
    setContent(data);
    setLoading(false);
  };

  useEffect(() => {
    loadContent();
  }, [path]);

  return (
    <section className="page">
      <PageHeader
        kicker="Archivo"
        title="Visualizador de contenido"
        description={path}
      />

      {loading ? (
        <Loading text="Cargando archivo..." />
      ) : (
        <article className="card">
          <div className="file-content">{content}</div>
        </article>
      )}
    </section>
  );
}

export default FileViewerPage;
