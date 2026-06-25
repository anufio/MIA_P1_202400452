import { useEffect, useState } from "react";
// Disks.jsx: contiene la página principal de administración de discos, mostrando un encabezado y un formulario para crear y eliminar discos. También incluye una lista de discos disponibles y el disco seleccionado actualmente, utilizando los hooks useDisks y useDiskSelection para manejar la carga de discos y la selección del disco.
const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api";

export default function Disks() {
  const [disks, setDisks] = useState([]);
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  const [name, setName] = useState("");
  const [size, setSize] = useState("");
  const [unit, setUnit] = useState("M");
  const [fit, setFit] = useState("FF");

  async function request(path, options = {}) {
    const res = await fetch(`${API_URL}${path}`, {
      ...options,
      headers: {
        "Content-Type": "application/json",
      },
    });

    const data = await res.json().catch(() => ({}));

    if (!res.ok || data.success === false) {
      throw new Error(data.error || data.message || "Error en la petición");
    }

    return data;
  }

  async function loadDisks() {
    try {
      setLoading(true);
      setMessage("");

      const data = await request("/disks");
      const list = data.disks || data.data || [];

      setDisks(Array.isArray(list) ? list : []);
    } catch (error) {
      setMessage("Error conectando con backend: " + error.message);
      setDisks([]);
    } finally {
      setLoading(false);
    }
  }

  async function createDisk(e) {
    e.preventDefault();

    if (!name.trim()) {
      setMessage("Ingrese el nombre del disco.");
      return;
    }

    if (!size || Number(size) <= 0) {
      setMessage("Ingrese un tamaño válido.");
      return;
    }

    try {
      setLoading(true);
      setMessage("");

      await request("/disks", {
        method: "POST",
        body: JSON.stringify({
          name: name.trim(),
          size: Number(size),
          unit,
          fit,
        }),
      });

      setName("");
      setSize("");
      setUnit("M");
      setFit("FF");

      setMessage("Disco creado correctamente.");
      await loadDisks();
    } catch (error) {
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  async function deleteDisk(path) {
    if (!window.confirm("¿Eliminar disco?")) return;

    try {
      setLoading(true);
      setMessage("");

      await request("/disks/delete", {
        method: "POST",
        body: JSON.stringify({ path }),
      });

      setMessage("Disco eliminado correctamente.");
      await loadDisks();
    } catch (error) {
      setMessage(error.message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadDisks();
  }, []);

  return (
    <div style={styles.page}>
      <h1 style={styles.title}>Discos</h1>

      <div style={styles.box}>
        <h2 style={styles.subtitle}>Crear nuevo disco</h2>

        <form onSubmit={createDisk} style={styles.form}>
          <label style={styles.label}>
            Nombre
            <input
              style={styles.input}
              value={name}
              placeholder="Ejemplo: disco1.dsk"
              onChange={(e) => setName(e.target.value)}
            />
          </label>

          <label style={styles.label}>
            Tamaño
            <input
              style={styles.input}
              type="number"
              value={size}
              placeholder="Ejemplo: 20"
              onChange={(e) => setSize(e.target.value)}
            />
          </label>

          <label style={styles.label}>
            Unidad
            <select style={styles.input} value={unit} onChange={(e) => setUnit(e.target.value)}>
              <option value="K">K</option>
              <option value="M">M</option>
            </select>
          </label>

          <label style={styles.label}>
            Fit
            <select style={styles.input} value={fit} onChange={(e) => setFit(e.target.value)}>
              <option value="FF">FF</option>
              <option value="BF">BF</option>
              <option value="WF">WF</option>
            </select>
          </label>

          <button style={styles.button} disabled={loading}>
            Crear disco
          </button>
        </form>

        {message && <p style={styles.message}>{message}</p>}
      </div>

      <div style={styles.box}>
        <h2 style={styles.subtitle}>Discos creados</h2>

        {loading && <p>Cargando...</p>}

        {!loading && disks.length === 0 && (
          <p>No hay discos creados todavía.</p>
        )}

        <div style={styles.grid}>
          {disks.map((disk) => (
            <div key={disk.id || disk.path} style={styles.card}>
              <h3>{disk.name}</h3>
              <p><strong>Ruta:</strong> {disk.path}</p>
              <p><strong>Tamaño:</strong> {disk.size} bytes</p>
              <p><strong>Fit:</strong> {disk.fit}</p>

              <button style={styles.deleteButton} onClick={() => deleteDisk(disk.path)}>
                Eliminar
              </button>
            </div>
          ))}
        </div>
      </div>
    </div>
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
    maxWidth: "900px",
  },
  form: {
    display: "grid",
    gap: "12px",
    maxWidth: "480px",
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
  deleteButton: {
    padding: "10px 12px",
    borderRadius: "10px",
    border: "none",
    background: "#8f2f4f",
    color: "white",
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
