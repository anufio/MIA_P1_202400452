import { SYMBOLS } from "../utils/constants";
import { formatBytes } from "../utils/formatters";
// BrowserItem.jsx: contiene el componente que representa un elemento del explorador de archivos, mostrando su nombre, ruta, tipo, permisos y tamaño. También maneja el evento de apertura del elemento al hacer clic en él.

function BrowserItem({ item, onOpen }) {
  const symbol = item.type === "folder" ? SYMBOLS.folder : SYMBOLS.file;
  const typeLabel = item.type === "folder" ? "Carpeta" : "Archivo";

  return (
    <article className="card item-card" onClick={() => onOpen?.(item)}>
      <div className="item-head">
        <div className="symbol-badge">{symbol}</div>

        <div>
          <div className="item-title">{item.name}</div>
          <div className="item-subtitle">{item.path}</div>
        </div>
      </div>

      <div className="item-subtitle" style={{ marginTop: 14 }}>
        Tipo: {typeLabel} · Permisos: {item.permissions || "---"} · Tamaño:{" "}
        {formatBytes(item.size)}
      </div>
    </article>
  );
}

export default BrowserItem;
