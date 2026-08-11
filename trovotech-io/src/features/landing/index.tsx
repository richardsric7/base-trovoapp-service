import Hero from "./components/Hero";
// import Steps from "./components/Steps";
import Compliance from "./components/Compliance";
import Market from "./components/Market";
import CTA from "./components/CTA";
import HowItWorks from "./components/HowItWorks";

const LandingPage = () => {
  return (
    <main className="w-full flex flex-col">
      <Hero />
      {/* <Steps /> */}
      <HowItWorks />
      <Compliance />
      <Market />
      <CTA />
    </main>
  );
};

export default LandingPage;
