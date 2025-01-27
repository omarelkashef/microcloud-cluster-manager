import { Icon } from "@canonical/react-components";
import { useQuery } from "@tanstack/react-query";
import React, { FC } from "react";
import { queryKeys } from "util/queryKeys";
import { fetchConfigurations } from "api/settings";

type Props = {
  clusterName: string;
};

const ClusterMetricsButton: FC<Props> = ({ clusterName }) => {
  const { data: configurations } = useQuery({
    queryKey: [queryKeys.configuration],
    queryFn: fetchConfigurations,
  });

  const baseUrl = configurations?.grafana_base_url?.value;

  if (!baseUrl) {
    return null;
  }

  return (
    <a
      className="p-segmented-control__button p-button u-no-margin--bottom has-icon"
      href={`${baseUrl}/lxd?orgId=1&var-job=${clusterName}`}
      target="_blank"
      rel="noopener noreferrer"
    >
      <Icon name="external-link" />
      <span>Metrics</span>
    </a>
  );
};

export default ClusterMetricsButton;
