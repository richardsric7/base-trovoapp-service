import React, { useState } from "react";
import { styled } from "styled-components";
import trovoLogo from "../utils/assets/images/trovo-logo.svg";
import hamburgerIcon from "../utils/assets/images/hamburgerIcon.svg";
import { theme } from "../utils/theme";
// import { NavSideDrawer } from "./navbarSidebar";

export const Navbar = () => {
  const [viewMobileNavbar, setViewMobileNavbar] = useState(false)
  const navbarLinks = [
    {
      name: "Home",
      link: "#home",
    },
    {
      name: "About Us ",
      link: '#features',
    },
    {
      name: "Products",
      link: "/#download",
    },
    {
      name: "Services",
      link: "/",
    },
    {
      name: "Documentation",
      link: "/",
    },
    {
      name: "Blog",
      link: "/",
    },
    {
      name: "Contact Us",
      link: "/#contact",
    },
  ];

  return (
    <NavLayout>
      <NavContainer>
        <LogoContainer>
          <LogoImage src={trovoLogo} alt="trovotech logo" />
        </LogoContainer>
        <NavItems>
          <NavLinks>
            {navbarLinks.map((item, index) => {
              return <LinkComponent option={item} key={index} />;
            })}
          </NavLinks>
          {/* <NavButtons href="#download">Download App</NavButtons> */}
        </NavItems>
        <HamburgerContainer onClick={()=> setViewMobileNavbar(true)}>
          <HamburgerIcon src={hamburgerIcon} />
        </HamburgerContainer>
      </NavContainer>
      {/* {viewMobileNavbar && <NavSideDrawer closeDrawer={()=>setViewMobileNavbar(false)} />} */}
    </NavLayout>
  );
};

const LinkComponent = ({ option }) => {
  return (
    <LinkDetails>
    
      <LinkText href={option.link}>{option.name}</LinkText>
    </LinkDetails>
  );
};

const LinkDetails = styled.div``;

const LinkText = styled.a`
  font-size: 16px;
  font-family: "Montserrat";
  font-style: normal;
  font-weight: 400;
  line-height: 28px;
  color: #1b1d21;
  text-decoration: none;
  &.active {
    font-weight: 600;
  }
  &:hover {
    border: none;
  }
`;

const NavLayout = styled.nav`
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 16px 80px;
  margin: auto;
  position: fixed;
  top: 0;
  right: 0;
  left: 0;
  background-color: ${theme.colors.white};
  z-index: 5;
  
  @media (min-width: 1450px) {
    max-width: 1440px;
    margin: auto
  }
  @media (max-width: 650px) {
    padding: 16px 30px;
  }
`;
const NavContainer = styled.div`
  display: flex;
  justify-content: space-between;
  width: 100%;
  align-items: center;
`;
const LogoContainer = styled.div`
  align-self: center;
`;
const LogoImage = styled.img`

@media (max-width: 650px) {
    width: 135px;
  }
`;
const NavItems = styled.div`
  display: flex;
  column-gap: 64px;
  align-items: center;

@media (max-width: 1100px) {
  column-gap: 30px;
}
@media (max-width: 1000px) {
    display: none;
  }
`;
const NavLinks = styled.div`
  display: flex;
  column-gap: 32px;

  @media (max-width: 1100px) {
    column-gap: 25px;
  }
`;
const NavButtons = styled.a`
  border-radius: 10px;
  border: 1px solid ${theme.colors.trovoBlue};
  background: ${theme.colors.trovoBlue};
  color: ${theme.colors.white};
  font-size: 20px;
  font-family: "Montserrat", sans-serif;
  font-weight: 600;
  line-height: 26px;
  padding: 12px 10px;
  width: 185px;
  display: block;
  text-align: center;
  text-decoration: none;
`;
const HamburgerContainer = styled.div`
  display: none;

  @media (max-width: 1000px) {
    display: block;
  }
`
const HamburgerIcon = styled.img`
  
`