import { Icon } from "@canonical/react-components";
import React, { FC } from "react";
import classnames from "classnames";

type Props = {
  uiUrl: string;
  className?: string;
};

const ClusterUiButton: FC<Props> = ({ uiUrl, className }) => {
  if (!uiUrl) {
    return null;
  }

  return (
    <a
      className={classnames("p-button u-no-margin--bottom has-icon", className)}
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
