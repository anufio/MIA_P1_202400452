// main.jsx: es el punto de entrada de la aplicación React, que renderiza el componente principal App dentro del elemento raíz del DOM. También importa los estilos globales, de diseño y de componentes para la aplicación.
import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.jsx";
import "./styles/global.css";
import "./styles/layout.css";
import "./styles/components.css";

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
