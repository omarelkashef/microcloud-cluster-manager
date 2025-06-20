import React, { FC } from "react";
import { Cluster } from "types/cluster";
import ConfigureClusterButton from "pages/clusters/actions/ConfigureClusterButton";
import RemoveClusterButton from "pages/clusters/actions/RemoveClusterButton";
import ClusterUiButton from "pages/clusters/actions/ClusterUiButton";
import ClusterMetricsButton from "pages/clusters/actions/ClusterMetricsButton";
import { List } from "@canonical/react-components";

type Props = {
  cluster: Cluster;
};

const ClusterActions: FC<Props> = ({ cluster }) => {
  return (
    <List
      inline
      className="actions-list"
      items={[
        <ConfigureClusterButton
          cluster={cluster}
          appearance="base"
          key="configure"
        />,
        <ClusterUiButton uiUrl={cluster.ui_url} appearance="base" key="ui" />,
        <ClusterMetricsButton
          clusterName={cluster.name}
          appearance="base"
          className="p-button--base"
          key="metrics"
        />,
        <RemoveClusterButton
          clusterName={cluster.name}
          appearance="base"
          key="remove"
        />,
      ]}
    />
  );
};

export default ClusterActions;
