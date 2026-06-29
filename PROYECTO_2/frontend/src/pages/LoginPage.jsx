import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getMountedPartitions } from "../api/diskApi";
import Alert from "../components/Alert";
import { useAuth } from "../hooks/useAuth";
import { validateLogin } from "../utils/validators";

function LoginPage() {
  const navigate = useNavigate();
  const { login, loading } = useAuth();

  const [mounted, setMounted] = useState([]);
  const [values, setValues] = useState({
    id: "",
    username: "root",
    password: "123"
  });

  const [error, setError] = useState("");

  useEffect(() => {
    async function loadMounted() {
      try {
        const data = await getMountedPartitions();
        const list = Array.isArray(data) ? data : [];

        setMounted(list);

        if (list.length > 0) {
          setValues((current) => ({
            ...current,
            id: current.id || list[0].id
          }));
        }
      } catch {
        setMounted([]);
      }
    }

    loadMounted();
  }, []);

  const handleChange = (event) => {
    setValues({
      ...values,
      [event.target.name]: event.target.value
    });
  };

  const handleSubmit = async (event) => {
    event.preventDefault();

    const cleanValues = {
      ...values,
      id: values.id.trim().toUpperCase(),
      username: values.username.trim()
    };

    const errors = validateLogin(cleanValues);

    if (Object.keys(errors).length > 0) {
      setError(Object.values(errors)[0]);
      return;
    }

    try {
      setError("");
      await login(cleanValues);
      navigate("/explorer");
    } catch (err) {
      setError(err.message || "No se pudo iniciar sesión.");
    }
  };

  return (
    <main className="login-wrap">
      <section className="card login-card">
        <div className="brand-mark">M</div>

        <h1 className="login-title">Iniciar sesión</h1>

        <p className="login-subtitle">
          Primero cree el disco, cree una partición, móntela y formatéela. Después
          inicie sesión con el ID de la partición montada.
        </p>

        <form className="form-grid" onSubmit={handleSubmit}>
          <div>
            <label>ID de partición</label>
            <select name="id" value={values.id} onChange={handleChange}>
              <option value="">Seleccione una partición montada</option>
              {mounted.map((part) => (
                <option key={part.id} value={part.id}>
                  {part.id} — {part.name}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label>Usuario</label>
            <input
              name="username"
              value={values.username}
              onChange={handleChange}
              placeholder="root"
            />
          </div>

          <div>
            <label>Contraseña</label>
            <input
              name="password"
              type="password"
              value={values.password}
              onChange={handleChange}
              placeholder="123"
            />
          </div>

          <Alert type="error">{error}</Alert>

          <button className="btn btn-primary" type="submit" disabled={loading}>
            {loading ? "Validando..." : "Entrar"}
          </button>

          <button
            className="btn btn-ghost"
            type="button"
            onClick={() => navigate("/disks")}
          >
            Volver a discos
          </button>
        </form>
      </section>
    </main>
  );
}

export default LoginPage;