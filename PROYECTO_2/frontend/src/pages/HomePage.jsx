import { Link } from "react-router-dom";
import { ROUTES } from "../utils/constants";

export default function HomePage() {
  return (
    <main style={styles.page}>
      <section style={styles.hero}>
        <p style={styles.kicker}>MIA Proyecto Fase 2</p>
        <h1 style={styles.title}>Visualizador web del sistema de archivos EXT2</h1>

        

        <div style={styles.actions}>
          <Link to={ROUTES.disks} style={styles.primaryLink}>
            Crear disco
          </Link>
          <Link to={ROUTES.partitions} style={styles.secondaryLink}>
            Ver particiones
          </Link>
        </div>
      </section>

      
    </main>
  );
}

const styles = {
  page: {
    minHeight: "100vh",
    padding: "32px",
    background: "#fff7fa",
    color: "#3b202b"
  },
  hero: {
    background: "white",
    border: "1px solid #e5bcc9",
    borderRadius: "18px",
    padding: "30px",
    maxWidth: "960px",
    marginBottom: "22px"
  },
  kicker: {
    margin: 0,
    color: "#b85c7a",
    fontWeight: "700"
  },
  title: {
    marginTop: "8px",
    marginBottom: "12px",
    fontSize: "32px"
  },
  text: {
    maxWidth: "760px",
    lineHeight: 1.6,
    color: "#765062"
  },
  actions: {
    display: "flex",
    gap: "12px",
    flexWrap: "wrap",
    marginTop: "22px"
  },
  primaryLink: {
    padding: "12px 18px",
    borderRadius: "12px",
    background: "#b85c7a",
    color: "white",
    textDecoration: "none",
    fontWeight: "700"
  },
  secondaryLink: {
    padding: "12px 18px",
    borderRadius: "12px",
    border: "1px solid #b85c7a",
    color: "#b85c7a",
    textDecoration: "none",
    fontWeight: "700"
  },
  steps: {
    display: "grid",
    gridTemplateColumns: "repeat(auto-fit, minmax(220px, 1fr))",
    gap: "14px",
    maxWidth: "960px"
  },
  step: {
    display: "grid",
    gap: "8px",
    background: "white",
    border: "1px solid #e5bcc9",
    borderRadius: "16px",
    padding: "18px"
  }
};