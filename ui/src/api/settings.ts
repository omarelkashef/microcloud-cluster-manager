import { LxdApiResponse } from "types/apiResponse";
import { ConfigPair, ManagerOptions, MemberOptions } from "types/config";
import { handleResponse } from "util/helpers";

export const fetchManagerConfigOptions = (): Promise<ManagerOptions> => {
  return new Promise((resolve, reject) => {
    fetch("/1.0/config")
      .then(handleResponse)
      .then((data) =>
        resolve((data as LxdApiResponse<ManagerOptions>).metadata),
      )
      .catch(reject);
  });
};

export const updateManagerConfigs = (config: ConfigPair): Promise<void> => {
  return new Promise((resolve, reject) => {
    fetch("/1.0/config", {
      method: "PATCH",
      body: JSON.stringify({
        config,
      }),
    })
      .then(handleResponse)
      .then(() => resolve())
      .catch(reject);
  });
};

export const fetchMemberConfigOptions = (): Promise<MemberOptions[]> => {
  return new Promise((resolve, reject) => {
    fetch("/1.0/member/config")
      .then(handleResponse)
      .then((data) =>
        resolve((data as LxdApiResponse<MemberOptions[]>).metadata),
      )
      .catch(reject);
  });
};

export const updateMemberConfigs = (
  member: string,
  config: ConfigPair,
): Promise<void> => {
  return new Promise((resolve, reject) => {
    fetch(`/1.0/member/${member}/config`, {
      method: "PATCH",
      body: JSON.stringify({
        config,
      }),
    })
      .then(handleResponse)
      .then(() => resolve())
      .catch(reject);
  });
};
