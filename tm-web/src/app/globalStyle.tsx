import { createGlobalStyle } from "styled-components";

const GlobalStyle = createGlobalStyle`
  @font-face {
    font-family: 'MatahariFont';
    src: url('../../public/fonts/Matahari-400Regular.ttf') format('ttf'),
         url('../../public/fonts/MatahariExtended-600ExtSemBd.ttf') format('ttf');
    font-weight: normal;
    font-style: normal;
  }
`;
