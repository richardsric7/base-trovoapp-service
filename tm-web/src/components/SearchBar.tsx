import styled from "styled-components";
// import BasicButton from "./button";
import { FC, ReactElement } from "react";
import { FiSearch } from "react-icons/fi";
import { FaSearch } from "react-icons/fa";

type ISearch = {
  placeholder?: string;
  rightItem?: ReactElement;
  labelText?: string;
  searchText?: string;
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
  value?: string;
  handleSearch?: () => void;
  className?: any;
  type?: any;
  customWidth?: string;
};
export const SearchBar: FC<ISearch & { customWidth?: string }> = ({
  placeholder,
  rightItem,
  labelText,
  searchText,
  onChange,
  value,
  handleSearch,
  className,
  type,
  customWidth,
}) => {
  return (
    <Wrapper className={className} customWidth={customWidth}>
      {labelText && <InputLabel>{labelText}</InputLabel>}
      <SearchContainer>
        <SearchInput
          value={value}
          placeholder="Start search..."
          onChange={onChange}
          type={type}
        />

        <SearchButton onClick={handleSearch}>
          <FaSearch />
        </SearchButton>
      </SearchContainer>
    </Wrapper>
  );
};
const Wrapper = styled.div<{ customWidth?: string }>`
  display: flex;
  flex-direction: column;
  width: ${({ customWidth }) => customWidth || "320px"};
`;
const InputLabel = styled.label`
  font-family: inherit;
  font-style: normal;
  font-weight: 400;
  font-size: 16px;
  line-height: 22px;
  margin-bottom: 8px;
  color: #3f3f3f;

  @media (max-width: 768px) {
    margin-bottom: 8px;
  }
`;
const SearchContainer = styled.div`
  display: flex;
  height: 38px;
  align-items: center;
  justify-content: space-between;
  background-color: #f2f6f9;
  border: none;
  width: 100%;
  border-radius: 8px;
  padding: 8px 4px 8px 12px;
  box-sizing: border-box;
`;
const InputWrapper = styled.div`
  // display: flex;
  // align-items: center;
  // border-radius: 16px;
  // padding: 8px 12px;
  @media (max-width: 768px) {
    // width: 100%;
  }
`;

const SearchButton = styled.button`
  width: 32px;
  height: 32px;
  padding: 8px;
  gap: 10px;
  border-radius: 8px;
  border: 0;
  font-family: inherit;
  background-color: #007cdf;
  color: #fff;
  // flex: 1;
  cursor: pointer;
  outline: none;
`;
const SearchInput = styled.input`
  background-color: transparent;
  border: none;
  outline: none;
  color: #00225a;
  font-family: inherit;

  ::placeholder {
    color: #00225a;
    font-weight: 400;
    font-size: 14px;
    line-height: 17px;
    font-family: inherit;
  }
`;
