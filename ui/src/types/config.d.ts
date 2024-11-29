export interface ConfigData {
  value: string;
  description: string;
  title: string;
}

export interface Configuration {
  [key: string]: ConfigData;
}
