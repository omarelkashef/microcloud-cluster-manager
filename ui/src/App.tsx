import React, { FC, lazy, Suspense } from "react";
import { Navigate, Route, Routes, useLocation } from "react-router-dom";
import NoMatch from "pages/NoMatch";
import Settings from "pages/Settings";
import { useAuth } from "context/auth";

const ClusterList = lazy(() => import("pages/clusters/ClusterList"));
const Login = lazy(() => import("pages/Login"));

const App: FC = () => {
  const { pathname } = useLocation();
  const { isAuthLoading, isAuthenticated } = useAuth();
  const isLoginPath = pathname === "/ui/login";

  if (!isAuthLoading && !isLoginPath && !isAuthenticated) {
    window.location.href = "/ui/login";
    return null;
  }

  if (!isAuthLoading && isLoginPath && isAuthenticated) {
    window.location.href = "/ui/sites";
    return null;
  }

  return (
    <Suspense fallback={<div>Loading</div>}>
      <Routes>
        <Route path="/" element={<Navigate to="/ui/sites" replace={true} />} />
        <Route
          path="/ui"
          element={<Navigate to="/ui/sites" replace={true} />}
        />
        <Route path="/ui/login" element={<Login />} />
        <Route path="/ui/sites" element={<ClusterList />} />
        <Route path="/ui/sites/:activeTab" element={<ClusterList />} />
        <Route path="/ui/settings" element={<Settings />} />
        <Route path="*" element={<NoMatch />} />
      </Routes>
    </Suspense>
  );
};

export default App;
