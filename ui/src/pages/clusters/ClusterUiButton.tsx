import { Icon } from "@canonical/react-components";
import React, { FC } from "react";

type Props = {
  uiUrl: string;
};

const ClusterUiButton: FC<Props> = ({ uiUrl }) => {
  if (!uiUrl) {
    return null;
  }

  return (
    <a
      className="p-segmented-control__button p-button u-no-margin--bottom has-icon"
      href={uiUrl}
      target="_blank"
      rel="noopener noreferrer"
    >
      <Icon name="external-link" />
      <span>LXD UI</span>
    </a>
  );
};

export default ClusterUiButton;
