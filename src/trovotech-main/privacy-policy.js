import React from "react";
import { Navbar } from "./navbar";
import { Footer } from "./footer";
import styled from "styled-components";

export const PrivacyPolicy = () => {
  return (
    <main>
      <Navbar />
      <StyledMain>
      <section class="page-header page-header-classic page-header-md">
        <div class="container">
          <div class="row">
            <div class="col-md-8 order-2 order-md-1 align-self-center p-static">
              <span
                class="page-header-title-border visible"
                // style="width: 103.039px;"
              ></span>
              <h1 data-title-border="">Privacy Policy</h1>
            </div>
          </div>
        </div>
      </section>
      <div role="main" class="main">
      <div class="container">
        <p 
        // style="text-align: center;"
        ><strong>TROVO APP PRIVACY POLICY</strong></p>
        <p>Please read this privacy policy (the &ldquo;<strong>Policy</strong>&rdquo;) carefully to understand how we
          use your personal information. By accessing or using the Trovo App Application (&ldquo;<strong>Trovo
            App</strong>&rdquo;), you shall be deemed to have accepted and agreed to the terms hereafter set forth.
          If you do not agree with this Policy, please do not use Trovo App.</p>
        <p>Trovotech Limited &nbsp;(&ldquo;<strong>Trovotech</strong>&rdquo;) a private company limited by shares, is a
          blockchain technology company that is hinged on the three cardinal points of blockchain solutions;
          scalability, speed, and user experience, to lead innovation in the Web3 space by building practical Blockchain
          solutions that solve real world problems. At Trovotech, we take privacy issues seriously at.</p>
        <p>This Privacy Policy informs you of the ways we ensure privacy and the confidentiality of Personal Data. We
          are compliant with applicable privacy laws in the countries in which we operate. This policy describes the
          information we gather, how we may use those Personal Data and the circumstances under which we may disclose
          such information to third-parties. It is your own responsibility to read the Privacy Policy carefully before
          you start to use Trovo App.</p>
        <p>This Policy may change from time to time. If there are any material changes to how your personal information
          is used, we will notify you of any changes. Your continued use of the App after we make changes is deemed to
          be acceptance of those changes, so please check the Policy periodically for updates.</p>
        <p><strong>1. &nbsp;Policy Orientation</strong></p>
        <ul>
          <li>
            <p>Your privacy is a right. We devote time and energy to ensure respect and protection of your rights.</p>
          </li>
          <li>
            <p>The protection we offer is technologically neutral and does not depend on any technological techniques
              used. The protection applies to the processing of data by automated means as well as manual processing</p>
          </li>
          <li>
            <p>Trovo App is a non-custodial wallet, your funds and private keys are your entire responsibility.</p>
          </li>
        </ul>
        <p><strong>2. &nbsp;Registration / Login</strong></p>
        <p>Trovo App requires you to register and log in.</p>
        <p><strong>3. &nbsp;Funds and Private Keys</strong></p>
        <p>Your funds and private keys will stay on your device.</p>
        <p><strong>4. Image</strong></p>
        <p>Your camera will only be used for reading QR codes. Camera images will never leave your device.</p>
        <p><strong>5. &nbsp;Information from Device used</strong></p>
        <p>We may collect information about the mobile device you use to access our mobile application, including the
          hardware model, operating system and version, unique device identifiers and mobile network information. This
          information will never be communicated to third-parties unless you provide prior specific consent.</p>
        {/* <p></p> */}
       <p><strong>6. &nbsp;Other information </strong></p>
        <p>Your email address, virtual currency addresses, mobile phone number and any other information you choose to
          provide will never be communicated to third-parties unless you provide prior specific consent.</p>
        <p><strong>7. &nbsp;Blockchain Transactions</strong></p>
        <p>Your blockchain transactions may be relayed through servers (&ldquo;nodes&rdquo;) and will be publicly
          visible due to the public nature of distributed ledger systems.</p>
        <p><strong>8. &nbsp;Secure Communication with our Servers</strong></p>
        <p>All of our servers support HTTPS.</p>
        <p><strong>9. &nbsp;Third-Party Servers Communication</strong></p>
        <p>We cannot guarantee the privacy of your Internet connection. Exchange rates, balances, transactions and other
          blockchain information may be read from or relayed to, third-party servers.&nbsp;</p>
        <p><strong>10. Aggregate Usage Statistics</strong></p>
        <p>We may collect Trovo App services usage information to improve function or user interface
          (&ldquo;<strong>UI</strong>&rdquo;), but will only use this information in an aggregated, anonymized fashion,
          and never in association with your name, email, or other personally identifying information.</p>
        <p><strong>11. Personally Identifying Information</strong></p>
        <p>You may choose to provide us with personally identifying information in order to participate in certain
          programs, activate features, or obtain other benefits but you will always be able to use the basic features of
          Trovo App without providing personally identifying information. If you provide us with personally
          identifying information then, notwithstanding any part of this policy, we may use that information to provide
          you with services or to improve the functioning or UI of Trovo App.</p>
        <p><strong>12. Confidentiality of User Information</strong></p>
        <p>Your personally identifying information will be kept strictly confidential and never provided to
          third-parties (other than in an aggregated, anonymized report such as the number of users per month). All
          Trovo App staff are bound by confidentiality agreements.</p>
        <p><strong>13. Reasons for collecting Personal Information</strong></p>
        <p>Personal information is data that can be used to identify you directly or indirectly, or to contact you. Our
          Privacy Policy covers all personal information that you voluntarily submit to us and that we obtain from our
          partners. This Privacy Policy does not apply to anonymized data, as it cannot be used to identify you.</p>
        <p>You may be asked to provide personal information anytime you are in contact with Trovo App.&nbsp;</p>
        <p>We may use your personal information in accordance with this Privacy Policy. Except as described in this
          Privacy Policy, Trovo App will not give, sell, rent or loan any personal Information to any third party.
        </p>
        <p><strong>14. The Information we collect</strong></p>
        <p>Any information that we collect is necessary either to provide you the services or to comply with legal
          regulations and legislation.</p>
        <p>You may refuse to provide us with some types of personal information requested, but in this case, we may not
          be able to ensure that the service works properly for you and that all services are available.</p>
        <p>Depending on the service you wish to use, we may collect the following types of information:</p>
        <ul>
          <li>
            <p>Full name, date of birth, nationality, place of birth, age, phone number, home address, and/or email.</p>
          </li>
          <li>
            <p>Tax number, passport number, driver&rsquo;s license details, national identity card details, photograph
              identification cards.</p>
          </li>
          <li>
            <p>Bank account details, payment card details, transaction history, trading data, and/or tax identification.
            </p>
          </li>
          <li>
            <p>Survey responses, information provided to our support team, public social networking posts,
              authentication data, security questions, user ID, other data collected via cookies and similar
              technologies.</p>
          </li>
        </ul>
        <p><strong>&nbsp;15. Use of Personal Information</strong></p>
        <p>The main purpose for collecting personal information is to provide you with a secure, and convenient
          experience. How this information can be used:</p>
        <ul>
          <li>
            <p>to create, develop, operate, deliver, and improve the services that you use on Trovo App App;</p>
          </li>
          <li>
            <p>to comply with the legal requirements and applicable legislation;</p>
          </li>
          <li>
            <p>for loss prevention and anti-fraud purposes.&nbsp;</p>
          </li>
        </ul>
        <p><strong>16. Information from external sources</strong></p>
        <p>We may obtain information about you from third-party sources as required or permitted by applicable law, such
          as public databases, credit bureaus, ID verification partners, resellers and channel partners, joint marketing
          partners, and social media platforms. We obtain such information to comply with our legal obligations, such as
          anti-money laundering and financing of terrorism laws. Our lawful basis for processing such data is compliance
          with legal obligations. In some cases, we may process additional data about you to ensure our Services are not
          used fraudulently or for other illicit activities.</p>
        <p><strong>17. Information Collected Automatically</strong>&nbsp;</p>
        <p>We may collect information about your computer, including where available your IP address, operating system,
          and browser type, for system administration, problem solving and service improvement. This is statistical data
          about our users&apos; actions and patterns and does not identify any individual.</p>
        <p><strong>18. Sharing Personal Information with Third-Parties </strong></p>
        <p>Trovo App will not share Personal Information of registered users with any third-party. However, we may
          share your information in the following circumstances:</p>
        <ul>
          <li>
            <p>With financial institutions who are our partners, to process payments you have authorized.</p>
          </li>
          <li>
            <p>With law enforcement or regulatory agencies, as may be required by law.</p>
          </li>
        </ul>
        <p><strong>&nbsp;19. Protection of Personal Information</strong></p>
        <p>We will not permit any third-parties to contact you directly on an unsolicited basis in relation to their own
          products or services. We do not sell, trade, or rent your personal identification information to others. You
          should never disclose your account password to unauthorized parties. We use certain security measures to help
          keep your personal information safe, but we cannot guarantee that these measures will stop any users trying to
          get around the privacy or security settings of &nbsp;Trovo App Application through unforeseen and/or
          illegal activity.</p>
        <p><strong>20. Request for &nbsp;Assistance</strong></p>
        <p>If you contact support, then as part of the assistance request, we may incidentally collect your personally
          identifying information, but we will endeavor to keep that information secure and confidential. Support
          information may be managed using a third-party service and further terms of that service may apply to your
          support request.</p>
        <p><strong>21. Queries</strong></p>
        <p>Any query with regards to this policy, the collection, use and disclosure of Personal Data by Trovo App
          operators or access to your Personal Data which is required by law to be disclosed should be directed to
          info@trovotech.io</p>
      </div>
            </div>
     </StyledMain>
 {/* </div> */}

      <Footer />
    </main>
  );
};

const StyledMain = styled.main`
  width: 80%;
  margin: auto;
  margin-top: 100px;
  font-family: "Montserrat";
`