import CommandConsole from "../components/CommandConsole";
import PageHeader from "../components/PageHeader";
// CommandPage.jsx: contiene la página principal de la consola de comandos, mostrando un encabezado y el componente CommandConsole para ejecutar comandos en el backend. También incluye una tabla con los comandos disponibles y su uso esperado.
function CommandPage() {
  return (
    <section className="page">
      <PageHeader
        kicker="Consola"
        title="Ejecución de comandos"
        description="Área para enviar comandos al backend y reflejar los cambios en el visualizador."
      />

      <CommandConsole />

      <article className="card">
        <div className="card-title">Comandos de Fase 2</div>

        <table className="table">
          <thead>
            <tr>
              <th>Comando</th>
              <th>Uso esperado</th>
            </tr>
          </thead>

          <tbody>
            <tr>
              <td>fdisk -add</td>
              <td>Agregar o quitar espacio a una partición existente.</td>
            </tr>

            <tr>
              <td>fdisk -delete</td>
              <td>Eliminar una partición con modo fast o full.</td>
            </tr>

            <tr>
              <td>remove</td>
              <td>Eliminar archivo o carpeta con validación de permisos.</td>
            </tr>

            <tr>
              <td>edit</td>
              <td>Modificar el contenido de un archivo.</td>
            </tr>

            <tr>
              <td>rename</td>
              <td>Cambiar el nombre de un archivo o carpeta.</td>
            </tr>

            <tr>
              <td>copy</td>
              <td>Copiar archivo o carpeta hacia otro destino.</td>
            </tr>

            <tr>
              <td>move</td>
              <td>Mover archivo o carpeta hacia otro destino.</td>
            </tr>
          </tbody>
        </table>
      </article>
    </section>
  );
}

export default CommandPage;
