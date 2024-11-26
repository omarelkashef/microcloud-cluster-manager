import { ConfigField, ConfigOptions, ConfigOptionsKeys } from "types/config";

export const configDescriptionToHtml = (input: string): string => {
  // special characters
  let result = input
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll("\n", "<br>");

  // replace admonition markup
  result = result
    .replaceAll("```", "")
    .replaceAll("{important}", "<b>Important</b>");

  // code blocks
  let count = 0;
  const maxCodeblockReplacementCount = 100; // avoid infinite loop
  while (result.includes("`") && count++ < maxCodeblockReplacementCount) {
    result = result.replace("`", "<code>").replace("`", "</code>");
  }

  return result;
};

export const getConfigMetadata = (
  configOptions?: ConfigOptions,
): Record<string, ConfigField> => {
  if (!configOptions) {
    return {};
  }

  const configMetadata: Record<string, ConfigField> = {};

  for (const entity in configOptions.configs) {
    const configEntityCategories =
      configOptions.configs[entity as ConfigOptionsKeys];

    for (const category in configEntityCategories) {
      const categoryKeys = configEntityCategories[category].keys;

      for (const config of categoryKeys) {
        const [key, value] = Object.entries(config)[0];
        configMetadata[key] = {
          ...value,
          category: category,
          key,
        };
      }
    }
  }

  return configMetadata;
};
