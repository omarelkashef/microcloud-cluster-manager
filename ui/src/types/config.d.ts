export type ConfigPair = Record<string, string | undefined>;

interface ManagerOptions {
  config: {
    [key: string]: string;
  };
}

interface MemberOptions {
  target?: string;
  config: {
    https_address: string;
    external_address?: string;
  };
}

export type ConfigField = ConfigOption & {
  category: string;
  key: string;
};

export interface ConfigOption {
  longdesc?: string;
  scope?: "global" | "local";
  shortdesc?: string;
  type: "bool" | "string" | "integer";
}

export interface ConfigOptionCategories {
  [category: string]: {
    keys: {
      [key: string]: ConfigOption;
    }[];
  };
}

export interface ConfigOptions {
  configs: {
    cluster: ConfigOptionCategories;
    member: ConfigOptionCategories;
  };
}

export type ConfigOptionsKeys = keyof ConfigOptions["configs"];
