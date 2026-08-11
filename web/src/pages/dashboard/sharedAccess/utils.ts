export const parsePermissionInput = (value: string) =>
  value
    .split(',')
    .map((entry) => entry.trim())
    .filter(Boolean);
