// App.jsx: es el archivo principal de la aplicación React, que configura el enrutamiento, los proveedores de contexto y la estructura general de la interfaz de usuario. Define rutas protegidas para usuarios autenticados y maneja la navegación entre diferentes páginas de la aplicación.
import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import Navbar from "./components/Navbar";
import Sidebar from "./components/Sidebar";
import { AuthProvider } from "./context/AuthContext";
import { DiskProvider } from "./context/DiskContext";
import { FileSystemProvider } from "./context/FileSystemContext";
import { useAuth } from "./hooks/useAuth";
import CommandPage from "./pages/CommandPage";
import Disks from "./pages/Disks";
import ExplorerPage from "./pages/ExplorerPage";
import FileViewerPage from "./pages/FileViewerPage";
import HomePage from "./pages/HomePage";
import LoginPage from "./pages/LoginPage";
import PartitionPage from "./pages/PartitionPage";
import ReportsPage from "./pages/ReportsPage";

function ProtectedRoute({ children }) {
  const { isAuthenticated } = useAuth();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return children;
}

function AppLayout() {
  return (
    <div className="app-shell">
      <Sidebar />

      <main className="content-area">
        <Navbar />

        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/disks" element={<Disks />} />
          <Route path="/partitions" element={<PartitionPage />} />
          <Route path="/explorer" element={<ExplorerPage />} />
          <Route path="/file-viewer" element={<FileViewerPage />} />
          <Route path="/reports" element={<ReportsPage />} />
          <Route path="/commands" element={<CommandPage />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </main>
    </div>
  );
}

function AppRoutes() {
  const { isAuthenticated } = useAuth();

  return (
    <Routes>
      <Route
        path="/login"
        element={
          isAuthenticated ? <Navigate to="/" replace /> : <LoginPage />
        }
      />

      <Route
        path="/*"
        element={
          <ProtectedRoute>
            <AppLayout />
          </ProtectedRoute>
        }
      />
    </Routes>
  );
}

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <DiskProvider>
          <FileSystemProvider>
            <AppRoutes />
          </FileSystemProvider>
        </DiskProvider>
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;