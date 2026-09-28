import React from "react";
import styled from "styled-components";

interface SourceDestinationTagsProps {
  isSwap: boolean;
  sourceAssetCode?: string;
  sourceNetwork?: string;
  destinationAssetCode?: string;
  destinationNetwork?: string;
}

// Shows which asset/network a payment left with (source) and which it
// arrived as (destination). For a plain payment the two sides are
// identical, so this collapses into a single combined tag; a row where
// they differ is a swap, shown as two distinct From/To tags.
const SourceDestinationTags: React.FC<SourceDestinationTagsProps> = ({
  isSwap,
  sourceAssetCode,
  sourceNetwork,
  destinationAssetCode,
  destinationNetwork,
}) => {
  if (!isSwap) {
    return (
      <TagsRow>
        <Tag tone="neutral">
          {destinationAssetCode || "--"} · {destinationNetwork || "base"}
        </Tag>
      </TagsRow>
    );
  }

  return (
    <TagsRow>
      <Tag tone="warning">
        From {sourceAssetCode || "--"} · {sourceNetwork || "base"}
      </Tag>
      <Tag tone="success">
        To {destinationAssetCode || "--"} · {destinationNetwork || "base"}
      </Tag>
    </TagsRow>
  );
};

export default SourceDestinationTags;

const TagsRow = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
`;

const toneColors: Record<string, { border: string; bg: string; text: string }> = {
  neutral: { border: "#007CDF", bg: "#ACD1EF", text: "#007CDF" },
  warning: { border: "#BE9B00", bg: "#FCE9CF", text: "#8A5A1F" },
  success: { border: "#00A859", bg: "#DCF4E6", text: "#1E7A46" },
};

const Tag = styled.span<{ tone: "neutral" | "warning" | "success" }>`
  border: 1px solid ${({ tone }) => toneColors[tone].border};
  background-color: ${({ tone }) => toneColors[tone].bg};
  color: ${({ tone }) => toneColors[tone].text};
  font-size: 11px;
  font-weight: 500;
  line-height: 14px;
  border-radius: 16px;
  padding: 3px 10px;
  display: inline-flex;
  align-items: center;
  white-space: nowrap;
`;
