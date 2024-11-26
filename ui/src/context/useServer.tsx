import { useQuery } from "@tanstack/react-query";
import { queryKeys } from "util/queryKeys";
import { UseQueryResult } from "@tanstack/react-query/src/types";
import { Server } from "types/server";
import { fetchServer } from "api/server";

export const useServer = (): UseQueryResult<Server> => {
  return useQuery({
    queryKey: [queryKeys.server],
    queryFn: fetchServer,
  });
};
