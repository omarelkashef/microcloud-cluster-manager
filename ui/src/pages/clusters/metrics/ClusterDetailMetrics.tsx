import { Cluster } from "types/cluster";
import { ClusterCpu } from "./ClusterCpu";
import { ClusterMemory } from "./ClusterMemory";
import { ClusterDisk } from "./ClusterDisk";
import { FC } from "react";

interface Props {
  cluster: Cluster;
}

const ClusterDetailMetrics: FC<Props> = ({ cluster }: Props) => {
  return (
    <div className="cluster-detail-metrics">
      <div className="meter-row">
        <span className="meter-row__title u-no-margin p-heading--5 u-no-padding">
          CPU
        </span>
        <ClusterCpu cluster={cluster} containerClassname="meter-row__metrics" />
      </div>
      <div className="meter-row">
        <span className="meter-row__title u-no-margin p-heading--5 u-no-padding">
          Memory
        </span>
        <ClusterMemory
          cluster={cluster}
          containerClassname="meter-row__metrics"
        />
      </div>
      <div className="meter-row">
        <span className="meter-row__title u-no-margin p-heading--5 u-no-padding">
          Disk
        </span>
        <ClusterDisk
          cluster={cluster}
          containerClassname="meter-row__metrics"
        />
      </div>
    </div>
  );
};

export default ClusterDetailMetrics;
