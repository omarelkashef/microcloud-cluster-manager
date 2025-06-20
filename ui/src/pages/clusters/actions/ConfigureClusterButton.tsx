import { Button, Icon } from "@canonical/react-components";
import React, { FC } from "react";
import usePanelParams from "context/usePanelParams";
import { Cluster } from "types/cluster";
import classnames from "classnames";

interface Props {
  cluster: Cluster;
  appearance?: string;
  className?: string;
}

const ConfigureClusterButton: FC<Props> = ({
  cluster,
  appearance = "",
  className,
}) => {
  const panelParams = usePanelParams();

  return (
    <Button
      appearance={appearance}
      className={classnames("u-no-margin--bottom", className)}
      hasIcon
      onClick={() => panelParams.openConfigureCluster(cluster.name)}
    >
      <Icon name="external-link" />
      <span>Configure</span>
    </Button>
  );
};

export default ConfigureClusterButton;
