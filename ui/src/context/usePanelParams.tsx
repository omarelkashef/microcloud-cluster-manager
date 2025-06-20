import { useSearchParams } from "react-router-dom";

export interface PanelHelper {
  panel: string | null;
  cluster: string | null;
  clear: () => void;
  openEnrolCluster: () => void;
  openConfigureCluster: (cluster?: string) => void;
}

export const panels = {
  enrolCluster: "enrol-cluster",
  configureCluster: "configure-cluster",
};

type ParamMap = Record<string, string>;

const usePanelParams = (): PanelHelper => {
  const [params, setParams] = useSearchParams();

  const craftResizeEvent = () => {
    setTimeout(() => window.dispatchEvent(new Event("resize")), 100);
  };

  const setPanelParams = (panel: string, args: ParamMap = {}) => {
    const newParams = new URLSearchParams();
    newParams.set("panel", panel);
    for (const [key, value] of Object.entries(args)) {
      if (value) {
        newParams.set(key, value);
      }
    }
    setParams(newParams);
    craftResizeEvent();
  };

  const clearParams = () => {
    const newParams = new URLSearchParams(params);
    // we only want to remove search params set when opening the panel
    // pre-existing search params should be kept e.g. params from the search bar
    newParams.delete("panel");
    newParams.delete("cluster");
    setParams(newParams);
    craftResizeEvent();
  };

  return {
    panel: params.get("panel"),
    cluster: params.get("cluster"),

    clear: () => {
      clearParams();
    },

    openEnrolCluster: () => {
      setPanelParams(panels.enrolCluster);
    },

    openConfigureCluster: (cluster?: string) => {
      const params = { cluster: cluster || "" };
      if (cluster) {
        params.cluster = cluster;
      }
      setPanelParams(panels.configureCluster, params);
    },
  };
};

export default usePanelParams;
