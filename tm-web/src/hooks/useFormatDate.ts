import { useCallback } from "react";

const useFormatDate = () => {
  return useCallback((dateString?: string) => {
    if (!dateString) return "N/A";

    const date = new Date(dateString);
    return date.toLocaleDateString("en-us", {
      day: "numeric",
      month: "short",
      year: "numeric",
    });
  }, []);
};

export default useFormatDate;
