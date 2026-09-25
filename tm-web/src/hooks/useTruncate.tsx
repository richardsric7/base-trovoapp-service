import React from "react";
import { useMemo } from "react";

const useTruncate = (text: string, maxLength: number) => {
  const truncatedText = useMemo(() => {
    if (text.length > maxLength) {
      return `${text.slice(0, maxLength)}...`;
    }
    return text;
  }, [text, maxLength]);

  return truncatedText;
};

interface TruncatedTextProps {
  text: string;
  maxLength: number;
}

const TruncatedText: React.FC<TruncatedTextProps> = ({ text, maxLength }) => {
  const truncated = useTruncate(text, maxLength);
  return <span>{truncated}</span>;
};

export default TruncatedText;
