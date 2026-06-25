import { useNavigate } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";
import { APP_INFO } from "../utils/constants";
// Navbar.jsx: contiene el componente que representa la barra de navegación superior, mostrando el nombre de la aplicación, el nombre del estudiante y su carnet. También maneja la autenticación del usuario, mostrando un botón para iniciar sesión o cerrar sesión según corresponda.

function Navbar() {
  const navigate = useNavigate();
  const { session, isAuthenticated, logout } = useAuth();

  const handleLogout = async () => {
    await logout();
    navigate("/login");
  };

  return (
    <header className="topbar">
      <div>
        <div className="topbar-title">{APP_INFO.name}</div>
        <div className="topbar-subtitle">
          {APP_INFO.student} · {APP_INFO.carnet}
        </div>
      </div>

      <div className="btn-row">
        {isAuthenticated ? (
          <>
            <span className="status-pill">
              {session.username} · {session.id}
            </span>

            <button className="btn btn-ghost" onClick={handleLogout}>
              Cerrar sesión
            </button>
          </>
        ) : (
          <button className="btn btn-primary" onClick={() => navigate("/login")}>
            Iniciar sesión
          </button>
        )}
      </div>
    </header>
  );
}

export default Navbar;
