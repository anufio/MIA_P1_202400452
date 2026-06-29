import { useState } from "react";
import Alert from "../components/Alert";
import DiskCard from "../components/DiskCard";
import Loading from "../components/Loading";
import PageHeader from "../components/PageHeader";
import { createDisk, deleteDisk } from "../api/diskApi";
import { useDiskSelection, useDisks } from "../hooks/useDisks";
import { validateDiskForm } from "../utils/validators";
// DiskPage.jsx: contiene la página principal de administración de discos, mostrando un encabezado y un formulario para crear y eliminar discos. También incluye una lista de discos disponibles y el disco seleccionado actualmente, utilizando los hooks useDisks y useDiskSelection para manejar la carga de discos y la selección del disco.
function DiskPage() {
  const { disks, loading, loadDisks } = useDisks();
  const { selectedDisk, setSelectedDisk } = useDiskSelection();

  const [form, setForm] = useState({
    name: "nuevo.dsk",
    path: "/home/anufio/mia/cali/nuevo.dsk",
    size: 10,
    unit: "M",
    fit: "WF"
  });

  const [message, setMessage] = useState("");
  const [messageType, setMessageType] = useState("info");

  const handleChange = (event) => {
    setForm({
      ...form,
      [event.target.name]: event.target.value
    });
  };

  const handleCreate = async () => {
    const errors = validateDiskForm(form);

    if (Object.keys(errors).length > 0) {
      setMessage(Object.values(errors)[0]);
      setMessageType("error");
      return;
    }

    const result = await createDisk(form);
    setMessage(result.message || "Solicitud enviada.");
    setMessageType(result.success === false ? "error" : "success");
    await loadDisks();
  };

  const handleDelete = async () => {
    if (!selectedDisk) {
      setMessage("Seleccione un disco antes de eliminar.");
      setMessageType("error");
      return;
    }

    const result = await deleteDisk(selectedDisk.path);
    setMessage(result.message || "Solicitud enviada.");
    setMessageType(result.success === false ? "error" : "success");
    await loadDisks();
  };

  return (
    <section className="page">
      <PageHeader
        kicker="Discos"
        title="Administración de discos"
        description="Cree, elimine y seleccione discos para trabajar con sus particiones."
      />

      <div className="grid grid-2">
        <article className="card">
          <div className="card-title">Crear disco</div>

          <div className="form-grid">
            <div>
              <label>Nombre</label>
              <input name="name" value={form.name} onChange={handleChange} />
            </div>

            <div>
              <label>Ruta</label>
              <input name="path" value={form.path} onChange={handleChange} />
            </div>

            <div className="form-row">
              <div>
                <label>Tamaño</label>
                <input
                  name="size"
                  type="number"
                  value={form.size}
                  onChange={handleChange}
                />
              </div>

              <div>
                <label>Unidad</label>
                <select name="unit" value={form.unit} onChange={handleChange}>
                  <option value="B">Bytes</option>
                  <option value="K">Kilobytes</option>
                  <option value="M">Megabytes</option>
                </select>
              </div>
            </div>

            <div>
              <label>Fit</label>
              <select name="fit" value={form.fit} onChange={handleChange}>
                <option value="FF">First Fit</option>
                <option value="BF">Best Fit</option>
                <option value="WF">Worst Fit</option>
              </select>
            </div>

            <div className="btn-row">
              <button className="btn btn-primary" onClick={handleCreate}>
                Crear disco
              </button>

              <button className="btn btn-danger" onClick={handleDelete}>
                Eliminar seleccionado
              </button>
            </div>

            <Alert type={messageType}>{message}</Alert>
          </div>
        </article>

        <article className="card">
          <div className="card-title">Disco seleccionado</div>

          {selectedDisk ? (
            <>
              <p className="card-muted">{selectedDisk.name}</p>
              <p className="card-muted">{selectedDisk.path}</p>
            </>
          ) : (
            <p className="card-muted">Seleccione un disco de la lista.</p>
          )}
        </article>
      </div>

      {loading ? (
        <Loading />
      ) : (
        <div className="grid grid-3">
          {disks.map((disk) => (
            <DiskCard
              key={disk.id}
              disk={disk}
              selected={selectedDisk?.id === disk.id}
              onSelect={setSelectedDisk}
            />
          ))}
        </div>
      )}
    </section>
  );
}

export default DiskPage;
