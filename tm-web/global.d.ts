declare module "*.png";
declare module "*.jpg";
declare module "*.gif";
declare module "*.pdf";
declare module "react-native-otp-textinput";
declare module "react-native-freshchat-sdk";
declare module "*.svg" {
  import React from "react";
  import { SvgProps } from "react-native-svg";
  const content: React.FC<SvgProps>;
  export default content;
}
