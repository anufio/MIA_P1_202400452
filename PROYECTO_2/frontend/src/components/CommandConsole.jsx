import { useEffect, useState } from "react";
import { executeCommand } from "../api/commandApi";
import { getMountedPartitions } from "../api/diskApi";
import { useAuth } from "../hooks/useAuth";

export default function CommandConsole() {
  const { session } = useAuth();
  const [mounted, setMounted] = useState([]);
  const [selectedId, setSelectedId] = useState(session?.partId || session?.id || "");
  const [command, setCommand] = useState("");
  const [output, setOutput] = useState("La salida del comando aparecerá aquí.");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  async function loadMounted() {
    try {
      const data = await getMountedPartitions();
      const list = Array.isArray(data) ? data : [];
      const sessionPartId = session?.partId || session?.id || "";

      setMounted(list);

      if (sessionPartId) {
        setSelectedId(sessionPartId);
        return;
      }

      if (!selectedId && list.length > 0) {
        setSelectedId(list[0].id);
      }
    } catch (error) {
      setMessage(error.message);
    }
  }

  async function handleExecute() {
    try {
      setLoading(true);
      setMessage("");

      const result = await executeCommand(command, selectedId);

      setOutput(result.output || result.message || "Comando ejecutado.");
      setCommand("");
      await loadMounted();
    } catch (error) {
      setMessage(error.message);
      setOutput("No se ejecutó el comando.");
    } finally {
      setLoading(false);
    }
  }

  function handleClear() {
    setCommand("");
    setOutput("La salida del comando aparecerá aquí.");
    setMessage("");
  }

  useEffect(() => {
    loadMounted();
  }, [session?.partId, session?.id]);

  return (
    <div className="card">
      <div className="form-grid">
        <div>
          <label>Sesión activa</label>
          <div className="status-pill">
            {session?.user || session?.username} · {session?.partId || session?.id}
          </div>
        </div>

        <div>
          <label>Partición para comandos de archivos</label>

          <select
            value={selectedId}
            onChange={(event) => setSelectedId(event.target.value)}
          >
            <option value="">Seleccione una partición montada</option>
            {mounted.map((part) => (
              <option key={part.id} value={part.id}>
                {part.id} — {part.name}
              </option>
            ))}
          </select>
        </div>

        <button className="btn btn-ghost" onClick={loadMounted}>
          Actualizar montadas
        </button>

        <div>
          <label>Comando</label>

          <textarea
            value={command}
            onChange={(event) => setCommand(event.target.value)}
            placeholder='Ejemplo: mkfile -path="/docs/a.txt" -cont="Hola mundo"'
          />
        </div>

        <div className="btn-row">
          <button
            className="btn btn-primary"
            onClick={handleExecute}
            disabled={loading}
          >
            {loading ? "Ejecutando..." : "Ejecutar comando"}
          </button>

          <button className="btn btn-ghost" onClick={handleClear}>
            Limpiar
          </button>
        </div>

        {message && (
          <div
            style={{
              padding: "12px",
              borderRadius: "10px",
              background: "#f5d7e1",
              color: "#3b202b"
            }}
          >
            {message}
          </div>
        )}

        <div className="console-box">{output}</div>

        <div className="card-muted">
          <p><strong>Ejemplos:</strong></p>
          <pre style={{ whiteSpace: "pre-wrap" }}>
{`fdisk -path="/ruta/disco1.dsk" -name="Part1" -add=1 -unit=M
fdisk -path="/ruta/disco1.dsk" -name="Part1" -delete=fast
mkdir -path="/docs"
mkfile -path="/docs/nota.txt" -cont="Hola mundo"
edit -path="/docs/nota.txt" -cont="Contenido actualizado"
rename -path="/docs/nota.txt" -name="tarea.txt"
copy -path="/docs/tarea.txt" -dest="/copia.txt"
move -path="/copia.txt" -dest="/docs/copia.txt"
remove -path="/docs/copia.txt"`}
          </pre>
        </div>
      </div>
    </div>
  );
}