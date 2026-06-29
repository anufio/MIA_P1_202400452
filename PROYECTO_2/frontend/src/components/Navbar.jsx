import { useNavigate } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";
import { APP_INFO } from "../utils/constants";

function Navbar() {
  const navigate = useNavigate();
  const { session, isAuthenticated, logout } = useAuth();

  const handleLogout = async () => {
    await logout();
    navigate("/");
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
              {session.user || session.username} · {session.partId || session.id}
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