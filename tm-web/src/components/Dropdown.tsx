"use client";
import { MouseEvent, ReactNode, useEffect, useRef, useState } from "react";
import { createPortal } from "react-dom";
import { FaAngleDown } from "react-icons/fa6";
import styled, { css } from "styled-components";
export const DropdownSelect = ({
  value,
  options,
  children,
  placeholder,
  onSelect,
  leftIcon,
  errorMessage,
  labelText,
  className,
  isTopAddon,
  isBottomAddon,
  labelColor = "#828282", // Default label color
  placeholderColor = "#00225A", // Default placeholder color

  backgroundColor = "#F3F4F6", // Default background color
  borderColor = "#E0E0E0", // Default border color
  iconColor = "#828282",
  selectColor = "#00225A",
  style,
  borderless = false,
  ...rest
}: {
  value: string;
  options: string[];
  children?: ReactNode;
  placeholder?: string;
  onSelect?: (item: string, index: number) => void;
  leftIcon?: ReactNode;
  errorMessage?: string;
  labelText?: string;
  className?: ReactNode;
  isTopAddon?: boolean;
  isBottomAddon?: boolean;
  labelColor?: string; // Add label color prop
  placeholderColor?: string;
  borderless?: boolean;
  backgroundColor?: string;
  borderColor?: string;
  selectColor?: string;
  iconColor?: string;

  style?: React.CSSProperties;
}) => {
  const node = useRef<HTMLDivElement>(null);
  const portalRef = useRef<HTMLDivElement | null>(null);
  const [show, setShow] = useState(false);
  const [coords, setCoords] = useState<{
    top: number;
    left: number;
    width: number;
  } | null>(null);

  const updateCoords = () => {
    if (!node.current) return;
    const rect = node.current.getBoundingClientRect();
    setCoords({ top: rect.bottom, left: rect.left, width: rect.width });
  };

  const clickOutside = (e: any) => {
    const target = e.target as Node;
    if (node.current?.contains(target)) return;
    if (portalRef.current && portalRef.current.contains(target)) return;
    setShow(false);
  };

  const handleSelect = (
    e: MouseEvent<HTMLButtonElement, globalThis.MouseEvent>,
    item: string,
    index: number
  ) => {
    e.preventDefault();
    setShow(false);
    if (onSelect) onSelect(item, index);
  };
  useEffect(() => {
    document.addEventListener("mousedown", clickOutside);
    return () => {
      document.removeEventListener("mousedown", clickOutside);
    };
  }, [show]);

  useEffect(() => {
    if (show) {
      updateCoords();
      window.addEventListener("scroll", updateCoords, true);
      window.addEventListener("resize", updateCoords);
    }
    return () => {
      window.removeEventListener("scroll", updateCoords, true);
      window.removeEventListener("resize", updateCoords);
    };
  }, [show]);
  return (
    <Wrapper>
      {labelText && (
        <InputLabel style={{ color: labelColor }}>{labelText}</InputLabel>
      )}
      <Container
        ref={node}
        {...rest}
        error={!!errorMessage}
        className={"inputContainer"}
        onClick={() => setShow(!show)}
        borderless={borderless}
        backgroundColor={backgroundColor}
        borderColor={borderColor}
      >
        <DropdownTextWrapper>
          {leftIcon}
          <DropDownTxt
            id={"dropDownTxt"}
            // isPlaceholder={value === "" && !!placeholder}
            style={{
              color: value ? selectColor : placeholderColor,
            }}
          >
            {!value && placeholder ? placeholder : value}
          </DropDownTxt>
        </DropdownTextWrapper>
        <FaAngleDown color={iconColor} />
        {/* render inline content placeholder only for layout; actual menu is portalled */}
      </Container>
      {show &&
        coords &&
        createPortal(
          <PortalContent
            ref={portalRef}
            show={show}
            style={{ top: coords.top, left: coords.left, width: coords.width }}
            id="content-portal"
          >
            {isTopAddon}
            {children}
            {options?.map((item, index) => (
              <DropDownItem
                key={index.toString()}
                onClick={(e) => handleSelect(e, item, index)}
              >
                {item}
              </DropDownItem>
            ))}
            {isBottomAddon}
          </PortalContent>,
          document.body
        )}
      {errorMessage && <ErrorMessage>{errorMessage}</ErrorMessage>}
    </Wrapper>
  );
};
const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
`;
const ErrorMessage = styled.p`
  font-size: 10px;
  color: red;
  margin-top: 8px;
  margin-bottom: 0;
  line-height: 14px;
`;
const Container = styled.div.attrs<{
  error?: boolean;
  borderless?: boolean;
  borderColor?: string;
  backgroundColor?: string;
}>((props) => ({
  borderless: undefined,
  borderColor: undefined,
  backgroundColor: undefined,
  error: undefined,
}))<{
  error: boolean;
  borderless?: boolean;
  borderColor?: string;
  backgroundColor: string;
}>`
  border-radius: 10px;
  padding: 8px 12px;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
  background-color: ${({ backgroundColor }) =>
    backgroundColor || "transparent"};
  border: ${({ borderless, borderColor }) =>
    borderless ? "none" : `1px solid ${borderColor || "#E0E0E0"}`};

  ${({ error }) =>
    error &&
    css`
      border: 1px solid red;
    `}
`;

const DropdownTextWrapper = styled.div`
  justify-content: space-between;
  display: flex;
  align-items: center;
  margin-right: 7px;
`;
const DropDownTxt = styled.p`
  font-style: normal;
  margin: 0;
  font-weight: 500;
  font-size: 14px;
  line-height: 24px;
  color: #828282;
  flex: 1;
`;

export const DropdownContent = styled.div<{ show: boolean }>`
  display: ${({ show }) => (show ? "block" : "none")};
  position: absolute;
  top: 40px;
  right: 0;
  background-color: #ffffff;
  overflow: hidden;
  /* min-width: 160px; */
  width: 100%;
  overflow-y: auto;
  max-height: 250px;
  box-shadow: 0px 0px 1px rgba(12, 26, 75, 0.2),
    0px 1px 3px rgba(50, 50, 71, 0.1);
  border-radius: 8px;
  z-index: 10;
  & button {
    color: #232735;
    font-size: 14px;
    line-height: 24px;
    padding: 8px 10px;
    text-decoration: none;
    display: block;
    cursor: pointer;
  }
  & button:hover {
    background-color: #f1f1f1;
  }
`;
export const DropDownItem = styled.button`
  background-color: transparent;
  border: none;
  text-align: left;
  width: 100%;
  cursor: pointer;
  font-family: inherit;
`;

const InputLabel = styled.label`
  font-weight: 500;
  font-size: 14px;
  line-height: 24px;
  margin-bottom: 8px;
  color: #00225a;

  @media (max-width: 768px) {
    margin-bottom: 8px;
  }
`;

const PortalContent = styled(DropdownContent)`
  position: fixed;
  z-index: 10000;
`;
