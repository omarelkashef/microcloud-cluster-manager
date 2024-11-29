import { LxdApiResponse } from "types/apiResponse";
import { Configuration } from "types/config";
import { handleResponse } from "util/helpers";

export const fetchConfigurations = (): Promise<Configuration> => {
  return new Promise((resolve, reject) => {
    fetch("/1.0/configuration")
      .then(handleResponse)
      .then((data) => resolve((data as LxdApiResponse<Configuration>).metadata))
      .catch(reject);
  });
};
