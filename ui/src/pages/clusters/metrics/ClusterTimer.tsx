import { useQuery } from "@tanstack/react-query";
import { fetchCluster } from "api/clusters";
import { FC, useEffect, useState } from "react";
import { Cluster } from "types/cluster";
import { getSecondsSinceLastHeartbeat } from "util/helpers";
import { queryKeys } from "util/queryKeys";

interface Props {
  cluster: Cluster;
}

const ClusterTimer: FC<Props> = ({ cluster }: Props) => {
  const [seconds, setSeconds] = useState(getSecondsSinceLastHeartbeat(cluster));

  const { data: data } = useQuery({
    queryKey: [queryKeys.clusters, cluster.name],
    queryFn: () => fetchCluster(cluster.name),
    refetchInterval: 60000,
    refetchIntervalInBackground: true,
  });

  useEffect(() => {
    const timerId = setInterval(() => {
      setSeconds((prev) => prev + 1);
    }, 1000);
    return () => clearInterval(timerId);
  }, []);

  useEffect(() => {
    setSeconds(getSecondsSinceLastHeartbeat(data as Cluster));
  }, [data]);

  const getFormattedTimeLeft = () => {
    const minutes = Math.floor(seconds / 60);
    const displaySeconds = seconds % 60;

    return `${minutes}:${displaySeconds < 10 ? `0${displaySeconds}` : displaySeconds}`;
  };

  return (
    <div className="cluster-detail-countdown">
      <div className="u-no-margin u-no-padding p-heading--3">
        {getFormattedTimeLeft()}
      </div>
      <div>Since last heartbeat</div>
    </div>
  );
};
export default ClusterTimer;
