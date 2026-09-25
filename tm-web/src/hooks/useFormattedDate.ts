import { useMemo } from "react";

const useFormattedDate = (dateString: string | undefined) => {
  const formattedDate = useMemo(() => {
    if (!dateString) return "N/A";

    try {
      const date = new Date(dateString);
      if (!isNaN(date.getTime())) {
        return date.toLocaleDateString("en-US", {
          year: "numeric",
          month: "short",
          day: "numeric",
          hour: "2-digit",
          minute: "2-digit",
        });
      }
      return "N/A";
    } catch (error) {
      return "N/A";
    }
  }, [dateString]);

  return formattedDate;
};

export default useFormattedDate;
