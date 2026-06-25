// PageHeader.jsx: contiene el componente que representa el encabezado de la página, mostrando un título, una descripción opcional y un conjunto de acciones opcionales. También maneja la visualización de un "kicker" opcional que aparece encima del título.
function PageHeader({ kicker, title, description, actions }) {
  return (
    <header className="page-header">
      <div>
        {kicker && <div className="page-kicker">{kicker}</div>}
        <h1 className="page-title">{title}</h1>
        {description && <p className="page-description">{description}</p>}
      </div>

      {actions && <div className="btn-row">{actions}</div>}
    </header>
  );
}

export default PageHeader;
