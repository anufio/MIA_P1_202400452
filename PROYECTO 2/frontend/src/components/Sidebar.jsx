import { NavLink } from "react-router-dom";
import { ROUTES, SYMBOLS } from "../utils/constants";

const links = [
  { to: ROUTES.home, label: "Inicio", symbol: SYMBOLS.home },
  { to: ROUTES.disks, label: "Discos", symbol: SYMBOLS.disk },
  { to: ROUTES.partitions, label: "Particiones", symbol: SYMBOLS.partition },
  { to: ROUTES.explorer, label: "Explorador", symbol: SYMBOLS.folder },
  { to: ROUTES.reports, label: "Reportes", symbol: SYMBOLS.report },
  { to: ROUTES.commands, label: "Comandos", symbol: SYMBOLS.console }
];

function Sidebar() {
  return (
    <aside className="sidebar">
      <div className="brand">
        <div className="brand-mark">M</div>
        <div className="brand-title">MIA Fase 2</div>
        <div className="brand-subtitle">Sistema de archivos EXT2</div>
      </div>

      <nav className="nav-menu">
        {links.map((link) => (
          <NavLink
            key={link.to}
            to={link.to}
            className={({ isActive }) =>
              isActive ? "nav-link active" : "nav-link"
            }
          >
            <span className="nav-symbol">{link.symbol}</span>
            <span>{link.label}</span>
          </NavLink>
        ))}
      </nav>
    </aside>
  );
}

export default Sidebar;
