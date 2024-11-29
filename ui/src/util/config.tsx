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
