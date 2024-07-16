import { FC } from "react";
import App from "./App";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Application } from "@canonical/react-components";

import Navigation from "./components/Navigation";
import { AuthProvider } from "context/auth";

const queryClient = new QueryClient();

const Root: FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <Application>
          <Navigation />
          <App />
        </Application>
      </AuthProvider>
    </QueryClientProvider>
  );
};

export default Root;
