import { useState } from "react";
import { useNavigate } from "react-router-dom";
import Alert from "../components/Alert";
import { useAuth } from "../hooks/useAuth";
import { validateLogin } from "../utils/validators";
// LoginPage.jsx: contiene la página principal de inicio de sesión, mostrando un formulario para ingresar el ID de partición, el nombre de usuario y la contraseña. También incluye un mensaje de error en caso de que los datos ingresados sean incorrectos o incompletos, y utiliza el hook useAuth para manejar la autenticación del usuario.
function LoginPage() {
  const navigate = useNavigate();
  const { login, loading } = useAuth();

  const [values, setValues] = useState({
    id: "523A",
    username: "root",
    password: "123"
  });

  const [error, setError] = useState("");

  const handleChange = (event) => {
    setValues({
      ...values,
      [event.target.name]: event.target.value
    });
  };

  const handleSubmit = async (event) => {
    event.preventDefault();

    const errors = validateLogin(values);

    if (Object.keys(errors).length > 0) {
      setError(Object.values(errors)[0]);
      return;
    }

    setError("");
    await login(values);
    navigate("/");
  };

  return (
    <main className="login-wrap">
      <section className="card login-card">
        <div className="brand-mark">M</div>

        <h1 className="login-title">Iniciar sesión</h1>

        <p className="login-subtitle">
          Acceso al sistema de archivos EXT2 desde la interfaz web.
        </p>

        <form className="form-grid" onSubmit={handleSubmit}>
          <div>
            <label>ID de partición</label>
            <input
              name="id"
              value={values.id}
              onChange={handleChange}
              placeholder="Ejemplo: 523A"
            />
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
        </form>
      </section>
    </main>
  );
}

export default LoginPage;
