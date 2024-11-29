import { FC } from "react";
import { AppStatus } from "@canonical/react-components";
import { useQuery } from "@tanstack/react-query";
import { queryKeys } from "util/queryKeys";
import { fetchConfigurations } from "api/settings";

interface Props {
  className?: string;
}

const StatusBar: FC<Props> = () => {
  const { data: configurations } = useQuery({
    queryKey: [queryKeys.configuration],
    queryFn: fetchConfigurations,
  });

  const version = configurations?.api_version?.value;

  return (
    <AppStatus className="status-bar" id="status-bar">
      <span className="server-version p-text--small">Version {version}</span>
    </AppStatus>
  );
};

export default StatusBar;
