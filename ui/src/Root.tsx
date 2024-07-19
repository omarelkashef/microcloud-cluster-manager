import { FC } from "react";
import App from "./App";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Application } from "@canonical/react-components";

import Navigation from "./components/Navigation";
import { AuthProvider } from "context/auth";
import StatusBar from "components/StatusBar";

const queryClient = new QueryClient();

const Root: FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <Application>
          <Navigation />
          <App />
          <StatusBar />
        </Application>
      </AuthProvider>
    </QueryClientProvider>
  );
};

export default Root;
