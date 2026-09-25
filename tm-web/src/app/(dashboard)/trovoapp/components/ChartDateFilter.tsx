"use client";
import React, { useState, useMemo } from "react";
import { Button, Card, Select, DatePicker, Menu } from "antd";
import styled from "styled-components";
import dayjs, { Dayjs } from "dayjs";
import { FaAngleDown } from "react-icons/fa6";
import SecondaryButton from "@/components/SecondaryButton";
import PrimaryButton from "@/components/PrimaryButton";

const { RangePicker } = DatePicker;
const { Option } = Select;

type DateRangeOption =
  | "Today"
  | "Yesterday"
  | "Last week"
  | "Last month"
  | "Last 3 months"
  | "Last year"
  | "All time"
  | "Custom";

type DisplayUnit = "hours" | "days" | "weeks" | "months" | "years";

interface ChartDateFilterProps {
  onApply: (
    range: DateRangeOption,
    unit: DisplayUnit,
    customDates?: [Dayjs, Dayjs]
  ) => void;
}

const ChartDateFilter: React.FC<ChartDateFilterProps> = ({ onApply }) => {
  const [panelVisible, setPanelVisible] = useState(false);
  const [selectedRange, setSelectedRange] = useState<DateRangeOption>("Today");
  const [displayUnit, setDisplayUnit] = useState<DisplayUnit>("hours");
  const [customDates, setCustomDates] = useState<[Dayjs, Dayjs] | null>(null);

  const availableDisplayUnits = useMemo(() => {
    switch (selectedRange) {
      case "Today":
      case "Yesterday":
        return ["hours"];
      case "Last week":
        return ["hours", "days"];
      case "Last month":
        return ["hours", "days", "weeks"];
      case "Last 3 months":
      case "Last year":
        return ["days", "weeks", "months"];
      case "All time":
        return ["days", "weeks", "months", "years"];
      case "Custom":
        return ["hours", "days", "weeks", "months", "years"];
      default:
        return ["days"];
    }
  }, [selectedRange]);

  const handleApply = () => {
    onApply(selectedRange, displayUnit, customDates || undefined);
    setPanelVisible(false);
  };

  const handleCancel = () => {
    setPanelVisible(false);
  };

  return (
    <Wrapper>
      <Button type="default" onClick={() => setPanelVisible((prev) => !prev)}>
        {selectedRange} <FaAngleDown color="#00225A" />
      </Button>

      {panelVisible && (
        <Panel>
          <TopRow>
            <StyledMenu
              mode="vertical"
              selectedKeys={[selectedRange]}
              onClick={(e) => {
                const key = e.key as DateRangeOption;
                setSelectedRange(key);
                setDisplayUnit(availableDisplayUnits[0] as DisplayUnit); // Reset unit on change
              }}
              items={[
                "Today",
                "Yesterday",
                "Last week",
                "Last month",
                "Last 3 months",
                "Last year",
                "All time",
                "Custom",
              ].map((label) => ({
                key: label,
                label,
              }))}
              style={{ width: 140 }}
            />
            <RightRow>
              <UnitGroup>
                <Text>Display chart in:</Text>
                <Divider />

                <UnitButtonContainer>
                  {availableDisplayUnits.map((unit) => (
                    <Button key={unit} value={unit}>
                      {unit}
                    </Button>
                  ))}
                </UnitButtonContainer>
              </UnitGroup>

              <PanelContent>
                {selectedRange === "Custom" && (
                  <RangePicker
                    style={{ width: "100%", marginBottom: 24 }}
                    onChange={(dates) => {
                      if (dates && dates[0] && dates[1]) {
                        setCustomDates([dates[0], dates[1]]);
                      } else {
                        setCustomDates(null);
                      }
                    }}
                  />
                )}
              </PanelContent>

              <Content>
                <Divider />
                <ButtonContainer>
                  <SecondaryButton
                    onClick={handleCancel}
                    buttonStyle={{ width: "146px" }}
                  >
                    Cancel
                  </SecondaryButton>
                  <PrimaryButton
                    onClick={handleApply}
                    buttonStyle={{ width: "146px" }}
                  >
                    Apply
                  </PrimaryButton>
                </ButtonContainer>
              </Content>
            </RightRow>
          </TopRow>
        </Panel>
      )}
    </Wrapper>
  );
};

export default ChartDateFilter;

const Wrapper = styled.div`
  position: relative;
  display: inline-block;
`;

const TopRow = styled.div`
  display: flex;

  gap: 12px;
`;

const Text = styled.p`
  font-weight: 600;
  font-style: SemiBold;
  font-size: 16px;
  line-height: 24px;
  color: #191919;
`;

const UnitGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Panel = styled(Card)`
  position: absolute;
  top: 45px;
  left: 0;
  width: 520px;
  display: flex;
  z-index: 1000;
  padding: 0;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
`;

const PanelContent = styled.div`
  flex: 1;
  padding-top: 10px;
`;

const UnitButtonContainer = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  align-items: center;
  gap: 10px;
`;
const RightRow = styled.div``;

const ButtonContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 10px;
`;

const Divider = styled.div`
  background-color: #e0e0e0;
  height: 1px;
  width: 100%;
`;

const Content = styled.div`
  margin-top: 200px;
`;

const StyledMenu = styled(Menu)`
  font-weight: 500;
  font-size: 16px;
  line-height: 24px;
  color: #191919;

  .ant-menu-item {
    font-family: inherit;
    font-size: inherit;
    font-size: 14px;
    color: inherit;
  }

  .ant-menu-item-selected {
    background-color: #f2f6f9 !important;
    font-weight: 600;
    border: 1px solid #62ace8;
    border-radius: 8px;
  }
`;
