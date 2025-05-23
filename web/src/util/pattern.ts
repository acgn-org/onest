export const ParseTextWithPattern = (
  text: string,
  regexp: RegExp,
  pattern: string,
) => {
  const matches = regexp.exec(text);
  const matchesObj = Object.fromEntries(
    (matches ?? []).map((v, i) => [String(i), v]),
  );
  return pattern.replace(/\$(\w+)|\$\{([^}]+)\}/g, (_, key1, key2) => {
    const key = key1 || key2;
    return matchesObj[key] ?? "";
  });
};
