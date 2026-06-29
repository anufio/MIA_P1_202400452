// Alert.jsx: contiene el componente de alerta que muestra mensajes de información, advertencia o error en la interfaz de usuario, dependiendo del tipo de alerta especificado.
function Alert({ type = "info", children }) {
  if (!children) {
    return null;
  }

  return <div className={`alert ${type}`}>{children}</div>;
}

export default Alert;
