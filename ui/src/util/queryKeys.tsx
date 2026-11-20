export const queryKeys = {
  clusters: "clusters",
  links: "links",
  tokens: "tokens",
  server: "server",
  configuration: "configuration",
  cluster: "cluster",
  identities: "identities",
};

export const allClusterLinksQueryKey = [
  queryKeys.clusters,
  { resource: queryKeys.links },
  "all",
] as const;
