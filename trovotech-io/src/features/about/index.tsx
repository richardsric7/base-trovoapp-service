import AboutHero from "./components/AboutHero";
import WhatWeBuild from "./components/WhatWeBuild";
import ProblemWeAreSolving from "./components/ProblemWeAreSolving";
import Vision from "./components/Vision";
import Values from "./components/Values";
import RegulatoryAlignment from "./components/RegulatoryAlignment";
// import Team from "./components/Team";
import AboutCTA from "./components/AboutCTA";

const AboutPage = () => {
  return (
    <main className="w-full flex flex-col pt-20">
      <AboutHero />

      <ProblemWeAreSolving />
      <WhatWeBuild />
      <Vision />
      <Values />
      <RegulatoryAlignment />

      {/* <Team /> */}
      <AboutCTA />
    </main>
  );
};

export default AboutPage;

