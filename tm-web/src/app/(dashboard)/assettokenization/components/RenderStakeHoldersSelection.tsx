import styled from "styled-components";
import StakeHoldersInput from "./StakeHoldersInput";

interface RenderStakeHoldersSelectionProps {
  title: string;
  description: string;
  placeholder: string;
  value: string;
  onChange: (value: string) => void;
  suggestions: string[];
  onSelect: (name: string) => void;
  selectedItems: string[];
  onRemove: () => void;
  hasOrganizations: boolean;
}

const RenderStakeHoldersSelection: React.FC<
  RenderStakeHoldersSelectionProps
> = ({
  title,
  description,
  placeholder,
  value,
  onChange,
  suggestions,
  onSelect,
  selectedItems,
  onRemove,
  hasOrganizations,
}) => {
  return (
    <Content>
      <Title>{title}</Title>
      <Text>{description}</Text>
      <Line />
      {hasOrganizations ? (
        <StakeHoldersInput
          placeholder={placeholder}
          value={value}
          onChange={onChange}
          suggestions={suggestions}
          onSelectSuggestion={onSelect}
          selectedItems={selectedItems}
          onRemoveItem={onRemove}
        />
      ) : (
        <EmptyState>
          No {title.toLowerCase()} organizations available yet. Please create
          one first.
        </EmptyState>
      )}
    </Content>
  );
};

export default RenderStakeHoldersSelection;

const Content = styled.div`
  display: flex;
  flex-direction: column;
  background-color: #f2f6f9;
  border-radius: 12px;
  padding: 16px;
  gap: 4px;
`;

const Title = styled.h2`
  font-weight: 600;
  font-size: 14px;
  line-height: 16px;
  letter-spacing: 0.1px;
  color: #191919;
`;

const Text = styled.p`
  font-weight: 400;
  font-size: 11px;
  line-height: 16px;
  letter-spacing: 0.1px;
  color: #828282;
`;

const Line = styled.div`
  background-color: #e0e0e0;
  width: 100%;
  height: 1px;
  margin-bottom: 8px;
`;

const EmptyState = styled.div`
  padding: 16px;
  background-color: #fff3cd;
  border: 1px solid #ffc107;
  border-radius: 8px;
  color: #856404;
  font-size: 12px;
  line-height: 18px;
  text-align: center;
`;
