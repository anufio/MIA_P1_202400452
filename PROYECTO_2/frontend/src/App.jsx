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

function RequireSession({ children }) {
  const { isAuthenticated, checkingSession } = useAuth();

  if (checkingSession) {
    return <section className="page">Validando sesión...</section>;
  }

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
          <Route path="/reports" element={<ReportsPage />} />
          <Route
            path="/explorer"
            element={
              <RequireSession>
                <ExplorerPage />
              </RequireSession>
            }
          />
          <Route
            path="/file-viewer"
            element={
              <RequireSession>
                <FileViewerPage />
              </RequireSession>
            }
          />
          <Route
            path="/commands"
            element={
              <RequireSession>
                <CommandPage />
              </RequireSession>
            }
          />
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
        element={isAuthenticated ? <Navigate to="/explorer" replace /> : <LoginPage />}
      />
      <Route path="/*" element={<AppLayout />} />
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