export const ParseStringInputToNumber = (s: string | number): number | null => {
  if (typeof s === "number") {
    return s;
  }
  const parsed = parseInt(s);
  if (!isNaN(parsed)) return parsed;
  return null;
};
