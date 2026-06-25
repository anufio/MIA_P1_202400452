import ResetButton from "../components/ResetButton";
// HomePage.jsx: contiene la página principal del proyecto, mostrando un encabezado y un botón para reiniciar el estado del sistema de archivos y los discos. También incluye una breve descripción del propósito de la página y su funcionalidad.
export default function HomePage() {
  return (
    <main
      style={{
        minHeight: "100vh",
        padding: "32px",
        background: "#fff7fa",
        color: "#3b202b",
      }}
    >
      <section
        style={{
          background: "white",
          border: "1px solid #e5bcc9",
          borderRadius: "18px",
          padding: "28px",
          maxWidth: "900px",
        }}
      >
        <h1 style={{ marginTop: 0 }}>MIA Proyecto 2</h1>

        <p>
          Panel principal para administrar discos, particiones, reportes y sistema
          de archivos EXT2.
        </p>

        <div style={{ marginTop: "24px" }}>
          <ResetButton />
        </div>

        <p style={{ marginTop: "16px", color: "#765062" }}>
          Use este botón antes de una prueba limpia o antes de iniciar una nueva
          demostración.
        </p>
      </section>
    </main>
  );
}
