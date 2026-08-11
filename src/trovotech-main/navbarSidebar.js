import { useNavigate } from "react-router-dom";
import { styled } from "styled-components";
// import { theme } from "../utils/theme";

export const NavSideDrawer = ({ closeDrawer }) => {
  const navigate = useNavigate();

  const navbarLinks = [
    {
      name: "Home",
      link: "/home",
    },
    {
      name: "Features ",
      link: "/features",
    },
    {
      name: "Download",
      link: "/download",
    },
    {
      name: "Help Center",
      link: "/help-center/account-deletion",
    },
    {
      name: "Contact Us",
      link: "/contact",
    },
  ];

  const handleNavigation = (path) => {
    navigate(path);
  };

  return (
    <>
      <DrawerContainer onClick={closeDrawer}></DrawerContainer>
      <DrawerSection onClick={(e) => e.stopPropagation()}>
        <NavLinks>
          <CancelIconContainer onClick={closeDrawer}>
            <CancelIcon>x</CancelIcon>
          </CancelIconContainer>
          {navbarLinks.map((item, index) => {
            return (
              <LinkComponent
                option={item}
                key={index}
                onClick={() => handleNavigation(item.link)}
              />
            );
          })}
        </NavLinks>
      </DrawerSection>
    </>
  );
};
const LinkComponent = ({ option, onClick }) => {
  return (
    <LinkDetails>
      <LinkText onClick={onClick}>{option.name}</LinkText>
    </LinkDetails>
  );
};
const DrawerContainer = styled.section`
  position: fixed;
  z-index: 2;
  top: 0;
  bottom: 0;
  right: 0;
  left: 0;
  display: flex;
  justify-content: flex-end;
  align-items: flex-start;
`;
const DrawerSection = styled.div`
  position: fixed;
  top: 0;
  bottom: 0;
  right: 0;
  z-index: 5;
  background: #ffffff;
  box-shadow: 0px 10px 24px rgba(0, 0, 0, 0.1);
  border-radius: 10px;
  width: 180px;
  height: 964px;
  padding: 29px 26px;
`;
const NavLinks = styled.div`
  display: flex;
  flex-direction: column;
  row-gap: 32px;

  @media (max-width: 1100px) {
    column-gap: 25px;
  }
`;

const LinkDetails = styled.div``;

const LinkText = styled.span`
  font-size: 16px;
  font-family: "Montserrat";
  font-style: normal;
  font-weight: 400;
  line-height: 28px;
  color: #1b1d21;
  text-decoration: none;
  cursor: pointer;
  &.active {
    font-weight: 600;
  }
  &:hover {
    border: none;
  }
`;
const CancelIconContainer = styled.div`
  display: flex;
  justify-content: flex-end;
`;
const CancelIcon = styled.div`
  cursor: pointer;
`;
