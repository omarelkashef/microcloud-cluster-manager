import { FC } from "react";
import { AppStatus } from "@canonical/react-components";

interface Props {
  className?: string;
}

const StatusBar: FC<Props> = () => {
  return (
    <AppStatus className="status-bar" id="status-bar">
      <span className="server-version p-text--small">Version 0.1</span>
    </AppStatus>
  );
};

export default StatusBar;
