// Loading.jsx: contiene el componente que representa un estado de carga, mostrando un mensaje de texto mientras se espera la carga de información.
function Loading({ text = "Cargando información..." }) {
  return <div className="empty-state">{text}</div>;
}

export default Loading;
