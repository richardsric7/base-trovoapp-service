import React from "react";
import chevronDown from "../utils/assets/images/chevron-down-icon.svg";
import chevronUp from "../utils/assets/images/chevron-up.svg";
import styled from "styled-components";
const Paginate = ({
  stepsPerPage,
  totalSteps,
  nextPage,
  showNextPage,
  showLastPage,
  previousPage,
}) => {
  const pageNumbers = [];

  for (let i = 1; i <= Math.ceil(totalSteps / stepsPerPage); i++) {
    pageNumbers.push(i);
  }

  return (
    <PaginationContainer >
      {showLastPage && (
        <div onClick={previousPage} className="page-number">
          <img src={chevronUp} alt="" />
        </div>
      )}
      {showNextPage && (
        <div onClick={nextPage} className="page-number">
          <img src={chevronDown} alt="" />
        </div>
      )}
    </PaginationContainer>
  );
};

export default Paginate;

const PaginationContainer = styled.div`
   display: flex;
   justify-content: center;

   div{
    cursor: pointer;
   }
`