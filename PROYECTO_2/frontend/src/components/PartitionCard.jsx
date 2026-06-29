import { SYMBOLS } from "../utils/constants";
import { formatBytes } from "../utils/formatters";
// PartitionCard.jsx: contiene el componente que representa una tarjeta de partición, mostrando su nombre, tipo, fit, tamaño y estado de montaje. También maneja los eventos de selección, montaje, desmontaje y formateo de la partición al hacer clic en los botones correspondientes.

function PartitionCard({
  partition,
  selected,
  onSelect,
  onMount,
  onUnmount,
  onFormat
}) {
  return (
    <article
      className="card item-card"
      onClick={() => onSelect?.(partition)}
      style={selected ? { borderColor: "var(--color-primary)" } : undefined}
    >
      <div className="item-head">
        <div className="symbol-badge">{SYMBOLS.partition}</div>

        <div>
          <div className="item-title">{partition.name}</div>
          <div className="item-subtitle">
            Tipo: {partition.type} · Fit: {partition.fit}
          </div>
        </div>
      </div>

      <div className="item-subtitle" style={{ marginTop: 14 }}>
        ID: {partition.id || "sin montar"} · Tamaño:{" "}
        {formatBytes(partition.size)}
      </div>

      <div className="btn-row" style={{ marginTop: 16 }}>
        {partition.mounted ? (
          <button
            className="btn btn-ghost"
            onClick={(event) => {
              event.stopPropagation();
              onUnmount?.(partition);
            }}
          >
            Desmontar
          </button>
        ) : (
          <button
            className="btn btn-primary"
            onClick={(event) => {
              event.stopPropagation();
              onMount?.(partition);
            }}
          >
            Montar
          </button>
        )}

        <button
          className="btn"
          onClick={(event) => {
            event.stopPropagation();
            onFormat?.(partition);
          }}
        >
          Formatear
        </button>
      </div>
    </article>
  );
}

export default PartitionCard;
