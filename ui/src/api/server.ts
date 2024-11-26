import { LxdApiResponse } from "types/apiResponse";
import { ConfigOptions } from "types/config";
import { Server } from "types/server";
import { handleResponse } from "util/helpers";

export const fetchServer = (): Promise<Server> => {
  return new Promise((resolve, reject) => {
    fetch("/1.0")
      .then(handleResponse)
      .then((data) => resolve((data as LxdApiResponse<Server>).metadata))
      .catch(reject);
  });
};

export const fetchConfigOptions = (): Promise<ConfigOptions | null> => {
  return new Promise((resolve, reject) => {
    fetch("/1.0/metadata/configuration")
      .then(handleResponse)
      .then((data) => resolve((data as LxdApiResponse<ConfigOptions>).metadata))
      .catch(reject);
  });
};
