import { SYMBOLS } from "../utils/constants";
import { formatBytes } from "../utils/formatters";
// DiskCard.jsx: contiene el componente que representa una tarjeta de disco, mostrando su nombre, ruta, tamaño y tipo de ajuste. También maneja el evento de selección del disco al hacer clic en él.

function DiskCard({ disk, selected, onSelect }) {
  return (
    <article
      className="card item-card"
      onClick={() => onSelect?.(disk)}
      style={selected ? { borderColor: "var(--color-primary)" } : undefined}
    >
      <div className="item-head">
        <div className="symbol-badge">{SYMBOLS.disk}</div>

        <div>
          <div className="item-title">{disk.name}</div>
          <div className="item-subtitle">{disk.path}</div>
        </div>
      </div>

      <div className="item-subtitle" style={{ marginTop: 14 }}>
        Tamaño: {formatBytes(disk.size)} · Fit: {disk.fit || "WF"}
      </div>
    </article>
  );
}

export default DiskCard;
