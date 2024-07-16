import { useQuery } from "@tanstack/react-query";
import { fetchClusters } from "api/clusters";
import { createContext, FC, ReactNode, useContext } from "react";
import { queryKeys } from "util/queryKeys";

interface ContextProps {
  isAuthenticated: boolean;
  isAuthLoading: boolean;
}

const initialState: ContextProps = {
  isAuthenticated: false,
  isAuthLoading: true,
};

export const AuthContext = createContext<ContextProps>(initialState);

interface ProviderProps {
  children: ReactNode;
}

export const AuthProvider: FC<ProviderProps> = ({ children }) => {
  // FIXME: this should query /1.0 when the endpoint is ready to check if the user is authenticated
  const { error, isLoading } = useQuery({
    queryKey: [queryKeys.clusters],
    queryFn: fetchClusters,
    retry: false,
  });

  const isAuthenticated = !isLoading && error?.message !== "not authorized";

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        isAuthLoading: isLoading,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export function useAuth() {
  return useContext(AuthContext);
}
