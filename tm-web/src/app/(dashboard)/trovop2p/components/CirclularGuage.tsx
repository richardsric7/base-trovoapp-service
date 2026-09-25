import React from "react";
import styled from "styled-components";

interface GaugeProps {
  value: number;
  size?: number;
  thickness?: number;
  primaryColor?: string;
  backgroundColor?: string;
  showValue?: boolean;
}

const CircularGauge = ({
  value,
  size = 200,
  thickness = 12,
  primaryColor = "#0066FF",
  backgroundColor = "#E6E6E6",
  showValue = true,
}: GaugeProps) => {
  // Ensure value is between 0 and 100
  const normalizedValue = Math.min(Math.max(value, 0), 100);

  // Calculating the SVG parameters
  const radius = (size - thickness) / 2;
  const center = size / 2;

  // Calculating the arc parameters
  const startAngle = -90; // Start from -90 degrees
  const endAngle = 90; // End at 90 degrees (270 degree arc)
  const angleRange = endAngle - startAngle;
  const progressAngle = startAngle + (angleRange * normalizedValue) / 100;

  // Creating the the arc paths
  const createArc = (angle: number) => {
    const radians = ((angle - 90) * Math.PI) / 180;
    const x = center + radius * Math.cos(radians);
    const y = center + radius * Math.sin(radians);
    const largeArc = Math.abs(angle - startAngle) > 180 ? 1 : 0;

    return [
      `M ${center + radius * Math.cos(((startAngle - 90) * Math.PI) / 180)}`,
      `${center + radius * Math.sin(((startAngle - 90) * Math.PI) / 180)}`,
      `A ${radius} ${radius} 0 ${largeArc} 1 ${x} ${y}`,
    ].join(" ");
  };

  const backgroundPath = createArc(endAngle);
  const progressPath = createArc(progressAngle);

  // Calculating the dot position
  const dotAngle = ((progressAngle - 90) * Math.PI) / 180;
  const dotX = center + radius * Math.cos(dotAngle);
  const dotY = center + radius * Math.sin(dotAngle);

  return (
    <GaugeContainer>
      <GaugeSVG $size={size} viewBox={`0 0 ${size} ${size}`}>
        {/* Background arc */}
        <path
          d={backgroundPath}
          fill="none"
          stroke={backgroundColor}
          strokeWidth={thickness}
          strokeLinecap="round"
        />

        {/* Progress arc */}
        <path
          d={progressPath}
          fill="none"
          stroke={primaryColor}
          strokeWidth={thickness}
          strokeLinecap="round"
          style={{ transition: "all 0.5s ease-out" }}
        />
        {/* Dot indicator with white border */}

        <circle cx={dotX} cy={dotY} r={thickness / 1 + 2} fill="white" />
        <circle cx={dotX} cy={dotY} r={thickness / 1.3} fill={primaryColor} />
      </GaugeSVG>

      {showValue && <ValueText>{normalizedValue}%</ValueText>}
    </GaugeContainer>
  );
};

export default CircularGauge;

const GaugeContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
`;

const GaugeSVG = styled.svg<{ $size: number }>`
  width: ${(props) => props.$size}px;
  height: ${(props) => props.$size}px;
`;

const ValueText = styled.span`
  font-size: 48px;
  font-weight: 700;
  line-height: 84px;
  text-align: center;
  color: #00225a;
  margin-top: -250px;
`;
