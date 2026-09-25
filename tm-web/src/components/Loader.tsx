import React from "react";
import styled from "styled-components";

const Loader = () => {
  return (
    <div>
      <LoadingText>
        <LoadingSpinner /> loading...
      </LoadingText>
    </div>
  );
};

export default Loader;

const LoadingText = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  color: #007cdf;
  font-weight: 500;
  margin-top: 10px;
`;

const LoadingSpinner = styled.div`
  width: 20px;
  height: 20px;
  border: 3px solid #007cdf;
  border-top: 3px solid transparent;
  border-radius: 50%;
  animation: spin 1s linear infinite;

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }
`;
