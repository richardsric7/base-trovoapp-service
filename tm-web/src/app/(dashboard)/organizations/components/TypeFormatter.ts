import { useCallback } from "react";

export const usePrettyType = () => {
  const prettyType = useCallback((type: string) => {
    switch (type) {
      case "ASSET_MANAGER":
        return "Asset Manager";
      case "RATING_AGENCY":
        return "Rating Agency";
      // add more as needed...
      default:
        return type
          .replace(/_/g, " ") // The /g flag replaces ALL underscores
          .toLowerCase()
          .replace(/(^\w|\s\w)/g, (m) => m.toUpperCase());
    }
  }, []);

  return prettyType;
};
