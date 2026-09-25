export const getInitials = (name?: string | null): string => {
  return name?.trim().charAt(0).toUpperCase() || "?";
};
